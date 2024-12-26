package utils

type Joke struct {
	ID       int
	Setup    string `json:"setup"`
	Delivery string `json:"delivery"`
}
