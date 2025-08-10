import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { Eye, EyeOff, User, Lock, AlertCircle, CheckCircle } from 'lucide-react';
import { useAuthStore } from '../stores/authStore';
import type {LoginRequest, UserType} from '../types/auth';
import { ApiError } from '../services/api';

interface LoginFormData {
  user_type: UserType;
  user_id: string;
  password: string;
}

const LoginForm: React.FC = () => {
  const { t } = useTranslation();
  const { login, isLoading, error, clearError } = useAuthStore();
  const [showPassword, setShowPassword] = useState(false);
  const [loginSuccess, setLoginSuccess] = useState(false);

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
    setError,
  } = useForm<LoginFormData>({
    defaultValues: {
      user_type: 'student',
      user_id: '',
      password: '',
    },
  });

  const userType = watch('user_type');

  const userTypeOptions: { value: UserType; label: string }[] = [
    { value: 'student', label: t('login.userTypes.student') },
    { value: 'teacher', label: t('login.userTypes.teacher') },
    { value: 'admin', label: t('login.userTypes.admin') },
  ];

  const getUserIdPlaceholder = (type: UserType): string => {
    switch (type) {
      case 'admin':
        return t('login.placeholders.admin');
      case 'teacher':
        return t('login.placeholders.teacher');
      case 'student':
        return t('login.placeholders.student');
      default:
        return '';
    }
  };

  const onSubmit = async (data: LoginFormData) => {
    try {
      clearError();
      await login(data as LoginRequest);
      setLoginSuccess(true);
      
      // Simulate redirect delay
      setTimeout(() => {
        // Here you would typically redirect to dashboard
        console.log('Redirecting to dashboard...');
      }, 1500);
    } catch (error) {
      setLoginSuccess(false);
      
      if (error instanceof ApiError) {
        // Map API error translation keys to form errors
        if (error.translationKey === 'error.invalid_credentials') {
          setError('user_id', { message: t('errors.invalid_credentials') });
          setError('password', { message: t('errors.invalid_credentials') });
        }
      }
    }
  };

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  return (
    <div className="w-full max-w-md mx-auto">
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        {/* User Type Selection */}
        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('login.userType')} *
          </label>
          <div className="grid grid-cols-3 gap-2">
            {userTypeOptions.map((option) => (
              <label
                key={option.value}
                className={`relative flex items-center justify-center p-3 border rounded-lg cursor-pointer transition-all duration-200 ${
                  userType === option.value
                    ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20 text-primary-700 dark:text-primary-300'
                    : 'border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:border-primary-300 dark:hover:border-primary-600'
                }`}
              >
                <input
                  type="radio"
                  value={option.value}
                  {...register('user_type', { required: t('login.validation.userTypeRequired') })}
                  className="sr-only"
                />
                <span className="text-sm font-medium">{option.label}</span>
                {userType === option.value && (
                  <CheckCircle className="absolute top-1 right-1 w-4 h-4 text-primary-500" />
                )}
              </label>
            ))}
          </div>
          {errors.user_type && (
            <p className="text-red-500 text-sm flex items-center space-x-1">
              <AlertCircle className="w-4 h-4" />
              <span>{errors.user_type.message}</span>
            </p>
          )}
        </div>

        {/* User ID Input */}
        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('login.userId')} *
          </label>
          <div className="relative">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <User className="h-5 w-5 text-gray-400" />
            </div>
            <input
              type="text"
              {...register('user_id', {
                required: t('login.validation.userIdRequired'),
              })}
              className={`block w-full pl-10 pr-3 py-2 border rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-800 dark:text-white dark:placeholder-gray-500 ${
                errors.user_id
                  ? 'border-red-300 dark:border-red-600'
                  : 'border-gray-300 dark:border-gray-600'
              }`}
              placeholder={getUserIdPlaceholder(userType)}
            />
          </div>
          {errors.user_id && (
            <p className="text-red-500 text-sm flex items-center space-x-1">
              <AlertCircle className="w-4 h-4" />
              <span>{errors.user_id.message}</span>
            </p>
          )}
        </div>

        {/* Password Input */}
        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
            {t('login.password')} *
          </label>
          <div className="relative">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Lock className="h-5 w-5 text-gray-400" />
            </div>
            <input
              type={showPassword ? 'text' : 'password'}
              {...register('password', {
                required: t('login.validation.passwordRequired'),
                minLength: {
                  value: 6,
                  message: t('login.validation.passwordMinLength', {length:6}),
                },
              })}
              className={`block w-full pl-10 pr-10 py-2 border rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-800 dark:text-white dark:placeholder-gray-500 ${
                errors.password
                  ? 'border-red-300 dark:border-red-600'
                  : 'border-gray-300 dark:border-gray-600'
              }`}
              placeholder={t('login.placeholders.password')}
            />
            <button
              type="button"
              onClick={togglePasswordVisibility}
              className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            >
              {showPassword ? (
                <EyeOff className="h-5 w-5" />
              ) : (
                <Eye className="h-5 w-5" />
              )}
            </button>
          </div>
          {errors.password && (
            <p className="text-red-500 text-sm flex items-center space-x-1">
              <AlertCircle className="w-4 h-4" />
              <span>{errors.password.message}</span>
            </p>
          )}
        </div>

        {/* Error Message */}
        {error && (
          <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
            <div className="flex items-center space-x-2">
              <AlertCircle className="h-5 w-5 text-red-500" />
              <p className="text-red-700 dark:text-red-400 text-sm">{error}</p>
            </div>
          </div>
        )}

        {/* Success Message */}
        {loginSuccess && (
          <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4">
            <div className="flex items-center space-x-2">
              <CheckCircle className="h-5 w-5 text-green-500" />
              <p className="text-green-700 dark:text-green-400 text-sm">
                {t('success.login_successful')}
              </p>
            </div>
          </div>
        )}

        {/* Submit Button */}
        <button
          type="submit"
          disabled={isLoading}
          className="w-full flex items-center justify-center px-4 py-2 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed dark:focus:ring-offset-gray-800 transition-all duration-200"
        >
          {isLoading ? (
            <div className="flex items-center space-x-2">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
              <span>{t('login.signingIn')}</span>
            </div>
          ) : (
            t('login.submit')
          )}
        </button>
      </form>
    </div>
  );
};

export default LoginForm;