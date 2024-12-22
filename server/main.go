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
const API_URL = "https://v2.jokeapi.dev/joke/Programming?amount=10&type=twopart"
const API_JOKES_PATH = "./jokes/apiJokes.json"
const CUSTOM_JOKES_PATH = "./jokes/customJokes.json"

var jokes map[int]utils.Joke

func main() {
	log.Println(`
		      ___           ___           ___           ___           ___                 
		     /\__\         /\  \         /\__\         /\  \         /\  \          ___   
		    /:/  /        /::\  \       /:/  /        /::\  \       /::\  \        /\  \  
		   /:/__/        /:/\:\  \     /:/__/        /:/\:\  \     /:/\:\  \       \:\  \ 
		  /::\  \ ___   /::\~\:\  \   /::\  \ ___   /::\~\:\  \   /::\~\:\  \      /::\__\
		 /:/\:\  /\__\ /:/\:\ \:\__\ /:/\:\  /\__\ /:/\:\ \:\__\ /:/\:\ \:\__\  __/:/\/__/
		 \/__\:\/:/  / \/__\:\/:/  / \/__\:\/:/  / \/__\:\/:/  / \/__\:\/:/  / /\/:/  /   
		      \::/  /       \::/  /       \::/  /       \::/  /       \::/  /  \::/__/    
		      /:/  /        /:/  /        /:/  /        /:/  /         \/__/    \:\__\    
		     /:/  /        /:/  /        /:/  /        /:/  /                    \/__/    
		     \/__/         \/__/         \/__/         \/__/                  (server)
	`)

	server, _ := CreateServer()
	log.Print("Listening on :", EXPOSED_PORT)
	server.ListenAndServe()
}

func CreateServer() (*http.Server, *http.ServeMux ){
	mux := http.NewServeMux()
	fetchApi := true

	if len(os.Args) > 1 && os.Args[1] == "--no-api" {
		log.Println("Using previously fetch jokes.")
		fetchApi = false
	}

	go buildDadabase(fetchApi)
	BuildHandlers(mux)

	server := &http.Server{
		Addr:    ":" + EXPOSED_PORT,
		Handler: mux,
	}

	return server, mux
}

func BuildHandlers(muxServer *http.ServeMux) {

	sendResponse := func(joke utils.Joke, w http.ResponseWriter) {
		jData, _ := json.Marshal(joke)
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(jData); err != nil {
			log.Println("Could not send response", err)
		}
	}

	muxServer.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Up")); err != nil {
			log.Println("Could not send response", err)
		}
	})

	muxServer.HandleFunc("/random", func(w http.ResponseWriter, r *http.Request) {
		sendResponse(GetRandomJoke(), w)
	})

	muxServer.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		jokeID, _ := strconv.Atoi(queryParams.Get("id"))
		if _, exists := jokes[jokeID]; exists{
			joke := utils.Joke{
				ID:       jokeID,
				Setup:    jokes[jokeID].Setup,
				Delivery: jokes[jokeID].Delivery,
			}
			sendResponse(joke, w)
		}else{
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 joke not found"))
			return
		}
		
	})

	muxServer.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		buildDadabase(true)
		if _, err := w.Write([]byte("Dadabase ready !")); err != nil {
			log.Println("Could not send response", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	muxServer.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		jData, _ := json.Marshal(jokes)
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(jData); err != nil {
			log.Println("Could not send response", err)
		}
	})
}

func GetRandomJoke() utils.Joke {
	var keys []int
	for key := range jokes {
		keys = append(keys, key)
	}

	randomIndex := rand.Intn(len(keys))

	randomKey := keys[randomIndex]
	return jokes[randomKey]
}

func buildDadabase(resetApiFile bool) {
	apiJokes := make(map[int]utils.Joke)

	if resetApiFile {
		apiJokes = fetchApi()
	} else {
		apiJokes = utils.ToMap(GetJokesFromFile(API_JOKES_PATH))
	}

	jokes = utils.Merge(apiJokes, getCustomJokes())
	log.Println("Dadabase ready")
}

func GetJokesFromFile(path string) []utils.Joke {
	var jokesArr []utils.Joke

	customJokesFile, err_openCustomJokesFile := os.Open(path)
	if err_openCustomJokesFile != nil {
		log.Fatal("Could not open file", err_openCustomJokesFile)
	}

	customJokesByte, err_readCustomJokesFile := io.ReadAll(customJokesFile)
	if err_readCustomJokesFile != nil {
		log.Fatal("Could not read file", err_readCustomJokesFile)
	}

	err_deserializeCustomJokes := json.Unmarshal(customJokesByte, &jokesArr)
	if err_deserializeCustomJokes != nil {
		log.Fatal("Could not deserialize file", err_deserializeCustomJokes)
	}

	return jokesArr
}

func getCustomJokes() map[int]utils.Joke {
	return utils.ToMap(GetJokesFromFile(CUSTOM_JOKES_PATH))
}

func fetchApi() map[int]utils.Joke {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	req, _ := http.NewRequest("GET", API_URL, nil)
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

	jokesArr := response.Results

	JsonApiJokes, err_serializeApiJokes := json.Marshal(jokesArr)
	if err_serializeApiJokes != nil {
		log.Println("Could not serialize API jokes", err_serializeApiJokes)
	}

	APIJokesFile, err_createFile := os.Create(API_JOKES_PATH)
	if err_createFile != nil {
		log.Println("Could not create API jokes file", err_createFile)
	}

	if _, err_writeFile := APIJokesFile.Write(JsonApiJokes); err != nil {
		log.Println("Could not write API jokes file", err_writeFile)
	}

	return utils.ToMap(jokesArr)
}
