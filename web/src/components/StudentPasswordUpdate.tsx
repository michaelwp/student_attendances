import React from 'react';
import { useTranslation } from 'react-i18next';
import { useForm } from 'react-hook-form';
import { useToast } from '../utils/toast-helpers';
import { studentDashboardApi } from '../services/api';
import type { PasswordUpdateData } from '../types/models';

interface PasswordFormData {
  old_password: string;
  new_password: string;
  confirm_password: string;
}

export const StudentPasswordUpdate: React.FC = () => {
  const { t } = useTranslation();
  const { showSuccess, showError } = useToast();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
    watch,
  } = useForm<PasswordFormData>();

  const newPassword = watch('new_password');

  const onSubmit = async (data: PasswordFormData) => {
    if (data.new_password !== data.confirm_password) {
      showError(t('common.error'), t('validation.passwords_must_match'));
      return;
    }

    try {
      const passwordData: PasswordUpdateData = {
        old_password: data.old_password,
        new_password: data.new_password,
      };
      
      await studentDashboardApi.updatePassword(passwordData);
      showSuccess(
        t('student_page.password_updated'),
        t('student_page.password_updated_success')
      );
      reset();
    } catch (error: unknown) {
      console.error('Failed to update password:', error);
      showError(
        t('common.error'),
        (error as Error)?.message || t('student_page.password_update_failed')
      );
    }
  };

  return (
    <div className="bg-white dark:bg-gray-800 shadow rounded-lg">
      <div className="px-4 py-5 sm:p-6">
        <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white mb-4">
          {t('student_page.change_password')}
        </h3>
        <p className="text-sm text-gray-500 dark:text-gray-400 mb-6">
          {t('student_page.change_password_desc')}
        </p>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6 max-w-lg">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('student_page.current_password')} *
            </label>
            <input
              type="password"
              {...register('old_password', { 
                required: t('validation.required'),
                minLength: { value: 6, message: t('validation.min_length', { length: 6 }) }
              })}
              className="block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder={t('student_page.current_password_placeholder')}
            />
            {errors.old_password && (
              <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.old_password.message}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('student_page.new_password')} *
            </label>
            <input
              type="password"
              {...register('new_password', { 
                required: t('validation.required'),
                minLength: { value: 8, message: t('validation.min_length', { length: 8 }) },
                pattern: {
                  value: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/,
                  message: t('validation.password_complexity')
                }
              })}
              className="block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder={t('student_page.new_password_placeholder')}
            />
            {errors.new_password && (
              <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.new_password.message}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('student_page.confirm_password')} *
            </label>
            <input
              type="password"
              {...register('confirm_password', { 
                required: t('validation.required'),
                validate: (value) => value === newPassword || t('validation.passwords_must_match')
              })}
              className="block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder={t('student_page.confirm_password_placeholder')}
            />
            {errors.confirm_password && (
              <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.confirm_password.message}</p>
            )}
          </div>

          <div className="bg-blue-50 dark:bg-blue-900/50 border border-blue-200 dark:border-blue-800 rounded-md p-4">
            <div className="flex">
              <div className="flex-shrink-0">
                <svg className="h-5 w-5 text-blue-400" viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
                </svg>
              </div>
              <div className="ml-3">
                <h3 className="text-sm font-medium text-blue-800 dark:text-blue-200">
                  {t('student_page.password_requirements')}
                </h3>
                <div className="mt-2 text-sm text-blue-700 dark:text-blue-300">
                  <ul className="list-disc pl-5 space-y-1">
                    <li>{t('validation.password_min_length', { length: 8 })}</li>
                    <li>{t('validation.password_uppercase')}</li>
                    <li>{t('validation.password_lowercase')}</li>
                    <li>{t('validation.password_number')}</li>
                      <li>{t('validation.password_special')}</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>

          <div className="flex justify-end">
            <button
              type="submit"
              disabled={isSubmitting}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isSubmitting && (
                <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
              )}
              {isSubmitting ? t('common.updating') : t('student_page.update_password')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default StudentPasswordUpdate;