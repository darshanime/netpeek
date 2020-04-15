package main

import "testing"

func TestGetFlagMap(t *testing.T) {
	tests := []struct {
		sPort           string
		dPort           string
		sHost           string
		dHost           string
		expectedFlagMap string
	}{
		{expectedFlagMap: "0000"},
		{expectedFlagMap: "0100", dPort: "8080"},
		{expectedFlagMap: "1000", sPort: "80"},
		{expectedFlagMap: "1100", sPort: "80", dPort: "8080"},
		{expectedFlagMap: "0010", sHost: "127.0.0.1"},
		{expectedFlagMap: "0001", dHost: "127.0.0.1"},
		{expectedFlagMap: "1010", sPort: "80", sHost: "127.0.0.1"},
		{expectedFlagMap: "1110", dPort: "81", sPort: "82", sHost: "127.0.0.1"},
		{expectedFlagMap: "1111", dPort: "81", sPort: "82", sHost: "127.0.0.1", dHost: "127.0.0.2"},
	}
	for ti, tc := range tests {
		sPort, dPort, sHost, dHost = &tc.sPort, &tc.dPort, &tc.sHost, &tc.dHost
		if receivedBPF := getFlagMap(); receivedBPF != tc.expectedFlagMap {
			t.Errorf("test %d: incorrect flag map, expected: %s, got: %s", ti, tc.expectedFlagMap, receivedBPF)
		}
	}
}

func TestGetBPFProgram(t *testing.T) {
	tests := []struct {
		sPort           string
		dPort           string
		sHost           string
		dHost           string
		expectedFlagMap string
	}{
		{expectedFlagMap: "tcp"},
		{expectedFlagMap: "tcp and (dst port 8080 or src port 8080)", dPort: "8080"},
		{expectedFlagMap: "tcp and (src port 80 or dst port 80)", sPort: "80"},
		{expectedFlagMap: "tcp and ((src port 80 and dst port 8080) or (dst port 80 and src port 8080))", sPort: "80", dPort: "8080"},
		{expectedFlagMap: "tcp and (src host 127.0.0.1 or dst host 127.0.0.1)", sHost: "127.0.0.1"},
		{expectedFlagMap: "tcp and (dst host 127.0.0.1 or src host 127.0.0.1)", dHost: "127.0.0.1"},
		{expectedFlagMap: "tcp and ((src port 80 and src host 127.0.0.1) or (dst port 80 and dst host 127.0.0.1))", sPort: "80", sHost: "127.0.0.1"},
		{expectedFlagMap: "tcp and ((src port 82 and dst port 81 and src host 127.0.0.1) or (dst port 82 and src port 81 and dst host 127.0.0.1))", dPort: "81", sPort: "82", sHost: "127.0.0.1"},
		{expectedFlagMap: "tcp and ((src port 82 and dst port 81 and src host 127.0.0.1 and dst host 127.0.0.2) or (dst port 82 and src port 81 and dst host 127.0.0.1 and src host 127.0.0.2))", dPort: "81", sPort: "82", sHost: "127.0.0.1", dHost: "127.0.0.2"},
	}
	for ti, tc := range tests {
		sPort, dPort, sHost, dHost = &tc.sPort, &tc.dPort, &tc.sHost, &tc.dHost
		if receivedBPF := getBPFProgram(); receivedBPF != tc.expectedFlagMap {
			t.Errorf("test %d: incorrect bpf, expected: %s, got: %s", ti, tc.expectedFlagMap, receivedBPF)
		}
	}
}
