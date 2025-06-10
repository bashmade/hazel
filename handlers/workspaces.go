package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/freekobie/hazel/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateWorkspace(c *gin.Context) {
	var input struct {
		Name        string    `json:"name" binding:"required"`
		Description string    `json:"description"`
		UserID      uuid.UUID `json:"userId" binding:"required,uuid"`
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ws := &models.Workspace{
		Name:        input.Name,
		Description: input.Description,
		User:        &models.User{Id: input.UserID},
	}

	err = h.wss.NewWorkspace(c.Request.Context(), ws)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ws)
}

func (h *Handler) GetWorkspace(c *gin.Context) {
	idStr := c.Param("id")

	if err := validate.Var(idStr, "uuid"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	ws, err := h.wss.GetWorkspace(c.Request.Context(), uuid.MustParse(idStr))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, ws)
}

// Get all the workspaces where a user has membership.
func (h *Handler) GetUserWorkspaces(c *gin.Context) {
	idStr := c.Param("id")
	fmt.Println(idStr)

	if err := validate.Var(idStr, "uuid"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id format"})
		return
	}

	workspaces, err := h.wss.GetWorkspace(c.Request.Context(), uuid.MustParse(idStr))
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, workspaces)
}
