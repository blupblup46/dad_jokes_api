# Requirements

Docker or go (v1.22.2)

# How does it work

At startup, the backend server makes a call to the joke API (https://v2.jokeapi.dev/) and stores the result in a file at `server/jokes/apiJokes.json`.
You can add your own jokes at `server/jokes/customJokes.json`. `server/jokes/apiJokes.json` is bound to change, do not add jokes here
Then, jokes from the API joke and custom jokes are merged. If the two files have a joke of same ID, the joke from the custom joke file will overwrite the one from the API joke.

# Startup

## Docker compose
At project root, run `docker compose up`
The client will be available at http://localhost:3000

## go
In `server/`, run `go run main.go`
You also have the possibility to start the server without running a first call to the jokes API with the option `--no-api`
You can add your own jokes at `server/jokes/customJokes.json`. 
`server/jokes/apiJokes.json` is bound to change, do not add jokes here

In `client/`, configure the parameters in config.json. The existing config should work out of the box.
You can also set the config as environment variables, `API_URL, API_PORT, CLIENT_PORT`
Then run `go run main.go`


# Client endpoints

#### /style
Serves the css file
#### /joke
Serves a html template filled with a random joke from the server
#### /joke?id=[id]
Serves a html template filled with the joke with id from the server.
#### /reset
Makes a call for /reset at server API
#### /search
Serves a html template with a form to make a server API call at /search?id=[filled id].


# Server endpoints

#### /health
Check if the server is up and running. If the server is OK, you will get a 200 OK response with "up" message
#### /random
get a random joke from the fusion of api and custom jokes
#### /search?id=[id]
search for a joke with id. If no joke exist with id, returns a 404 not found response
#### /reset
makes a call at the joke API and updates the `server/jokes/apiJokes.json` file with the new jokes. Then, custom jokes are added 
#### /all
returns the fusion between custom and api jokes