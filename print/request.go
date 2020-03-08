package print

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const requestTemplate string = `{{ .Method }} {{ .URL }}
{{range $key, $value := .Headers }}{{ $key }}: {{ $value }}
{{end}}{{ .Body }}
~~~~~~~~~~~~`

type requestContext struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

func Request(req *http.Request) {
	resultsTmpl, err := template.New("Meanings").Parse(requestTemplate)
	if err != nil {
		panic(fmt.Sprintf("cannot init template - %s\n", err.Error()))
	}
	reqContext := getRequestContext(req)
	err = resultsTmpl.Execute(os.Stdout, reqContext)
	if err != nil {
		panic(fmt.Sprintf("cannot render template - %s\n", err.Error()))
	}
}

func getRequestContext(req *http.Request) requestContext {
	defer req.Body.Close()
	rc := requestContext{Method: req.Method, URL: req.Host + req.URL.RequestURI(), Headers: map[string]string{}}
	rcBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(fmt.Sprintf("cannot read req.Body - %s\n", err.Error()))
	}
	rc.Body = string(rcBody)
	for key, val := range req.Header {
		rc.Headers[key] = strings.Join(val, ",")
	}
	return rc
}
