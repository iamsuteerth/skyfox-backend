package models

type Movie struct {
	MovieId     string `json:"movieId"`
	Name        string `json:"name"`
	Duration    string `json:"duration"`
	Plot        string `json:"plot"`
	ImdbRating  string `json:"imdbRating"`
	MoviePoster string `json:"moviePoster"`
	Genre       string `json:"genre"`
}
