package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"server/utils"
)

const EXPOSED_PORT = "8080"
const API_ROOT_URL = "https://icanhazdadjoke.com"
const API_BATCH_SUFFIX = "/search"

var jokes map[int]string

func main() {
	buildDadabase()

	log.Print("Listening on :", EXPOSED_PORT)
	err := http.ListenAndServe(":"+EXPOSED_PORT, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func buildDadabase() {
	log.Print("calling API")

	client := &http.Client{}
	req, _ := http.NewRequest("GET", API_ROOT_URL+API_BATCH_SUFFIX, nil)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Could not fetch the API: ", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Could not read response from API: ", err)
	}

	var response utils.BatchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal("Could not deserialize JSON: ", err)
	}

	jokes = utils.ToJokesMap(response.Results)
	log.Print(jokes)

}
