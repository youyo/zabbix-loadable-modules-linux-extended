package main

import (
	"bufio"
	"errors"
	"os"

	"github.com/drael/GOnetstat"
	g2z "gopkg.in/cavaliercoder/g2z.v3"
)

func init() {
	g2z.RegisterUint64Item("linux_extended.netstat.count", "LISTEN,tcp", linuxExtendedNetstatCount)
	g2z.RegisterDiscoveryItem("linux_extended.swap.discovery", "", linuxExtendedSwapDiscovery)
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

func linuxExtendedSwapDiscovery(request *g2z.AgentRequest) (d g2z.DiscoveryData, err error) {
	var l []string
	f, err := os.Open("/proc/swaps")
	defer f.Close()
	if err != nil {
		return
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		l = append(l, s.Text())
	}
	if len(l) >= 2 {
		d = append(d, g2z.DiscoveryItem{
			"SWAP": "true",
		})
	}
	return
}

func main() {}
