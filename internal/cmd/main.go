package main

import (
	"github.com/joho/godotenv"
	"internal.com/cmd/api"
	"internal.com/cmd/api/handlers"
	dbmanager "internal.com/internal/platform/database"
	"internal.com/internal/service"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	database, err := dbmanager.NewDatabase().Connect()
	if err != nil {
		log.Fatalf("[CONNECT DATABASE FAIL]: %s", err.Error())
	}

	internalRepository := service.NewServiceRepository(database)

	internalService := service.NewInternalService(internalRepository)
	err = internalService.InsertDefaultService()
	if err != nil {
		log.Fatalf("[INSERT DEFAULT FAIL]: %s", err.Error())
	}

	internalHandler := handlers.NewInternalHandler(internalService)

	err = api.New(internalHandler)
	if err != nil {
		log.Fatalf("[START SERVER FAIL]: %s", err.Error())
	}
}
