// Buggeon - SelfHosted service for bug and task tracking
// Copyright (C) 2026 DEVE corp.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"bugtracker/config"
	"bugtracker/internal/db"
	"bugtracker/internal/handlers"
	"bugtracker/internal/middleware"
	"bugtracker/internal/repositories"
	"bugtracker/internal/services"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error while config loading")
	}

	mongoHost := os.Getenv("DB_HOST")
	mongoPort := os.Getenv("DB_PORT")
	mongoUser := os.Getenv("DB_USER")
	mongoPassword := os.Getenv("DB_PASSWORD")

	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", mongoUser, mongoPassword, mongoHost, mongoPort)

	if err := db.Connect(connectionString); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Access-Control-Allow-Origin"},
	}))

	userRepo := repositories.NewUserRepo()
	projectRepo := repositories.NewProjectRepo()
	memberRepo := repositories.NewMemberRepo()
	boardRepo := repositories.NewBoardRepo()
	cardRepo := repositories.NewCardRepo()

	boardService := services.NewBoardService(boardRepo)
	cardService := services.NewCardService(cardRepo)
	memberService := services.NewMemberService(memberRepo)
	projectService := services.NewProjectService(projectRepo, memberService)
	tokenService := services.NewTokenService(config.LoadConfig())
	userService := services.NewUserService(userRepo, tokenService)
	systemService := services.NewSystemService(userRepo)

	projectHandler := handlers.NewProjectHandler(projectService, cardService, boardService, memberService)
	userHandler := handlers.NewUserHandler(userService)
	systemHandler := handlers.NewSystemHandler(systemService)

	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	protected := server.Group("/api")
	protected.Use(authMiddleware.AuthRequired())
	{
		server.GET("/api/projects", func(ctx *gin.Context) {
			projectHandler.GetProjects(ctx)
		})

		server.GET("/api/projects/:project_id", func(ctx *gin.Context) {
			projectHandler.GetProjects(ctx)
		})

		server.POST("/api/projects", func(ctx *gin.Context) {
			projectHandler.CreateProject(ctx)
		})

		server.GET("/api/projects/members/:member_id", func(ctx *gin.Context) {
			projectHandler.GetMember(ctx)
		})

		server.GET("/api/projects/members", func(ctx *gin.Context) {
			projectHandler.GetMembers(ctx)
		})

		server.POST("/api/projects/members", func(ctx *gin.Context) {
			projectHandler.DeleteMember(ctx)
		})

		server.POST("/api/projects/boards", func(ctx *gin.Context) {
			projectHandler.CreateBoard(ctx)
		})

		server.GET("/api/projects/boards", func(ctx *gin.Context) {
			projectHandler.GetBoards(ctx)
		})

		server.GET("/api/projects/boards/:board_id", func(ctx *gin.Context) {
			projectHandler.GetBoard(ctx)
		})

		server.POST("/api/projects/boards", func(ctx *gin.Context) {
			projectHandler.GetBoard(ctx)
		})

		server.DELETE("/api/projects/boards", func(ctx *gin.Context) {
			projectHandler.DeleteBoard(ctx)
		})

		server.POST("/api/projects/boards/cards", func(ctx *gin.Context) {
			projectHandler.CreateCard(ctx)
		})

		server.GET("/api/projects/boards/cards", func(ctx *gin.Context) {
			projectHandler.GetCards(ctx)
		})

		server.GET("/api/projects/boards/cards/:card_id", func(ctx *gin.Context) {
			projectHandler.DeleteBoard(ctx)
		})

		server.GET("/api/getAllUsers", func(ctx *gin.Context) {
			systemHandler.GetAllUsers(ctx)
		})

		server.DELETE("/api/projects/boards/cards", func(ctx *gin.Context) {
			projectHandler.DeleteBoard(ctx)
		})
	}

	server.POST("/user/registration", func(ctx *gin.Context) {
		userHandler.Registration(ctx)
	})
	server.POST("user/login", func(ctx *gin.Context) {
		userHandler.Login(ctx)
	})

	server.GET("/test", handlers.Test)

	server.Run(":" + os.Getenv("PORT"))

}
