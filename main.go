package main

import (
	"bytes"
	"flag"
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
			return true
		}
	}
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

func download_url(url string, url_channel chan<- []byte) {
	//download the URL and send the contents back dwn the channel
	txt := fmt.Sprintf("Downloading %s", url)
	fmt.Println(txt)
	resp, err := http.Get(url)
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

	url_channel <- contents

}

func main() {
	ipaddress_file_out := flag.String("fileout", "ip-blacklist.txt", "A list of IP addresses")
	flag.Parse()

	processed_ips := make([]string, 0)
	url_channel := make(chan []byte)

	for _, url := range blacklist_urls {
		prefix_url := url_prefix(url)
		go download_url(prefix_url, url_channel)
		ip_addresses := match_ip(string(<-url_channel))
		for _, value := range ip_addresses {
			exists := ip_address_in_slice(value, processed_ips)
			if exists == false {
				processed_ips = append(processed_ips, value)
			}
		}
	}

	fmt.Printf("Processed")
	for _, value := range processed_ips {
		fmt.Printf("%s\n", value)
	}
	//open a file and output the addresses.
	file_handle, err := os.Create(*ipaddress_file_out)
	if err != nil {
		fmt.Println(err)
		file_handle.Close()
		return
	}
	//write the addresses
	for _, value := range processed_ips {
		fmt.Fprintln(file_handle, value)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	//close the file
	cat := file_handle.Close()
	if cat != nil {
		fmt.Println(err)
		return
	}

}
