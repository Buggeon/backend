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
	"bugtracker/graph"
	s3storage "bugtracker/internal/S3Storage"
	"bugtracker/internal/db"
	"bugtracker/internal/handlers"
	"bugtracker/internal/middleware"
	"bugtracker/internal/models"
	"bugtracker/internal/repositories"
	"bugtracker/internal/security"
	"bugtracker/internal/services"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createAdmin(tokenService *services.TokenService) error {

	admin_password, err := security.HashPassword(os.Getenv("ADMIN_PASSWORD"))

	if err != nil {
		return err
	}

	admin_login := "@" + os.Getenv("ADMIN_LOGIN")
	admin_name := os.Getenv("ADMIN_NAME")
	admin_email := os.Getenv("ADMIN_EMAIL")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	users := db.GetCollection("users")

	count, err := users.CountDocuments(ctx, bson.M{"role": "admin"})

	if err != nil {
		return err
	}

	if count == 0 {

		admin := models.User{
			ID:        primitive.NewObjectID(),
			Name:      admin_name,
			Login:     admin_login,
			Password:  admin_password,
			Email:     admin_email,
			Role:      "admin",
			CreatedAt: time.Now(),
		}

		refreshToken, err := tokenService.GenerateRefreshToken(&admin)

		if err != nil {
			return err
		}

		admin.RefreshTokens = []string{refreshToken}

		_, err = users.InsertOne(context.TODO(), admin)

		return err
	}

	return nil

}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error while config loading")
	}

	mongoHost := os.Getenv("DB_HOST")
	mongoPort := os.Getenv("DB_PORT")
	mongoUser := os.Getenv("DB_USER")
	mongoPassword := os.Getenv("DB_PASSWORD")

	connectionString := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin",
		mongoUser, mongoPassword, mongoHost, mongoPort)

	if err := db.Connect(connectionString); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	corsConfig := cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
			"HEAD",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	server := gin.Default()
	server.Use(corsConfig)

	s3Storage := s3storage.NewS3Storage()

	userRepo := repositories.NewUserRepo()
	projectRepo := repositories.NewProjectRepo()
	memberRepo := repositories.NewMemberRepo()
	boardRepo := repositories.NewBoardRepo()
	cardRepo := repositories.NewCardRepo()
	messageRepo := repositories.NewMessageRepo()
	schemaRepo := repositories.NewSchemaRepo()

	boardService := services.NewBoardService(boardRepo, projectRepo)
	cardService := services.NewCardService(cardRepo, boardRepo)
	memberService := services.NewMemberService(memberRepo, projectRepo)
	schemaService := services.NewSchemaService(projectRepo, schemaRepo, s3Storage)
	projectService := services.NewProjectService(projectRepo, memberService, memberRepo, s3Storage)
	tokenService := services.NewTokenService(config.LoadConfig())
	userService := services.NewUserService(userRepo, tokenService)
	systemService := services.NewSystemService(userRepo)
	messageService := services.NewMessageService(messageRepo, cardRepo)

	projectHandler := handlers.NewProjectHandler(projectService, schemaService, cardService, boardService, memberService, messageService)
	userHandler := handlers.NewUserHandler(userService)
	systemHandler := handlers.NewSystemHandler(systemService)

	authMiddleware := middleware.NewAuthMiddleware(tokenService)

	if err := createAdmin(tokenService); err != nil {
		fmt.Println("Failed to create admin")
	}

	graphqlResolver := &graph.Resolver{
		ProjectService: projectService,
		BoardService:   boardService,
		CardService:    cardService,
		UserService:    userService,
		MemberService:  memberService,
		MessageService: messageService,
	}

	gqlHandler := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: graphqlResolver},
		),
	)

	setupRoutes(server, projectHandler, userHandler, systemHandler, authMiddleware, gqlHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s", port)
	server.Run("127.0.0.1:" + port)
}

func setupRoutes(
	server *gin.Engine,
	projectHandler *handlers.ProjectHandler,
	userHandler *handlers.UserHandler,
	systemHandler *handlers.SystemHandler,
	authMiddleware *middleware.AuthMiddleware,
	gqlHandler *handler.Server,
) {
	server.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	server.POST("/auth/register", userHandler.Registration)
	server.POST("/auth/login", userHandler.Login)
	server.POST("/auth/refreshtoken", userHandler.RefreshAccessToken)
	server.GET("/test", handlers.Test)

	api := server.Group("/api")
	api.Use(authMiddleware.AuthRequired())
	{
		projects := api.Group("/projects")
		{
			projects.GET("", projectHandler.GetProjects)
			projects.GET("/:project_id", projectHandler.GetProject)
			projects.POST("", projectHandler.CreateProject)
			projects.DELETE("/:project_id", projectHandler.DeleteProject)
			projects.POST("/:project_id/schemas", projectHandler.AddProjectSchema)
			projects.PATCH("/:project_id/logo", projectHandler.SetProjectLogo)

			schemas := projects.Group("/:project_id/schemas")
			{
				schemas.POST("", projectHandler.AddProjectSchema)
				schemas.GET("/:_id", projectHandler.GetMember)
				schemas.POST("", projectHandler.CreateMember)
				schemas.DELETE("/:_id", projectHandler.DeleteMember)
			}

			members := projects.Group("/:project_id/members")
			{
				members.GET("", projectHandler.GetMembers)
				members.GET("/:member_id", projectHandler.GetMember)
				members.POST("", projectHandler.CreateMember)
				members.DELETE("/:member_id", projectHandler.DeleteMember)
			}

			boards := projects.Group("/:project_id/boards")
			{
				boards.GET("", projectHandler.GetBoards)
				boards.GET("/:board_id", projectHandler.GetBoard)
				boards.POST("", projectHandler.CreateBoard)
				boards.DELETE("/:board_id", projectHandler.DeleteBoard)

				cards := boards.Group("/:board_id/cards")
				{
					cards.GET("", projectHandler.GetCards)
					cards.GET("/:card_id", projectHandler.GetCard)
					cards.POST("", projectHandler.CreateCard)
					cards.DELETE("/:card_id", projectHandler.DeleteCard)
					cards.PUT("/:card_id/updatelocation", projectHandler.UpdateCardLocation)

					messages := cards.Group("/:card_id/messages")
					{
						messages.POST("", projectHandler.NewMessage)
					}
				}
			}
		}

		queryGroup := api.Group("/query")
		queryGroup.POST("", gin.WrapH(gqlHandler))

		api.GET("/users", systemHandler.GetAllUsers)
	}
}
