package main

import (
	"net/http"
	"testing"
	"net/http/httptest"
	
)


func compareResult(a, b []User, fieldComp string) bool {
	
		if len(a) != len(b){
			return false
		}
		result := true
		for i := range a {
			switch fieldComp {
			case "Id":
				if (a)[i].Id != (b)[i].Id{
					result = false
				}
			case "Age":
				if (a)[i].Age  != (b)[i].Age {
					result = false
				}
			}
		}
		return result
	}
	

func Cases() []TestCase {

	cases := []TestCase{
		TestCase {
			AccessToken: "d41d8cd98f00b204e9800998ecf8427e",
			Request: &SearchRequest {
				Limit: 5,
				Offset: 0,
				OrderBy: 0,
				OrderField: "",
				Query: "",
			},
			IsError: false,
			Response: &SearchResponse {
				Users: 
					[]User{
						User{
							Id: 0,
						},
						User{
							Id: 1,
						},
						User{
							Id: 2,
						},
						User{
							Id: 3,
						},
						User{
							Id: 4,
						},
					},
				NextPage: true,
			},
		},
		TestCase {
			AccessToken: "d41d8cd98f00b204e9800998ecf8427e",
			Request: &SearchRequest {
				Limit: 30,
				Offset: 0,
				OrderBy: 1,
				OrderField: "Id",
				Query: "",	
			},
			IsError: false,
			Response: &SearchResponse {
				Users: []User {User{Id:0,}, User{Id:1,}, User{Id:2,}, User{Id:3,}, User{Id:4,}, User{Id:5,}, User{Id:6,}, User{Id:7,}, 
				User{Id:8,}, User{Id:9,}, User{Id:10,}, User{Id:11,}, User{Id:12,}, User{Id:13,}, User{Id:14,}, 
				User{Id:15,}, User{Id:16,}, User{Id:17,}, User{Id:18,}, User{Id:19,}, User{Id:20,}, User{Id:21,}, 
				User{Id:22,}, User{Id:23,}, User{Id:24,}, },
				NextPage: true,

			},
		},	
		TestCase {
			AccessToken: "d41d8cd98fwqe00b204e9800",
			Request: &SearchRequest {
				Limit: 30,
				Offset: 0,
				OrderBy: 1,
				OrderField: "Id",
				Query: "",
			},
			IsError:  true,
			Response: &SearchResponse {
				Users: 
					[]User{},
			},
		},
		
		TestCase {
			AccessToken: "d41d8cd98f00b204e9800998ecf8427e",
			Request: &SearchRequest {
				Limit: 5,
				Offset: 0,
				OrderBy: -1,
				OrderField: "Id",
				Query: "",
			},
			IsError: false,
			Response: &SearchResponse{
				Users: 
					[]User{
						User{
							Id: 34,
						},
						User{
							Id: 33,
						},
						User{
							Id: 32,
						},
						User{
							Id: 31,
						},
						User{
							Id: 30,
						},
					},
					NextPage: true,					
				},
		},

		TestCase {
			AccessToken: "d41d8cd98f00b204e9800998ecf8427e",
			Request: &SearchRequest {
				Limit: 5,
				Offset: 0,
				OrderBy: 1,
				OrderField: "",
				Query: "",
			},
			IsError: false,
			Response: &SearchResponse{
			Users: 
				[]User{
					User{
						Id: 15,
					},
					User{
						Id: 16,
					},
					User{
						Id: 19,
					},
					User{
						Id: 22,
					},
					User{
						Id: 5,
					},
				},
				NextPage: true,
			},
		},
		
		TestCase {
			AccessToken: "d41d8cd98f00b204e9800998ecf8427e",
			Request: &SearchRequest {
				Limit: 1,
				Offset: 0,
				OrderBy: 1,
				OrderField: "Warreho",
				Query: "Boyd Wolf",
			},
			IsError: true,
			Response: &SearchResponse{
				Users: 
					[]User{},
					NextPage: false,
				},
		},	
		TestCase {
			AccessToken: "d41d8cd98f00b204e9800998ecf8427e",
			Request: &SearchRequest {
				Limit: 1,
				Offset: -1,
				OrderBy: 0,
				OrderField: "",
				Query: "",
			},
			IsError: true,
			Response: &SearchResponse{
				Users: 
					[]User{},
					NextPage: false,
				},
		},	
		TestCase {
			AccessToken: "d41d8cd98f00b204e9800998ecf8427e",
			Request: &SearchRequest {
				Limit: 1,
				Offset: 0,
				OrderBy: 0,
				OrderField: "",
				Query: "new",
			},
			IsError: false,
			Response: &SearchResponse{
				Users: 
					[]User{},
					NextPage: false,
				},
		},	
		TestCase {
			AccessToken: "d41d8cd98f00b204e9800998ecf8427e",
			Request: &SearchRequest {
				Limit: 1,
				Offset: 0,
				OrderBy: 0,
				OrderField: "",
				Query: "",
			},
			IsError: false,
			Response: &SearchResponse{
				Users: 
					[]User{User{Id:0}},
					NextPage: true,
				},
		},
		TestCase {
			AccessToken: "d41d8cd98f00b204e9800998ecf8427e",
			Request: &SearchRequest {
				Limit: 1,
				Offset: 0,
				OrderBy: 1,
				OrderField: "Id",
				Query: "Boyd Wolf",
			},
			IsError: false,
			Response: &SearchResponse{
				Users: 
					[]User{User{Id:0,},  },
					NextPage: false,
				},
		},
		TestCase {
			AccessToken: "d41d8cd98f00b204e9800998ecf8427e",
			Request: &SearchRequest {
				Limit: 1,
				Offset: 0,
				OrderBy: 1,
				OrderField: "Id",
				Query: "",
			},
			IsError: false,
			Response: &SearchResponse{
				Users: 
					[]User{User{Id:0,},  },
					NextPage: true,
				},
		},
		TestCase {
			AccessToken: "d41d8cd98f00b204e9800998ecf8427e",
			Request: &SearchRequest {
				Limit: 1,
				Offset: 0,
				OrderBy: -1,
				OrderField: "Age",
				Query: "",
			},
			IsError: false,
			Response: &SearchResponse{
				Users: 
					[]User{User{Id:32,},  },
					NextPage: true,
				},
		},
		TestCase {
			AccessToken: "d41d8cd98f00b204e9800998ecf8427e",
			Request: &SearchRequest {
				Limit: -1,
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
		},	

	}	
	
	return cases
}



func funcIsError(t *testing.T, res *SearchResponse, err error, caseIter int, i TestCase ) {

	if err != nil && !i.IsError {
		t.Errorf("[%d] unexpected error: %#v", caseIter, err)
	}
	if err == nil && i.IsError {
		t.Errorf("[%d] expected error, got nil", caseIter)
	}
	if err == nil && !compareResult(i.Response.Users, res.Users, "Id") {
		t.Errorf("\n[%d] wrong result, expected \n\t%#v, got \n\t%#v", caseIter, i.Response, res.Users)
	}
	if err == nil && i.Response.NextPage != res.NextPage {
		t.Errorf("\n[%d] Incorrect NextPage \n\t%#v, got \n\t%#v", caseIter, i.Response, res.Users)
	}

}
// сюда писать тесты

func TestClientServer(t *testing.T) {

	cases := Cases()
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	
	
	for caseIter, i := range cases {
		sc := &SearchClient{}
		sc.AccessToken = i.AccessToken
		sc.URL = ts.URL
		
		res, err := sc.FindUsers(i.Request)
		funcIsError(t, res, err, caseIter, i)
	}
	ts.Close()

	i := TestCase {
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
	

	searchServerSleep := httptest.NewServer(http.HandlerFunc(SearchServerSleep))
	sc := &SearchClient{}
	sc.AccessToken = i.AccessToken
	sc.URL = searchServerSleep.URL
	res, err := sc.FindUsers(i.Request)
	funcIsError(t, res, err, 0, i)
	searchServerSleep.Close()

	searchServerStattusInternal := httptest.NewServer(http.HandlerFunc(SearchServerStattusInternal))
	sc.URL = searchServerStattusInternal.URL
	res, err = sc.FindUsers(i.Request)
	funcIsError(t, res, err, 0, i)
	searchServerStattusInternal.Close()

	searchServerLocal := httptest.NewServer(http.HandlerFunc(SearchServerLocal))
	sc.URL = searchServerLocal.URL
	res, err = sc.FindUsers(i.Request)
	funcIsError(t, res, err, 0, i)
	searchServerLocal.Close()

	searchServerBadRequest := httptest.NewServer(http.HandlerFunc(SearchServerBadRequest))
	sc.URL = searchServerBadRequest.URL
	res, err = sc.FindUsers(i.Request)
	funcIsError(t, res, err, 0, i)
	searchServerBadRequest.Close()

	searchServerBadJson := httptest.NewServer(http.HandlerFunc(SearchServerBadJson))
	sc.URL = searchServerBadJson.URL
	res, err = sc.FindUsers(i.Request)
	funcIsError(t, res, err, 0, i)
	searchServerBadJson.Close()

	searchServerBadJs := httptest.NewServer(http.HandlerFunc(SearchServerBadJs))
	sc.URL = searchServerBadJs.URL
	res, err = sc.FindUsers(i.Request)
	funcIsError(t, res, err, 0, i)
	searchServerBadJs.Close()

	searchServerBadReq := httptest.NewServer(http.HandlerFunc(SearchServerBadReq))
	sc.URL = searchServerBadReq.URL
	res, err = sc.FindUsers(i.Request)
	funcIsError(t, res, err, 0, i)
	searchServerBadReq.Close()

	searchServerErrorServerUnknown := httptest.NewServer(http.HandlerFunc(SearchServerErrorServerUnknown))
	sc.URL = searchServerErrorServerUnknown.URL
	res, err = sc.FindUsers(i.Request)
	funcIsError(t, res, err, 0, i)
	searchServerErrorServerUnknown.Close()
	
}