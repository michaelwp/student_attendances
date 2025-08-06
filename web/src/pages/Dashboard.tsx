import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../stores/authStore';
import { dashboardApi } from '../services/api';
import type { DashboardStats } from '../types/models';

export const Dashboard: React.FC = () => {
  const { t } = useTranslation();
  const { user } = useAuthStore();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchStats();
  }, []);

  const fetchStats = async () => {
    try {
      setLoading(true);
      const response = await dashboardApi.getStats();
      setStats(response.data || null);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch stats');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 dark:bg-red-900/50 border border-red-200 dark:border-red-700 rounded-md p-4">
        <div className="flex">
          <div className="flex-shrink-0">
            <svg className="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
          </div>
          <div className="ml-3">
            <h3 className="text-sm font-medium text-red-800 dark:text-red-200">
              {t('error.failed_to_load')}
            </h3>
            <p className="text-sm text-red-700 dark:text-red-300 mt-1">{error}</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
          {t('dashboard.welcome')}, {user?.type === 'admin' ? t('dashboard.admin') : ''}
        </h1>
        <p className="mt-1 text-sm text-gray-600 dark:text-gray-300">
          {t('dashboard.overview')}
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <DetailedStatsCard
          title={t('dashboard.stats.total_teachers')}
          total={stats?.total_teachers || 0}
          active={stats?.active_teachers || 0}
          inactive={stats?.inactive_teachers || 0}
          icon="👨‍🏫"
          color="blue"
        />
        <DetailedStatsCard
          title={t('dashboard.stats.total_students')}
          total={stats?.total_students || 0}
          active={stats?.active_students || 0}
          inactive={stats?.inactive_students || 0}
          icon="👨‍🎓"
          color="green"
        />
        <StatsCard
          title={t('dashboard.stats.total_classes')}
          value={stats?.total_classes || 0}
          icon="🏫"
          color="purple"
        />
        <DetailedStatsCard
          title={t('dashboard.stats.total_admins')}
          total={stats?.total_admins || 0}
          active={stats?.active_admins || 0}
          inactive={stats?.inactive_admins || 0}
          icon="👤"
          color="orange"
        />
      </div>

      {/* Today's Attendance */}
      {(stats?.total_attendance_today || 0) > 0 && (
        <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
          <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
            {t('dashboard.today_attendance')}
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="text-center">
              <div className="text-2xl font-bold text-green-600 dark:text-green-400">
                {stats?.present_today || 0}
              </div>
              <div className="text-sm text-gray-600 dark:text-gray-300">
                {t('dashboard.present')}
              </div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-red-600 dark:text-red-400">
                {stats?.absent_today || 0}
              </div>
              <div className="text-sm text-gray-600 dark:text-gray-300">
                {t('dashboard.absent')}
              </div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-yellow-600 dark:text-yellow-400">
                {stats?.late_today || 0}
              </div>
              <div className="text-sm text-gray-600 dark:text-gray-300">
                {t('dashboard.late')}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Quick Actions */}
      <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
        <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
          {t('dashboard.quick_actions')}
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <QuickActionButton
            title={t('dashboard.manage_teachers')}
            description={t('dashboard.manage_teachers_desc')}
            icon="👨‍🏫"
            href="#/teachers"
          />
          <QuickActionButton
            title={t('dashboard.manage_students')}
            description={t('dashboard.manage_students_desc')}
            icon="👨‍🎓"
            href="#/students"
          />
          <QuickActionButton
            title={t('dashboard.manage_classes')}
            description={t('dashboard.manage_classes_desc')}
            icon="🏫"
            href="#/classes"
          />
          <QuickActionButton
            title={t('dashboard.attendance')}
            description={t('dashboard.attendance_desc')}
            icon="📊"
            href="#/attendance"
          />
        </div>
      </div>
    </div>
  );
};

interface StatsCardProps {
  title: string;
  value: number;
  icon: string;
  color: 'blue' | 'green' | 'purple' | 'orange';
}

const StatsCard: React.FC<StatsCardProps> = ({ title, value, icon, color }) => {
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
            {value.toLocaleString()}
          </div>
        </div>
      </div>
    </div>
  );
};

interface DetailedStatsCardProps {
  title: string;
  total: number;
  active: number;
  inactive: number;
  icon: string;
  color: 'blue' | 'green' | 'purple' | 'orange';
}

const DetailedStatsCard: React.FC<DetailedStatsCardProps> = ({ title, total, active, inactive, icon, color }) => {
  const { t } = useTranslation();
  
  const colorClasses = {
    blue: 'bg-blue-50 dark:bg-blue-900/50 text-blue-600 dark:text-blue-400',
    green: 'bg-green-50 dark:bg-green-900/50 text-green-600 dark:text-green-400',
    purple: 'bg-purple-50 dark:bg-purple-900/50 text-purple-600 dark:text-purple-400',
    orange: 'bg-orange-50 dark:bg-orange-900/50 text-orange-600 dark:text-orange-400',
  };

  return (
    <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
      <div className="flex items-start">
        <div className={`flex-shrink-0 p-3 rounded-lg ${colorClasses[color]}`}>
          <span className="text-2xl">{icon}</span>
        </div>
        <div className="ml-4 flex-1">
          <div className="text-sm font-medium text-gray-600 dark:text-gray-300">
            {title}
          </div>
          <div className="text-2xl font-bold text-gray-900 dark:text-white">
            {total.toLocaleString()}
          </div>
          <div className="flex items-center space-x-4 mt-2 text-xs">
            <div className="flex items-center">
              <div className="w-2 h-2 bg-green-500 rounded-full mr-1"></div>
              <span className="text-gray-600 dark:text-gray-400">
                {t('dashboard.active')}: {active}
              </span>
            </div>
            <div className="flex items-center">
              <div className="w-2 h-2 bg-red-500 rounded-full mr-1"></div>
              <span className="text-gray-600 dark:text-gray-400">
                {t('dashboard.inactive')}: {inactive}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

interface QuickActionButtonProps {
  title: string;
  description: string;
  icon: string;
  href: string;
}

const QuickActionButton: React.FC<QuickActionButtonProps> = ({
  title,
  description,
  icon,
  href,
}) => {
  const handleClick = () => {
    window.location.hash = href;
  };

  return (
    <button
      onClick={handleClick}
      className="text-left p-4 border border-gray-200 dark:border-gray-600 rounded-lg hover:border-blue-500 hover:shadow-md transition-all duration-200 group"
    >
      <div className="flex items-start space-x-3">
        <div className="text-2xl">{icon}</div>
        <div className="flex-1 min-w-0">
          <h3 className="text-sm font-medium text-gray-900 dark:text-white group-hover:text-blue-600 dark:group-hover:text-blue-400">
            {title}
          </h3>
          <p className="text-xs text-gray-600 dark:text-gray-300 mt-1">
            {description}
          </p>
        </div>
      </div>
    </button>
  );
};

export default Dashboard;