package api

import (
	"net/http"

	"github.com/npub1337/facemate/internal/face"
	"github.com/npub1337/facemate/internal/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	faceService    face.Service
	storageService storage.Service
}

func NewHandler(faceService face.Service, storageService storage.Service) *Handler {
	return &Handler{
		faceService:    faceService,
		storageService: storageService,
	}
}

type TrainRequest struct {
	Image    string `json:"image" binding:"required"`     // base64 encoded image
	PersonID string `json:"person_id" binding:"required"` // unique identifier for the person
}

type CompareResponse struct {
	PersonID   string  `json:"person_id,omitempty"`
	Similarity float32 `json:"similarity,omitempty"`
	Error      string  `json:"error,omitempty"`
}

func (h *Handler) Train(c *gin.Context) {
	var req TrainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert base64 to bytes
	imageData := []byte(req.Image)

	// Train the face
	if err := h.faceService.Train(imageData, req.PersonID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Face trained successfully"})
}

func (h *Handler) Compare(c *gin.Context) {
	var req struct {
		Image string `json:"image" binding:"required"` // base64 encoded image
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert base64 to bytes
	imageData := []byte(req.Image)

	// Compare the face
	personID, similarity, err := h.faceService.Compare(imageData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CompareResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CompareResponse{
		PersonID:   personID,
		Similarity: similarity,
	})
}
