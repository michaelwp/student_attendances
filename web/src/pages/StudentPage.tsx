import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useToast } from '../utils/toast-helpers';
import { studentDashboardApi } from '../services/api';
import StudentPasswordUpdate from '../components/StudentPasswordUpdate';
import StudentProfileView from '../components/StudentProfileView';
import StudentAbsentRequestForm from '../components/StudentAbsentRequestForm';
import StudentAbsentRequestList from '../components/StudentAbsentRequestList';
import type { StudentProfile } from '../types/models';

export const StudentPage: React.FC = () => {
  const { t } = useTranslation();
  const { showError } = useToast();
  const [activeTab, setActiveTab] = useState('profile');
  const [profile, setProfile] = useState<StudentProfile | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    try {
      setLoading(true);
      const response = await studentDashboardApi.getProfile();
      setProfile(response.data as StudentProfile);
    } catch (error: unknown) {
      console.error('Failed to fetch profile:', error);
      showError(
        t('common.error'),
        (error as Error)?.message || t('student_dashboard.profile_fetch_failed')
      );
    } finally {
      setLoading(false);
    }
  };

  const tabs = [
    { id: 'profile', name: t('student_page.profile_tab'), icon: '👤' },
    { id: 'password', name: t('student_page.password_tab'), icon: '🔒' },
    { id: 'absent_requests', name: t('student_page.absent_requests_tab'), icon: '📝' },
  ];

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!profile) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">
            {t('student_dashboard.profile_not_found')}
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
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
          {t('student_page.welcome')}, {profile.first_name} {profile.last_name}!
        </h1>
        <p className="mt-1 text-sm text-gray-600 dark:text-gray-300">
          {t('student_page.title')}
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatsCard
          title={t('student_page.student_id')}
          value={profile.student_id}
          icon="📚"
          color="blue"
          isText={true}
        />
        <StatsCard
          title={t('student_page.class')}
          value={profile.class_name || 'N/A'}
          icon="🏫"
          color="green"
          isText={true}
        />
        <StatsCard
          title={t('student_page.attendance_rate')}
          value={profile.attendance_stats?.attendance_rate
              ? `${profile.attendance_stats.attendance_rate.toFixed(1)}%`
              : 'N/A'}
          icon="📊"
          color="purple"
          isText={true}
        />
        <StatsCard
          title={t('student_page.total_attendance')}
          value={profile.attendance_stats?.total_days || 0}
          icon="✅"
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
          <StudentProfileView profile={profile} />
        )}

        {activeTab === 'password' && (
          <StudentPasswordUpdate />
        )}

        {activeTab === 'absent_requests' && (
          <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
            <div className="xl:col-span-1">
              <StudentAbsentRequestForm onRequestCreated={() => {}} />
            </div>
            <div className="xl:col-span-2">
              <StudentAbsentRequestList />
            </div>
          </div>
        )}
      </div>
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

export default StudentPage;