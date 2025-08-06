import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import type {AuthState, LoginRequest, User} from '../types/auth';
import { apiService, ApiError } from '../services/api';

interface AuthStore extends AuthState {
  initializeAuth: () => void;
}

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,

      initializeAuth: () => {
        const { user } = get();
        
        // Check if token is expired
        if (user && user.expires_at * 1000 < Date.now()) {
          // Token is expired, clear auth state
          set({
            user: null,
            isAuthenticated: false,
            error: null,
          });
          return;
        }

        // If user exists and token is valid, set auth token for API service
        if (user?.token) {
          apiService.setAuthToken(user.token);
          set({ isAuthenticated: true });
        }
      },

      login: async (credentials: LoginRequest) => {
        set({ isLoading: true, error: null });

        try {
          const response = await apiService.login(credentials);
          
          const user: User = {
            id: response.user_id,
            type: response.user_type,
            token: response.token,
            expires_at: response.expires_at,
          };

          // Set token for future API calls
          apiService.setAuthToken(response.token);

          set({
            user,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          });
        } catch (error) {
          let errorMessage = 'Login failed';
          
          if (error instanceof ApiError) {
            errorMessage = error.message;
          } else if (error instanceof Error) {
            errorMessage = error.message;
          }

          set({
            user: null,
            isAuthenticated: false,
            isLoading: false,
            error: errorMessage,
          });
          throw error;
        }
      },

      logout: async () => {
        set({ isLoading: true, error: null });

        try {
          // Call logout API if user is authenticated
          const { isAuthenticated } = get();
          if (isAuthenticated) {
            await apiService.logout();
          }
        } catch (error) {
          console.warn('Logout API call failed:', error);
          // Continue with local logout even if API call fails
        } finally {
          // Clear local auth state regardless of API call result
          set({
            user: null,
            isAuthenticated: false,
            isLoading: false,
            error: null,
          });
        }
      },

      clearError: () => {
        set({ error: null });
      },
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
      onRehydrateStorage: () => (state) => {
        // Initialize auth after rehydration
        if (state) {
          state.initializeAuth();
        }
      },
    }
  )
);