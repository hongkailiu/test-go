package info

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// Info represents information
type Info struct {
	Version string    `json:"version"`
	Ips     []string  `json:"ips"`
	Now     time.Time `json:"now"`
}

// GetInfo returns the required information
func GetInfo(version string) *Info {

	i := Info{}
	i.Version = version
	i.Ips = getIPs()
	i.Now = time.Now()
	return &i
}

func getIPs() []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Error(err.Error())
		return []string{err.Error()}
	}
	// handle err
	var result []string
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			result = append(result, err.Error())
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if len(ip) != 0 {
				result = append(result, ip.String())
			}
		}
	}
	return result
}
