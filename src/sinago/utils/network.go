package utils

import (
	"net"
	"time"
	"strconv"
	"net/http"
	"fmt"
)

func IsReachable(protocol string,ip string, port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout(protocol, ip+":"+strconv.Itoa(port), timeout)
	if err != nil {
		return false
	} else {
		defer conn.Close()
		return true
	}
}

func GetInteralIP() string {
	var ip string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		ip = "Unknown Address"
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return ip
}

func HeaderToArray(header http.Header) (res []string) {
	for name, values := range header {
		for _, value := range values {
			res = append(res, fmt.Sprintf("%s: %s", name, value))
		}
	}
	return
}
