package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"net/url"
	"net/http/httputil"
	smpl "github.com/benmeehan/go-smpl"
)

type CompileRequest struct {
	Code     string `json:"code"`
	Input    string `json:"input"`
	Language string `json:"language"`
}

type LoadBalancer struct {
	LanguageServers map[string][]string
	NextServerIndex map[string]int
}

func (lb *LoadBalancer) handleCompile(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal server error")
		log.Printf("Failed to read request body: %v", err)
		return
	}

	// Create a new buffer and set it as the request body
	requestBuffer := bytes.NewBuffer(body)
	r.Body = ioutil.NopCloser(requestBuffer)

	var compileReq CompileRequest
	err = json.Unmarshal(body, &compileReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error occurred during parsing of the JSON request body")
		log.Printf("Failed to parse JSON request body: %v", err)
		return
	}

	// Get the list of servers based on the language
	servers := lb.LanguageServers[compileReq.Language]
	if servers == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid language")
		log.Printf("Invalid language: %s", compileReq.Language)
		return
	}

	// Get the next server index from the map
	nextServerIndex := lb.NextServerIndex[compileReq.Language]
	if nextServerIndex >= len(servers) {
		nextServerIndex = 0 // Reset to the first server if index exceeds the server list size
	}
	serverURL := servers[nextServerIndex]

	// Update the next server index for the language
	lb.NextServerIndex[compileReq.Language] = (nextServerIndex + 1) % len(servers)

	// Parse the server URL
	u, err := url.Parse(serverURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal server error")
		log.Printf("Failed to parse server URL: %v", err)
		return
	}
	u.Path = "/compile"

	// Create a reverse proxy with a custom Director function
	proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = u.Scheme
		req.URL.Host = u.Host
		req.URL.Path = u.Path
		req.Host = u.Host
	}}
	proxy.ServeHTTP(w, r)
}

func main() {
	var lb LoadBalancer

	var s smpl.Configuration
	s.InitializeFromFile("./urls.smpl")

	lb.LanguageServers = s.Geta
	lb.NextServerIndex = make(map[string]int)

	ticker := time.NewTicker(14 * time.Minute)
	go pingServers(lb.LanguageServers,ticker)

	http.HandleFunc("/compile", lb.handleCompile)
	fmt.Println("Load Balancer listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func pingServers(servers map[string][string], ticker *time.ticker){
	// Ping the servers periodically to see if they are alive
	for {
		select {
		case <-ticker.C:
			for k,v:=range servers{
				_, err := http.post(v)
				if err != nil {
					fmt.Printf("Error pinging %s server with url %s: %v\n", k, v, err)
					return
				}
			}
		}
	}
}