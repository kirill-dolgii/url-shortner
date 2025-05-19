package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/kirill-dolgii/url-shortner/internal/config/dbconfig"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

func main() {
	// CLI флаги
	up := flag.Bool("up", false, "apply all up migrations")
	down := flag.Bool("down", false, "rollback last migration")
	force := flag.Int("force", -1, "force database to specific version")
	showVersion := flag.Bool("version", false, "show current migration version")
	flag.Parse()

	// Загрузка конфигурации из .env
	cfg, err := dbconfig.LoadDBConfig()
	if err != nil {
		log.Fatalf("failed to load DB config: %v", err)
	}

	// Подключение к базе
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("failed to initialize migrator: %v", err)
	}

	// --force N
	if *force >= 0 {
		if err := m.Force(*force); err != nil {
			log.Fatalf("failed to force version: %v", err)
		}
		fmt.Printf("✅ Forced to version %d\n", *force)
		return
	}

	// --version
	if *showVersion {
		v, dirty, err := m.Version()
		if err == migrate.ErrNilVersion {
			fmt.Println("📦 No migration has been applied yet.")
			return
		}
		if err != nil {
			log.Fatalf("failed to get version: %v", err)
		}
		fmt.Printf("📦 Current version: %d (dirty: %v)\n", v, dirty)
		return
	}

	// --down
	if *down {
		v, _, _ := m.Version()
		if v == 0 {
			fmt.Println("📭 No migration to roll back.")
			return
		}
		if err := m.Steps(-1); err != nil {
			log.Fatalf("⛔ Rollback failed: %v", err)
		}
		fmt.Println("⏪ Rolled back 1 migration.")
		return
	}

	// --up
	if *up {
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				fmt.Println("✅ No new migrations to apply.")
			} else {
				log.Fatalf("migration failed: %v", err)
			}
		} else {
			fmt.Println("✅ All migrations applied successfully.")
		}
		return
	}

	// Если флагов нет
	fmt.Println("⚠️  Please specify one of the flags: --up, --down, --force=N, --version")
}
