package main

import (
	"log"
	"io"
	"net/url"
	"net/http"
	"github.com/gorilla/mux"
)

var (
	baseURL = "https://cooltra.electricfeel.net"
	vehiclesEndpoint = "/integrator/v1/vehicles"
	accessToken = "Bearer 0fb6f9fffe309680c17d6fb7203cded9a39fc5b865f36d0763211e70a9948c58"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/ping", logger(Ping))
	r.HandleFunc("/vehicles", logger(Vehicles))

	// TODO: Add TLS and Read timeout.
	serve := &http.Server{
		Addr: ":8080",
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
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
