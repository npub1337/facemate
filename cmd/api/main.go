package main

import (
	"log"
	"net/http"

	"github.com/npub1337/facemate/internal/api"
	"github.com/npub1337/facemate/internal/face"
	"github.com/npub1337/facemate/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize face recognition service
	faceService, err := face.NewService()
	if err != nil {
		log.Fatalf("Failed to initialize face recognition service: %v", err)
	}

	// Initialize storage service
	storageService := storage.NewService()

	// Initialize router
	router := gin.Default()

	// Initialize API handlers
	handler := api.NewHandler(faceService, storageService)

	// Register routes
	router.POST("/train", handler.Train)
	router.POST("/compare", handler.Compare)

	// Start server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
