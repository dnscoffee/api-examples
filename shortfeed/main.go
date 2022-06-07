package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const defaultEndpoint = "https://api.dns.coffee/api/v0/feeds/domains/old/date"
const dateFmt = "2006-01-02"

var (
	endpoint    = flag.String("endpoint", defaultEndpoint, "endpoint to use")
	apiKey      = flag.String("apikey", "", "DNS coffee API Key. REQUIRED")
	date        = flag.String("date", time.Now().Format(dateFmt), "the date to fetch YYYY-MM-DD")
	minLen      = flag.Uint("len", 5, "minimum domain length to return")
	ignoredTLDs = flag.String("ingore-tlds", "", "comma separated list of TLDs to ignore")
)

// JSONResponse JSON-API root data object
type JSONResponse struct {
	Error *JSONError  `json:"error,omitempty"`
	Data  *FeedDomain `json:"data,omitempty"`
}

// JSONError JSON-API error object
type JSONError struct {
	ID     string `json:"-"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type FeedDomain struct {
	Change  string    `json:"change,omitempty"`
	Date    time.Time `json:"date"`
	Domains []string  `json:"domains"`
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var ignoredTLDMap map[string]bool

func getFeed(dateStr string) []string {
	fullPath := fmt.Sprintf("%s/%s", *endpoint, dateStr)
	req, err := http.NewRequest(http.MethodGet, fullPath, nil)
	check(err)
	req.Header.Set("X-API-Key", *apiKey)
	resp, err := http.DefaultClient.Do(req)
	check(err)
	defer resp.Body.Close()
	var data JSONResponse
	if resp.StatusCode == 200 {
		// good
		err = json.NewDecoder(resp.Body).Decode(&data)
		check(err)
	} else {
		// bad
		log.Printf("HTTP request error, response code: %d", resp.StatusCode)
		err = json.NewDecoder(resp.Body).Decode(&data)
		check(err)
		if data.Error != nil {
			log.Printf("Response Error: %+v", data.Error)
		}
		return nil
	}
	return data.Data.Domains
}

func filter(size uint, domains []string) {
	for _, domain := range domains {
		l := uint(len(domain))
		if l <= size {
			tld := domain[strings.LastIndex(domain, ".")+1:]
			if !ignoredTLDMap[tld] {
				fmt.Printf("%d\t%s\n", len(domain), domain)
			}
		}
	}
}

func main() {
	flag.Parse()
	if len(*apiKey) < 10 {
		log.Fatal("invalid API Key")
	}
	ignoredTLDMap = make(map[string]bool)
	for _, tld := range strings.Split(strings.ToUpper(*ignoredTLDs), ",") {
		ignoredTLDMap[tld] = true
	}
	domains := getFeed(*date)
	//log.Printf("Got %d domains", len(domains))
	filter(*minLen, domains)

}
