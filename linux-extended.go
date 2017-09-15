package main

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/drael/GOnetstat"
	g2z "gopkg.in/cavaliercoder/g2z.v3"
)

func init() {
	g2z.RegisterUint64Item("linux_extended.netstat.count", "LISTEN,tcp", linuxExtendedNetstatCount)
	g2z.RegisterDiscoveryItem("linux_extended.swap.discovery", "", linuxExtendedSwapDiscovery)
	g2z.RegisterUint64Item("linux_extended.swap.size", "/swap,used", linuxExtendedSwapSize)
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
	f, err := os.Open("/proc/swaps")
	defer f.Close()
	if err != nil {
		return
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		if !strings.HasPrefix(s.Text(), "Filename") {
			device := strings.Fields(s.Text())[0]
			d = append(d, g2z.DiscoveryItem{
				"DEVICE": device,
			})
		}
	}
	return
}

func linuxExtendedSwapSize(request *g2z.AgentRequest) (value uint64, err error) {
	device := ""
	unit := "used"
	switch len(request.Params) {
	case 2:
		if request.Params[0] != "" {
			device = request.Params[0]
		}
		if request.Params[1] != "" {
			unit = request.Params[1]
		}
	case 1:
		if request.Params[0] != "" {
			device = request.Params[0]
		}
	default:
		err = errors.New("Undefine device")
		return
	}

	var l []string
	f, err := os.Open("/proc/swaps")
	defer f.Close()
	if err != nil {
		return
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		if strings.HasPrefix(s.Text(), device) {
			l = strings.Fields(s.Text())
		}
	}
	total, err := strconv.ParseUint(l[2], 10, 64)
	if err != nil {
		return
	}
	used, err := strconv.ParseUint(l[3], 10, 64)
	if err != nil {
		return
	}

	switch unit {
	case "used":
		value = used
	case "total":
		value = total
	case "free":
		value = total - used
	case "pfree":
		value = used / total * 100
	}
	return
}

func main() {}
