package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

// User representa un usuario del sistema
type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Role        string    `json:"role"` // admin, editor, viewer
	HouseholdID string    `json:"household_id"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Household representa el grupo familiar/domiciliario
type Household struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category representa una categoría de transacción
type Category struct {
	ID          string    `json:"id"`
	HouseholdID string    `json:"household_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // income, expense, transfer
	Color       string    `json:"color"`
	Icon        string    `json:"icon"`
	IsCustom    bool      `json:"is_custom"`
	CreatedAt   time.Time `json:"created_at"`
}

// Transaction representa una transacción financiera
type Transaction struct {
	ID              string    `json:"id"`
	HouseholdID     string    `json:"household_id"`
	UserID          string    `json:"user_id"`
	CategoryID      string    `json:"category_id"`
	Type            string    `json:"type"` // income, expense, transfer
	Description     string    `json:"description"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	TransactionDate time.Time `json:"transaction_date"`
	PaymentMethod   string    `json:"payment_method"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	// Relaciones
	Category *Category `json:"category,omitempty"`
	User     *User     `json:"user,omitempty"`
}

// CurrencyExchange representa compra/venta de divisas
type CurrencyExchange struct {
	ID               string    `json:"id"`
	HouseholdID      string    `json:"household_id"`
	UserID           string    `json:"user_id"`
	OperationType    string    `json:"operation_type"` // buy, sell
	AmountLocal      float64   `json:"amount_local"`
	AmountForeign    float64   `json:"amount_foreign"`
	ExchangeRate     float64   `json:"exchange_rate"`
	ForeignCurrency  string    `json:"foreign_currency"`
	ExchangeDate     time.Time `json:"exchange_date"`
	Location         string    `json:"location"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
}

// RecurringExpense representa impuestos/servicios recurrentes
type RecurringExpense struct {
	ID          string    `json:"id"`
	HouseholdID string    `json:"household_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	CategoryID  *string   `json:"category_id"`
	Frequency   string    `json:"frequency"` // monthly, quarterly, annual, bimonthly
	DueDay      *int      `json:"due_day"`
	DueMonth    *int      `json:"due_month"`
	Status      string    `json:"status"` // active, inactive, paid
	LastPaidDate *time.Time `json:"last_paid_date"`
	NextDueDate time.Time `json:"next_due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ServicePayment representa un pago de servicio/impuesto
type ServicePayment struct {
	ID               string    `json:"id"`
	HouseholdID      string    `json:"household_id"`
	UserID           string    `json:"user_id"`
	RecurringExpenseID *string  `json:"recurring_expense_id"`
	ServiceName      string    `json:"service_name"`
	Amount           float64   `json:"amount"`
	Currency         string    `json:"currency"`
	PaymentDate      time.Time `json:"payment_date"`
	PaymentMethod    string    `json:"payment_method"`
	ReferenceNumber  string    `json:"reference_number"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
}

// AccountBalance representa saldo en cuenta
type AccountBalance struct {
	ID          string    `json:"id"`
	HouseholdID string    `json:"household_id"`
	Currency    string    `json:"currency"`
	Balance     float64   `json:"balance"`
	AccountType string    `json:"account_type"` // cash, savings, checking, dollar_savings
	AccountName string    `json:"account_name"`
	LastUpdated time.Time `json:"last_updated"`
}

// Budget representa un presupuesto por categoría
type Budget struct {
	ID              string    `json:"id"`
	HouseholdID     string    `json:"household_id"`
	CategoryID      string    `json:"category_id"`
	Month           int       `json:"month"`
	Year            int       `json:"year"`
	BudgetedAmount  float64   `json:"budgeted_amount"`
	SpentAmount     float64   `json:"spent_amount"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ============================================================================
// DTOs (Data Transfer Objects) para API
// ============================================================================

// CreateTransactionRequest DTO
type CreateTransactionRequest struct {
	CategoryID      string    `json:"category_id" binding:"required"`
	Type            string    `json:"type" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	Amount          float64   `json:"amount" binding:"required,gt=0"`
	Currency        string    `json:"currency" binding:"required"`
	TransactionDate time.Time `json:"transaction_date" binding:"required"`
	PaymentMethod   string    `json:"payment_method"`
	Notes           string    `json:"notes"`
}

// UpdateTransactionRequest DTO
type UpdateTransactionRequest struct {
	CategoryID      string     `json:"category_id"`
	Description     string     `json:"description"`
	Amount          *float64   `json:"amount"`
	Currency        string     `json:"currency"`
	TransactionDate *time.Time `json:"transaction_date"`
	PaymentMethod   string     `json:"payment_method"`
	Notes           string     `json:"notes"`
}

// CreateCurrencyExchangeRequest DTO
type CreateCurrencyExchangeRequest struct {
	OperationType   string    `json:"operation_type" binding:"required"`
	AmountLocal     float64   `json:"amount_local" binding:"required,gt=0"`
	AmountForeign   float64   `json:"amount_foreign" binding:"required,gt=0"`
	ExchangeRate    float64   `json:"exchange_rate" binding:"required,gt=0"`
	ForeignCurrency string    `json:"foreign_currency" binding:"required"`
	ExchangeDate    time.Time `json:"exchange_date" binding:"required"`
	Location        string    `json:"location"`
	Notes           string    `json:"notes"`
}

// CreateCategoryRequest DTO
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
}

// CreateRecurringExpenseRequest DTO
type CreateRecurringExpenseRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Currency    string `json:"currency"`
	CategoryID  *string `json:"category_id"`
	Frequency   string `json:"frequency" binding:"required"`
	DueDay      *int   `json:"due_day"`
	DueMonth    *int   `json:"due_month"`
}

// LoginRequest DTO
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest DTO
type RegisterRequest struct {
	Email        string `json:"email" binding:"required,email"`
	Name         string `json:"name" binding:"required"`
	Password     string `json:"password" binding:"required,min=8"`
	HouseholdID  string `json:"household_id" binding:"required"`
}

// AuthResponse DTO
type AuthResponse struct {
	Token     string `json:"token"`
	User      User   `json:"user"`
	Household Household `json:"household"`
}

// ============================================================================
// Report Models
// ============================================================================

// SummaryReport contiene un resumen financiero
type SummaryReport struct {
	TotalIncome      float64       `json:"total_income"`
	TotalExpense     float64       `json:"total_expense"`
	NetBalance       float64       `json:"net_balance"`
	Period           string        `json:"period"`
	CurrencyBalances map[string]float64 `json:"currency_balances"`
}

// MonthlyBreakdown contiene desglose mensual
type MonthlyBreakdown struct {
	Month            int                     `json:"month"`
	Year             int                     `json:"year"`
	Income           float64                 `json:"income"`
	Expense          float64                 `json:"expense"`
	Balance          float64                 `json:"balance"`
	CategoryBreakdown map[string]CategoryStat `json:"category_breakdown"`
}

// CategoryStat estadísticas por categoría
type CategoryStat struct {
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Percentage   float64 `json:"percentage"`
}

// CurrencyBalanceReport saldo por divisa
type CurrencyBalanceReport struct {
	ARS      float64 `json:"ars"`
	USD      float64 `json:"usd"`
	TotalARS float64 `json:"total_ars"` // Conversión de todo a ARS
	TotalUSD float64 `json:"total_usd"` // Conversión de todo a USD
}
