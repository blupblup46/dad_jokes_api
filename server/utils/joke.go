package utils

type BatchResponse struct {
	Results []Joke `json:"results"`
}

type CustomJokes struct {
	Jokes []Joke `json:"jokes"`
}

type Joke struct {
	ID   string `json:"id"`
	Joke string `json:"joke"`
}

func ToJokesMap(jokesArr []Joke, isFromApi bool) map[string]string {
	res := make(map[string]string)
	var prefix string

	if isFromApi {
		prefix = "API_"
	}

	for _, joke := range jokesArr {
		res[prefix+joke.ID] = joke.Joke
	}
	return res
}

func Merge(jokes1 map[string]string, jokes2 map[string]string) map[string]string {
	for k, v := range jokes2 {
		jokes1[k] = v
	}

	return jokes1
}
