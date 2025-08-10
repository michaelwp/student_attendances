import React from 'react';
import { useTranslation } from 'react-i18next';
import type { TeacherProfile } from '../types/models';

interface TeacherProfileViewProps {
  profile: TeacherProfile;
}

export const TeacherProfileView: React.FC<TeacherProfileViewProps> = ({ profile }) => {
  const { t } = useTranslation();

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <div className="bg-white dark:bg-gray-800 shadow rounded-lg">
      <div className="px-4 py-5 sm:p-6">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <div className="w-8 h-8 bg-blue-100 dark:bg-blue-900/50 rounded-lg flex items-center justify-center">
                <span className="text-blue-600 dark:text-blue-400">👨‍🏫</span>
              </div>
            </div>
            <div className="ml-3">
              <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white">
                {t('teacher_page.profile_information')}
              </h3>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                {t('teacher_page.profile_description')}
              </p>
            </div>
          </div>
        </div>

        <div className="border-t border-gray-200 dark:border-gray-700 pt-4">
          <dl className="grid grid-cols-1 gap-x-4 gap-y-6 sm:grid-cols-2">
            <div>
              <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                {t('teacher_page.teacher_id')}
              </dt>
              <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                {profile.teacher_id}
              </dd>
            </div>

            <div>
              <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                {t('teacher_page.full_name')}
              </dt>
              <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                {profile.first_name} {profile.last_name}
              </dd>
            </div>

            <div>
              <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                {t('teacher_page.email')}
              </dt>
              <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                {profile.email}
              </dd>
            </div>

            <div>
              <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                {t('teacher_page.phone')}
              </dt>
              <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                {profile.phone}
              </dd>
            </div>

            <div>
              <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                {t('teacher_page.status')}
              </dt>
              <dd className="mt-1 text-sm">
                <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                  profile.is_active
                    ? 'bg-green-100 text-green-800 dark:bg-green-900/50 dark:text-green-300'
                    : 'bg-red-100 text-red-800 dark:bg-red-900/50 dark:text-red-300'
                }`}>
                  <span className="mr-1">
                    {profile.is_active ? '✅' : '❌'}
                  </span>
                  {profile.is_active ? t('common.active') : t('common.inactive')}
                </span>
              </dd>
            </div>

            <div>
              <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">
                {t('teacher_page.member_since')}
              </dt>
              <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                {formatDate(profile.created_at)}
              </dd>
            </div>
          </dl>
        </div>

        {/* Classes Section */}
        {profile.classes && profile.classes.length > 0 && (
          <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
            <h4 className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-3">
              {t('teacher_page.assigned_classes')}
            </h4>
            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
              {profile.classes.map((classItem) => (
                <div key={classItem.id} className="bg-gray-50 dark:bg-gray-700 rounded-lg p-3">
                  <div className="flex items-center justify-between">
                    <div>
                      <h5 className="text-sm font-medium text-gray-900 dark:text-white">
                        {classItem.name}
                      </h5>
                      <p className="text-xs text-gray-500 dark:text-gray-400">
                        {classItem.description}
                      </p>
                    </div>
                    <div className="flex-shrink-0">
                      <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900/50 dark:text-blue-300">
                        🏫
                      </span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Quick Stats */}
        <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
          <h4 className="text-sm font-medium text-gray-500 dark:text-gray-400 mb-3">
            {t('teacher_page.quick_stats')}
          </h4>
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-3">
            <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-3">
              <div className="text-2xl font-bold text-blue-600 dark:text-blue-400">
                {profile.classes?.length || 0}
              </div>
              <div className="text-xs text-blue-600 dark:text-blue-400">
                {t('teacher_page.total_classes')}
              </div>
            </div>
            <div className="bg-green-50 dark:bg-green-900/20 rounded-lg p-3">
              <div className="text-2xl font-bold text-green-600 dark:text-green-400">
                {profile.total_students || 0}
              </div>
              <div className="text-xs text-green-600 dark:text-green-400">
                {t('teacher_page.total_students')}
              </div>
            </div>
            <div className="bg-orange-50 dark:bg-orange-900/20 rounded-lg p-3">
              <div className="text-2xl font-bold text-orange-600 dark:text-orange-400">
                {profile.pending_requests || 0}
              </div>
              <div className="text-xs text-orange-600 dark:text-orange-400">
                {t('teacher_page.pending_requests')}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default TeacherProfileView;