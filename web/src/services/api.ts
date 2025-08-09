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
    // Convert headers to Record<string, string> to ensure type safety
    const headers = {
      ...this.defaultHeaders,
      ...(options.headers as Record<string, string> || {}),
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
    
    // Don't set Content-Type for FormData - browser sets it automatically with boundary
    const isFormData = options.body instanceof FormData;
    
    // Convert headers to Record<string, string> to ensure type safety
    const defaultHeaders: Record<string, string> = {
      ...this.defaultHeaders,
      ...(options.headers as Record<string, string> || {}),
    };
    
    // Only set Content-Type for non-FormData requests
    if (!isFormData) {
      defaultHeaders['Content-Type'] = 'application/json';
    }
    
    const defaultOptions: RequestInit = {
      headers: defaultHeaders,
      credentials: 'include', // Include cookies for authentication
      ...options,
    };

    try {
      console.log('Making request to:', url);
      console.log('Request options:', {
        method: defaultOptions.method || 'GET',
        headers: defaultOptions.headers,
        bodyType: defaultOptions.body ? typeof defaultOptions.body : 'none',
        isFormData: isFormData
      });
      
      const response = await fetch(url, defaultOptions);
      
      console.log('Response status:', response.status, response.statusText);
      
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
      `/teachers/all${query ? `?${query}` : ''}`
    );
  },
  getById: (id: number) =>
    apiService.request<ApiResponse<Teacher>>(`/teachers/record-id/${id}`),
  getByTeacherId: (teacherId: string) =>
    apiService.request<ApiResponse<Teacher>>(`/teachers/teacher-id/${teacherId}`),
  create: (data: TeacherFormData) =>
    apiService.request<ApiResponse<Teacher>>('/teachers', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  update: (id: number, data: Partial<TeacherFormData>) =>
    apiService.request<ApiResponse<Teacher>>(`/teachers/record-id/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  delete: (id: number) =>
    apiService.request<ApiResponse<null>>(`/teachers/record-id/${id}`, {
      method: 'DELETE',
    }),
  updatePassword: (teacherId: string, data: PasswordUpdateData) =>
    apiService.request<ApiResponse<null>>(`/teachers/teacher-id/${teacherId}/password`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  resetPassword: (teacherId: string) =>
    apiService.request<PasswordResetResponse>(`/teachers/teacher-id/${teacherId}/reset-password`, {
      method: 'PUT',
    }),
};

// Students API
export const studentsApi = {
  getAll: (params?: { limit?: number; offset?: number }) => {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.offset) searchParams.set('offset', params.offset.toString());
    const query = searchParams.toString();
    return apiService.request<PaginatedResponse<Student>>(
      `/students/all${query ? `?${query}` : ''}`
    );
  },
  getById: (id: number) =>
    apiService.request<ApiResponse<Student>>(`/students/record-id/${id}`),
  getByStudentId: (studentId: string) =>
    apiService.request<ApiResponse<Student>>(`/students/student-id/${studentId}`),
  getByClassId: (classId: number) =>
    apiService.request<ApiResponse<Student[]>>(`/students/class-id/${classId}`),
  create: (data: StudentFormData) =>
    apiService.request<ApiResponse<Student>>('/students', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  update: (id: number, data: Partial<StudentFormData>) =>
    apiService.request<ApiResponse<Student>>(`/students/record-id/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  delete: (id: number) =>
    apiService.request<ApiResponse<null>>(`/students/record-id/${id}`, {
      method: 'DELETE',
    }),
  updatePassword: (studentId: string, data: PasswordUpdateData) =>
    apiService.request<ApiResponse<null>>(`/students/student-id/${studentId}/password`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  resetPassword: (studentId: string) =>
    apiService.request<PasswordResetResponse>(`/students/student-id/${studentId}/reset-password`, {
      method: 'PUT',
    }),
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
  getAll: (params?: { limit?: number; offset?: number; class_id?: number | string; date?: string; status?: string }) => {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.offset) searchParams.set('offset', params.offset.toString());
    if (params?.class_id) searchParams.set('class_id', params.class_id.toString());
    if (params?.date) searchParams.set('date', params.date);
    if (params?.status) searchParams.set('status', params.status);
    const query = searchParams.toString();
    return apiService.request<PaginatedResponse<Attendance>>(
      `/attendances/all${query ? `?${query}` : ''}`
    );
  },
  getById: (id: number) =>
    apiService.request<ApiResponse<Attendance>>(`/attendances/attendances-id/${id}`),
  create: (data: AttendanceFormData) =>
    apiService.request<ApiResponse<Attendance>>('/attendances', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  update: (id: number, data: Partial<AttendanceFormData>) =>
    apiService.request<ApiResponse<Attendance>>(`/attendances/attendances-id/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  delete: (id: number) =>
    apiService.request<ApiResponse<null>>(`/attendances/attendances-id/${id}`, {
      method: 'DELETE',
    }),
  getByStudentId: (studentId: string, params?: { date?: string; limit?: number; offset?: number }) => {
    const searchParams = new URLSearchParams();
    if (params?.date) searchParams.set('date', params.date);
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.offset) searchParams.set('offset', params.offset.toString());
    const query = searchParams.toString();
    return apiService.request<ApiResponse<Attendance[]>>(
      `/attendances/student-id/${studentId}${query ? `?${query}` : ''}`
    );
  },
  getByClassId: (classId: number, params?: { date?: string; limit?: number; offset?: number }) => {
    const searchParams = new URLSearchParams();
    if (params?.date) searchParams.set('date', params.date);
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.offset) searchParams.set('offset', params.offset.toString());
    const query = searchParams.toString();
    return apiService.request<ApiResponse<Attendance[]>>(
      `/attendances/class-id/${classId}${query ? `?${query}` : ''}`
    );
  },
  getByDateRange: (startDate: string, endDate: string, params?: { limit?: number; offset?: number }) => {
    const searchParams = new URLSearchParams();
    searchParams.set('start_date', startDate);
    searchParams.set('end_date', endDate);
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.offset) searchParams.set('offset', params.offset.toString());
    const query = searchParams.toString();
    return apiService.request<ApiResponse<Attendance[]>>(
      `/attendances/date-range?${query}`
    );
  },
};

// Admins API
export const adminsApi = {
  getAll: (params?: { limit?: number; offset?: number }) => {
    const searchParams = new URLSearchParams();
    if (params?.limit) searchParams.set('limit', params.limit.toString());
    if (params?.offset) searchParams.set('offset', params.offset.toString());
    const query = searchParams.toString();
    return apiService.request<PaginatedResponse<Admin>>(
      `/admins/all${query ? `?${query}` : ''}`
    );
  },
  getById: (id: number) =>
    apiService.request<ApiResponse<Admin>>(`/admins/admin-id/${id}`),
  getByEmail: (email: string) =>
    apiService.request<ApiResponse<Admin>>(`/admins/email/${encodeURIComponent(email)}`),
  create: (data: AdminFormData) =>
    apiService.request<ApiResponse<Admin>>('/admins', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  update: (id: number, data: Partial<AdminFormData>) =>
    apiService.request<ApiResponse<Admin>>(`/admins/admin-id/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  delete: (id: number) =>
    apiService.request<ApiResponse<null>>(`/admins/admin-id/${id}`, {
      method: 'DELETE',
    }),
  updatePassword: (id: number, data: PasswordUpdateData) =>
    apiService.request<ApiResponse<null>>(`/admins/admin-id/${id}/password`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  setActiveStatus: (id: number, isActive: boolean) =>
    apiService.request<ApiResponse<null>>(`/admins/admin-id/${id}/status`, {
      method: 'PUT',
      body: JSON.stringify({ is_active: isActive }),
    }),
};

// Dashboard Stats API
export const dashboardApi = {
  getStats: () =>
    apiService.request<ApiResponse<DashboardStats>>('/admins/stats'),
};

// Student Attendance Marking API (public endpoint, no auth required)
export const studentAttendanceApi = {
  markAttendance: (data: { student_id: string; password: string }) =>
    apiService.request<{ student_name: string; message: string }>('/attendance/mark', {
      method: 'POST',
      body: JSON.stringify(data),
    }),
};