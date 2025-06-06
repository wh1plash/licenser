package main

import (
	"fmt"
	"licenser/server/api"
	"licenser/server/store"
	"log"
	"os"
	"strconv"
	"time"

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

	redisClient := store.NewRedisClient(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	ttl, err := time.ParseDuration(os.Getenv("REDIS_TTL"))
	if err != nil {
		log.Fatal("error to rarse duration")
	}
	dbCache := store.NewChachedStore(db, redisClient, ttl)

	app := fiber.New()
	appHandler := api.NewAppHandler(dbCache)

	app.Get("/app", appHandler.HandleGetApp)
	app.Post("/app", appHandler.HandleInsertApp)
	app.Get("/apps", appHandler.HandleGetAppList)

	log.Fatal(app.Listen(os.Getenv("LISTEN_ADDR")))

}

func mustLoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
