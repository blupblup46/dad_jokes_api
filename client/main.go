package main

import (
	"client/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

const CONFIG_PATH = "config.json"
const JOKE_TEMPLATE_PATH = "./html_files/joke.html"
const RESET_TEMPALTE_PATH = "./html_files/reset.html"
const SEARCH_TEMPLATE_PATH = "./html_files/search.html"

var config utils.Config
var logErr = log.New(os.Stderr, "", 0)

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
	server, _ := CreateServer()
	log.Print("Listening on :", config.ExposePort)
	if err := server.ListenAndServe(); err != nil {
		logErr.Printf("Could not start server at port %s, %s", config.ExposePort, err)
	}
}

// load config from the config.json file
func loadConfig() utils.Config {

	configFile, err_openConfigFile := os.Open(CONFIG_PATH)
	if err_openConfigFile != nil {
		logErr.Print("Could not open config file", err_openConfigFile)
	} else {
		configByte, err_readConfigFile := io.ReadAll(configFile)
		if err_readConfigFile != nil {
			logErr.Print("Could not read config file", err_readConfigFile)
		} else {
			err_deserializeConfig := json.Unmarshal(configByte, &config)
			if err_deserializeConfig != nil {
				logErr.Print("Could not deserialize config file", err_deserializeConfig)
			}
		}
	}

	if os.Getenv("API_URL") != "" {
		config.ApiUrl = os.Getenv("API_URL")
	}

	if os.Getenv("API_PORT") != "" {
		config.ApiPort = os.Getenv("API_PORT")
	}

	if os.Getenv("CLIENT_PORT") != "" {
		config.ExposePort = os.Getenv("CLIENT_PORT")
	}

	return config
}

func CreateServer() (*http.Server, *http.ServeMux) {
	mux := http.NewServeMux()
	BuildHandlers(mux)

	server := &http.Server{
		Addr:    ":" + config.ExposePort,
		Handler: mux,
	}

	return server, mux
}

/*
Build handlers
muxServer: the server to attacch handlers to
*/
func BuildHandlers(muxServer *http.ServeMux) {

	// serve the css file for html. Probably not the correct way to do that
	muxServer.HandleFunc("/style", func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join("./html_files/style.css")
		http.ServeFile(w, r, filePath)
	})

	// convenience endpoint that redirect to /joke
	muxServer.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://localhost:"+config.ExposePort+"/joke", http.StatusMovedPermanently)
	})

	// Serves a html template filled with a random joke from the server
	// With parameter ?id=[id], serves a html template filled with the joke with id from the server.
	// If id = 0, serves a random joke
	muxServer.HandleFunc("/joke", func(w http.ResponseWriter, r *http.Request) {
		var joke utils.Joke
		var statusCode int

		queryParams := r.URL.Query()
		jokeId, _ := strconv.Atoi(queryParams.Get("id"))

		if jokeId == 0 {
			joke, statusCode = fetchApi("/random")
		} else {
			path := fmt.Sprint("/search?id=", jokeId)
			joke, statusCode = fetchApi(path)
		}

		if statusCode == http.StatusNotFound {
			logErr.Printf("Joke #%d not found\n", jokeId)
			joke = utils.Joke{ID: -1, Setup: "Joke #" + strconv.Itoa(jokeId) + " not found", Delivery: ""}
		} else if statusCode != http.StatusOK {
			logErr.Print("Could not fetch API, status code:", statusCode)
			joke = utils.Joke{ID: -1, Setup: "Could not fetch API", Delivery: ""}
		}

		tmpl, err := template.ParseFiles(JOKE_TEMPLATE_PATH)
		if err != nil {
			logErr.Print("Error parsing template", JOKE_TEMPLATE_PATH)
		}
		templateVars := map[string]any{"joke": joke, "port": config.ExposePort}
		tmpl.Execute(w, templateVars)

	})

	//Makes a call for /reset at server API
	muxServer.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		statusCode := resetRequest()

		resetMessage := "Dadabase reseted !"
		if statusCode != http.StatusOK {
			resetMessage = "Could not reset jokes API dadabase"
			errMessage := fmt.Sprint(resetMessage, statusCode)
			logErr.Print(errMessage)
		}

		tmpl, err := template.ParseFiles(RESET_TEMPALTE_PATH)

		if err != nil {
			logErr.Print("Error parsing template", RESET_TEMPALTE_PATH)
		}

		templateVars := map[string]any{"resetMessage": resetMessage, "port": config.ExposePort}
		tmpl.Execute(w, templateVars)
	})

	// Serves a html template with a form to make a server API call at /search?id=[filled id].
	muxServer.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(SEARCH_TEMPLATE_PATH)

		if err != nil {
			logErr.Print("Error parsing template", SEARCH_TEMPLATE_PATH)
		}

		tmpl.Execute(w, config.ExposePort)
	})

}

// makes a http get request at url
// url = the server to fetch
func getRequest(url string) *http.Response {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		logErr.Print("Error creating request: ", err)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logErr.Print("Could not fetch the API: ", err)
	}

	return resp
}

// fetch the server API
// path = path to fetch (/joke, /reset, ...)
func fetchApi(path string) (utils.Joke, int) {
	url := fmt.Sprintf("%s:%s%s", config.ApiUrl, config.ApiPort, path)
	log.Println(url)
	resp := getRequest(url)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logErr.Print("Could not read response from API: ", err)
	}
	var response utils.Joke

	if resp.StatusCode == http.StatusOK {
		if err := json.Unmarshal(body, &response); err != nil {
			logErr.Print("Could not deserialize JSON: ", err)
		}
	}

	return response, resp.StatusCode
}

// makes a call at /reset endpoint at server API
func resetRequest() int {
	url := fmt.Sprintf("%s:%s%s", config.ApiUrl, config.ApiPort, "/reset")
	resp := getRequest(url)
	defer resp.Body.Close()

	return resp.StatusCode
}
