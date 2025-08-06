// Base model interface
export interface BaseModel {
  id: number;
  created_at: string;
  updated_at: string;
}

// Teacher model
export interface Teacher extends BaseModel {
  teacher_id: string;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  photo_path?: string;
}

// Student model
export interface Student extends BaseModel {
  student_id: string;
  classes_id: number;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  photo_path?: string;
}

// Class model
export interface Class extends BaseModel {
  name: string;
  homeroom_teacher: string;
  description: string;
}

// Admin model
export interface Admin extends BaseModel {
  email: string;
  last_login?: string;
  is_active: boolean;
}

// Attendance model
export interface Attendance extends BaseModel {
  student_id: string;
  class_id: number;
  date: string;
  status: AttendanceStatus;
  description?: string;
}

// Attendance status enum
export type AttendanceStatus = 'present' | 'absent' | 'late';

// API Response wrappers
export interface ApiResponse<T> {
  translate_key: string;
  message: string;
  data?: T;
  total?: number;
  limit?: number;
  offset?: number;
}

export interface PaginatedResponse<T> {
  translate_key: string;
  message: string;
  data: T[];
  total: number;
  limit: number;
  offset: number;
}

// Form data interfaces for creating/updating
export interface TeacherFormData {
  teacher_id: string;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  password?: string;
}

export interface StudentFormData {
  student_id: string;
  classes_id: number;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  password?: string;
}

export interface ClassFormData {
  name: string;
  homeroom_teacher: string;
  description: string;
}

export interface AdminFormData {
  email: string;
  password?: string;
  is_active: boolean;
}

export interface AttendanceFormData {
  student_id: string;
  class_id: number;
  date: string;
  status: AttendanceStatus;
  description?: string;
}

// Password update interfaces
export interface PasswordUpdateData {
  old_password: string;
  new_password: string;
}

export interface PasswordResetResponse {
  translate_key: string;
  message: string;
  newPassword: string;
}

// Photo upload response
export interface PhotoUploadResponse {
  translate_key: string;
  message: string;
  path: string;
}

export interface PhotoUrlResponse {
  translate_key: string;
  message: string;
  url: string;
}

// Statistics for dashboard
export interface DashboardStats {
  total_teachers: number;
  total_students: number;
  total_classes: number;
  total_admins: number;
  total_attendance_today?: number;
  present_today?: number;
  absent_today?: number;
  late_today?: number;
}