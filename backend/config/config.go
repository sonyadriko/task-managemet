package config

import "os"

type Config struct {
	JWTSecret      string
	JWTExpiration  int
	ServerPort     string
	AllowedOrigins []string
}

func GetConfig() *Config {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-this-in-production"
	}

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173"
	}

	return &Config{
		JWTSecret:      jwtSecret,
		JWTExpiration:  24, // hours
		ServerPort:     serverPort,
		AllowedOrigins: []string{allowedOrigins},
	}
}
