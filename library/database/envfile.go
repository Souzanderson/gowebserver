package database

import (
	"encoding/json"
	"log"
	"os"
)

func init() {
	StartEnvironment()
}

var env_data map[string]string

func getEnvData() map[string]string {
	env_data = read_file_env()

	return env_data
}

func read_file_env() map[string]string {
	env_location_file := "./production.json"

	if os.Getenv("GO_ENV") == "release" {
		env_location_file = "./production.json"
	}
	if os.Getenv("GO_ENV") == "development" {
		env_location_file = "./development.json"
	}

	data, err := os.ReadFile(env_location_file)

	if err != nil {
		log.Fatalf("Erro ao abrir environment!")
	}

	values := string(data)
	var x map[string]string
	json.Unmarshal([]byte(values), &x)

	log.Println("[ENVIRONMENT] Environment location => ", env_location_file)

	return x
}

func StartEnvironment() {
	for k, v := range getEnvData() {
		os.Setenv(k, v)
	}
}

func GetEnvKey(key string) string {
	return env_data[key]
}
