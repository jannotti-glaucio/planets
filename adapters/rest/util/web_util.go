package util

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/go-playground/form/v4"
)

//Parser reponsable for sing a struct
func Parser(request *http.Request, data interface{}) {
	if request != nil && request.Body != nil {
		request.ParseMultipartForm(128)

		//If request form is empty try parse json
		if len(request.Form) < 1 {
			decodeJSON(request.Body, data)
		} else {
			decodeForm(request.Form, data)
		}
	}
}

func decodeJSON(r io.Reader, obj interface{}) (err error) {
	decoder := json.NewDecoder(r)
	decoder.UseNumber()
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&obj)
	return
}

func decodeForm(data url.Values, obj interface{}) (err error) {
	decoder := form.NewDecoder()
	decoder.SetNamespacePrefix("[")
	decoder.SetNamespaceSuffix("]")
	decoder.SetTagName("json")
	err = decoder.Decode(&obj, data)
	return
}
