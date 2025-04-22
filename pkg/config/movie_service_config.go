package config

type MovieServiceConfig struct {
	BaseURL string
	APIKey  string
}

func GetMovieServiceConfig() MovieServiceConfig {
	return MovieServiceConfig{
		BaseURL: getEnvOrDefault("MOVIE_SERVICE_URL", "http://localhost:4567"),
		APIKey:  getEnvOrDefault("MOVIE_SERVICE_API_KEY", "test"),
	}
}
