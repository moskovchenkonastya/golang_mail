package main

import (
	"net/http"
	"sort"
	"strings"
	"strconv"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"encoding/xml"
	"io"
	"time"
	"net/http/httptest"
	
)
type TestCase struct{
	Request      *SearchRequest
	Response     *SearchResponse
	IsError      bool
	AccessToken  string
}

const (
	CorrectAccessToken = "d41d8cd98f00b204e9800998ecf8427e"
)

const xmlData = "dataset.xml"

func SearchParams(r *http.Request) (*SearchRequest, error) {

	var err error = nil 
	params := &SearchRequest{}
	params.Limit, err = strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		params.Limit = 25
	}

	params.Offset, err = strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		params.Offset = 0;
	}

	params.Query = r.FormValue("query")
	params.OrderField = r.FormValue("order_field")
	if params.OrderField != "" && strings.Index("Id_Name_Age", params.OrderField) < 0 {
		err = fmt.Errorf(ErrorBadOrderField)
		return nil, err
	}

	params.OrderBy, err = strconv.Atoi(r.FormValue("order_by"))
	if err != nil {
		params.OrderBy = 0
	}
	return params, err
}

func readData() *[]User {

	readDataXML, err := ioutil.ReadFile(xmlData)
	if err != nil {
		err = fmt.Errorf("Can not read xmlData")
		panic(err)
	}

	input := bytes.NewReader(readDataXML)
	decoder := xml.NewDecoder(input)

	users := make([]User, 0)
	var user User
	var firstName, lastName string

	for {

		tok, tokenErr := decoder.Token()
		if tokenErr != nil && tokenErr != io.EOF {
			fmt.Println("error happend", tokenErr)
			break
		} else if tokenErr == io.EOF {
			break
		}
		if tok == nil {
			fmt.Println("t is nil break")
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			switch tok.Name.Local {
				case "id":
					if err := decoder.DecodeElement(&user.Id, &tok); err != nil {
						fmt.Println("error happend", err)
					}	
				case "age":
					if err := decoder.DecodeElement(&user.Age, &tok); err != nil {
						fmt.Println("error happend", err)
					}
				case "first_name":
					if err := decoder.DecodeElement(&firstName, &tok); err != nil {
						fmt.Println("error happend", err)
					}				
				case "last_name":
					if err := decoder.DecodeElement(&lastName, &tok); err != nil {
						fmt.Println("error happend", err)
					}								
				case "about":
					if err := decoder.DecodeElement(&user.About, &tok); err != nil {
						fmt.Println("error happend", err)
					}								
				case "gender":
					if err := decoder.DecodeElement(&user.Gender, &tok); err != nil {
						fmt.Println("error happend", err)
					}		
				case "row":
					user = User{}
					firstName = ""
					lastName = ""
			}

		case xml.EndElement:
			if tok.Name.Local == "row" {
				user.Name = firstName + " " + lastName
				users = append(users, user)
			}
		}
	}
	return &users
}

func sortData(params *SearchRequest, data []User) *[]User {

	filtrData := make([]User, 0)
	cur := 1

	field := "Name"
	if params.OrderField != "" {
		field = params.OrderField
	}

	if params.OrderBy != OrderByAsIs {
		sort.Slice(data, func(i, j int) bool {
			switch field {
			case "Id":
				return (data[j].Id - data[i].Id) * params.OrderBy >= 0
			case "Name":
				return strings.Compare(data[i].Name, data[j].Name) * params.OrderBy < 0
			case "Age":
				return (data[j].Age - data[i].Age) * params.OrderBy >= 0
			}
			return true
		})
	}

	for _, user := range data {

		skipValue := false
		if params.Query != "" {
			if strings.Index(user.Name + "_" + user.About, params.Query) < 0 {
				skipValue = true
			} 
		}

		if !skipValue {
			cur++
			
			if params.Offset + 1 < cur {
				filtrData = append(filtrData, user)
			}

			if (params.Limit + params.Offset) < cur {
				break
			}
		}
	}
	return &filtrData

}

func SearchServer(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("AccessToken") != CorrectAccessToken {

		w.WriteHeader(http.StatusUnauthorized)

	} else {
		//var err error
		params, err := SearchParams(r)
		var result []byte

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errResp := SearchErrorResponse{}
			errResp.Error = fmt.Sprint(err)
			result, _  = json.Marshal(errResp)
				
		} else {

			data := *readData()
			filtrData := *sortData(params, data)
			result, _  = json.Marshal(filtrData)
		}

		w.Write([]byte(result))
	}
	
}

func SearchServerSleep(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 6)
}

func SearchServerStattusInternal(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func SearchServerLocal(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(187)
	w.Header().Add("Location", "")
}

func SearchServerBadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"Id\":''2'2}"))	
}

func SearchServerBadJson(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{\"Id\":''2'2}"))
}

func SearchServerBadReq(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("{\"Error\":\"TEST ERROR MESSAGE\"}"))
}

func SearchServerErrorServerUnknown(w http.ResponseWriter, r *http.Request) {
	item := TestCase{
	  AccessToken: "d41d8cd98f00b204e9800998ecf8427e",
	  Request: &SearchRequest {
		Limit: 1,
		Offset: 0,
		OrderBy: 0,
		OrderField: "",
		Query: "",
	  },
	  IsError: true,
	  Response: &SearchResponse{
		Users: 
		  []User{},
	  },

	}
	errorUnknownSever := httptest.NewServer(http.HandlerFunc(SearchServerErrorServerUnknown))
	searchClient := new(SearchClient)
	searchClient.AccessToken = ""
	searchClient.URL = errorUnknownSever.URL
	_, err := searchClient.FindUsers(item.Request)
	if err != nil {
	  errorUnknownSever.Close()
	} else {
	  errorUnknownSever.Close()
	}
	return
}

func SearchServerBadJs(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
}


