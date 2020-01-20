package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var blacklist_urls = []string{
	"raw.githubusercontent.com/ktsaou/blocklist-ipsets/master/firehol_level1.netset",
	"lists.blocklist.de/lists/all.txt",
	"https://raw.githubusercontent.com/stamparm/ipsum/master/levels/3.txt",
	"http://cinsscore.com/list/ci-badguys.txt",
}

func match_ip(pattern string) []string {
	//match ip addresses from string pattern and return slice of ips as string
	re := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(?:/\d{1,2}|)`)
	result := re.FindAllString(pattern, -1)
	return result
}

func ip_address_in_slice(ip_address string, ips []string) bool {
	//find a string in slice return boolean
	for _, val := range ips {
		if val == ip_address {
			fmt.Printf("True")
			return true
		}
	}
	fmt.Printf("False")
	return false
}

func url_prefix(url string) string {
	//prepend an appropriate http protocol if url doesn't have one
	is_http := strings.HasPrefix(url, "http://")
	is_https := strings.HasPrefix(url, "https://")
	prefix_url := ""

	if is_http == false && is_https == false {
		fmt.Println("Not Http or Https")
		var buffer bytes.Buffer
		buffer.WriteString("http://")
		buffer.WriteString(url)
		prefix_url = buffer.String()
	} else {
		prefix_url = url
	}
	return prefix_url
}

func main() {
	processed_ips := make([]string, 0)

	for _, url := range blacklist_urls {
		prefix_url := url_prefix(url)
		resp, err := http.Get(prefix_url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		ip_addresses := match_ip(string(contents))

		for _, value := range ip_addresses {
			processed_ips = append(processed_ips, value)
		}
	}

	fmt.Printf("Processed")
	for _, value := range processed_ips {
		fmt.Printf("%s\n", value)
	}

}