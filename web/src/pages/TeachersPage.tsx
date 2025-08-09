import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useForm, useWatch } from 'react-hook-form';
import DataTable from '../components/DataTable';
import type { Column } from '../components/DataTable';
import Modal, { ConfirmModal } from '../components/Modal';
import ResetPasswordConfirmModal from '../components/ResetPasswordConfirmModal';
import PasswordResetModal from '../components/PasswordResetModal';
import UpdatePasswordModal from '../components/UpdatePasswordModal';
import { useToast } from '../components/Toast';
import { teachersApi } from '../services/api';
import type { Teacher, TeacherFormData } from '../types/models';
import { passwordValidationRules, validatePassword, getPasswordStrength, isValidEmail } from '../utils/validation';

export const TeachersPage: React.FC = () => {
  const { t } = useTranslation();
  const { showSuccess, showError } = useToast();
  const [teachers, setTeachers] = useState<Teacher[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [confirmModal, setConfirmModal] = useState<{
    isOpen: boolean;
    teacher: Teacher | null;
  }>({ isOpen: false, teacher: null });
  const [resetPasswordModal, setResetPasswordModal] = useState<{
    isOpen: boolean;
    teacher: Teacher | null;
  }>({ isOpen: false, teacher: null });
  const [passwordResetResult, setPasswordResetResult] = useState<{
    isOpen: boolean;
    newPassword: string;
    entityName: string;
  }>({ isOpen: false, newPassword: '', entityName: '' });
  const [updatePasswordModal, setUpdatePasswordModal] = useState<{
    isOpen: boolean;
    teacher: Teacher | null;
  }>({ isOpen: false, teacher: null });
  const [editingTeacher, setEditingTeacher] = useState<Teacher | null>(null);
  const [editingTeacherId, setEditingTeacherId] = useState<number | null>(null);
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
    control,
  } = useForm<TeacherFormData>();

  const password = useWatch({ control, name: 'password' }) || '';
  const retypePassword = useWatch({ control, name: 'retype_password' }) || '';
  const [passwordValidationErrors, setPasswordValidationErrors] = useState<string[]>([]);
  const [passwordStrength, setPasswordStrength] = useState(0);

  useEffect(() => {
    if (password && !editingTeacherId) {
      const errors = validatePassword(password);
      setPasswordValidationErrors(errors);
      setPasswordStrength(getPasswordStrength(password));
    } else {
      setPasswordValidationErrors([]);
      setPasswordStrength(0);
    }
  }, [password, editingTeacherId]);

  useEffect(() => {
    fetchTeachers();
  }, [pagination.current, pagination.pageSize]);

  // Debug effect to track editingTeacher changes
  useEffect(() => {
    console.log('editingTeacher changed to:', editingTeacher);
  }, [editingTeacher]);

  // Debug effect to track editingTeacherId changes
  useEffect(() => {
    console.log('editingTeacherId changed to:', editingTeacherId);
  }, [editingTeacherId]);

  const fetchTeachers = async () => {
    try {
      setLoading(true);
      const response = await teachersApi.getAll({
        limit: pagination.pageSize,
        offset: (pagination.current - 1) * pagination.pageSize,
      });
      setTeachers(response.data || []);
      setPagination(prev => ({
        ...prev,
        total: response.total || 0,
      }));
    } catch (error: any) {
      console.error('Failed to fetch teachers:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingTeacher(null);
    setEditingTeacherId(null);
    reset();
    setModalOpen(true);
  };

  const handleEdit = (teacher: Teacher) => {
    console.log('handleEdit called with teacher:', teacher);
    console.log('editingTeacher before setting:', editingTeacher);
    setEditingTeacher(teacher);
    setEditingTeacherId(teacher.id);
    console.log('editingTeacher after setting (may not be immediately updated due to async state):', teacher);
    
    // Use reset() with values instead of individual setValue calls
    // This ensures the form is properly initialized for editing
    reset({
      teacher_id: teacher.teacher_id,
      first_name: teacher.first_name,
      last_name: teacher.last_name,
      email: teacher.email,
      phone: teacher.phone,
      // Don't include password fields for editing
    });
    
    setModalOpen(true);
    console.log('Modal opened, editingTeacher should be:', teacher);
  };

  const handleDelete = (teacher: Teacher) => {
    setConfirmModal({ isOpen: true, teacher });
  };

  const handleResetPassword = (teacher: Teacher) => {
    setResetPasswordModal({ isOpen: true, teacher });
  };

  const handleUpdatePassword = (teacher: Teacher) => {
    setUpdatePasswordModal({ isOpen: true, teacher });
  };


  const handleConfirmDelete = async () => {
    if (!confirmModal.teacher) return;
    
    try {
      await teachersApi.delete(confirmModal.teacher.id);
      showSuccess(
        t('teachers.deleted_successfully'),
        t('teachers.teacher_deleted_message')
      );
      await fetchTeachers();
      setConfirmModal({ isOpen: false, teacher: null });
    } catch (error: any) {
      console.error('Failed to delete teacher:', error);
      showError(
        t('common.error'),
        error?.message || t('teachers.delete_failed')
      );
    }
  };

  const handleConfirmResetPassword = async () => {
    if (!resetPasswordModal.teacher) return;
    
    try {
      const response = await teachersApi.resetPassword(resetPasswordModal.teacher.teacher_id);
      const entityName = `${resetPasswordModal.teacher.first_name} ${resetPasswordModal.teacher.last_name}`;
      
      setResetPasswordModal({ isOpen: false, teacher: null });
      setPasswordResetResult({
        isOpen: true,
        newPassword: response.newPassword,
        entityName,
      });
      
      await fetchTeachers();
    } catch (error: any) {
      console.error('Failed to reset teacher password:', error);
      showError(
        t('common.error'),
        error?.message || t('teachers.password_reset_failed')
      );
    }
  };


  const onSubmit = async (data: TeacherFormData) => {
    console.log('=== onSubmit CALLED ===');
    // Use editingTeacherId as the reliable source of truth for edit mode
    const isEditMode = editingTeacherId !== null;
    const currentEditingId = editingTeacherId;
    
    try {
      // Remove retype_password before sending to API
      const { retype_password, ...apiData } = data;
      
      if (isEditMode && currentEditingId) {
        console.log('Updating teacher with ID:', currentEditingId);
        await teachersApi.update(currentEditingId, apiData);
        showSuccess(
          t('teachers.updated_successfully'),
          t('teachers.teacher_updated_message')
        );
      } else {
        console.log('Creating new teacher');
        await teachersApi.create(apiData);
        showSuccess(
          t('teachers.created_successfully'),
          t('teachers.teacher_created_message')
        );
      }
      
      await fetchTeachers();
      setModalOpen(false);
      setEditingTeacher(null); // Reset editing state after successful operation
      setEditingTeacherId(null);
      reset();
    } catch (error: any) {
      console.error('Failed to save teacher:', error);
      showError(
        t('common.error'),
        error?.message || t('teachers.save_failed')
      );
    }
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

  const columns: Column<Teacher>[] = [
    {
      key: 'teacher_id',
      title: t('teachers.teacher_id'),
      width: '32',
    },
    {
      key: 'first_name',
      title: t('teachers.first_name'),
      render: (_, record) => `${record.first_name} ${record.last_name}`,
    },
    {
      key: 'email',
      title: t('teachers.email'),
    },
    {
      key: 'phone',
      title: t('teachers.phone'),
    },
    {
      key: 'is_active',
      title: t('teachers.status'),
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
        form="teacher-form"
        disabled={isSubmitting}
        className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {isSubmitting ? (
          <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
        ) : null}
        {editingTeacherId ? t('common.update') : t('common.create')}
      </button>
      <button
        type="button"
        onClick={() => {
          setModalOpen(false);
          setEditingTeacher(null);
          setEditingTeacherId(null);
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
            {t('teachers.title')}
          </h1>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-300">
            {t('teachers.description')}
          </p>
        </div>
        <button
          onClick={handleCreate}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          <svg className="-ml-1 mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          {t('teachers.add_teacher')}
        </button>
      </div>

      {/* Data Table */}
      <DataTable
        data={teachers}
        columns={columns}
        loading={loading}
        onEdit={handleEdit}
        onDelete={handleDelete}
        customActions={(teacher) => (
          <>
            <button
              onClick={() => handleUpdatePassword(teacher)}
              className="text-green-600 hover:text-green-900 dark:text-green-400 dark:hover:text-green-300 mr-3"
              title={t('passwords.update_password')}
            >
              🔧
            </button>
            <button
              onClick={() => handleResetPassword(teacher)}
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
          setEditingTeacher(null);
          setEditingTeacherId(null);
          reset();
        }}
        title={editingTeacherId ? t('teachers.edit_teacher') : t('teachers.add_teacher')}
        size="lg"
        footer={modalFooter}
      >
        <form id="teacher-form" onSubmit={(e) => {
          console.log('Form submission triggered');
          console.log('Form errors:', errors);
          console.log('Is submitting:', isSubmitting);
          handleSubmit(onSubmit)(e);
        }} className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('teachers.teacher_id')} *
              </label>
              <input
                type="text"
                {...register('teacher_id', { required: t('validation.required') })}
                className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder={t('teachers.teacher_id_placeholder')}
              />
              {errors.teacher_id && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.teacher_id.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('teachers.email')} *
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
                placeholder={t('teachers.email_placeholder')}
              />
              {errors.email && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.email.message}</p>
              )}
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('teachers.first_name')} *
              </label>
              <input
                type="text"
                {...register('first_name', { required: t('validation.required') })}
                className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder={t('teachers.first_name_placeholder')}
              />
              {errors.first_name && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.first_name.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('teachers.last_name')} *
              </label>
              <input
                type="text"
                {...register('last_name', { required: t('validation.required') })}
                className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder={t('teachers.last_name_placeholder')}
              />
              {errors.last_name && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.last_name.message}</p>
              )}
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('teachers.phone')} *
              </label>
              <input
                type="tel"
                {...register('phone', { required: t('validation.required') })}
                className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder={t('teachers.phone_placeholder')}
              />
              {errors.phone && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.phone.message}</p>
              )}
            </div>

          </div>

          {!editingTeacherId && (
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {t('teachers.password')} *
                </label>
                <input
                  type="password"
                  {...register('password', { 
                    required: !editingTeacherId ? t('validation.required') : false,
                    validate: (value) => {
                      if (!editingTeacherId && value) {
                        const errors = validatePassword(value);
                        return errors.length === 0 || errors.join('; ');
                      }
                      return true;
                    },
                  })}
                  className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder={t('teachers.password_placeholder')}
                />
                
                {password && !editingTeacherId && (
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
                      const isValid = rule.test(password);
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
                  {t('teachers.retype_password')} *
                </label>
                <input
                  type="password"
                  {...register('retype_password', { 
                    required: !editingTeacherId ? t('validation.required') : false,
                    validate: (value) => {
                      if (!editingTeacherId && value !== password) {
                        return t('validation.passwords_must_match');
                      }
                      return true;
                    },
                  })}
                  className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder={t('teachers.retype_password_placeholder')}
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
        onClose={() => setConfirmModal({ isOpen: false, teacher: null })}
        onConfirm={handleConfirmDelete}
        title={t('teachers.delete_teacher')}
        message={t('teachers.delete_confirm', { 
          name: confirmModal.teacher ? `${confirmModal.teacher.first_name} ${confirmModal.teacher.last_name}` : '' 
        })}
        confirmText={t('common.delete')}
        type="danger"
      />

      {/* Reset Password Modal */}
      <ResetPasswordConfirmModal
        isOpen={resetPasswordModal.isOpen}
        onClose={() => setResetPasswordModal({ isOpen: false, teacher: null })}
        onConfirm={handleConfirmResetPassword}
        title={t('common.reset_password')}
        message={resetPasswordModal.teacher ? 
          t('teachers.reset_password_confirm', { 
            name: `${resetPasswordModal.teacher.first_name} ${resetPasswordModal.teacher.last_name}` 
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
        onClose={() => setUpdatePasswordModal({ isOpen: false, teacher: null })}
        userType="teacher"
        userId={updatePasswordModal.teacher?.teacher_id || ''}
        userDisplayName={updatePasswordModal.teacher ? 
          `${updatePasswordModal.teacher.first_name} ${updatePasswordModal.teacher.last_name}` : ''
        }
      />
    </div>
  );
};

export default TeachersPage;