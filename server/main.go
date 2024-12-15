package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"server/utils"
	"strconv"
)

const EXPOSED_PORT = "8080"
const API_ROOT_URL = "https://icanhazdadjoke.com"
const API_BATCH_SUFFIX = "/search"

var jokes map[int]string

func main() {
	go buildDadabase()
	buildHandlers()

	log.Print("Listening on :", EXPOSED_PORT)
	err := http.ListenAndServe(":"+EXPOSED_PORT, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func buildHandlers() {

	sendResponse := func(joke utils.Joke, w http.ResponseWriter) {
		jData, _ := json.Marshal(joke)
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(jData); err != nil {
			log.Println("Could send response", err)
		}
	}

	http.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {
		sendResponse(getRandomJoke(), w)
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		jokeID, _ := strconv.Atoi(queryParams.Get("id"))
		joke := utils.Joke{
			ID:   strconv.Itoa(jokeID),
			Joke: jokes[jokeID],
		}
		sendResponse(joke, w)
	})

	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		buildDadabase()
		if _, err := w.Write([]byte("Dadabase ready !")); err != nil {
			log.Println("Could send response", err)
		}
	})
}

func getRandomJoke() utils.Joke {
	index := rand.Intn(len(jokes))
	return utils.Joke{
		ID:   strconv.Itoa(index),
		Joke: jokes[index],
	}
}

func buildDadabase() {
	jokes = utils.Merge(fetchApi(), getCustomJokes())
	log.Println("Dadabase ready")
}

func getCustomJokes() []string {
	var customJokes []string

	customJokesFile, err_openCustomJokesFile := os.Open("./jokes/customJokes.json")
	if err_openCustomJokesFile != nil {
		log.Fatal("Could not open custom jokes file", err_openCustomJokesFile)
	}

	customJokesByte, err_readCustomJokesFile := io.ReadAll(customJokesFile)
	if err_readCustomJokesFile != nil {
		log.Fatal("Could not read custom jokes file", err_readCustomJokesFile)
	}

	err_deserializeCustomJokes := json.Unmarshal(customJokesByte, &customJokes)
	if err_deserializeCustomJokes != nil {
		log.Fatal("Could not deserialize custom jokes file", err_deserializeCustomJokes)
	}

	return customJokes
}

func fetchApi() []string {
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

	jokesArr := utils.ToArray(response.Results)

	JsonApiJokes, err_serializeApiJokes := json.Marshal(jokesArr)
	if err_serializeApiJokes != nil{
		log.Println("Could not serialize API jokes", err_serializeApiJokes)
	}


	APIJokesFile, err_createFile := os.Create("./jokes/API_jokes.json")
	if err_createFile != nil{
		log.Println("Could not create API jokes file", err_createFile)
	}

	if _, err_writeFile := APIJokesFile.Write(JsonApiJokes); err != nil {
		log.Println("Could not write API jokes file", err_writeFile)
	}

	return jokesArr
}
