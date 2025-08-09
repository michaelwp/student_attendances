import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useForm, useWatch } from 'react-hook-form';
import DataTable from '../components/DataTable';
import type { Column } from '../components/DataTable';
import Modal, { ConfirmModal } from '../components/Modal';
import UpdatePasswordModal from '../components/UpdatePasswordModal';
import { useToast } from '../components/Toast';
import { adminsApi } from '../services/api';
import type { Admin, AdminFormData } from '../types/models';
import { isValidEmail, passwordValidationRules, validatePassword, getPasswordStrength } from '../utils/validation';

export const AdminsPage: React.FC = () => {
  const { t } = useTranslation();
  const { showSuccess, showError } = useToast();
  const [admins, setAdmins] = useState<Admin[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [confirmModal, setConfirmModal] = useState<{
    isOpen: boolean;
    admin: Admin | null;
  }>({ isOpen: false, admin: null });
  const [updatePasswordModal, setUpdatePasswordModal] = useState<{
    isOpen: boolean;
    admin: Admin | null;
  }>({ isOpen: false, admin: null });
  const [editingAdmin, setEditingAdmin] = useState<Admin | null>(null);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [passwordStrength, setPasswordStrength] = useState(0);
  const [passwordValidationErrors, setPasswordValidationErrors] = useState<string[]>([]);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
    setValue,
    control,
  } = useForm<AdminFormData>();

  const password = useWatch({
    control,
    name: 'password',
    defaultValue: '',
  });

  const retypePassword = useWatch({
    control,
    name: 'retype_password',
    defaultValue: '',
  });

  useEffect(() => {
    fetchAdmins();
  }, [pagination.current, pagination.pageSize]);

  // Password validation effect
  useEffect(() => {
    if (password) {
      setPasswordValidationErrors(validatePassword(password));
      setPasswordStrength(getPasswordStrength(password));
    } else {
      setPasswordValidationErrors([]);
      setPasswordStrength(0);
    }
  }, [password]);

  const fetchAdmins = async () => {
    try {
      setLoading(true);
      const response = await adminsApi.getAll({
        limit: pagination.pageSize,
        offset: (pagination.current - 1) * pagination.pageSize,
      });
      setAdmins(response.data || []);
      setPagination(prev => ({
        ...prev,
        total: response.total || 0,
      }));
    } catch (error: any) {
      console.error('Failed to fetch admins:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = () => {
    setEditingAdmin(null);
    reset();
    setPasswordValidationErrors([]);
    setPasswordStrength(0);
    setModalOpen(true);
  };

  const handleEdit = (admin: Admin) => {
    setEditingAdmin(admin);
    setValue('email', admin.email);
    setValue('is_active', admin.is_active);
    setModalOpen(true);
  };

  const handleDelete = (admin: Admin) => {
    setConfirmModal({ isOpen: true, admin });
  };

  const handleUpdatePassword = (admin: Admin) => {
    setUpdatePasswordModal({ isOpen: true, admin });
  };

  const handleConfirmDelete = async () => {
    if (!confirmModal.admin) return;
    
    try {
      await adminsApi.delete(confirmModal.admin.id);
      showSuccess(
        t('admins.deleted_successfully'),
        t('admins.admin_deleted_message')
      );
      await fetchAdmins();
      setConfirmModal({ isOpen: false, admin: null });
    } catch (error: any) {
      console.error('Failed to delete admin:', error);
      showError(
        t('common.error'),
        error?.message || t('admins.delete_failed')
      );
    }
  };

  const onSubmit = async (data: AdminFormData) => {
    try {
      // Remove retype_password before sending to API
      const { retype_password, ...apiData } = data;
      
      if (editingAdmin) {
        await adminsApi.update(editingAdmin.id, apiData);
        showSuccess(
          t('admins.updated_successfully'),
          t('admins.admin_updated_message')
        );
      } else {
        await adminsApi.create(apiData);
        showSuccess(
          t('admins.created_successfully'),
          t('admins.admin_created_message')
        );
      }
      await fetchAdmins();
      setModalOpen(false);
      reset();
    } catch (error: any) {
      console.error('Failed to save admin:', error);
      showError(
        t('common.error'),
        error?.message || t('admins.save_failed')
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

  const columns: Column<Admin>[] = [
    {
      key: 'email',
      title: t('admins.email'),
      width: '32',
    },
    {
      key: 'is_active',
      title: t('admins.status'),
      render: (value) => getStatusBadge(value),
    },
    {
      key: 'last_login',
      title: t('admins.last_login'),
      render: (value) => value ? new Date(value).toLocaleDateString() : t('common.never'),
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
        form="admin-form"
        disabled={isSubmitting}
        className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {isSubmitting ? (
          <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
        ) : null}
        {editingAdmin ? t('common.update') : t('common.create')}
      </button>
      <button
        type="button"
        onClick={() => setModalOpen(false)}
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
            {t('admins.title')}
          </h1>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-300">
            {t('admins.description')}
          </p>
        </div>
        <button
          onClick={handleCreate}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          <svg className="-ml-1 mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          {t('admins.add_admin')}
        </button>
      </div>

      {/* Data Table */}
      <DataTable
        data={admins}
        columns={columns}
        loading={loading}
        onEdit={handleEdit}
        onDelete={handleDelete}
        customActions={(admin) => (
          <>
            <button
              onClick={() => handleUpdatePassword(admin)}
              className="text-green-600 hover:text-green-900 dark:text-green-400 dark:hover:text-green-300 mr-3"
              title={t('passwords.update_password')}
            >
              🔧
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
        onClose={() => setModalOpen(false)}
        title={editingAdmin ? t('admins.edit_admin') : t('admins.add_admin')}
        size="lg"
        footer={modalFooter}
      >
        <form id="admin-form" onSubmit={handleSubmit(onSubmit)} className="space-y-4">

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              {t('admins.email')} *
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
              placeholder={t('admins.email_placeholder')}
              disabled={!!editingAdmin}
            />
            {errors.email && (
              <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.email.message}</p>
            )}
            {editingAdmin && (
              <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                {t('admins.email_readonly')}
              </p>
            )}
          </div>

          {!editingAdmin && (
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {t('admins.password')} *
                </label>
                <input
                  type="password"
                  {...register('password', { 
                    required: !editingAdmin ? t('validation.required') : false,
                    validate: (value) => {
                      if (!editingAdmin && value) {
                        const errors = validatePassword(value);
                        return errors.length === 0 || errors.join('; ');
                      }
                      return true;
                    },
                  })}
                  className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder={t('admins.password_placeholder')}
                />

                {/* Password Strength Indicator */}
                {!editingAdmin && password && (
                  <div className="mt-2">
                    <div className="flex items-center justify-between text-xs">
                      <span className="text-gray-500 dark:text-gray-400">Password Strength</span>
                      <span className={`font-medium ${
                        passwordStrength < 40 ? 'text-red-500' : 
                        passwordStrength < 80 ? 'text-yellow-500' : 'text-green-500'
                      }`}>
                        {passwordStrength < 40 ? 'Weak' : 
                         passwordStrength < 80 ? 'Medium' : 'Strong'}
                      </span>
                    </div>
                    <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-1.5 mt-1">
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

                {/* Password Validation Rules */}
                {!editingAdmin && passwordValidationErrors.length > 0 && (
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

              {/* Retype Password Field */}
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  {t('common.confirm_password')} *
                </label>
                <input
                  type="password"
                  {...register('retype_password', { 
                    required: !editingAdmin ? t('validation.required') : false,
                    validate: (value) => {
                      if (!editingAdmin && value !== password) {
                        return t('validation.passwords_must_match');
                      }
                      return true;
                    },
                  })}
                  className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder={t('common.confirm_password')}
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

          <div className="flex items-center">
            <input
              type="checkbox"
              {...register('is_active')}
              className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 dark:border-gray-600 rounded"
              defaultChecked={true}
            />
            <label className="ml-2 block text-sm text-gray-900 dark:text-gray-300">
              {t('admins.is_active')}
            </label>
          </div>
        </form>
      </Modal>

      {/* Delete Confirmation Modal */}
      <ConfirmModal
        isOpen={confirmModal.isOpen}
        onClose={() => setConfirmModal({ isOpen: false, admin: null })}
        onConfirm={handleConfirmDelete}
        title={t('admins.delete_admin')}
        message={t('admins.delete_confirm', { 
          name: confirmModal.admin ? confirmModal.admin.email : '' 
        })}
        confirmText={t('common.delete')}
        type="danger"
      />

      {/* Update Password Modal */}
      <UpdatePasswordModal
        isOpen={updatePasswordModal.isOpen}
        onClose={() => setUpdatePasswordModal({ isOpen: false, admin: null })}
        userType="admin"
        userId={updatePasswordModal.admin?.id || 0}
        userDisplayName={updatePasswordModal.admin?.email || ''}
      />
    </div>
  );
};

export default AdminsPage;