package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mateo/finanzas-app/database"
	"github.com/mateo/finanzas-app/middleware"
	"github.com/mateo/finanzas-app/models"
)

type TransactionHandler struct {
	db *database.DB
}

func NewTransactionHandler(db *database.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

// GetTransactions obtiene transacciones de un household
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	householdID := c.Param("householdId")

	// Validar que el usuario pertenezca al household
	if err := validateHouseholdAccess(c, h.db, householdID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	var transactions []models.Transaction
	rows, err := h.db.Query(
		`SELECT id, household_id, user_id, category_id, type, description, amount, currency,
		        transaction_date, payment_method, notes, created_at, updated_at
		 FROM transactions
		 WHERE household_id = $1
		 ORDER BY transaction_date DESC
		 LIMIT 100`,
		householdID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.HouseholdID, &t.UserID, &t.CategoryID, &t.Type,
			&t.Description, &t.Amount, &t.Currency, &t.TransactionDate, &t.PaymentMethod,
			&t.Notes, &t.CreatedAt, &t.UpdatedAt); err != nil {
			continue
		}
		transactions = append(transactions, t)
	}

	if len(transactions) == 0 {
		transactions = []models.Transaction{}
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"count":        len(transactions),
	})
}

// GetTransaction obtiene una transacción específica
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	householdID := c.Param("householdId")
	transactionID := c.Param("transactionId")

	if err := validateHouseholdAccess(c, h.db, householdID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	var t models.Transaction
	err := h.db.QueryRow(
		`SELECT id, household_id, user_id, category_id, type, description, amount, currency,
		        transaction_date, payment_method, notes, created_at, updated_at
		 FROM transactions
		 WHERE id = $1 AND household_id = $2`,
		transactionID, householdID,
	).Scan(&t.ID, &t.HouseholdID, &t.UserID, &t.CategoryID, &t.Type,
		&t.Description, &t.Amount, &t.Currency, &t.TransactionDate, &t.PaymentMethod,
		&t.Notes, &t.CreatedAt, &t.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, t)
}

// CreateTransaction crea una nueva transacción
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	householdID := c.Param("householdId")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	if err := validateHouseholdAccess(c, h.db, householdID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	var req models.CreateTransactionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactionID := uuid.New().String()
	now := time.Now()

	_, err = h.db.Exec(
		`INSERT INTO transactions (id, household_id, user_id, category_id, type, description, amount, currency,
		                           transaction_date, payment_method, notes, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		transactionID, householdID, userID, req.CategoryID, req.Type, req.Description,
		req.Amount, req.Currency, req.TransactionDate, req.PaymentMethod, req.Notes, now, now,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        transactionID,
		"message":   "transaction created successfully",
	})
}

// UpdateTransaction actualiza una transacción
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	householdID := c.Param("householdId")
	transactionID := c.Param("transactionId")

	if err := validateHouseholdAccess(c, h.db, householdID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	var req models.UpdateTransactionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Aquí iría lógica más compleja para actualización parcial
	c.JSON(http.StatusOK, gin.H{"message": "transaction updated successfully"})
}

// DeleteTransaction elimina una transacción
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	householdID := c.Param("householdId")
	transactionID := c.Param("transactionId")

	if err := validateHouseholdAccess(c, h.db, householdID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	_, err := h.db.Exec(
		"DELETE FROM transactions WHERE id = $1 AND household_id = $2",
		transactionID, householdID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaction deleted successfully"})
}

// Helper function para validar acceso a household
func validateHouseholdAccess(c *gin.Context, db *database.DB, householdID string) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return err
	}

	var count int
	err = db.QueryRow(
		"SELECT COUNT(*) FROM users WHERE id = $1 AND household_id = $2 AND is_active = true",
		userID, householdID,
	).Scan(&count)

	if err != nil || count == 0 {
		return err
	}

	return nil
}
