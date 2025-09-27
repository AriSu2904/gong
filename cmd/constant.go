package cmd

import "github.com/charmbracelet/huh"

var dbDriver = map[string]string{
	"postgresql": "github.com/jackc/pgx/v5/stdlib",
	"mysql":      "github.com/go-sql-driver/mysql",
	"sqlite":     "github.com/mattn/go-sqlite3",
	"mongodb":    "go.mongodb.org/mongo-driver",
}

var dbOptions = []huh.Option[string]{
	huh.NewOption("PostgreSQL", "postgresql"),
	huh.NewOption("MySQL", "mysql"),
	huh.NewOption("SQLite", "sqlite"),
	huh.NewOption("MongoDB", "mongodb"),
	huh.NewOption("None", "N/A"),
}
