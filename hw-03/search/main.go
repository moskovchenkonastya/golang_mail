package main

import (
	"net/http"
	_ "net/http/pprof"
	"strconv"
	//json "encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	//"regexp"
	"strings"
	// "log"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

var (
	_ *jlexer.Lexer
	_ *jwriter.Writer
)




// подключение профайлера
/*
{   "browsers":["Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36","LG-LX550 AU-MIC-LX550/2.0 MMP/2.0 Profile/MIDP-2.0 Configuration/CLDC-1.1","Mozilla/5.0 (Android; Linux armv7l; rv:10.0.1) Gecko/20100101 Firefox/10.0.1 Fennec/10.0.1","Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; MATBJS; rv:11.0) like Gecko"],
	"company":"Flashpoint",
	"country":"Dominican Republic",
	"email":"JonathanMorris@Muxo.edu",
	"job":"Programmer Analyst #{N}",
	"name":"Sharon Crawford",
	"phone":"176-88-49"}
*/

type Users struct {
	Browsers []string
	// company string
	// country string
	Email string
	// job string
	Name string
	// phone string
}


func easyjsonUnmarshalJSON(in *jlexer.Lexer, out *Users) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Users) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonUnmarshalJSON(&r, v)
	return r.Error()
}

func cpuhogger() {
	var acc uint64
	for {
		acc += 1
		if acc&1 == 0 {
			acc <<= 1
		}
	}
}

func main() {
	go http.ListenAndServe("0.0.0.0:8080", nil)
	cpuhogger()
}


//var rAndr = regexp.MustCompile("Android")
//var rMSIE = regexp.MustCompile("MSIE")

func FastSearch(out io.Writer) {

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	//var r = regexp.MustCompile("@")
	seenBrowsers := make(map[string]bool, 1000)
	uniqueBrowsers := 0
	foundUsers := ""

	lines := strings.Split(string(fileContents), "\n")

	var user Users

	//users := make([]map[string]interface{}, 0)
	for i, line := range lines {
		//user.UnmarshalJSON := make(map[string]interface{})
		// fmt.Printf("%v %v\n", err, line)
		err := user.UnmarshalJSON([]byte(line))
		if err != nil {
			panic(err)
		}
		//users = append(users, user)
	//}

	//for i, user := range users {

		isAndroid := false
		isMSIE := false

		//browsers, ok := user["browsers"].([]interface{})
		//if !ok {
			// log.Println("cant cast browsers")
		//	continue
		//}

		for _, browser := range user.Browsers {
			//browser, ok := browserRaw .(string)
			//if !ok {
				// log.Println("cant cast browser to string")
			//	continue
			//}
			tmpIsAndroid := strings.Contains(browser, "Android")
			tmpIsMSIE := strings.Contains(browser, "MSIE")

			if tmpIsAndroid || tmpIsMSIE {
				
				if tmpIsAndroid {
					isAndroid = true
				}

				if tmpIsMSIE {
					isMSIE = true
				}

				_, ok := seenBrowsers[browser]
				if !ok {
					// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
					seenBrowsers[browser] = true
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.Replace(user.Email, "@", " [at] ", -1)
		//foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
		foundUsers += "[" + strconv.Itoa(i) +  "] " + user.Name + " <" + email + ">\n"
	}

	fmt.Fprintln(out, "found users:\n" + foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
