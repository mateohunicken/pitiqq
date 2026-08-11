// Tipos para la aplicación Finanzas Domésticas

export interface User {
  id: string;
  email: string;
  name: string;
  role: 'admin' | 'editor' | 'viewer';
  household_id: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Household {
  id: string;
  name: string;
  description?: string;
  currency: string;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: string;
  household_id: string;
  name: string;
  description?: string;
  type: 'income' | 'expense' | 'transfer';
  color?: string;
  icon?: string;
  is_custom: boolean;
  created_at: string;
}

export interface Transaction {
  id: string;
  household_id: string;
  user_id: string;
  category_id: string;
  type: 'income' | 'expense' | 'transfer';
  description: string;
  amount: number;
  currency: string;
  transaction_date: string;
  payment_method?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
  category?: Category;
  user?: User;
}

export interface CurrencyExchange {
  id: string;
  household_id: string;
  user_id: string;
  operation_type: 'buy' | 'sell';
  amount_local: number;
  amount_foreign: number;
  exchange_rate: number;
  foreign_currency: string;
  exchange_date: string;
  location?: string;
  notes?: string;
  created_at: string;
}

export interface RecurringExpense {
  id: string;
  household_id: string;
  name: string;
  description?: string;
  amount: number;
  currency: string;
  category_id?: string;
  frequency: 'monthly' | 'quarterly' | 'annual' | 'bimonthly';
  due_day?: number;
  due_month?: number;
  status: 'active' | 'inactive' | 'paid';
  last_paid_date?: string;
  next_due_date: string;
  created_at: string;
  updated_at: string;
}

export interface SummaryReport {
  total_income: number;
  total_expense: number;
  net_balance: number;
  period: string;
  currency_balances: Record<string, number>;
}

export interface MonthlyBreakdown {
  month: number;
  year: number;
  income: number;
  expense: number;
  balance: number;
  category_breakdown: Record<string, CategoryStat>;
}

export interface CategoryStat {
  category_name: string;
  amount: number;
  percentage: number;
}

export interface CurrencyBalanceReport {
  ars: number;
  usd: number;
  total_ars: number;
  total_usd: number;
}

export interface AuthResponse {
  token: string;
  user: User;
  household: Household;
}

export interface ApiError {
  error: string;
}
