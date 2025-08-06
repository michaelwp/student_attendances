export type UserType = 'admin' | 'teacher' | 'student';

export interface LoginRequest {
  user_type: UserType;
  user_id: string;
  password: string;
}

export interface LoginResponse {
  translate_key: string;
  message: string;
  token: string;
  user_type: UserType;
  user_id: string;
  expires_at: number;
}

export interface ApiError {
  translate_key: string;
  error: string;
}

export interface User {
  id: string;
  type: UserType;
  token: string;
  expires_at: number;
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
  clearError: () => void;
}