import type { TeacherFormData, StudentFormData, AdminFormData } from '../types/models';

// Helper functions to transform form data for API calls

export const prepareTeacherDataForAPI = (data: TeacherFormData, isUpdate: boolean = false): any => {
  const { retype_password, password, ...apiData } = data;
  
  // Handle phone field - send null if empty for backend compatibility (*string type)
  const phone = data.phone && data.phone.trim() !== '' ? data.phone.trim() : null;
  
  // Create a properly typed object with all required fields
  const cleanData = {
    teacher_id: apiData.teacher_id,
    first_name: apiData.first_name,
    last_name: apiData.last_name,
    email: apiData.email,
    phone: phone, // Send null or string value, not empty string
    // Only include password for creation, not updates
    ...(password && !isUpdate && { password: password })
  };
  
  console.log('Preparing teacher data for API:', { original: data, cleaned: cleanData, isUpdate });
  
  return cleanData;
};

export const prepareStudentDataForAPI = (data: StudentFormData, isUpdate: boolean = false): any => {
  const { retype_password, password, ...apiData } = data;
  
  // Handle phone field - send null if empty for backend compatibility (*string type)
  const phone = data.phone && data.phone.trim() !== '' ? data.phone.trim() : null;
  
  // Create a properly typed object with all required fields
  const cleanData = {
    student_id: apiData.student_id,
    classes_id: apiData.classes_id,
    first_name: apiData.first_name,
    last_name: apiData.last_name,
    email: apiData.email,
    phone: phone, // Send null or string value, not empty string
    // Only include password for creation, not updates
    ...(password && !isUpdate && { password: password })
  };
  
  console.log('Preparing student data for API:', { original: data, cleaned: cleanData, isUpdate });
  
  return cleanData;
};

export const prepareAdminDataForAPI = (data: AdminFormData) => {
  const { retype_password, ...apiData } = data;
  return apiData;
};