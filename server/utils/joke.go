package utils

type BatchResponse struct {
	Results []Joke `json:"jokes"`
}

type JokesArr struct {
	Jokes []Joke `json:"jokes"`
}

type Joke struct {
	ID       int
	Setup    string `json:"setup"`
	Delivery string `json:"delivery"`
}

// turns a Joke array to a map with key as joke.ID
// jokesArr = the joke array to turn into a map
func ToMap(jokesArr []Joke) map[int]Joke {

	res := make(map[int]Joke)

	for _, joke := range jokesArr {
		res[joke.ID] = joke
	}

	return res
}

// merge to jokes map,
// if the joke from the second map has an ID existing in the first map, the joke in the second map will override the joke in the first map
func Merge(jokes1 map[int]Joke, jokes2 map[int]Joke) map[int]Joke {

	for _, v := range jokes2 {
		jokes1[v.ID] = v
	}

	return jokes1
}
