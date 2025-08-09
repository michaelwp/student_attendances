import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useForm } from 'react-hook-form';
import DataTable from '../components/DataTable';
import type { Column } from '../components/DataTable';
import Modal, { ConfirmModal } from '../components/Modal';
import ResetPasswordConfirmModal from '../components/ResetPasswordConfirmModal';
import PasswordResetModal from '../components/PasswordResetModal';
import UpdatePasswordModal from '../components/UpdatePasswordModal';
import { useToast } from '../components/Toast';
import { studentsApi, classesApi } from '../services/api';
import type { Student, StudentFormData, Class } from '../types/models';
import { passwordValidationRules, validatePassword, getPasswordStrength, isValidEmail } from '../utils/validation';

export const StudentsPage: React.FC = () => {
  const { t } = useTranslation();
  const { showSuccess, showError } = useToast();
  const [students, setStudents] = useState<Student[]>([]);
  const [classes, setClasses] = useState<Class[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [confirmModal, setConfirmModal] = useState<{
    isOpen: boolean;
    student: Student | null;
  }>({ isOpen: false, student: null });
  const [resetPasswordModal, setResetPasswordModal] = useState<{
    isOpen: boolean;
    student: Student | null;
  }>({ isOpen: false, student: null });
  const [passwordResetResult, setPasswordResetResult] = useState<{
    isOpen: boolean;
    newPassword: string;
    entityName: string;
  }>({ isOpen: false, newPassword: '', entityName: '' });
  const [updatePasswordModal, setUpdatePasswordModal] = useState<{
    isOpen: boolean;
    student: Student | null;
  }>({ isOpen: false, student: null });
  const [editingStudent, setEditingStudent] = useState<Student | null>(null);
  const [editingStudentId, setEditingStudentId] = useState<number | null>(null);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
    watch,
  } = useForm<StudentFormData>();
  
  const password = watch('password');
  const retypePassword = watch('retype_password');
  const [passwordStrength, setPasswordStrength] = useState(0);
  const [passwordValidationErrors, setPasswordValidationErrors] = useState<string[]>([]);
  
  useEffect(() => {
    if (password) {
      setPasswordValidationErrors(validatePassword(password));
    } else {
      setPasswordValidationErrors([]);
    }
  }, [password]);
  
  useEffect(() => {
    if (password) {
      setPasswordStrength(getPasswordStrength(password));
    } else {
      setPasswordStrength(0);
    }
  }, [password]);

  useEffect(() => {
    fetchStudents();
    fetchClasses();
  }, [pagination.current, pagination.pageSize]);

  const fetchStudents = async () => {
    try {
      setLoading(true);
      const response = await studentsApi.getAll({
        limit: pagination.pageSize,
        offset: (pagination.current - 1) * pagination.pageSize,
      });
      setStudents(response.data || []);
      setPagination(prev => ({
        ...prev,
        total: response.total || 0,
      }));
    } catch (error: any) {
      console.error('Failed to fetch students:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchClasses = async () => {
    try {
      const response = await classesApi.getAll({ limit: 100, offset: 0 });
      setClasses(response.data || []);
    } catch (error: any) {
      console.error('Failed to fetch classes:', error);
    }
  };

  const handleCreate = () => {
    setEditingStudent(null);
    setEditingStudentId(null);
    reset();
    setModalOpen(true);
  };

  const handleEdit = (student: Student) => {
    console.log('handleEdit called with student:', student);
    console.log('editingStudent before setting:', editingStudent);
    setEditingStudent(student);
    setEditingStudentId(student.id);
    console.log('editingStudent after setting (may not be immediately updated due to async state):', student);
    
    // Use reset() with values instead of individual setValue calls
    // This ensures the form is properly initialized for editing
    reset({
      student_id: student.student_id,
      classes_id: student.classes_id,
      first_name: student.first_name,
      last_name: student.last_name,
      email: student.email,
      phone: student.phone,
      // Don't include password fields for editing
    });
    
    setModalOpen(true);
    console.log('Modal opened, editingStudent should be:', student);
  };

  const handleDelete = (student: Student) => {
    setConfirmModal({ isOpen: true, student });
  };

  const handleResetPassword = (student: Student) => {
    setResetPasswordModal({ isOpen: true, student });
  };

  const handleUpdatePassword = (student: Student) => {
    setUpdatePasswordModal({ isOpen: true, student });
  };


  const handleConfirmDelete = async () => {
    if (!confirmModal.student) return;
    
    try {
      await studentsApi.delete(confirmModal.student.id);
      showSuccess(
        t('students.deleted_successfully'),
        t('students.student_deleted_message')
      );
      await fetchStudents();
      setConfirmModal({ isOpen: false, student: null });
    } catch (error: any) {
      console.error('Failed to delete student:', error);
      showError(
        t('common.error'),
        error?.message || t('students.delete_failed')
      );
    }
  };

  const handleConfirmResetPassword = async () => {
    if (!resetPasswordModal.student) return;
    
    try {
      const response = await studentsApi.resetPassword(resetPasswordModal.student.student_id);
      const entityName = `${resetPasswordModal.student.first_name} ${resetPasswordModal.student.last_name}`;
      
      setResetPasswordModal({ isOpen: false, student: null });
      setPasswordResetResult({
        isOpen: true,
        newPassword: response.newPassword,
        entityName,
      });
      
      await fetchStudents();
    } catch (error: any) {
      console.error('Failed to reset student password:', error);
      showError(
        t('common.error'),
        error?.message || t('students.password_reset_failed')
      );
    }
  };


  const onSubmit = async (data: StudentFormData) => {
    console.log('=== onSubmit CALLED ===');
    // Use editingStudentId as the reliable source of truth for edit mode
    const isEditMode = editingStudentId !== null;
    const currentEditingId = editingStudentId;

    try {
      // Remove retype_password before sending to API
      const { retype_password, ...apiData } = data;

      if (isEditMode && currentEditingId) {
        console.log('Updating student with ID:', currentEditingId);
        await studentsApi.update(currentEditingId, apiData);
        showSuccess(
          t('students.updated_successfully'),
          t('students.student_updated_message')
        );
      } else {
        console.log('Creating new student');
        await studentsApi.create(apiData);
        showSuccess(
          t('students.created_successfully'),
          t('students.student_created_message')
        );
      }

      await fetchStudents();
      setModalOpen(false);
      setEditingStudent(null); // Reset editing state after successful operation
      setEditingStudentId(null);
      reset();
    } catch (error: any) {
      console.error('Failed to save student:', error);
      showError(
        t('common.error'),
        error?.message || t('students.save_failed')
      );
    }
  };

  const getClassName = (classId: number) => {
    const classItem = classes.find(c => c.id === classId);
    return classItem ? classItem.name : `Class ${classId}`;
  };

  const getStatusBadge = (isActive: boolean) => {
    return isActive ? (
      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-900/50 dark:text-green-300">
        {t('common.active')}
      </span>
    ) : (
      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800 dark:bg-red-900/50 dark:text-red-300">
        {t('common.inactive')}
      </span>
    );
  };

  const columns: Column<Student>[] = [
    {
      key: 'student_id',
      title: t('students.student_id'),
      width: '32',
    },
    {
      key: 'first_name',
      title: t('students.name'),
      render: (_, record) => `${record.first_name} ${record.last_name}`,
    },
    {
      key: 'classes_id',
      title: t('students.class'),
      render: (value) => getClassName(value),
    },
    {
      key: 'email',
      title: t('students.email'),
    },
    {
      key: 'phone',
      title: t('students.phone'),
    },
    {
      key: 'is_active',
      title: t('students.status'),
      render: (value) => getStatusBadge(value),
    },
    {
      key: 'created_at',
      title: t('common.created_at'),
      render: (value) => new Date(value).toLocaleDateString(),
    },
  ];

  const modalFooter = (
    <>
      <button
        type="submit"
        form="student-form"
        disabled={isSubmitting}
        className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {isSubmitting ? (
          <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
        ) : null}
        {editingStudentId ? t('common.update') : t('common.create')}
      </button>
      <button
        type="button"
        onClick={() => {
          setModalOpen(false);
          setEditingStudent(null);
          setEditingStudentId(null);
          reset();
        }}
        disabled={isSubmitting}
        className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 dark:border-gray-600 shadow-sm px-4 py-2 bg-white dark:bg-gray-800 text-base font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {t('common.cancel')}
      </button>
    </>
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            {t('students.title')}
          </h1>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-300">
            {t('students.description')}
          </p>
        </div>
        <button
          onClick={handleCreate}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          <svg className="-ml-1 mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          {t('students.add_student')}
        </button>
      </div>

      {/* Data Table */}
      <DataTable
        data={students}
        columns={columns}
        loading={loading}
        onEdit={handleEdit}
        onDelete={handleDelete}
        customActions={(student) => (
          <>
            <button
              onClick={() => handleUpdatePassword(student)}
              className="text-green-600 hover:text-green-900 dark:text-green-400 dark:hover:text-green-300 mr-3"
              title={t('passwords.update_password')}
            >
              🔧
            </button>
            <button
              onClick={() => handleResetPassword(student)}
              className="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300 mr-3"
              title={t('common.reset_password')}
            >
              🔑
            </button>
          </>
        )}
        pagination={{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: pagination.total,
          onChange: (page, pageSize) => {
            setPagination(prev => ({
              ...prev,
              current: page,
              pageSize,
            }));
          },
        }}
      />

      {/* Create/Edit Modal */}
      <Modal
        isOpen={modalOpen}
        onClose={() => {
          console.log('Modal closing, resetting editing state');
          setModalOpen(false);
          setEditingStudent(null);
          setEditingStudentId(null);
          reset();
        }}
        title={editingStudentId ? t('students.edit_student') : t('students.add_student')}
        size="lg"
        footer={modalFooter}
      >
        <form id="student-form" onSubmit={(e) => {
          console.log('Form submission triggered');
          console.log('Form errors:', errors);
          console.log('Is submitting:', isSubmitting);
          handleSubmit(onSubmit)(e);
        }} className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('students.student_id')} *
              </label>
              <input
                type="text"
                {...register('student_id', { required: t('validation.required') })}
                className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder={t('students.student_id_placeholder')}
              />
              {errors.student_id && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.student_id.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('students.class')} *
              </label>
              <select
                {...register('classes_id', { 
                  required: t('validation.required'),
                  valueAsNumber: true,
                })}
                className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="">{t('students.select_class')}</option>
                {classes.map((classItem) => (
                  <option key={classItem.id} value={classItem.id}>
                    {classItem.name}
                  </option>
                ))}
              </select>
              {errors.classes_id && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.classes_id.message}</p>
              )}
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('students.first_name')} *
              </label>
              <input
                type="text"
                {...register('first_name', { required: t('validation.required') })}
                className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder={t('students.first_name_placeholder')}
              />
              {errors.first_name && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.first_name.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('students.last_name')} *
              </label>
              <input
                type="text"
                {...register('last_name', { required: t('validation.required') })}
                className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder={t('students.last_name_placeholder')}
              />
              {errors.last_name && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.last_name.message}</p>
              )}
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('students.email')} *
              </label>
              <input
                type="email"
                {...register('email', { 
                  required: t('validation.required'),
                  validate: (value) => {
                    if (!isValidEmail(value)) {
                      return t('validation.invalid_email');
                    }
                    return true;
                  },
                })}
                className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder={t('students.email_placeholder')}
              />
              {errors.email && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.email.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('students.phone')} *
              </label>
              <input
                type="tel"
                {...register('phone', { required: t('validation.required') })}
                className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder={t('students.phone_placeholder')}
              />
              {errors.phone && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.phone.message}</p>
              )}
            </div>
          </div>

          {!editingStudentId && (
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {t('students.password')} *
                </label>
                <input
                  type="password"
                  {...register('password', { 
                    required: !editingStudentId ? t('validation.required') : false,
                    validate: (value) => {
                      if (!editingStudentId && value) {
                        const errors = validatePassword(value);
                        return errors.length === 0 || errors.join('; ');
                      }
                      return true;
                    },
                  })}
                  className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder={t('students.password_placeholder')}
                />
                
                {password && !editingStudentId && (
                  <div className="mt-2">
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-gray-500">Password Strength</span>
                      <span className={`font-medium ${
                        passwordStrength < 40 ? 'text-red-500' : 
                        passwordStrength < 80 ? 'text-yellow-500' : 'text-green-500'
                      }`}>
                        {passwordStrength < 40 ? 'Weak' : 
                         passwordStrength < 80 ? 'Medium' : 'Strong'}
                      </span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-1.5 mt-1">
                      <div 
                        className={`h-1.5 rounded-full transition-all duration-300 ${
                          passwordStrength < 40 ? 'bg-red-500' : 
                          passwordStrength < 80 ? 'bg-yellow-500' : 'bg-green-500'
                        }`}
                        style={{ width: `${passwordStrength}%` }}
                      ></div>
                    </div>
                  </div>
                )}
                
                {passwordValidationErrors.length > 0 && (
                  <div className="mt-2 space-y-1">
                    {passwordValidationRules.map((rule, index) => {
                      const isValid = password ? rule.test(password) : false;
                      return (
                        <div key={index} className={`flex items-center text-xs ${
                          isValid ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'
                        }`}>
                          <span className="mr-2">{isValid ? '✓' : '✗'}</span>
                          <span>{rule.message}</span>
                        </div>
                      );
                    })}
                  </div>
                )}
                
                {errors.password && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.password.message}</p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {t('students.retype_password')} *
                </label>
                <input
                  type="password"
                  {...register('retype_password', { 
                    required: !editingStudentId ? t('validation.required') : false,
                    validate: (value) => {
                      if (!editingStudentId && value !== password) {
                        return t('validation.passwords_must_match');
                      }
                      return true;
                    },
                  })}
                  className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder={t('students.retype_password_placeholder')}
                />
                {retypePassword && retypePassword !== password && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">
                    {t('validation.passwords_must_match')}
                  </p>
                )}
                {errors.retype_password && (
                  <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.retype_password.message}</p>
                )}
              </div>
            </div>
          )}
        </form>
      </Modal>

      {/* Delete Confirmation Modal */}
      <ConfirmModal
        isOpen={confirmModal.isOpen}
        onClose={() => setConfirmModal({ isOpen: false, student: null })}
        onConfirm={handleConfirmDelete}
        title={t('students.delete_student')}
        message={t('students.delete_confirm', { 
          name: confirmModal.student ? `${confirmModal.student.first_name} ${confirmModal.student.last_name}` : '' 
        })}
        confirmText={t('common.delete')}
        type="danger"
      />

      {/* Reset Password Modal */}
      <ResetPasswordConfirmModal
        isOpen={resetPasswordModal.isOpen}
        onClose={() => setResetPasswordModal({ isOpen: false, student: null })}
        onConfirm={handleConfirmResetPassword}
        title={t('common.reset_password')}
        message={resetPasswordModal.student ? 
          t('students.reset_password_confirm', { 
            name: `${resetPasswordModal.student.first_name} ${resetPasswordModal.student.last_name}` 
          }) : ''
        }
      />

      {/* Password Reset Result Modal */}
      <PasswordResetModal
        isOpen={passwordResetResult.isOpen}
        onClose={() => setPasswordResetResult({ isOpen: false, newPassword: '', entityName: '' })}
        newPassword={passwordResetResult.newPassword}
        entityName={passwordResetResult.entityName}
      />

      {/* Update Password Modal */}
      <UpdatePasswordModal
        isOpen={updatePasswordModal.isOpen}
        onClose={() => setUpdatePasswordModal({ isOpen: false, student: null })}
        userType="student"
        userId={updatePasswordModal.student?.student_id || ''}
        userDisplayName={updatePasswordModal.student ? 
          `${updatePasswordModal.student.first_name} ${updatePasswordModal.student.last_name}` : ''
        }
      />
    </div>
  );
};

export default StudentsPage;