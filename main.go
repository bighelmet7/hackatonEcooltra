package main

import (
	"bytes"
	"strconv"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Vehicle payload.
type Vehicle struct {
	ID       string    `json:"id"`
	Position []float64 `json:"position"`
	Range    int64     `json:"range"`
}

var (
	baseURL          = "http://ecooltra.arnaugarcia.com"
	vehiclesEndpoint = "/vehicles.json"
)

func main() {

	// TODO:
	// /vehicle/<id>/ returns a single obj, print the available perimeter.
	r := mux.NewRouter()
	r.HandleFunc("/ping", logger(Ping))
	r.HandleFunc("/api/vehicles", logger(Vehicles))

	// TODO: Add TLS and Read timeout.
	serve := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Println("Server running...")
	log.Fatal(serve.ListenAndServe())
}

func logger(f func(w http.ResponseWriter, req *http.Request)) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("[%s] %s - %s", req.RemoteAddr, req.RequestURI, req.Method)
		f(w, req)
	}
}

// Ping probes the alives status of the services.
func Ping(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "This method it's not supported.", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Pong\n")
}

func groupBy(vehicles []Vehicle, maxMeters, threshold int64) []Vehicle {
	var result []Vehicle
	// TODO: this is just for critical status only, should be more dynamic status.
	for _, vehicle := range vehicles {
		if vehicle.Range <= int64(maxMeters/threshold) {
			result = append(result, vehicle)
		}
	}
	return result
}

type Data struct {
	Results []Vehicle `json:"results"`
	Total int `json:"total"`
}

// Vehicles nothing special, yet.
func Vehicles(w http.ResponseWriter, req *http.Request) {
	var (
		maxMeters int64 = 70000
		threshold int64 = 1
	)
	if req.Method != http.MethodGet {
		http.Error(w, "This method it's not supported.", http.StatusMethodNotAllowed)
		return
	}
	fullRange, ok := req.URL.Query()["maxMeters"]
	if ok {
		fullRangeInt, err := strconv.Atoi(fullRange[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		maxMeters = int64(fullRangeInt)
	}
	thresh, ok := req.URL.Query()["threshold"]
	if ok {
		threshInt, err := strconv.Atoi(thresh[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if threshInt <= 0  {
			http.Error(w, "threshold should be a positive integer.", http.StatusInternalServerError)
			return
		}
		threshold = int64(threshInt)
	}
	cli := &http.Client{}
	u, err := url.Parse(baseURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.Path = vehiclesEndpoint
	nReq, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	nReq.Header.Add("Content-Type", "application/json")
	resp, err := cli.Do(nReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var vehicles []Vehicle
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &vehicles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vehiclesGroup := groupBy(vehicles, maxMeters, threshold)
	results := Data{
		Results: vehiclesGroup,
		Total: len(vehiclesGroup),
	}
	b, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	reader := bytes.NewReader(b)

	w.Header().Set("Access-Control-Allow-Origin", "*")
   	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, reader)
}
