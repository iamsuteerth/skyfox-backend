package config

type PaymentServiceConfig struct {
	BaseURL string
	APIKey  string
}

func GetPaymentServiceConfig() PaymentServiceConfig {
	return PaymentServiceConfig{
		BaseURL: getEnvOrDefault("PAYMENT_GATEWAY_URL", "http://localhost:8082"),
		APIKey:  getEnvOrDefault("PAYMENT_GATEWAY_API_KEY", ""),
	}
}
