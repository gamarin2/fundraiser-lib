package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

const BaseUrl = "http://localhost:8545"

var ListeningPort int
var ContractAddress string

func main() {
	flag.IntVar(&ListeningPort, "port", 8080, "listening port")
	flag.StringVar(&ContractAddress, "contract", "0xa4028F2aec0ad18964e368338E5268FebB4F5423", "contract address")
	flag.Parse()

	fmt.Println("Listening on port", ListeningPort)
	fmt.Println("Contract address", ContractAddress)

	// hard coded calls
	http.HandleFunc("/eth/weiPerAtom", weiPerAtom)
	http.HandleFunc("/eth/totalAtom", totalAtom)
	http.HandleFunc("/eth/totalWei", totalWei)
	http.HandleFunc("/eth/numDonations", numDonations)
	http.HandleFunc("/eth/isActive", isActive)
	http.ListenAndServe(fmt.Sprintf(":%d", ListeningPort), nil)
}

const (
	// function signatures
	SIG_WeiPerAtom   = "0x574a5e31"
	SIG_TotalAtom    = "0x615fa416"
	SIG_TotalWei     = "0x6d98e9fc"
	SIG_NumDonations = "0x5e9a1849"
	SIG_IsActive     = "0x22f3e2d4"
)

func isActive(w http.ResponseWriter, r *http.Request) {
	makeRequest(SIG_IsActive, w, r)
}

func weiPerAtom(w http.ResponseWriter, r *http.Request) {
	makeRequest(SIG_WeiPerAtom, w, r)
}

func totalAtom(w http.ResponseWriter, r *http.Request) {
	makeRequest(SIG_TotalAtom, w, r)
}

func totalWei(w http.ResponseWriter, r *http.Request) {
	makeRequest(SIG_TotalWei, w, r)
}

func numDonations(w http.ResponseWriter, r *http.Request) {
	makeRequest(SIG_NumDonations, w, r)
}

func makeRequest(contractData string, w http.ResponseWriter, r *http.Request) {
	// Make json rpc request
	params := CallParams{ContractAddress, contractData}
	jsonR := JSONRPC{
		ID:      "1",
		JSONRPC: "2.0",
		Method:  "eth_call",
		Params:  []CallParams{params},
	}
	rString, _ := json.Marshal(jsonR)

	resp, err := http.Post(BaseUrl, "application/json", bytes.NewBuffer(rString))
	if err != nil {
		httpError(w, err.Error())
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		httpError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8600")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Content-Length, Cache-Control, cf-connecting-ip")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
	w.Write(b)
}

type JSONRPC struct {
	ID      string       `json:"id"`
	JSONRPC string       `json:"jsonrpc"`
	Method  string       `json:"method"`
	Params  []CallParams `json:"params"`
}

type CallParams struct {
	To   string `json:"to"`
	Data string `json:"data"`
}

func httpError(w http.ResponseWriter, errStr string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(errStr + "\n"))
}
