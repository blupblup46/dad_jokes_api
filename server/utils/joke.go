package utils


type BatchResponse struct {
	Results []Joke `json:"jokes"`
}

type JokesArr struct {
	Jokes []Joke `json:"jokes"`
}

type Joke struct {
	ID 		 int
	Setup    string `json:"setup"`
	Delivery string `json:"delivery"`
}

func ToMap(jokesArr []Joke) map[int]Joke {

	res := make(map[int]Joke)

	for _, joke := range jokesArr {
		res[joke.ID] = joke
	}

	return res
}

func Merge(jokes1 map[int]Joke, jokes2 map[int]Joke) map[int]Joke {

	for _, v := range jokes2 {
		jokes1[v.ID] = v
	}

	return jokes1
}
