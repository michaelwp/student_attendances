import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { useToast } from '../utils/toast-helpers';
import { teacherDashboardApi } from '../services/api';
import type { PasswordUpdateData } from '../types/models';

interface PasswordFormData {
  old_password: string;
  new_password: string;
  confirm_password: string;
}

export const TeacherPasswordUpdate: React.FC = () => {
  const { t } = useTranslation();
  const { showSuccess, showError } = useToast();
  const [loading, setLoading] = useState(false);
  
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    watch,
  } = useForm<PasswordFormData>();

  const newPassword = watch('new_password');

  const onSubmit = async (data: PasswordFormData) => {
    if (data.new_password !== data.confirm_password) {
      showError(t('common.error'), t('teacher_page.password_mismatch'));
      return;
    }

    try {
      setLoading(true);
      const passwordData: PasswordUpdateData = {
        old_password: data.old_password,
        new_password: data.new_password,
      };
      
      await teacherDashboardApi.updatePassword(passwordData);
      
      showSuccess(
        t('common.success'),
        t('teacher_page.password_updated_successfully')
      );
      
      reset();
    } catch (error: unknown) {
      console.error('Password update error:', error);
      showError(
        t('common.error'),
        (error as Error)?.message || t('teacher_page.password_update_failed')
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white dark:bg-gray-800 shadow rounded-lg">
      <div className="px-4 py-5 sm:p-6">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <div className="w-8 h-8 bg-orange-100 dark:bg-orange-900/50 rounded-lg flex items-center justify-center">
                <span className="text-orange-600 dark:text-orange-400">🔒</span>
              </div>
            </div>
            <div className="ml-3">
              <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white">
                {t('teacher_page.change_password')}
              </h3>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                {t('teacher_page.change_password_description')}
              </p>
            </div>
          </div>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          {/* Current Password */}
          <div>
            <label htmlFor="old_password" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              {t('teacher_page.current_password')}
            </label>
            <div className="mt-1">
              <input
                {...register('old_password', {
                  required: t('teacher_page.current_password_required'),
                })}
                type="password"
                id="old_password"
                className="appearance-none block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                placeholder={t('teacher_page.enter_current_password')}
              />
              {errors.old_password && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">
                  {errors.old_password.message}
                </p>
              )}
            </div>
          </div>

          {/* New Password */}
          <div>
            <label htmlFor="new_password" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              {t('teacher_page.new_password')}
            </label>
            <div className="mt-1">
              <input
                {...register('new_password', {
                  required: t('teacher_page.new_password_required'),
                  minLength: {
                    value: 8,
                    message: t('teacher_page.password_min_length', { length: 8 })
                  }
                })}
                type="password"
                id="new_password"
                className="appearance-none block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                placeholder={t('teacher_page.enter_new_password')}
              />
              {errors.new_password && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">
                  {errors.new_password.message}
                </p>
              )}
            </div>
          </div>

          {/* Confirm Password */}
          <div>
            <label htmlFor="confirm_password" className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              {t('teacher_page.confirm_new_password')}
            </label>
            <div className="mt-1">
              <input
                {...register('confirm_password', {
                  required: t('teacher_page.confirm_password_required'),
                  validate: (value) =>
                    value === newPassword || t('teacher_page.password_mismatch')
                })}
                type="password"
                id="confirm_password"
                className="appearance-none block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                placeholder={t('teacher_page.confirm_new_password')}
              />
              {errors.confirm_password && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">
                  {errors.confirm_password.message}
                </p>
              )}
            </div>
          </div>

          {/* Password Requirements */}
          <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-700 rounded-md p-4">
            <div className="flex">
              <div className="flex-shrink-0">
                <svg className="h-5 w-5 text-blue-400" viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
                </svg>
              </div>
              <div className="ml-3">
                <h3 className="text-sm font-medium text-blue-800 dark:text-blue-300">
                  {t('teacher_page.password_requirements')}
                </h3>
                <div className="mt-2 text-sm text-blue-700 dark:text-blue-300">
                  <ul className="list-disc space-y-1 pl-5">
                    <li>{t('teacher_page.password_requirement_length')}</li>
                    <li>{t('teacher_page.password_requirement_security')}</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>

          {/* Submit Button */}
          <div className="flex justify-end">
            <button
              type="submit"
              disabled={loading}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? (
                <>
                  <svg className="animate-spin -ml-1 mr-3 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  {t('teacher_page.updating_password')}
                </>
              ) : (
                t('teacher_page.update_password')
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default TeacherPasswordUpdate;