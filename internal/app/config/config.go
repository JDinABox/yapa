package config

import (
	"flag"
	"os"
)

// Config for the app.
type Config struct {
	// Address to listen to
	// Default: :8080
	Address string `json:"address"`
	// Database Path
	// Default: ./data.sqlite3
	DatabasePath string `json:"database_path"`
	// Domain used for redirects
	// Example: example.com
	Domain string `json:"domain"`
}

type option func(*Config) error

func New(options ...option) *Config {
	config := &Config{
		Address:      ":8080",
		DatabasePath: "./data.sqlite3",
	}

	for _, opt := range options {
		err := opt(config)
		if err != nil {
			panic(err)
		}
	}

	return config
}

func FromFlags() option {
	return func(c *Config) error {
		var (
			address      string
			domain       string
			databasePath string
		)

		flag.StringVar(&address, "address", "", "The address to listen on.")
		flag.StringVar(&databasePath, "db", "", "The path to the database file.")
		flag.StringVar(&domain, "domain", "", "The domain of the server.")
		flag.Parse()

		if address != "" {
			c.Address = address
		}
		if domain != "" {
			c.Domain = domain
		}
		if databasePath != "" {
			c.DatabasePath = databasePath
		}

		return nil
	}
}

func FromEnv() option {
	return func(c *Config) error {
		if address := os.Getenv("ADDRESS"); address != "" {
			c.Address = address
		}
		if databasePath := os.Getenv("DATABASE_PATH"); databasePath != "" {
			c.DatabasePath = databasePath
		}
		if domain := os.Getenv("DOMAIN"); domain != "" {
			c.Domain = domain
		}

		return nil
	}
}

func Defaults() option {
	return func(c *Config) error {
		c.Address = "0.0.0.0:8080"
		c.DatabasePath = "./data.sqlite3"
		c.Domain = ""
		return nil
	}
}
