import axios, { AxiosInstance } from 'axios';
import * as Types from '@types/index';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Agregar token a cada request
    this.client.interceptors.request.use((config) => {
      const token = localStorage.getItem('auth_token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Manejar errores globales
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('auth_token');
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }

  // ==================== Auth ====================

  async register(
    email: string,
    name: string,
    password: string,
    householdName: string
  ): Promise<Types.AuthResponse> {
    const response = await this.client.post('/auth/register', {
      email,
      name,
      password,
      household_name: householdName,
    });
    return response.data;
  }

  async login(email: string, password: string): Promise<Types.AuthResponse> {
    const response = await this.client.post('/auth/login', {
      email,
      password,
    });
    return response.data;
  }

  // ==================== Transacciones ====================

  async getTransactions(householdId: string): Promise<Types.Transaction[]> {
    const response = await this.client.get(
      `/households/${householdId}/transactions`
    );
    return response.data.transactions;
  }

  async getTransaction(
    householdId: string,
    transactionId: string
  ): Promise<Types.Transaction> {
    const response = await this.client.get(
      `/households/${householdId}/transactions/${transactionId}`
    );
    return response.data;
  }

  async createTransaction(
    householdId: string,
    data: Partial<Types.Transaction>
  ): Promise<{ id: string }> {
    const response = await this.client.post(
      `/households/${householdId}/transactions`,
      data
    );
    return response.data;
  }

  async updateTransaction(
    householdId: string,
    transactionId: string,
    data: Partial<Types.Transaction>
  ): Promise<void> {
    await this.client.put(
      `/households/${householdId}/transactions/${transactionId}`,
      data
    );
  }

  async deleteTransaction(
    householdId: string,
    transactionId: string
  ): Promise<void> {
    await this.client.delete(
      `/households/${householdId}/transactions/${transactionId}`
    );
  }

  // ==================== Divisas ====================

  async getCurrencyExchanges(
    householdId: string
  ): Promise<Types.CurrencyExchange[]> {
    const response = await this.client.get(
      `/households/${householdId}/currency-exchanges`
    );
    return response.data;
  }

  async createCurrencyExchange(
    householdId: string,
    data: Partial<Types.CurrencyExchange>
  ): Promise<{ id: string }> {
    const response = await this.client.post(
      `/households/${householdId}/currency-exchanges`,
      data
    );
    return response.data;
  }

  // ==================== Categorías ====================

  async getCategories(householdId: string): Promise<Types.Category[]> {
    const response = await this.client.get(
      `/households/${householdId}/categories`
    );
    return response.data;
  }

  async createCategory(
    householdId: string,
    data: Partial<Types.Category>
  ): Promise<{ id: string }> {
    const response = await this.client.post(
      `/households/${householdId}/categories`,
      data
    );
    return response.data;
  }

  // ==================== Servicios ====================

  async getRecurringExpenses(
    householdId: string
  ): Promise<Types.RecurringExpense[]> {
    const response = await this.client.get(
      `/households/${householdId}/recurring-expenses`
    );
    return response.data;
  }

  async createRecurringExpense(
    householdId: string,
    data: Partial<Types.RecurringExpense>
  ): Promise<{ id: string }> {
    const response = await this.client.post(
      `/households/${householdId}/recurring-expenses`,
      data
    );
    return response.data;
  }

  // ==================== Reportes ====================

  async getSummaryReport(householdId: string): Promise<Types.SummaryReport> {
    const response = await this.client.get(
      `/households/${householdId}/reports/summary`
    );
    return response.data;
  }

  async getMonthlyBreakdown(
    householdId: string
  ): Promise<Types.MonthlyBreakdown[]> {
    const response = await this.client.get(
      `/households/${householdId}/reports/monthly-breakdown`
    );
    return response.data;
  }

  async getCurrencyBalance(
    householdId: string
  ): Promise<Types.CurrencyBalanceReport> {
    const response = await this.client.get(
      `/households/${householdId}/reports/currency-balance`
    );
    return response.data;
  }

  // ==================== Usuarios ====================

  async getCurrentUser(): Promise<Types.User> {
    const response = await this.client.get('/me');
    return response.data;
  }

  async getHouseholdUsers(householdId: string): Promise<Types.User[]> {
    const response = await this.client.get(
      `/households/${householdId}/users`
    );
    return response.data;
  }
}

export const apiClient = new ApiClient();
