package main

import (
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

// Flags needed m8
var (
	baseURL          = "http://ecooltra.arnaugarcia.com"
	vehiclesEndpoint = "/vehicles.json"
	geoEndpoint = "/test.geojson"
	accessToken      = "Bearer 0fb6f9fffe309680c17d6fb7203cded9a39fc5b865f36d0763211e70a9948c58"
	maxMeters        = 65000
)

func main() {

	// TODO:
	// /vehicle/<id>/ returns a single obj, print the available perimeter.
	r := mux.NewRouter()
	r.HandleFunc("/ping", logger(Ping))
	r.HandleFunc("/api/vehicles", logger(Vehicles))
	r.HandleFunc("/api/geo", logger(Points))

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

func groupBy(vehicles []Vehicle) []Vehicle {
	var result []Vehicle
	// TODO: this is just for critical status only, should be more dynamic status.
	for _, vehicle := range vehicles {
		// critical is interpreted as the 25% of the total.
		if vehicle.Range <= int64(maxMeters/4) {
			result = append(result, vehicle)
		}
	}
	return result
}

// Vehicles nothing special, yet.
// TODO: add critic zones and other zones.
func Vehicles(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "This method it's not supported.", http.StatusMethodNotAllowed)
		return
	}
	cli := &http.Client{}
	u, err := url.Parse(baseURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.Path = vehiclesEndpoint
	q := u.Query()
	q.Set("system_id", "barcelona")
	u.RawQuery = q.Encode()
	nReq, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	nReq.Header.Add("Content-Type", "application/json")
	nReq.Header.Add("Authorization", accessToken)
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
	vehiclesGroup := groupBy(vehicles)
	b, err := json.Marshal(vehiclesGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func Points(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "This method it's not supported.", http.StatusMethodNotAllowed)
		return
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.Path = geoEndpoint
	resp, err := http.Get(u.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	w.Header().Set("Access-Control-Allow-Origin", "*")
   	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	io.Copy(w, resp.Body)
}
