package response

import (
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
)

type MovieResponse struct {
	ImdbId      string `json:"imdbid"`
	Title       string `json:"title"`
	RunTime     string `json:"runtime"`
	Plot        string `json:"plot"`
	ImdbRating  string `json:"imdbRating"`
	MoviePoster string `json:"poster"`
	Genre       string `json:"genre"`
}

func NewMovieResponse(movie *models.Movie) MovieResponse {
	return MovieResponse{
		ImdbId:      movie.MovieId,
		Title:       movie.Name,
		RunTime:     movie.Duration,
		Plot:        movie.Plot,
		ImdbRating:  movie.ImdbRating,
		MoviePoster: movie.MoviePoster,
		Genre:       movie.Genre,
	}
}
