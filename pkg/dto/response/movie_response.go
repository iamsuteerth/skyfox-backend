package response

import (
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
)

type MovieResponse struct {
	MovieId     string `json:"movieId"`
	Name        string `json:"name"`
	RunTime     string `json:"runtime"`
	Plot        string `json:"plot"`
	ImdbRating  string `json:"imdbRating"`
	MoviePoster string `json:"poster"`
	Genre       string `json:"genre"`
}

func NewMovieResponse(movie *models.Movie) MovieResponse {
	return MovieResponse{
		MovieId:     movie.MovieId,
		Name:        movie.Name,
		RunTime:     movie.Duration,
		Plot:        movie.Plot,
		ImdbRating:  movie.ImdbRating,
		MoviePoster: movie.MoviePoster,
		Genre:       movie.Genre,
	}
}
