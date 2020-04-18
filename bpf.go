package main

import "fmt"

func getBPFProgram() string {
	flagMap := getFlagMap()
	switch flagMap {
	case "0000":
		return "tcp"
	case "0001":
		return fmt.Sprintf("tcp and (dst host %s or src host %s)", *dHost, *dHost)
	case "0010":
		return fmt.Sprintf("tcp and (src host %s or dst host %s)", *sHost, *sHost)
	case "0100":
		return fmt.Sprintf("tcp and (dst port %s or src port %s)", *dPort, *dPort)
	case "1000":
		return fmt.Sprintf("tcp and (src port %s or dst port %s)", *sPort, *sPort)
	case "1001":
		return fmt.Sprintf("tcp and ((src port %s and dst host %s) or (dst port %s and src host %s))", *sHost, *dHost, *sHost, *dHost)
	case "1010":
		return fmt.Sprintf("tcp and ((src port %s and src host %s) or (dst port %s and dst host %s))", *sPort, *sHost, *sPort, *sHost)
	case "0011":
		return fmt.Sprintf("tcp and ((src host %s and dst host %s) or (dst host %s and src host %s))", *sHost, *dHost, *sHost, *dHost)
	case "0101":
		return fmt.Sprintf("tcp and ((dst port %s and dst host %s) or (src port %s and src host %s))", *dPort, *dHost, *dPort, *dHost)
	case "0110":
		return fmt.Sprintf("tcp and ((dst port %s and src host %s) or (src port %s and dst host %s))", *dPort, *sHost, *dPort, *sHost)
	case "1100":
		return fmt.Sprintf("tcp and ((src port %s and dst port %s) or (dst port %s and src port %s))", *sPort, *dPort, *sPort, *dPort)
	case "0111":
		return fmt.Sprintf("tcp and ((dst port %s and src host %s and dst host %s) or (src port %s and dst host %s and src host %s))", *dPort, *sHost, *dHost, *dPort, *sHost, *dHost)
	case "1011":
		return fmt.Sprintf("tcp and ((src port %s and src host %s and dst host %s) or (dst port %s and dst host %s and src host %s))", *sPort, *sHost, *dHost, *sPort, *sHost, *dHost)
	case "1101":
		return fmt.Sprintf("tcp and ((src port %s and dst port %s and dst host %s) or (dst port %s and src port %s and src host %s))", *sPort, *dPort, *dHost, *sPort, *dPort, *dHost)
	case "1110":
		return fmt.Sprintf("tcp and ((src port %s and dst port %s and src host %s) or (dst port %s and src port %s and dst host %s))", *sPort, *dPort, *sHost, *sPort, *dPort, *sHost)
	case "1111":
		return fmt.Sprintf("tcp and ((src port %s and dst port %s and src host %s and dst host %s) or (dst port %s and src port %s and dst host %s and src host %s))", *sPort, *dPort, *sHost, *dHost, *sPort, *dPort, *sHost, *dHost)
	default:
		panic("could not build bpf program")
	}
}

func getFlagMap() string {
	isPresent := func(opt *string) string {
		if *opt != "" {
			return "1"
		}
		return "0"
	}
	return isPresent(sPort) + isPresent(dPort) + isPresent(sHost) + isPresent(dHost)
}
