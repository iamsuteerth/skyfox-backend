package movieservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/iamsuteerth/skyfox-backend/pkg/config"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"
)

type MovieServiceResponse struct {
	ImdbId      string `json:"imdbid"`
	Title       string `json:"title"`
	RunTime     string `json:"runtime"`
	Plot        string `json:"plot"`
	ImdbRating  string `json:"imdbRating"`
	MoviePoster string `json:"poster"`
	Genre       string `json:"genre"`
}

func (m MovieServiceResponse) ToMovie() (*models.Movie, error) {
	runtime := strings.Split(m.RunTime, " ")[0]
	duration, err := time.ParseDuration(runtime + "m")

	if err != nil {
		log.Error().Err(err).Str("movie", m.Title).Msg("Failed to parse runtime of movie")
		return nil, utils.NewInternalServerError("MOVIE_CONVERSION_FAILED", "Failed to process movie runtime", err)
	}

	movie := &models.Movie{
		MovieId:     m.ImdbId,
		Name:        m.Title,
		Duration:    duration.String(),
		Plot:        m.Plot,
		ImdbRating:  m.ImdbRating,
		MoviePoster: m.MoviePoster,
		Genre:       m.Genre,
	}
	return movie, nil
}

type MovieService interface {
	GetMovieById(ctx context.Context, id string) (*models.Movie, error)
	GetAllMovies(ctx context.Context) ([]*models.Movie, error)
}

type movieService struct {
	config config.MovieServiceConfig
}

func NewMovieService(cfg config.MovieServiceConfig) MovieService {
	return &movieService{
		config: cfg,
	}
}

func (s *movieService) GetMovieById(ctx context.Context, id string) (*models.Movie, error) {
	url := fmt.Sprintf("%s/movies/%s", s.config.BaseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, utils.NewInternalServerError("REQUEST_CREATION_FAILED", "Failed to create request", err)
	}

	req.Header.Add("x-api-key", s.config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, utils.NewInternalServerError("MOVIE_SERVICE_ERROR", "Failed to connect to movie service", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewInternalServerError(
			"MOVIE_SERVICE_ERROR",
			fmt.Sprintf("Movie service returned status code %d", resp.StatusCode),
			fmt.Errorf("non-OK status: %d", resp.StatusCode),
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.NewInternalServerError("RESPONSE_READ_ERROR", "Failed to read response body", err)
	}

	var movieResponse MovieServiceResponse
	if err := json.Unmarshal(body, &movieResponse); err != nil {
		return nil, utils.NewInternalServerError("JSON_PARSE_ERROR", "Failed to parse movie data", err)
	}

	movie, err := movieResponse.ToMovie()
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (s *movieService) GetAllMovies(ctx context.Context) ([]*models.Movie, error) {
	url := fmt.Sprintf("%s/movies", s.config.BaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, utils.NewInternalServerError("REQUEST_CREATION_FAILED", "Failed to create request", err)
	}

	req.Header.Add("x-api-key", s.config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, utils.NewInternalServerError("MOVIE_SERVICE_ERROR", "Failed to connect to movie service", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewInternalServerError(
			"MOVIE_SERVICE_ERROR",
			fmt.Sprintf("Movie service returned status code %d", resp.StatusCode),
			fmt.Errorf("non-OK status: %d", resp.StatusCode),
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.NewInternalServerError("RESPONSE_READ_ERROR", "Failed to read response body", err)
	}

	var moviesResponse []MovieServiceResponse
	if err := json.Unmarshal(body, &moviesResponse); err != nil {
		return nil, utils.NewInternalServerError("JSON_PARSE_ERROR", "Failed to parse movies data", err)
	}

	movies := make([]*models.Movie, 0, len(moviesResponse))
	for _, movieResp := range moviesResponse {
		movie, err := movieResp.ToMovie()
		if err != nil {
			log.Error().Err(err).Str("movie", movieResp.Title).Msg("Skipping movie due to conversion error")
			continue
		}
		movies = append(movies, movie)
	}

	return movies, nil
}
