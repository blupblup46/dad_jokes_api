package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"server/utils"
	"strconv"
	"testing"
)

func TestServerStarts(t *testing.T) {
	_, handler := CreateServer()

	ts := httptest.NewServer(handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

}

func TestApiJokesAreWritten(t *testing.T) {
	_, handler := CreateServer()

	ts := httptest.NewServer(handler)
	defer ts.Close()

	file, err_openCustomJokesFile := os.Open(API_JOKES_PATH)
	if err_openCustomJokesFile != nil {
		t.Errorf("Could not open file %s", err_openCustomJokesFile)
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		t.Errorf("Error getting file stats for %s: %v", API_JOKES_PATH, err)
	}
	if stat.Size() == 0 {
		t.Error("File is empty")
	}

}

func TestJokesAreRead(t *testing.T) {
	jokes := GetJokesFromFile(API_JOKES_PATH)

	if jokes[0].Setup == "" {
		t.Error("Jokes are not correctly read")
	}
}

func TestSearchJoke(t *testing.T) {
	_, handler := CreateServer()

	ts := httptest.NewServer(handler)
	defer ts.Close()
	
	buildDadabase(false)
	joke := GetRandomJoke()

	resp, err := http.Get(ts.URL + "/search?id=" + strconv.Itoa(joke.ID))
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil{
		t.Fatal("Could not read response from API: ", err)
	}
	if resp.StatusCode == http.StatusNotFound{
		t.Fatalf("Joke #%s not found", strconv.Itoa(joke.ID))
	}

	var serverJoke utils.Joke
	if err := json.Unmarshal(body, &serverJoke); err != nil {
		t.Fatal("Could not deserialize JSON: ", err)
	}

	if serverJoke.ID != joke.ID{
		t.Error("joke from file and joke from server are not the same")
	}

}


func TestDatabaseReset(t *testing.T){

	getAllJokes := func(ts *httptest.Server)(map[int]utils.Joke){
		resp, err := http.Get(ts.URL + "/all")

		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
	
		if err != nil{
			t.Fatal("Could not read response from API: ", err)
		}
		var serverJokes map[int]utils.Joke
		if err := json.Unmarshal(body, &serverJokes); err != nil {
			t.Fatal("Could not deserialize JSON: ", err)
		}
		return serverJokes
	}
	_, handler := CreateServer()

	ts := httptest.NewServer(handler)
	defer ts.Close()

	firstBatch := getAllJokes(ts)
	http.Get(ts.URL + "/reset")
	secondBatch := getAllJokes(ts)

	if reflect.DeepEqual(firstBatch, secondBatch) {
		t.Error("Jokes did not reset")
	}
}