package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func main() {
	url := flag.String("u", "example", "url to check")
	flag.Parse()

	resp, err := http.Get("http://data.iana.org/TLD/tlds-alpha-by-domain.txt")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	list, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(list), "\n")

	for i := 1; i < len(lines)-1; i++ {
		newUrl := fmt.Sprintf("%s.%s", *url, strings.ToLower(lines[i]))

		hostsChannel := make(chan []string)
		urlChannel := make(chan string)

		go func() {
			hosts, err := net.LookupHost(newUrl)
			if err == nil {
				hostsChannel <- hosts
				urlChannel <- newUrl
			}
			close(hostsChannel)
		}()

		for host := range hostsChannel {
			fmt.Printf("%s: %s\n", <-urlChannel, host)
		}
	}
}
