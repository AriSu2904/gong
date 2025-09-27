package cmd

type TemplateData struct {
	ProjectName    string
	DatabaseDriver string
}

const envTemplate = `DB_SOURCE=postgresql://postgres:<your_password>@localhost:5432/dbname?sslmode=disable`

const cfgTemplateDb = `package config
import(
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DbSource		string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		log.Fatal("DB_SOURCE env variable not set")
	}

	return &Config{
		DbSource: dbSource,
	}, nil
}
`

const cfgTemplate = `package config
import(
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	return &Config{}, nil
}
`

const basicTemplate = `package main
import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)

	port := ":8080"

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	log.Println("Server running on port ", port)
}
`

func setCfgTemplate(isUnknownDriver bool) string {
	if isUnknownDriver {
		return cfgTemplate
	}

	return cfgTemplateDb
}

func setDbTemplate(isUnknownDriver bool) string {
	if isUnknownDriver {
		return "package database\n"
	}

	return ""
}

func setEnvTemplate(isUnknownDriver bool) string {
	if isUnknownDriver {
		return ""
	}

	return envTemplate
}

func setMainTemplate(isUnknownDriver bool) string {
	if isUnknownDriver {
		return basicTemplate
	}

	return ""
}
