import { create } from 'zustand';
import * as Types from '@types/index';
import { apiClient } from '@services/api';

interface AuthStore {
  user: Types.User | null;
  household: Types.Household | null;
  token: string | null;
  isLoading: boolean;
  error: string | null;
  isAuthenticated: boolean;

  // Actions
  login: (email: string, password: string) => Promise<void>;
  register: (
    email: string,
    name: string,
    password: string,
    householdName: string
  ) => Promise<void>;
  logout: () => void;
  checkAuth: () => void;
  clearError: () => void;
}

export const useAuthStore = create<AuthStore>((set) => ({
  user: null,
  household: null,
  token: null,
  isLoading: false,
  error: null,
  isAuthenticated: false,

  login: async (email: string, password: string) => {
    set({ isLoading: true, error: null });
    try {
      const response = await apiClient.login(email, password);
      localStorage.setItem('auth_token', response.token);
      set({
        token: response.token,
        user: response.user,
        household: response.household,
        isAuthenticated: true,
        isLoading: false,
      });
    } catch (error: any) {
      set({
        error: error.response?.data?.error || 'Login failed',
        isLoading: false,
      });
      throw error;
    }
  },

  register: async (
    email: string,
    name: string,
    password: string,
    householdName: string
  ) => {
    set({ isLoading: true, error: null });
    try {
      const response = await apiClient.register(
        email,
        name,
        password,
        householdName
      );
      localStorage.setItem('auth_token', response.token);
      set({
        token: response.token,
        user: response.user,
        household: response.household,
        isAuthenticated: true,
        isLoading: false,
      });
    } catch (error: any) {
      set({
        error: error.response?.data?.error || 'Registration failed',
        isLoading: false,
      });
      throw error;
    }
  },

  logout: () => {
    localStorage.removeItem('auth_token');
    set({
      user: null,
      household: null,
      token: null,
      isAuthenticated: false,
    });
  },

  checkAuth: () => {
    const token = localStorage.getItem('auth_token');
    if (!token) {
      set({ isAuthenticated: false });
      return;
    }

    set({ token, isAuthenticated: true });
  },

  clearError: () => set({ error: null }),
}));
