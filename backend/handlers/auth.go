package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/mateo/finanzas-app/database"
	"github.com/mateo/finanzas-app/models"
)

type AuthHandler struct {
	db *database.DB
}

func NewAuthHandler(db *database.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Register crea un nuevo usuario y household
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email         string `json:"email" binding:"required,email"`
		Name          string `json:"name" binding:"required"`
		Password      string `json:"password" binding:"required,min=8"`
		HouseholdName string `json:"household_name" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar que el email no exista
	var existingUser string
	err := h.db.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Crear household
	householdID := uuid.New().String()
	_, err = h.db.Exec(
		"INSERT INTO households (id, name, currency) VALUES ($1, $2, $3)",
		householdID, req.HouseholdName, "ARS",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create household"})
		return
	}

	// Crear categorías por defecto
	if err := h.db.CreateDefaultCategories(householdID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create default categories"})
		return
	}

	// Hash la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Crear usuario con rol admin (primer usuario del household)
	userID := uuid.New().String()
	_, err = h.db.Exec(
		`INSERT INTO users (id, email, name, password_hash, role, household_id, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, true)`,
		userID, req.Email, req.Name, string(hashedPassword), "admin", householdID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// Generar JWT
	user := models.User{
		ID:          userID,
		Email:       req.Email,
		Name:        req.Name,
		Role:        "admin",
		HouseholdID: householdID,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	household := models.Household{
		ID:       householdID,
		Name:     req.HouseholdName,
		Currency: "ARS",
	}

	token, _ := generateJWT(user, "your-secret-key-change-in-production")

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token:     token,
		User:      user,
		Household: household,
	})
}

// Login autentica un usuario
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buscar usuario
	var user models.User
	var passwordHash string
	err := h.db.QueryRow(
		`SELECT id, email, name, role, household_id, is_active, created_at, updated_at, password_hash
		 FROM users WHERE email = $1 AND is_active = true`,
		req.Email,
	).Scan(&user.ID, &user.Email, &user.Name, &user.Role, &user.HouseholdID, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt, &passwordHash)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Validar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Buscar household
	var household models.Household
	err = h.db.QueryRow(
		"SELECT id, name, description, currency, created_at, updated_at FROM households WHERE id = $1",
		user.HouseholdID,
	).Scan(&household.ID, &household.Name, &household.Description, &household.Currency,
		&household.CreatedAt, &household.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch household"})
		return
	}

	// Generar JWT
	token, _ := generateJWT(user, "your-secret-key-change-in-production")

	c.JSON(http.StatusOK, models.AuthResponse{
		Token:     token,
		User:      user,
		Household: household,
	})
}

// generateJWT crea un JWT token
func generateJWT(user models.User, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":          user.ID,
		"email":        user.Email,
		"name":         user.Name,
		"role":         user.Role,
		"household_id": user.HouseholdID,
		"exp":          time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 días
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
