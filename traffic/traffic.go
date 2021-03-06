package main

import (
	"log"
	"bytes"
	"errors"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"github.com/gorilla/mux"
)

var (	
	authURL = "https://com-shi-va.barcelona.cat/api/auth"
	baseURL = "https://api-com-shi-va.barcelona.cat"
	afectacionsEndpoint = "/afectacions/"
)

func main() {
	log.Println("Setting multiplexor")
	m := mux.NewRouter()
	m.HandleFunc("/ping", logger(Ping))
	m.HandleFunc("/api/traffic", logger(Traffic))

	server := &http.Server{
		Addr:":8081",
		Handler: m,
	}
	log.Println("Server running...")
	log.Fatalln(server.ListenAndServe())
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

type AccessToken struct {
	OK int `json:"ok"`
	Token string `json:"access_token"`
	TokenType string `json:"token_type"`
}

func getAuthToken() (token string, err error) {
	resp, err := http.Get(authURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	accessToken := AccessToken{}
	if err := json.Unmarshal(body, &accessToken); err != nil {
		return "", err
	}
	if accessToken.OK != 1 {
		return "", errors.New("Error getting the access token.")
	}
	return accessToken.Token, nil
}

// TODO: update del token si esta caducado.
// avoid 403 error when requesting for a token.
func Traffic(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return 
	}
	auth, err := getAuthToken()
	if err != nil {
		http.Error(w, "Error while getting the TrafficToken", http.StatusInternalServerError)
		return
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.Path = afectacionsEndpoint
	q := u.Query()
	q.Set("access_token", auth)
	q.Set("token_type", "Bearer")
	u.RawQuery = q.Encode()
	log.Println(u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Access-Control-Allow-Origin", "*")
   	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)

	reader := bytes.NewReader(body)
	io.Copy(w, reader)
}
