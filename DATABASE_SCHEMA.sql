-- ============================================================================
-- Finanzas Domésticas - Schema PostgreSQL
-- ============================================================================
-- Aplicación multi-usuario para gestión de finanzas domésticas
-- Soporta: ingresos, egresos, compra/venta de dólares, impuestos, servicios

-- Extensiones
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- 1. USUARIOS Y AUTENTICACIÓN
-- ============================================================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'viewer', -- 'admin', 'editor', 'viewer'
    household_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE households (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS', -- Moneda principal
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE users ADD CONSTRAINT fk_users_household
    FOREIGN KEY (household_id) REFERENCES households(id) ON DELETE CASCADE;

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_household_id ON users(household_id);

-- ============================================================================
-- 2. CATEGORÍAS Y SUBCATEGORÍAS
-- ============================================================================

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    household_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL, -- 'income', 'expense', 'transfer'
    color VARCHAR(7), -- Color hexadecimal para UI
    icon VARCHAR(50), -- Nombre del ícono (para frontend)
    is_custom BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (household_id) REFERENCES households(id) ON DELETE CASCADE,
    UNIQUE(household_id, name)
);

CREATE INDEX idx_categories_household ON categories(household_id);
CREATE INDEX idx_categories_type ON categories(type);

-- ============================================================================
-- 3. TRANSACCIONES (Ingresos, Egresos, Servicios, Impuestos)
-- ============================================================================

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    household_id UUID NOT NULL,
    user_id UUID NOT NULL,
    category_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL, -- 'income', 'expense', 'transfer'
    description VARCHAR(500) NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    transaction_date DATE NOT NULL,
    payment_method VARCHAR(50), -- 'cash', 'bank_transfer', 'credit_card', 'debit_card'
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (household_id) REFERENCES households(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE INDEX idx_transactions_household ON transactions(household_id);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);
CREATE INDEX idx_transactions_category ON transactions(category_id);
CREATE INDEX idx_transactions_user ON transactions(user_id);
CREATE INDEX idx_transactions_type ON transactions(type);

-- ============================================================================
-- 4. OPERACIONES DE DIVISAS (Compra/Venta de Dólares)
-- ============================================================================

CREATE TABLE currency_exchanges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    household_id UUID NOT NULL,
    user_id UUID NOT NULL,
    operation_type VARCHAR(20) NOT NULL, -- 'buy', 'sell'
    amount_local NUMERIC(15, 2) NOT NULL, -- Cantidad en moneda local (ARS)
    amount_foreign NUMERIC(15, 2) NOT NULL, -- Cantidad en moneda extranjera (USD)
    exchange_rate NUMERIC(10, 4) NOT NULL, -- Tasa de cambio utilizada
    foreign_currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    exchange_date DATE NOT NULL,
    location VARCHAR(100), -- Dónde se realizó el cambio
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (household_id) REFERENCES households(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_currency_exchanges_household ON currency_exchanges(household_id);
CREATE INDEX idx_currency_exchanges_date ON currency_exchanges(exchange_date);
CREATE INDEX idx_currency_exchanges_type ON currency_exchanges(operation_type);

-- ============================================================================
-- 5. IMPUESTOS Y SERVICIOS (Recurrentes)
-- ============================================================================

CREATE TABLE recurring_expenses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    household_id UUID NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    amount NUMERIC(15, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    category_id UUID,
    frequency VARCHAR(20) NOT NULL, -- 'monthly', 'quarterly', 'annual', 'bimonthly'
    due_day INT, -- Día del mes en que vence
    due_month INT, -- Mes (para anuales)
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- 'active', 'inactive', 'paid'
    last_paid_date DATE,
    next_due_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (household_id) REFERENCES households(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);

CREATE INDEX idx_recurring_expenses_household ON recurring_expenses(household_id);
CREATE INDEX idx_recurring_expenses_due_date ON recurring_expenses(next_due_date);
CREATE INDEX idx_recurring_expenses_status ON recurring_expenses(status);

-- ============================================================================
-- 6. PAGOS DE SERVICIOS E IMPUESTOS (Historial)
-- ============================================================================

CREATE TABLE service_payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    household_id UUID NOT NULL,
    user_id UUID NOT NULL,
    recurring_expense_id UUID,
    service_name VARCHAR(150) NOT NULL,
    amount NUMERIC(15, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    payment_date DATE NOT NULL,
    payment_method VARCHAR(50),
    reference_number VARCHAR(100), -- Número de comprobante o referencia
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (household_id) REFERENCES households(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (recurring_expense_id) REFERENCES recurring_expenses(id) ON DELETE SET NULL
);

CREATE INDEX idx_service_payments_household ON service_payments(household_id);
CREATE INDEX idx_service_payments_date ON service_payments(payment_date);

-- ============================================================================
-- 7. SALDOS DE CUENTAS (Snapshot de efectivo)
-- ============================================================================

CREATE TABLE account_balances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    household_id UUID NOT NULL,
    currency VARCHAR(3) NOT NULL,
    balance NUMERIC(15, 2) NOT NULL,
    account_type VARCHAR(50), -- 'cash', 'savings', 'checking', 'dollar_savings'
    account_name VARCHAR(100),
    last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (household_id) REFERENCES households(id) ON DELETE CASCADE,
    UNIQUE(household_id, currency, account_type)
);

CREATE INDEX idx_account_balances_household ON account_balances(household_id);

-- ============================================================================
-- 8. PRESUPUESTOS (Opcional pero útil)
-- ============================================================================

CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    household_id UUID NOT NULL,
    category_id UUID NOT NULL,
    month INT NOT NULL, -- 1-12
    year INT NOT NULL,
    budgeted_amount NUMERIC(15, 2) NOT NULL,
    spent_amount NUMERIC(15, 2) DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (household_id) REFERENCES households(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    UNIQUE(household_id, category_id, month, year)
);

CREATE INDEX idx_budgets_household ON budgets(household_id);
CREATE INDEX idx_budgets_period ON budgets(year, month);

-- ============================================================================
-- 9. AUDITORÍA
-- ============================================================================

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    household_id UUID NOT NULL,
    user_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (household_id) REFERENCES households(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_audit_logs_household ON audit_logs(household_id);
CREATE INDEX idx_audit_logs_date ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);

-- ============================================================================
-- Procedimientos Almacenados Útiles
-- ============================================================================

-- Función para actualizar saldo de cuenta automáticamente
CREATE OR REPLACE FUNCTION update_account_balance()
RETURNS TRIGGER AS $$
BEGIN
    -- Actualizar saldo basado en transacciones
    -- Se ejecutaría después de insertar/actualizar transacciones
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Función para actualizar last_updated en households
CREATE OR REPLACE FUNCTION update_household_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_households_updated_at
    BEFORE UPDATE ON households
    FOR EACH ROW
    EXECUTE FUNCTION update_household_timestamp();

-- ============================================================================
-- Datos Iniciales (Categorías por defecto)
-- ============================================================================

-- Se poplarán después de crear un household en la aplicación
-- Este schema es genérico y escalable
