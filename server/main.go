package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"server/utils"
)

const EXPOSED_PORT = "8080"
const API_ROOT_URL = "https://icanhazdadjoke.com"
const API_BATCH_SUFFIX = "/search"

var jokes map[string]string

func main() {
	go buildDadabase()

	log.Print("Listening on :", EXPOSED_PORT)
	err := http.ListenAndServe(":"+EXPOSED_PORT, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func buildDadabase(){
	jokes = utils.Merge(fetchApi(), getCustomJokes())
	log.Println("Dadabase ready")
}

func getCustomJokes() map[string]string {
	var customJokes utils.CustomJokes

	customJokesFile, err_openCustomJokesFile := os.Open("./jokes/customJokes.json")
	if  err_openCustomJokesFile != nil{
		log.Fatal("Could not open custom jokes file", err_openCustomJokesFile)
	}

	customJokesByte, err_readCustomJokesFile := io.ReadAll(customJokesFile)
	if  err_readCustomJokesFile != nil{
		log.Fatal("Could not read custom jokes file", err_readCustomJokesFile)
	}

	err_deserializeCustomJokes := json.Unmarshal(customJokesByte, &customJokes)
	if  err_deserializeCustomJokes != nil{
		log.Fatal("Could not deserialize custom jokes file", err_deserializeCustomJokes)
	}

	return utils.ToJokesMap(customJokes.Jokes, false)
}

func fetchApi() map[string]string {
	client := &http.Client{}
	defer client.CloseIdleConnections()
	

	req, _ := http.NewRequest("GET", API_ROOT_URL+API_BATCH_SUFFIX, nil)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Could not fetch the API: ", err)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal("Could not read response from API: ", err)
	}

	var response utils.BatchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal("Could not deserialize JSON: ", err)
	}

	jokes = utils.ToJokesMap(response.Results, true)

	JsonApiJokes, err_serializeApiJokes := json.Marshal(jokes)
	APIJokesFile, err_createFile := os.Create("./jokes/API_jokes.json")

	if err_createFile == nil && err_serializeApiJokes == nil{
		APIJokesFile.Write(JsonApiJokes)
	}

	return jokes
}
