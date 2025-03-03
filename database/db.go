package database

import (
	"database/sql"
	"fmt"
	"log"
	"go_server_monitor/config"
	
	_ "github.com/go-sql-driver/mysql"       // MySQL
	_ "github.com/lib/pq"                    // PostgreSQL
	_ "github.com/mattn/go-sqlite3"          // SQLite
	_ "github.com/denisenkom/go-mssqldb"     // MSSQL
	_ "github.com/godror/godror"             // Oracle
)

var DB *sql.DB

func ConnectDB(cfg *config.Config) error {
	if cfg.DBEnabled == 0 {
		log.Println("Databse Decativated.")
		return nil
	}

	var dsn string
	var err error

	switch cfg.DBType {
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
		DB, err = sql.Open("postgres", dsn)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
		DB, err = sql.Open("mysql", dsn)
	case "sqlite":
		DB, err = sql.Open("sqlite3", cfg.DBName)
	case "mssql":
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
		DB, err = sql.Open("sqlserver", dsn)
	case "oracle":
		dsn = fmt.Sprintf("%s/%s@%s:%s/%s",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
		DB, err = sql.Open("godror", dsn)
	default:
		return fmt.Errorf("Unsupported DB, please notify me at: \nhttps://github.com/Gambitdutku/Go-server-monitor-/issues: %s", cfg.DBType)

	}

	if err != nil {
		return fmt.Errorf("Database connection error: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("Database ping error: %w", err)
	}

	log.Printf("%s Connected to database", cfg.DBType)
	return nil
}

