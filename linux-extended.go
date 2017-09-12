package main

import (
	"errors"

	"github.com/drael/GOnetstat"
	g2z "gopkg.in/cavaliercoder/g2z.v3"
)

func init() {
	g2z.RegisterUint64Item("linux_extended.netstat.count", "LISTEN,tcp", linuxExtendedNetstatCount)
}

func linuxExtendedNetstatCount(request *g2z.AgentRequest) (value uint64, err error) {
	state := "LISTEN"
	protocol := "tcp"
	switch len(request.Params) {
	case 2:
		if request.Params[0] != "" {
			state = request.Params[0]
		}
		if request.Params[1] != "" {
			protocol = request.Params[1]
		}
	case 1:
		if request.Params[0] != "" {
			state = request.Params[0]
		}
	default:
	}

	var d []GOnetstat.Process
	switch protocol {
	case "tcp":
		d = GOnetstat.Tcp()
	case "udp":
		d = GOnetstat.Udp()
	case "tcp6":
		d = GOnetstat.Tcp6()
	case "udp6":
		d = GOnetstat.Udp6()
	default:
		err = errors.New("Not match protocol")
		return
	}
	switch state {
	case
		"LISTEN",
		"ESTABLISHED",
		"TIME_WAIT",
		"LISTENING",
		"SYN_SENT",
		"SYN_RECEIVED",
		"FIN_WAIT_1",
		"FIN_WAIT_2",
		"CLOSE_WAIT",
		"CLOSING",
		"LAST_ACK",
		"CLOSED":
		value = countUp(d, state)
	default:
		err = errors.New("Not match State")
		return
	}
	return
}

func countUp(data []GOnetstat.Process, state string) (value uint64) {
	for _, v := range data {
		if v.State == state {
			value++
		}
	}
	return
}

func main() {}
