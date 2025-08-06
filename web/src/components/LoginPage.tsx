import React from 'react';
import { useTranslation } from 'react-i18next';
import { GraduationCap } from 'lucide-react';
import LoginForm from './LoginForm';
import LanguageSwitcher from './LanguageSwitcher';
import ThemeToggle from './ThemeToggle';

const LoginPage: React.FC = () => {
  const { t } = useTranslation();

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800 transition-colors duration-300">
      {/* Header with controls */}
      <header className="absolute top-0 right-0 p-6">
        <div className="flex items-center space-x-4">
          <LanguageSwitcher />
          <ThemeToggle />
        </div>
      </header>

      {/* Main content */}
      <div className="flex items-center justify-center min-h-screen p-4">
        <div className="w-full max-w-md">
          {/* Logo and title */}
          <div className="text-center mb-8">
            <div className="flex items-center justify-center mb-6">
              <div className="bg-primary-600 dark:bg-primary-500 rounded-full p-4 shadow-lg">
                <GraduationCap className="w-12 h-12 text-white" />
              </div>
            </div>
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2 animate-fade-in">
              {t('app.title')}
            </h1>
            <p className="text-gray-600 dark:text-gray-400 animate-fade-in">
              {t('app.subtitle')}
            </p>
          </div>

          {/* Login card */}
          <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-700 p-8 animate-slide-up">
            <div className="mb-6">
              <h2 className="text-2xl font-semibold text-gray-900 dark:text-white text-center">
                {t('login.title')}
              </h2>
            </div>

            <LoginForm />
          </div>

          {/* Footer */}
          <div className="text-center mt-8 text-sm text-gray-500 dark:text-gray-400">
            <p>© 2024 Student Attendance System. All rights reserved.</p>
          </div>
        </div>
      </div>

      {/* Background pattern */}
      <div className="absolute inset-0 -z-10 overflow-hidden">
        <div className="absolute top-1/4 left-1/4 w-64 h-64 bg-primary-200 dark:bg-primary-900 rounded-full opacity-20 animate-pulse"></div>
        <div className="absolute bottom-1/4 right-1/4 w-48 h-48 bg-indigo-200 dark:bg-indigo-900 rounded-full opacity-20 animate-pulse delay-1000"></div>
        <div className="absolute top-3/4 left-1/3 w-32 h-32 bg-blue-200 dark:bg-blue-900 rounded-full opacity-20 animate-pulse delay-2000"></div>
      </div>
    </div>
  );
};

export default LoginPage;