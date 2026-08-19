package config

import "os"

type Config struct {
	Port                string
	DatabaseURL         string
	JWTSecret           string
	SMTPHost            string
	SMTPPort            string
	SMTPUser            string
	SMTPPass            string
	AppBaseURL          string
	RappelJours         string
	SchedulerIntervalle string
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() Config {
	return Config{
		Port:                getenv("PORT", "8086"),
		DatabaseURL:         getenv("DATABASE_URL", "postgres://postgres:zak@localhost:5433/nomorewaste?sslmode=disable"),
		JWTSecret:           getenv("JWT_SECRET", "dev-secret-change-me"),
		SMTPHost:            getenv("SMTP_HOST", ""),
		SMTPPort:            getenv("SMTP_PORT", "587"),
		SMTPUser:            getenv("SMTP_USER", ""),
		SMTPPass:            getenv("SMTP_PASS", ""),
		AppBaseURL:          getenv("APP_BASE_URL", "http://localhost:8086"),
		RappelJours:         getenv("RAPPEL_JOURS", "30"),
		SchedulerIntervalle: getenv("SCHEDULER_INTERVALLE", "24h"),
	}
}
