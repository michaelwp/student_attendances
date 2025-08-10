import React, { useState, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { useToast } from '../utils/toast-helpers';
import { useAuthStore } from '../stores/authStore';
import { useDarkMode } from '../hooks/useDarkMode';
import { teacherDashboardApi } from '../services/api';
import TeacherPasswordUpdate from '../components/TeacherPasswordUpdate';
import TeacherProfileView from '../components/TeacherProfileView';
import TeacherAbsentRequestList from '../components/TeacherAbsentRequestList';
import type { TeacherAbsentRequestListHandle } from '../components/TeacherAbsentRequestList';
import type { TeacherProfile } from '../types/models';

export const TeacherPage: React.FC = () => {
  const { t, i18n } = useTranslation();
  const { showError } = useToast();
  const { logout } = useAuthStore();
  const { isDark, toggleDarkMode } = useDarkMode();
  const [activeTab, setActiveTab] = useState('profile');
  const [profile, setProfile] = useState<TeacherProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const absentRequestListRef = useRef<TeacherAbsentRequestListHandle>(null);

  const handleLogout = async () => {
    try {
      await logout();
    } catch (error) {
      console.error('Logout error:', error);
    }
  };

  const toggleLanguage = () => {
    const newLang = i18n.language === 'en' ? 'id' : 'en';
    i18n.changeLanguage(newLang);
  };

  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    try {
      setLoading(true);
      const response = await teacherDashboardApi.getProfile();
      setProfile(response.data as TeacherProfile);
    } catch (error: unknown) {
      console.error('Failed to fetch profile:', error);
      showError(
        t('common.error'),
        (error as Error)?.message || t('teacher_dashboard.profile_fetch_failed')
      );
    } finally {
      setLoading(false);
    }
  };

  const handleRequestStatusChanged = () => {
    // Refresh the absent request list when a request status is changed
    if (absentRequestListRef.current) {
      absentRequestListRef.current.refreshRequests();
    }
  };

  const tabs = [
    { id: 'profile', name: t('teacher_page.profile_tab'), icon: '👤' },
    { id: 'password', name: t('teacher_page.password_tab'), icon: '🔒' },
    { id: 'absent_requests', name: t('teacher_page.absent_requests_tab'), icon: '📝' },
  ];

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!profile) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">
            {t('teacher_dashboard.profile_not_found')}
          </h2>
          <button
            onClick={fetchProfile}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            {t('common.retry')}
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      {/* Header */}
      <header className="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-semibold text-gray-900 dark:text-white">
                {t('app.title')} - {t('teacher_page.teacher_dashboard')}
              </h1>
            </div>
            
            <div className="flex items-center space-x-4">
              {/* Language Toggle */}
              <button
                onClick={toggleLanguage}
                className="p-2 text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white rounded-md hover:bg-gray-100 dark:hover:bg-gray-700"
                title={t('button.change_language')}
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129" />
                </svg>
              </button>
              
              {/* Dark Mode Toggle */}
              <button
                onClick={toggleDarkMode}
                className="p-2 text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white rounded-md hover:bg-gray-100 dark:hover:bg-gray-700"
                title={t('button.toggle_dark_mode')}
              >
                {isDark ? (
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                  </svg>
                ) : (
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                  </svg>
                )}
              </button>
              
              {/* Logout Button */}
              <button
                onClick={handleLogout}
                className="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
              >
                <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
                {t('button.logout')}
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div className="space-y-6">
          {/* Header */}
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
              {t('teacher_page.welcome')}, {profile.first_name} {profile.last_name}!
            </h1>
            <p className="mt-1 text-sm text-gray-600 dark:text-gray-300">
              {t('teacher_page.title')}
            </p>
          </div>

          {/* Stats Cards */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <StatsCard
              title={t('teacher_page.teacher_id')}
              value={profile.teacher_id}
              icon="👨‍🏫"
              color="blue"
              isText={true}
            />
            <StatsCard
              title={t('teacher_page.classes')}
              value={profile.classes?.length || 0}
              icon="🏫"
              color="green"
            />
            <StatsCard
              title={t('teacher_page.total_students')}
              value={profile.total_students || 0}
              icon="👨‍🎓"
              color="purple"
            />
            <StatsCard
              title={t('teacher_page.pending_requests')}
              value={profile.pending_requests || 0}
              icon="📝"
              color="orange"
            />
          </div>

          {/* Navigation Tabs */}
          <div className="border-b border-gray-200 dark:border-gray-700">
            <nav className="-mb-px flex space-x-8">
              {tabs.map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id)}
                  className={`py-2 px-1 border-b-2 font-medium text-sm ${
                    activeTab === tab.id
                      ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
                  }`}
                >
                  <span className="mr-2">{tab.icon}</span>
                  {tab.name}
                </button>
              ))}
            </nav>
          </div>

          {/* Tab Content */}
          <div>
            {activeTab === 'profile' && (
              <TeacherProfileView profile={profile} />
            )}

            {activeTab === 'password' && (
              <TeacherPasswordUpdate />
            )}

            {activeTab === 'absent_requests' && (
              <TeacherAbsentRequestList 
                ref={absentRequestListRef}
                onRequestStatusChanged={handleRequestStatusChanged}
              />
            )}
          </div>
        </div>
      </main>
    </div>
  );
};

interface StatsCardProps {
  title: string;
  value: number | string;
  icon: string;
  color: 'blue' | 'green' | 'purple' | 'orange';
  isText?: boolean;
}

const StatsCard: React.FC<StatsCardProps> = ({ title, value, icon, color, isText = false }) => {
  const colorClasses = {
    blue: 'bg-blue-50 dark:bg-blue-900/50 text-blue-600 dark:text-blue-400',
    green: 'bg-green-50 dark:bg-green-900/50 text-green-600 dark:text-green-400',
    purple: 'bg-purple-50 dark:bg-purple-900/50 text-purple-600 dark:text-purple-400',
    orange: 'bg-orange-50 dark:bg-orange-900/50 text-orange-600 dark:text-orange-400',
  };

  return (
    <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
      <div className="flex items-center">
        <div className={`flex-shrink-0 p-3 rounded-lg ${colorClasses[color]}`}>
          <span className="text-2xl">{icon}</span>
        </div>
        <div className="ml-4">
          <div className="text-sm font-medium text-gray-600 dark:text-gray-300">
            {title}
          </div>
          <div className="text-2xl font-bold text-gray-900 dark:text-white">
            {isText ? value : typeof value === 'number' ? value.toLocaleString() : value}
          </div>
        </div>
      </div>
    </div>
  );
};

export default TeacherPage;