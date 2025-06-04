package main

import (
	"fmt"
	"licenser/server/api"
	"licenser/server/store"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	mustLoadEnvVariables()
}

func main() {
	port, _ := strconv.Atoi(os.Getenv("PG_PORT"))
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("PG_HOST"), port, os.Getenv("PG_USER"), os.Getenv("PG_PASS"), os.Getenv("PG_DB_NAME"))
	db, err := store.NewPostgresStore(connStr)
	if err != nil {
		log.Fatal("error to connect to Posgres database")
		return
	}

	if err := db.Init(); err != nil {
		log.Fatal("error to create tables", "error", err.Error())
		return
	}

	if err := db.CreateApp(); err != nil {
		fmt.Println(err)
	}

	app := fiber.New()
	appHandler := api.NewAppHandler(db)

	app.Get("/app", appHandler.HandleGetApp)
	app.Post("/app", appHandler.HandleInsertApp)

	log.Fatal(app.Listen(":9080"))

}

func mustLoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
