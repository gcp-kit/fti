package config

import "os"

func GetProjectID(cfg *Config) string {
	if cfg.FirestoreProjectOnEmulator != "" {
		_ = os.Setenv("FIRESTORE_EMULATOR_HOST", cfg.FirestoreEmulatorHost)
		return cfg.FirestoreProjectOnEmulator
	}

	id, ok := os.LookupEnv("GCP_PROJECT")
	if ok {
		return id
	}

	id, ok = os.LookupEnv("GOOGLE_CLOUD_PROJECT")
	if ok {
		return id
	}

	return ""
}
