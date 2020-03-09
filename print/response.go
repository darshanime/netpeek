package print

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/darshanime/netpeek/stats"
)

const outputTemplate string = `{{ .ReqMethod }} {{ .ReqURL }}
{{range $key, $value := .ReqHeaders }}{{ $key }}: {{ $value }}
{{end}}
{{ .ReqBody }}
~~~~~~~~~~~~
{{range $key, $value := .RespHeaders }}{{ $key }}: {{ $value }}
{{end}}
{{ .RespBody }}
####

{{range $pkt := .PacketInfo }}{{ if eq $pkt.Dir "client->server" }}-->{{ else }}<--{{ end }}{{ if $pkt.FIN }}FIN {{ else }}...{{ end }}{{ if $pkt.SYN }}SYN {{ else }}...{{ end }}{{ if $pkt.RST }}RST {{ else }}...{{ end }}{{ if $pkt.PSH }}PSH {{ else }}...{{ end }}{{ if $pkt.ACK }}ACK {{ else }}...{{ end }}{{ if $pkt.URG }}URG {{ else }}...{{ end }}{{ if $pkt.ECE }}ECE {{ else }}...{{ end }}{{ if $pkt.CWR }}CWR {{ else }}...{{ end }}{{ if $pkt.NS }}NS {{ else }}...{{ end }} | CL - {{ $pkt.CaptureLength }} | TS - {{ $pkt.Timestamp }}
{{end}}`

type templateContext struct {
	ReqMethod  string
	ReqURL     template.HTML
	ReqHeaders map[string]string
	ReqBody    template.HTML

	RespHeaders map[string]string
	RespBody    template.HTML

	PacketInfo []stats.PacketInfo
}

func Response(req *http.Request, resp *http.Response, pktInfo []stats.PacketInfo) {
	resultsTmpl, err := template.New("Meanings").Parse(outputTemplate)
	if err != nil {
		panic(fmt.Sprintf("cannot init response template - %s\n", err.Error()))
	}
	templateContext := getContext(req, resp, pktInfo)
	var output strings.Builder
	err = resultsTmpl.Execute(&output, templateContext)
	if err != nil {
		panic(fmt.Sprintf("cannot render response template - %s\n", err.Error()))
	}
	fmt.Println(output.String())
}

func getContext(req *http.Request, resp *http.Response, pktInfo []stats.PacketInfo) templateContext {
	tc := templateContext{
		ReqMethod:   req.Method,
		ReqURL:      template.HTML(req.URL.String()),
		ReqHeaders:  map[string]string{},
		RespHeaders: map[string]string{},
		PacketInfo:  pktInfo,
	}

	rcBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(fmt.Sprintf("cannot read req.Body - %s\n", err.Error()))
	}
	defer req.Body.Close()
	tc.ReqBody = template.HTML(rcBody)

	for key, val := range req.Header {
		tc.ReqHeaders[key] = strings.Join(val, ",")
	}

	rcBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("cannot read resp.Body - %s\n", err.Error()))
	}
	defer resp.Body.Close()
	if val, ok := resp.Header["Accept"]; ok && strings.HasPrefix(val[0], "text") {
		tc.RespBody = template.HTML(rcBody)
	} else {
		tc.RespBody = template.HTML("<REDACTED> non text content type")
	}

	for key, val := range req.Header {
		tc.RespHeaders[key] = strings.Join(val, ",")
	}
	return tc
}
