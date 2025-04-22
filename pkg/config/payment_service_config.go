package config

type PaymentServiceConfig struct {
	BaseURL string
	APIKey  string
}

func GetPaymentServiceConfig() PaymentServiceConfig {
	return PaymentServiceConfig{
		BaseURL: getEnvOrDefault("PAYMENT_SERVICE_URL", "http://localhost:8082"),
		APIKey:  getEnvOrDefault("PAYMENT_SERVICE_API_KEY", ""),
	}
}
