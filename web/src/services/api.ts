import type {LoginRequest, LoginResponse, ApiError as ApiErrorType} from '../types/auth';
import type {
  Teacher,
  Student,
  Class,
  Admin,
  Attendance,
  ApiResponse,
  PaginatedResponse,
  TeacherFormData,
  StudentFormData,
  ClassFormData,
  AdminFormData,
  AttendanceFormData,
  PasswordUpdateData,
  PasswordResetResponse,
  PhotoUploadResponse,
  PhotoUrlResponse,
  DashboardStats,
} from '../types/models';

const API_BASE_URL = 'http://localhost:8080/api/v1';

class ApiError extends Error {
  status: number;
  statusText: string;
  translationKey: string;

  constructor(
    status: number,
    statusText: string,
    translationKey: string,
    message: string
  ) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.statusText = statusText;
    this.translationKey = translationKey;
  }
}

class ApiService {
  private defaultHeaders: Record<string, string> = {};


  async login(credentials: LoginRequest): Promise<LoginResponse> {
    try {
      const response = await this.request<LoginResponse>('/auth/login', {
        method: 'POST',
        body: JSON.stringify(credentials),
      });
      return response;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  }

  async logout(): Promise<void> {
    try {
      await this.authenticatedRequest('/auth/logout', {
        method: 'POST',
      });
    } catch (error) {
      console.error('Logout error:', error);
      throw error;
    }
  }

  // Utility method to set authorization header for subsequent requests
  setAuthToken(token: string): void {
    this.defaultHeaders = {
      ...this.defaultHeaders,
      Authorization: `Bearer ${token}`,
    };
  }

  // Override request method to include auth headers
  private async authenticatedRequest<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const headers = {
      ...this.defaultHeaders,
      ...options.headers,
    };

    return this.request<T>(endpoint, {
      ...options,
      headers,
    });
  }

  // Make request public for API modules
  public async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    
    const defaultOptions: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...this.defaultHeaders,
        ...options.headers,
      },
      credentials: 'include', // Include cookies for authentication
      ...options,
    };

    try {
      const response = await fetch(url, defaultOptions);
      
      if (!response.ok) {
        let errorData: ApiErrorType;
        
        try {
          errorData = await response.json();
        } catch {
          // If JSON parsing fails, create a generic error
          throw new ApiError(
            response.status,
            response.statusText,
            'error.network_error',
            `HTTP ${response.status}: ${response.statusText}`
          );
        }

        throw new ApiError(
          response.status,
          response.statusText,
          errorData.translate_key || 'error.api_error',
          errorData.error || 'An error occurred'
        );
      }

      const data = await response.json();
      return data as T;
    } catch (error) {
      if (error instanceof ApiError) {
        throw error;
      }
      
      // Network error or other fetch errors
      if (error instanceof TypeError && error.message.includes('fetch')) {
        throw new ApiError(
          0,
          'Network Error',
          'error.network_error',
          'Unable to connect to the server. Please check your internet connection.'
        );
      }
      
      throw new ApiError(
        500,
        'Unknown Error',
        'error.unknown',
        error instanceof Error ? error.message : 'An unknown error occurred'
      );
    }
  }
}

export const apiService = new ApiService();
export { ApiError };

// Teachers API
export const teachersApi = {
  getAll: (params?: { limit?: number; offset?: number }) => {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.offset) searchParams.set('offset', params.offset.toString());
    const query = searchParams.toString();
    return apiService.request<PaginatedResponse<Teacher>>(
      `/teachers${query ? `?${query}` : ''}`
    );
  },
  getById: (id: number) =>
    apiService.request<ApiResponse<Teacher>>(`/teachers/${id}`),
  getByTeacherId: (teacherId: string) =>
    apiService.request<ApiResponse<Teacher>>(`/teachers/teacher/${teacherId}`),
  create: (data: TeacherFormData) =>
    apiService.request<ApiResponse<Teacher>>('/teachers', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  update: (id: number, data: Partial<TeacherFormData>) =>
    apiService.request<ApiResponse<Teacher>>(`/teachers/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  delete: (id: number) =>
    apiService.request<ApiResponse<null>>(`/teachers/${id}`, {
      method: 'DELETE',
    }),
  updatePassword: (teacherId: string, data: PasswordUpdateData) =>
    apiService.request<ApiResponse<null>>(`/teachers/${teacherId}/password`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  resetPassword: (id: number) =>
    apiService.request<PasswordResetResponse>(`/teachers/${id}/reset-password`, {
      method: 'POST',
    }),
  uploadPhoto: (id: number, file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return apiService.request<PhotoUploadResponse>(`/teachers/${id}/photo`, {
      method: 'POST',
      body: formData,
    });
  },
  getPhotoUrl: (id: number) =>
    apiService.request<PhotoUrlResponse>(`/teachers/${id}/photo-url`),
};

// Students API
export const studentsApi = {
  getAll: (params?: { limit?: number; offset?: number }) => {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.offset) searchParams.set('offset', params.offset.toString());
    const query = searchParams.toString();
    return apiService.request<PaginatedResponse<Student>>(
      `/students${query ? `?${query}` : ''}`
    );
  },
  getById: (id: number) =>
    apiService.request<ApiResponse<Student>>(`/students/${id}`),
  getByStudentId: (studentId: string) =>
    apiService.request<ApiResponse<Student>>(`/students/student/${studentId}`),
  getByClassId: (classId: number) =>
    apiService.request<ApiResponse<Student[]>>(`/students/class/${classId}`),
  create: (data: StudentFormData) =>
    apiService.request<ApiResponse<Student>>('/students', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  update: (id: number, data: Partial<StudentFormData>) =>
    apiService.request<ApiResponse<Student>>(`/students/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  delete: (id: number) =>
    apiService.request<ApiResponse<null>>(`/students/${id}`, {
      method: 'DELETE',
    }),
  updatePassword: (studentId: string, data: PasswordUpdateData) =>
    apiService.request<ApiResponse<null>>(`/students/${studentId}/password`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  resetPassword: (id: number) =>
    apiService.request<PasswordResetResponse>(`/students/${id}/reset-password`, {
      method: 'POST',
    }),
  uploadPhoto: (id: number, file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return apiService.request<PhotoUploadResponse>(`/students/${id}/photo`, {
      method: 'POST',
      body: formData,
    });
  },
  getPhotoUrl: (id: number) =>
    apiService.request<PhotoUrlResponse>(`/students/${id}/photo-url`),
};

// Classes API
export const classesApi = {
  getAll: (params?: { limit?: number; offset?: number }) => {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.offset) searchParams.set('offset', params.offset.toString());
    const query = searchParams.toString();
    return apiService.request<PaginatedResponse<Class>>(
      `/classes${query ? `?${query}` : ''}`
    );
  },
  getById: (id: number) =>
    apiService.request<ApiResponse<Class>>(`/classes/${id}`),
  create: (data: ClassFormData) =>
    apiService.request<ApiResponse<Class>>('/classes', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  update: (id: number, data: Partial<ClassFormData>) =>
    apiService.request<ApiResponse<Class>>(`/classes/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  delete: (id: number) =>
    apiService.request<ApiResponse<null>>(`/classes/${id}`, {
      method: 'DELETE',
    }),
};

// Attendance API
export const attendanceApi = {
  getAll: (params?: { limit?: number; offset?: number; date?: string; class_id?: number; student_id?: string }) => {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.offset) searchParams.set('offset', params.offset.toString());
    if (params?.date) searchParams.set('date', params.date);
    if (params?.class_id) searchParams.set('class_id', params.class_id.toString());
    if (params?.student_id) searchParams.set('student_id', params.student_id);
    const query = searchParams.toString();
    return apiService.request<PaginatedResponse<Attendance>>(
      `/attendance${query ? `?${query}` : ''}`
    );
  },
  getById: (id: number) =>
    apiService.request<ApiResponse<Attendance>>(`/attendance/${id}`),
  create: (data: AttendanceFormData) =>
    apiService.request<ApiResponse<Attendance>>('/attendance', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  update: (id: number, data: Partial<AttendanceFormData>) =>
    apiService.request<ApiResponse<Attendance>>(`/attendance/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  delete: (id: number) =>
    apiService.request<ApiResponse<null>>(`/attendance/${id}`, {
      method: 'DELETE',
    }),
  getByStudentId: (studentId: string, date?: string) => {
    const searchParams = new URLSearchParams();
    if (date) searchParams.set('date', date);
    const query = searchParams.toString();
    return apiService.request<ApiResponse<Attendance[]>>(
      `/attendance/student/${studentId}${query ? `?${query}` : ''}`
    );
  },
  getByClassId: (classId: number, date?: string) => {
    const searchParams = new URLSearchParams();
    if (date) searchParams.set('date', date);
    const query = searchParams.toString();
    return apiService.request<ApiResponse<Attendance[]>>(
      `/attendance/class/${classId}${query ? `?${query}` : ''}`
    );
  },
  getByDate: (date: string) =>
    apiService.request<ApiResponse<Attendance[]>>(`/attendance/date/${date}`),
};

// Admins API
export const adminsApi = {
  getAll: (params?: { limit?: number; offset?: number }) => {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.offset) searchParams.set('offset', params.offset.toString());
    const query = searchParams.toString();
    return apiService.request<PaginatedResponse<Admin>>(
      `/admins${query ? `?${query}` : ''}`
    );
  },
  getById: (id: number) =>
    apiService.request<ApiResponse<Admin>>(`/admins/${id}`),
  getByEmail: (email: string) =>
    apiService.request<ApiResponse<Admin>>(`/admins/email/${encodeURIComponent(email)}`),
  create: (data: AdminFormData) =>
    apiService.request<ApiResponse<Admin>>('/admins', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  update: (id: number, data: Partial<AdminFormData>) =>
    apiService.request<ApiResponse<Admin>>(`/admins/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  delete: (id: number) =>
    apiService.request<ApiResponse<null>>(`/admins/${id}`, {
      method: 'DELETE',
    }),
  updatePassword: (email: string, data: PasswordUpdateData) =>
    apiService.request<ApiResponse<null>>(`/admins/${encodeURIComponent(email)}/password`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  setActiveStatus: (id: number, isActive: boolean) =>
    apiService.request<ApiResponse<null>>(`/admins/${id}/status`, {
      method: 'PUT',
      body: JSON.stringify({ is_active: isActive }),
    }),
};

// Dashboard Stats API
export const dashboardApi = {
  getStats: () =>
    apiService.request<ApiResponse<DashboardStats>>('/admins/stats'),
};