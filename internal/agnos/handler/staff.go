package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/Markikie/agnos/internal/agnos/api/request"
	"github.com/Markikie/agnos/internal/agnos/api/response"
	"github.com/Markikie/agnos/internal/agnos/service"
)

type StaffHandler struct {
	staffService service.StaffService
	jwtSecret    string
}

type Claims struct {
	StaffID  string `json:"staff_id"`
	Username string `json:"username"`
	Hospital string `json:"hospital"`
	jwt.RegisteredClaims
}

func NewStaffHandler(
	staffService service.StaffService,
) StaffHandler {
	return StaffHandler{
		staffService: staffService,
		jwtSecret:    "your-secret-key", // In production, use environment variable
	}
}

func (h *StaffHandler) CreateStaff(c *gin.Context) {
	var req request.StaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" || req.Hospital == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, password, and hospital are required"})
		return
	}

	staff, err := h.staffService.CreateStaff(req.Username, req.Password, req.Hospital)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Staff created successfully",
		"staff_id": staff.ID,
	})
}

func (h *StaffHandler) Login(c *gin.Context) {
	var req request.LoginStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get hospital from query parameter
	hospital := c.Query("hospital")
	if hospital == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hospital parameter is required"})
		return
	}

	staff, err := h.staffService.Login(req.Username, req.Password, hospital)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token
	claims := &Claims{
		StaffID:  staff.ID.String(),
		Username: staff.Username,
		Hospital: staff.Hospital,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, response.LoginStaffResponse{
		AccessToken: tokenString,
	})
}
