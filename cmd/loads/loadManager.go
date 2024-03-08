package main

import (
	"context"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/database"
)

func main() {
	config.InitConfig()
	db := database.GetDB()

	managerRepo := repository.InitManagerRepo(db)

	manager := models.ManagerBase{
		Login:    "testLogin",
		Password: "Password1",
	}

	id, err := managerRepo.Create(context.Background(), manager)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Manager with id %d crated", id)
}
