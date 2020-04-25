package print

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/gopacket"
)

func RequestToString(req *http.Request) string {
	var str strings.Builder
	str.WriteString(req.Method + " " + req.Host + req.URL.String() + "\n")
	for key, val := range req.Header {
		str.WriteString(key + ": " + strings.Join(val, ",") + "\n")
	}
	str.WriteString("\n")
	rcBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(fmt.Sprintf("cannot read resp.Body - %s\n", err.Error()))
	}
	defer req.Body.Close()
	str.Write(rcBody)
	return str.String()
}

func ConnDir(netflow, tcpflow gopacket.Flow) string {
	return netflow.Src().String() + ":" + tcpflow.Src().String() + "-->" + netflow.Dst().String() + ":" + tcpflow.Dst().String()
}
