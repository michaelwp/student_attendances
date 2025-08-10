import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useForm } from 'react-hook-form';
import { useToast } from '../utils/toast-helpers';
import { absentRequestApi } from '../services/api';
import type { AbsentRequestFormData } from '../types/models';

interface StudentAbsentRequestFormProps {
  onRequestCreated: () => void;
}

export const StudentAbsentRequestForm: React.FC<StudentAbsentRequestFormProps> = ({ onRequestCreated }) => {
  const { t } = useTranslation();
  const { showSuccess, showError } = useToast();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<AbsentRequestFormData>();

  const onSubmit = async (data: AbsentRequestFormData) => {
    try {
      setIsSubmitting(true);
      await absentRequestApi.create(data);
      showSuccess(
        t('student_page.request_created'),
        t('student_page.request_created_success')
      );
      reset();
      onRequestCreated();
    } catch (error: unknown) {
      console.error('Failed to create absent request:', error);
      showError(
        t('common.error'),
        (error as Error)?.message || t('student_page.request_create_failed')
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="bg-white dark:bg-gray-800 shadow rounded-lg">
      <div className="px-4 py-5 sm:p-6">
        <div className="flex items-center mb-4">
          <div className="flex-shrink-0">
            <div className="w-8 h-8 bg-blue-100 dark:bg-blue-900/50 rounded-lg flex items-center justify-center">
              <span className="text-blue-600 dark:text-blue-400">📝</span>
            </div>
          </div>
          <div className="ml-3">
            <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white">
              {t('student_page.create_absent_request')}
            </h3>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              {t('student_page.create_absent_request_desc')}
            </p>
          </div>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('student_page.request_date')} *
            </label>
            <input
              type="date"
              {...register('request_date', { 
                required: t('validation.required'),
                validate: (value) => {
                  const today = new Date().toISOString().split('T')[0];
                  return value >= today || t('student_page.date_future_only');
                }
              })}
              className="block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
            {errors.request_date && (
              <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.request_date.message}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('student_page.reason')} *
            </label>
            <textarea
              rows={4}
              {...register('reason', { 
                required: t('validation.required'),
                minLength: { value: 10, message: t('validation.min_length', { length: 10 }) },
                maxLength: { value: 500, message: t('validation.max_length', { length: 500 }) }
              })}
              className="block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder={t('student_page.reason_placeholder')}
            />
            {errors.reason && (
              <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.reason.message}</p>
            )}
          </div>

          <div className="bg-yellow-50 dark:bg-yellow-900/50 border border-yellow-200 dark:border-yellow-800 rounded-md p-3">
            <div className="flex">
              <div className="flex-shrink-0">
                <svg className="h-5 w-5 text-yellow-400" viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                </svg>
              </div>
              <div className="ml-3">
                <h3 className="text-sm font-medium text-yellow-800 dark:text-yellow-200">
                  {t('student_page.request_notice')}
                </h3>
                <div className="mt-2 text-sm text-yellow-700 dark:text-yellow-300">
                  <p>{t('student_page.request_notice_desc')}</p>
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
              {isSubmitting ? t('common.submitting') : t('student_page.submit_request')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default StudentAbsentRequestForm;