package cui

import (
	"testing"
)

func TestGetConnectionNameFromLine(t *testing.T) {
	line := "192.168.0.104:49369           172.217.160.174:80"
	output := getConnectionNameFromLine(line)
	if output != "->192.168.0.104:49369->172.217.160.174:80" {
		t.Fatalf("received unexpected output from getConnectionNameFromLine, expected: %s, got: %s\n", "->192.168.0.104:49369->172.217.160.174:8", output)
	}
}

func TestGetRequestNameFromLine(t *testing.T) {
	line := "1                             /"
	output := getRequestNameFromLine(line)
	if output != "" {
		t.Fatalf("received unexpected output from getConnectionNameFromLine, expected: %s, got: %s\n", "->192.168.0.104:49369->172.217.160.174:8", output)
	}
}
