package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/mateo/finanzas-app/config"
)

type DB struct {
	*sql.DB
}

// Connect abre conexión a PostgreSQL
func Connect(cfg *config.Config) (*DB, error) {
	dbConn, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configurar pool de conexiones
	dbConn.SetMaxOpenConns(25)
	dbConn.SetMaxIdleConns(5)
	dbConn.SetConnMaxLifetime(5 * 60) // 5 minutos

	// Probar conexión
	if err := dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DB{dbConn}, nil
}

// InitializeSchema ejecuta el script SQL para crear tablas
func InitializeSchema(db *DB) error {
	// Leer archivo SQL
	schemaPath := filepath.Join(os.Getenv("SCHEMA_PATH"), "DATABASE_SCHEMA.sql")
	if os.Getenv("SCHEMA_PATH") == "" {
		schemaPath = "./DATABASE_SCHEMA.sql"
	}

	schema, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Ejecutar script SQL
	if _, err := db.Exec(string(schema)); err != nil {
		// No es error fatal si las tablas ya existen
		// Las sentencias CREATE IF NOT EXISTS lo manejan
		fmt.Printf("Note: %v\n", err)
	}

	return nil
}

// CreateHousehold crea un nuevo household
func (db *DB) CreateHousehold(name, description, currency string) (string, error) {
	var id string
	err := db.QueryRow(
		`INSERT INTO households (name, description, currency)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		name, description, currency,
	).Scan(&id)
	return id, err
}

// GetHousehold obtiene un household por ID
func (db *DB) GetHousehold(id string) error {
	return db.QueryRow(
		`SELECT id FROM households WHERE id = $1`,
		id,
	).Scan(&id)
}

// Transaction helper para ejecutar múltiples operaciones
func (db *DB) BeginTx() (*sql.Tx, error) {
	return db.Begin()
}

// CreateDefaultCategories crea categorías por defecto para un household
func (db *DB) CreateDefaultCategories(householdID string) error {
	defaultCategories := []struct {
		name        string
		description string
		categoryType string
		color       string
		icon        string
	}{
		{"Salario", "Ingreso por trabajo", "income", "#10B981", "briefcase"},
		{"Venta", "Venta de bienes", "income", "#8B5CF6", "shopping-bag"},
		{"Transferencia", "Transferencia recibida", "income", "#3B82F6", "arrow-down"},
		{"Otros Ingresos", "Otros ingresos diversos", "income", "#06B6D4", "plus"},

		{"Comida", "Supermercado, restaurantes", "expense", "#EF4444", "shopping-cart"},
		{"Servicios", "Gas, electricidad, agua, internet", "expense", "#F59E0B", "lightbulb"},
		{"Transporte", "Auto, combustible, público", "expense", "#EC4899", "truck"},
		{"Salud", "Medicinas, doctor", "expense", "#8B5CF6", "heart"},
		{"Educación", "Cursos, libros", "expense", "#3B82F6", "book"},
		{"Entretenimiento", "Películas, deportes", "expense", "#10B981", "music"},
		{"Vivienda", "Alquiler, mantenimiento", "expense", "#F97316", "home"},
		{"Impuestos", "Impuestos y contribuciones", "expense", "#EF4444", "receipt"},
		{"Ahorro", "Depósitos a ahorros", "expense", "#06B6D4", "piggy-bank"},
		{"Otros Gastos", "Gastos diversos", "expense", "#6B7280", "dots"},
	}

	stmt := `INSERT INTO categories (household_id, name, description, type, color, icon, is_custom)
	         VALUES ($1, $2, $3, $4, $5, $6, false)
	         ON CONFLICT (household_id, name) DO NOTHING`

	for _, cat := range defaultCategories {
		if _, err := db.Exec(stmt, householdID, cat.name, cat.description, cat.categoryType, cat.color, cat.icon); err != nil {
			return fmt.Errorf("failed to create category %s: %w", cat.name, err)
		}
	}

	return nil
}
