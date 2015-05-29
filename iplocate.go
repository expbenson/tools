package main

import (
	"flag"
	. "fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
)

func httpGet(ip string) string {
	response, err := http.Get("http://ip.taobao.com/service/getIpInfo.php?ip=" + ip)
	if err != nil {
		Println("Http request error.")
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// handle error
	}

	return string(body)
}

func main() {
	flag.Usage = func() {
		Println("Usage: iplocate <ip or host>")
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		Println("Missing Argument(IP or Host).")
		os.Exit(1)
	}
	ipOrHost := args[0]
	var ip string
	if m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", ipOrHost); m {
		ip = ipOrHost
	} else {
		addrs, err := net.LookupHost(ipOrHost)
		if err != nil {
			Println("Error: ", err.Error())
			os.Exit(1)
		}
		if len(addrs) == 0 {
			Println("No IP")
			os.Exit(1)
		}
		ip = addrs[0]
	}
	response := httpGet(ip)
	// Println(response)
	js, err := simplejson.NewJson([]byte(response))
	if err != nil {
		Println("Error", err.Error())
		os.Exit(1)
	}
	code, _ := js.Get("code").Int()
	if code != 0 {
		Println("Error: code is not 0.")
	} else {
		data := js.Get("data")
		Printf("IP: %s\nCountry: %s\nArea: %s\nRegion: %s\nCity: %s\nISP: %s\n", data.Get("ip").MustString(), data.Get("country").MustString(), data.Get("area").MustString(), data.Get("region").MustString(), data.Get("city").MustString(), data.Get("isp").MustString())
	}
}
