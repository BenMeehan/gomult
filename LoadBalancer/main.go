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

	// Select a server from the list (using a simple round-robin strategy)
	serverURL := servers[0]

	u, err := url.Parse(serverURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal server error")
		log.Printf("Failed to parse server URL: %v", err)
		return
	}
	u.Path = "/compile"

	// Proxy the request to the selected server URL
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}

func main() {
	var lb LoadBalancer

	var s smpl.Configuration
	s.InitializeFromFile("./urls.smpl")

	lb.LanguageServers=s.Geta

	http.HandleFunc("/compile", lb.handleCompile)
	fmt.Println("Load Balancer listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}