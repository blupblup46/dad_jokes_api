package utils

type BatchResponse struct {
    Results      []Joke  `json:"results"`
}

type Joke struct {
    ID   string `json:"id"`
    Joke string `json:"joke"`
}

func ToJokesMap(jokesArr []Joke) map[int]string{
	res := make(map[int]string)
	
	for index, joke := range jokesArr{
		res[index] = joke.Joke
	}
	return res
}