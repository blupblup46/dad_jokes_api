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

func ToArray(jokesArr []Joke) []string {
	res := make([]string, 0, len(jokesArr))

	for _, joke := range jokesArr {
		res = append(res, joke.Joke)
	}
	return res
}

func Merge(jokes1 []string, jokes2 []string) map[int]string{
	res := make(map[int]string)
	ind := 0

	for _, v := range jokes1 {
		res[ind] = v
		ind++
	}

	for _, v := range jokes2 {
		res[ind] = v
		ind++
	}

	return res
}
