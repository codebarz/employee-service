package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/codebarz/employee-service/database"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {

	logger := log.New(os.Stdout, "employee-service ", log.LstdFlags)

	godotenv.Load()

	r := chi.NewRouter()

	dbURL := os.Getenv("DB_URL")

	db := database.NewDatabase(logger)

	dbCfg := database.Config{DatabaseURL: dbURL}

	conn, err := db.OpenConnection(dbCfg)

	if err != nil {
		fmtErr := fmt.Sprintf("Err connecting to postgres DB. [ERROR]:%v", err)
		log.Fatal(fmtErr)
		os.Exit(1)
	}

	connErr := conn.Ping()

	if connErr != nil {
		log.Fatal(connErr)
	}

	log.Println("DB connection successful")

	if err := db.Migrate(dbCfg); err != nil {
		log.Fatal("Migration err", err)
		os.Exit(1)
	}

	defer func() {
		// db.l..Log("Closing DB connection", err)
		conn.Close()
	}()

	http.ListenAndServe(":9090", r)
}
