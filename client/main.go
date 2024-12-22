package main

import (
	"client/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

const CONFIG_PATH = "config.json"
var config utils.Config

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
		     \/__/         \/__/         \/__/         \/__/                  (client)
	`)
	loadConfig()
	fmt.Println(config)
	server, _ := CreateServer()
	log.Print("Listening on :", config.ExposePort)
	server.ListenAndServe()
}

func loadConfig() utils.Config{
	configFile, err_openConfigFile := os.Open(CONFIG_PATH)
	if err_openConfigFile != nil {
		log.Fatal("Could not open config file", err_openConfigFile)
	}

	configByte, err_readConfigFile := io.ReadAll(configFile)
	if err_readConfigFile != nil {
		log.Fatal("Could not read file", err_readConfigFile)
	}

	err_deserializeConfig := json.Unmarshal(configByte, &config)
	if err_deserializeConfig != nil {
		log.Fatal("Could not deserialize file", err_deserializeConfig)
	}

	return config
}

func CreateServer() (*http.Server, *http.ServeMux ){
	mux := http.NewServeMux()

	BuildHandlers(mux)

	server := &http.Server{
		Addr:    ":" + config.ExposePort,
		Handler: mux,
	}

	return server, mux
}


func BuildHandlers(muxServer *http.ServeMux) {
	muxServer.HandleFunc("/joke", func(w http.ResponseWriter, r *http.Request) {
		var joke utils.Joke
		var statusCode int

		queryParams := r.URL.Query()
		jokeId, _ := strconv.Atoi(queryParams.Get("id"))

		if jokeId == 0{
			joke, statusCode = fetchApi("/random")
		}else{
			path := fmt.Sprint("/search?id=", jokeId)
			joke, statusCode = fetchApi(path)
		}

		if statusCode == http.StatusNotFound {
			log.Fatalf("Joke #%d not found", jokeId)
		}
		if statusCode != http.StatusOK {
			log.Fatal("Could not fetch API, status code:",statusCode)
		}

		tmpl, err := template.ParseFiles("./html_files/joke.html")
		if err != nil {
			log.Fatal("Error parsing template ./html_files/joke.html")
		}
		tmpl.Execute(w, joke)
	})

	muxServer.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		statusCode := resetRequest()

		resetMessage := "Dadabase reseted !"
		if statusCode != http.StatusOK {
			log.Fatal("Could not reset jokes API dadabase, status code:",statusCode)
			resetMessage = "Could not reset jokes API dadabase"
		}

		tmpl, err := template.ParseFiles("./html_files/reset.html")
		if err != nil {
			log.Fatal("Error parsing template ./html_files/reset.html")
		}

		tmpl.Execute(w, resetMessage)
	})


}

func getRequest(url string) *http.Response{
	client := &http.Client{}
	defer client.CloseIdleConnections()


	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Could not fetch the API: ", err)
	}

	return resp
}

func fetchApi(path string) (utils.Joke, int){
	url := fmt.Sprintf("%s:%s%s", config.ApiUrl, config.ApiPort, path)
	resp := getRequest(url)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Could not read response from API: ", err)
	}

	var response utils.Joke
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal("Could not deserialize JSON: ", err)
	}

	return response, resp.StatusCode

}

func resetRequest() int {
	url := fmt.Sprintf("%s:%s%s", config.ApiUrl, config.ApiPort, "/reset")
	resp := getRequest(url)
	defer resp.Body.Close()

	return resp.StatusCode
}