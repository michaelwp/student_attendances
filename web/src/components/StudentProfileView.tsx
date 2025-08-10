import React from 'react';
import { useTranslation } from 'react-i18next';
import type { StudentProfile } from '../types/models';

interface StudentProfileViewProps {
  profile: StudentProfile;
}

export const StudentProfileView: React.FC<StudentProfileViewProps> = ({ profile }) => {
  const { t } = useTranslation();

  return (
    <div className="bg-white dark:bg-gray-800 shadow rounded-lg">
      <div className="px-4 py-5 sm:p-6">
        <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white mb-4">
          {t('student_page.profile_info')}
        </h3>
        <p className="text-sm text-gray-500 dark:text-gray-400 mb-6">
          {t('student_page.profile_info_desc')}
        </p>

        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {t('student_page.student_id')}
            </label>
            <div className="mt-1 p-3 bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-md">
              <span className="text-sm text-gray-900 dark:text-gray-100 font-medium">
                {profile.student_id}
              </span>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {t('student_page.class')}
            </label>
            <div className="mt-1 p-3 bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-md">
              <span className="text-sm text-gray-900 dark:text-gray-100">
                {profile.class_name || 'Not assigned'}
              </span>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {t('student_page.first_name')}
            </label>
            <div className="mt-1 p-3 bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-md">
              <span className="text-sm text-gray-900 dark:text-gray-100">
                {profile.first_name}
              </span>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {t('student_page.last_name')}
            </label>
            <div className="mt-1 p-3 bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-md">
              <span className="text-sm text-gray-900 dark:text-gray-100">
                {profile.last_name}
              </span>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {t('student_page.email')}
            </label>
            <div className="mt-1 p-3 bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-md">
              <span className="text-sm text-gray-900 dark:text-gray-100">
                {profile.email}
              </span>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              {t('student_page.phone')}
            </label>
            <div className="mt-1 p-3 bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-md">
              <span className="text-sm text-gray-900 dark:text-gray-100">
                {profile.phone}
              </span>
            </div>
          </div>
        </div>

        <div className="mt-8 grid grid-cols-1 gap-6 sm:grid-cols-4">
          <div className="bg-green-50 dark:bg-green-900/50 p-4 rounded-lg">
            <div className="text-2xl font-bold text-green-600 dark:text-green-400">
              {profile.attendance_stats?.present_days || 0}
            </div>
            <div className="text-sm text-green-600 dark:text-green-400 font-medium">
              {t('student_page.present_days')}
            </div>
          </div>

            <div className="bg-blue-50 dark:bg-blue-900/50 p-4 rounded-lg">
                <div className="text-2xl font-bold text-blue-600 dark:text-blue-400">
                    {profile.attendance_stats?.excused_days || 0}
                </div>
                <div className="text-sm text-blue-600 dark:text-blue-400 font-medium">
                    {t('student_page.excused_days')}
                </div>
            </div>

          <div className="bg-red-50 dark:bg-red-900/50 p-4 rounded-lg">
            <div className="text-2xl font-bold text-red-600 dark:text-red-400">
              {profile.attendance_stats?.absent_days || 0}
            </div>
            <div className="text-sm text-red-600 dark:text-red-400 font-medium">
              {t('student_page.absent_days')}
            </div>
          </div>

          <div className="bg-yellow-50 dark:bg-yellow-900/50 p-4 rounded-lg">
            <div className="text-2xl font-bold text-yellow-600 dark:text-yellow-400">
              {profile.attendance_stats?.absent_days || 0}
            </div>
            <div className="text-sm text-yellow-600 dark:text-yellow-400 font-medium">
              {t('student_page.late_days')}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default StudentProfileView;