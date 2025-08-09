import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useForm } from 'react-hook-form';
import Modal from './Modal';
import { useToast } from '../utils/toast-helpers';
import { teachersApi, studentsApi, adminsApi } from '../services/api';
import { passwordValidationRules, validatePassword, getPasswordStrength } from '../utils/validation';

interface UpdatePasswordModalProps {
  isOpen: boolean;
  onClose: () => void;
  userType: 'teacher' | 'student' | 'admin';
  userId: string | number; // teacher_id/student_id/admin_id
  userDisplayName: string;
}

interface PasswordUpdateFormData {
  old_password: string;
  new_password: string;
  confirm_password: string;
}

const UpdatePasswordModal: React.FC<UpdatePasswordModalProps> = ({
  isOpen,
  onClose,
  userType,
  userId,
  userDisplayName,
}) => {
  const { t } = useTranslation();
  const { showSuccess, showError } = useToast();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState(0);
  const [passwordValidationErrors, setPasswordValidationErrors] = useState<string[]>([]);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    watch,
  } = useForm<PasswordUpdateFormData>();

  const newPassword = watch('new_password');
  const confirmPassword = watch('confirm_password');

  React.useEffect(() => {
    if (newPassword) {
      setPasswordValidationErrors(validatePassword(newPassword));
      setPasswordStrength(getPasswordStrength(newPassword));
    } else {
      setPasswordValidationErrors([]);
      setPasswordStrength(0);
    }
  }, [newPassword]);

  const onSubmit = async (data: PasswordUpdateFormData) => {
    setIsSubmitting(true);
    
    try {
      const updateData = {
        old_password: data.old_password,
        new_password: data.new_password,
      };

      switch (userType) {
        case 'teacher':
          await teachersApi.updatePassword(userId as string, updateData);
          break;
        case 'student':
          await studentsApi.updatePassword(userId as string, updateData);
          break;
        case 'admin':
          await adminsApi.updatePassword(userId as number, updateData);
          break;
        default:
          throw new Error('Invalid user type');
      }

      showSuccess(
        t('passwords.updated_successfully'),
        t('passwords.password_updated_message', { name: userDisplayName })
      );

      reset();
      onClose();
    } catch (error: unknown) {
      console.error('Failed to update password:', error);
      showError(
        t('common.error'),
        (error as Error)?.message || t('passwords.update_failed')
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    reset();
    setPasswordValidationErrors([]);
    setPasswordStrength(0);
    onClose();
  };

  const modalFooter = (
    <>
      <button
        type="submit"
        form="update-password-form"
        disabled={isSubmitting || passwordValidationErrors.length > 0}
        className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {isSubmitting ? (
          <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
        ) : null}
        {t('passwords.update_password')}
      </button>
      <button
        type="button"
        onClick={handleClose}
        disabled={isSubmitting}
        className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 dark:border-gray-600 shadow-sm px-4 py-2 bg-white dark:bg-gray-800 text-base font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {t('common.cancel')}
      </button>
    </>
  );

  return (
    <Modal
      isOpen={isOpen}
      onClose={handleClose}
      title={t('passwords.update_password_for', { name: userDisplayName })}
      size="md"
      footer={modalFooter}
    >
      <form id="update-password-form" onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {/* Current Password */}
        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('passwords.current_password')} *
          </label>
          <input
            type="password"
            {...register('old_password', { 
              required: t('validation.required'),
            })}
            className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            placeholder={t('passwords.current_password_placeholder')}
          />
          {errors.old_password && (
            <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.old_password.message}</p>
          )}
        </div>

        {/* New Password */}
        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('passwords.new_password')} *
          </label>
          <input
            type="password"
            {...register('new_password', { 
              required: t('validation.required'),
              validate: (value) => {
                if (value) {
                  const errors = validatePassword(value);
                  return errors.length === 0 || errors.join('; ');
                }
                return true;
              },
            })}
            className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            placeholder={t('passwords.new_password_placeholder')}
          />

          {/* Password Strength Indicator */}
          {newPassword && (
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
          {passwordValidationErrors.length > 0 && (
            <div className="mt-2 space-y-1">
              {passwordValidationRules.map((rule, index) => {
                const isValid = newPassword ? rule.test(newPassword) : false;
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

          {errors.new_password && (
            <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.new_password.message}</p>
          )}
        </div>

        {/* Confirm New Password */}
        <div>
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('passwords.confirm_new_password')} *
          </label>
          <input
            type="password"
            {...register('confirm_password', { 
              required: t('validation.required'),
              validate: (value) => {
                if (value !== newPassword) {
                  return t('validation.passwords_must_match');
                }
                return true;
              },
            })}
            className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            placeholder={t('passwords.confirm_new_password_placeholder')}
          />
          {confirmPassword && confirmPassword !== newPassword && (
            <p className="mt-1 text-sm text-red-600 dark:text-red-400">
              {t('validation.passwords_must_match')}
            </p>
          )}
          {errors.confirm_password && (
            <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.confirm_password.message}</p>
          )}
        </div>

        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-md p-3">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-blue-400" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <p className="text-sm text-blue-700 dark:text-blue-300">
                {t('passwords.security_notice')}
              </p>
            </div>
          </div>
        </div>
      </form>
    </Modal>
  );
};

export default UpdatePasswordModal;