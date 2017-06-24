package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var urlTests = []struct {
	Url   string
	Valid bool
}{
	{"not an url", false},
	{"http://not.google.form.url.com/somepath", false},
	{"https://docs.google.com/forms/d/1Z-_ewPdnRmbqYjSTk3_hEovJw-CCAOlHm8dOKiAOfvc/edit", false},                                 // private edit url
	{"https://docs.google.com/a/balsamiq.com/forms/d/e/1FAIpQLSfnPMNxYkElNjWoC-sAdxza2JXACLxQRmbPvfE0nJoBcFMLaw/viewform", true}, // company account
	{"https://docs.google.com/forms/d/e/1FAIpQLSdgHnrWlSYKOa6x26roquuVg5Z--4kUvnHh16rEPjgRpulJNQ/viewform", true},                // personal account
}

func TestURL(t *testing.T) {
	for _, urlTest := range urlTests {
		err := CheckURL(urlTest.Url)
		if err == nil != urlTest.Valid {
			t.Error(
				"URL",
				urlTest.Url,
				"EXPECTED",
				urlTest.Valid,
				"GOT",
				err == nil,
			)
		}
	}
}

var formTests = []struct {
	Url    string
	Status int
	Body   string
}{
	{
		Url:    "https://docs.google.com/forms/d/e/1FAIpQLSdjZK2A_L9zUprCxOJdvvjexmNmxwmZCN6vMmTXAIZhJqUg3w/viewform",
		Status: http.StatusOK,
		Body: `
		{
			"title": "Short Test",
			"header": "Short Test",
			"desc": "description",
			"path": "/forms",
			"action": "e/1FAIpQLSdjZK2A_L9zUprCxOJdvvjexmNmxwmZCN6vMmTXAIZhJqUg3w",
			"fields": [
				{
					"id": 10920109,
					"label": "Short",
					"desc": "Short Description",
					"typeid": 1,
					"widgets": [
						{
							"id": "499896788",
							"required": false
						}
					]
				}
			]
		}
		`,
	},
}

func TestFormHandler(t *testing.T) {
	for _, formTest := range formTests {
		req, err := http.NewRequest("GET", "/formdress?url="+formTest.Url, nil)
		if err != nil {
			log.Fatal(err)
		}

		rr := httptest.NewRecorder()

		FormDressHandler(rr, req)

		res := rr.Result()
		if res.StatusCode != formTest.Status {
			t.Error(
				"Status EXPECTED",
				formTest.Status,
				"GOT",
				res.StatusCode,
			)
		}

		resJSON := &map[string]interface{}{}
		body, _ := ioutil.ReadAll(res.Body)
		if err := json.Unmarshal(body, resJSON); err != nil {
			log.Fatal(err)
		}

		formJSON := &map[string]interface{}{}
		if err := json.Unmarshal([]byte(formTest.Body), formJSON); err != nil {
			log.Fatal(err)
		}

		// Make sure the random attribute value match
		(*formJSON)["fbzx"] = (*resJSON)["fbzx"]

		if !reflect.DeepEqual(resJSON, formJSON) {
			resString, _ := json.MarshalIndent(resJSON, "", "  ")
			formString, _ := json.MarshalIndent(formJSON, "", "  ")
			t.Errorf("FOR: %s\nEXPECTED\n%s\nGOT\n%s", formTest.Url, formString, resString)
		}
	}
}
