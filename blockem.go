package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
)

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
	//prepend an appropriate http protocol if url doesn't have one defined
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

func download_url(url string, url_channel chan []byte) {
	//download the URL and send the contents back down the channel
	txt := fmt.Sprintf("Downloading %s", url)
	fmt.Println(txt)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)

	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}
	url_channel <- contents
}

func create_processed_ip_slice(blacklist_urls []string) []string {

	processed_ips := make([]string, 0)
	url_channel := make(chan []byte)

	for _, url := range blacklist_urls {
		prefix_url := url_prefix(url)
		go download_url(prefix_url, url_channel)

		ip_addresses := match_ip(string(<-url_channel))
		for _, value := range ip_addresses {
			//Make sure we don't have this address already and that it is a parse-able IP address
			exists := ip_address_in_slice(value, processed_ips)
			if exists == false {
				fmt.Printf("Testing: %s\n", value)
				parse := net.ParseIP(value)
				if parse != nil {
					fmt.Printf("Valid: %s\n", value)
					processed_ips = append(processed_ips, value)
				} else {
					fmt.Printf("%s Not Valid\n", value)
				}
			}
		}

	}

	return processed_ips
}

var blacklist_urls_default = []string{
	"raw.githubusercontent.com/ktsaou/blocklist-ipsets/master/firehol_level1.netset",
	"https://raw.githubusercontent.com/stamparm/ipsum/master/levels/3.txt",
	"http://cinsscore.com/list/ci-badguys.txt",
}

func main() {
	ipaddress_file_out := flag.String("fileout", "badbots-blockem-ip-list.out", "Outfile for processed ip addresses")
	blacklist_urls_file_in := flag.String("blacklist_urls", "NotSet", "Input file containing comma seperated URLs")
	flag.Parse()

	var blacklist_urls = make([]string, 0)

	if *blacklist_urls_file_in == "NotSet" {
		blacklist_urls = blacklist_urls_default
	} else {

		file_bytes, err := ioutil.ReadFile(*blacklist_urls_file_in)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		blacklist_urls = strings.Split(string(file_bytes), "\n")
	}

	processed_ips := create_processed_ip_slice(blacklist_urls)

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
	file_err := file_handle.Close()
	if file_err != nil {
		fmt.Println(err)
		return
	}

}
