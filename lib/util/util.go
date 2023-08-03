package util

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"net"
	"os"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func GetInternalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				strip := ipnet.IP.String()
				if strings.HasPrefix(strip, "10.") || strings.HasPrefix(strip, "172.16.") {
					return strip, nil
				}
			}
		}
	}

	return "", errors.New("no internal ip found")
}

func Hostname() (string, error) {
	hostname, err := os.Hostname()
	return hostname, err
}

func IdcName() string {
	hostName, err := Hostname()
	if err != nil {
		return "default"
	} else {
		strList := strings.Split(hostName, ".")
		if len(strList) == 5 {
			return strList[2]
		} else {
			return "default"
		}

	}

}

func StructToJson(data interface{}) string {
	if result, err := json.Marshal(data); err != nil {
		return ""
	} else {
		return string(result)
	}

}

func ToJsonStr(payload any) string {
	bytes, _ := json.Marshal(payload)
	return string(bytes)
}
