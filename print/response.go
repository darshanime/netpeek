package print

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const outputTemplate string = `{{ .ReqMethod }} {{ .ReqURL }}
{{range $key, $value := .ReqHeaders }}{{ $key }}: {{ $value }}
{{end}}
{{ .ReqBody }}
~~~~~~~~~~~~
{{range $key, $value := .RespHeaders }}{{ $key }}: {{ $value }}
{{end}}
{{ .RespBody }}`

type templateContext struct {
	ReqMethod  string
	ReqURL     string
	ReqHeaders map[string]string
	ReqBody    string

	RespHeaders map[string]string
	RespBody    string
}

func Response(req *http.Request, resp *http.Response) {
	resultsTmpl, err := template.New("Meanings").Parse(outputTemplate)
	if err != nil {
		panic(fmt.Sprintf("cannot init response template - %s\n", err.Error()))
	}
	templateContext := getContext(req, resp)
	err = resultsTmpl.Execute(os.Stdout, templateContext)
	if err != nil {
		panic(fmt.Sprintf("cannot render response template - %s\n", err.Error()))
	}
}

func getContext(req *http.Request, resp *http.Response) templateContext {
	tc := templateContext{
		ReqMethod:   req.Method,
		ReqURL:      req.URL.String(),
		ReqHeaders:  map[string]string{},
		RespHeaders: map[string]string{},
	}

	rcBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(fmt.Sprintf("cannot read req.Body - %s\n", err.Error()))
	}
	defer req.Body.Close()
	tc.ReqBody = string(rcBody)

	for key, val := range req.Header {
		tc.ReqHeaders[key] = strings.Join(val, ",")
	}

	rcBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("cannot read resp.Body - %s\n", err.Error()))
	}
	defer resp.Body.Close()
	tc.RespBody = string(rcBody)

	for key, val := range req.Header {
		tc.RespHeaders[key] = strings.Join(val, ",")
	}
	return tc
}
