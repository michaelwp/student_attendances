import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useForm } from 'react-hook-form';
import DataTable from '../components/DataTable';
import type { Column } from '../components/DataTable';
import Modal, { ConfirmModal } from '../components/Modal';
import { useToast } from '../utils/toast-helpers';
import { classesApi, teachersApi } from '../services/api';
import type { Class, ClassFormData, Teacher } from '../types/models';

export const ClassesPage: React.FC = () => {
  const { t } = useTranslation();
  const { showSuccess, showError } = useToast();
  const [classes, setClasses] = useState<Class[]>([]);
  const [teachers, setTeachers] = useState<Teacher[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [confirmModal, setConfirmModal] = useState<{
    isOpen: boolean;
    classItem: Class | null;
  }>({ isOpen: false, classItem: null });
  const [editingClass, setEditingClass] = useState<Class | null>(null);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
    setValue,
  } = useForm<ClassFormData>();

  useEffect(() => {
    fetchClasses();
    fetchTeachers();
  }, [pagination.current, pagination.pageSize]);

  const fetchClasses = async () => {
    try {
      setLoading(true);
      const response = await classesApi.getAll({
        limit: pagination.pageSize,
        offset: (pagination.current - 1) * pagination.pageSize,
      });
      setClasses(response.data || []);
      setPagination(prev => ({
        ...prev,
        total: response.total || 0,
      }));
    } catch (error: unknown) {
      console.error('Failed to fetch classes:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchTeachers = async () => {
    try {
      const response = await teachersApi.getAll({ limit: 100, offset: 0 });
      setTeachers(response.data || []);
    } catch (error: unknown) {
      console.error('Failed to fetch teachers:', error);
    }
  };

  const handleCreate = () => {
    setEditingClass(null);
    reset();
    setModalOpen(true);
  };

  const handleEdit = (classItem: Class) => {
    setEditingClass(classItem);
    setValue('name', classItem.name);
    setValue('homeroom_teacher', classItem.homeroom_teacher);
    setValue('description', classItem.description);
    setModalOpen(true);
  };

  const handleDelete = (classItem: Class) => {
    setConfirmModal({ isOpen: true, classItem });
  };

  const handleConfirmDelete = async () => {
    if (!confirmModal.classItem) return;
    
    try {
      await classesApi.delete(confirmModal.classItem.id);
      showSuccess(
        t('classes.deleted_successfully'),
        t('classes.class_deleted_message')
      );
      await fetchClasses();
      setConfirmModal({ isOpen: false, classItem: null });
    } catch (error: unknown) {
      console.error('Failed to delete class:', error);
      showError(
        t('common.error'),
        (error as Error)?.message || t('classes.delete_failed')
      );
    }
  };

  const onSubmit = async (data: ClassFormData) => {
    try {
      if (editingClass) {
        await classesApi.update(editingClass.id, data);
        showSuccess(
          t('classes.updated_successfully'),
          t('classes.class_updated_message')
        );
      } else {
        await classesApi.create(data);
        showSuccess(
          t('classes.created_successfully'),
          t('classes.class_created_message')
        );
      }
      await fetchClasses();
      setModalOpen(false);
      reset();
    } catch (error: unknown) {
      console.error('Failed to save class:', error);
      showError(
        t('common.error'),
        (error as Error)?.message || t('classes.save_failed')
      );
    }
  };

  const getTeacherName = (teacherId: string) => {
    const teacher = teachers.find(t => t.teacher_id === teacherId);
    return teacher ? `${teacher.first_name} ${teacher.last_name}` : teacherId;
  };

  const columns: Column<Class>[] = [
    {
      key: 'name',
      title: t('classes.name'),
      width: '32',
    },
    {
      key: 'homeroom_teacher',
      title: t('classes.homeroom_teacher'),
      render: (value) => getTeacherName(value as string),
    },
    {
      key: 'description',
      title: t('classes.description'),
      render: (value) => (value as string) || '-',
    },
    {
      key: 'created_at',
      title: t('common.created_at'),
      render: (value) => new Date(value as string).toLocaleDateString(),
    },
  ];

  const modalFooter = (
    <>
      <button
        type="submit"
        form="class-form"
        disabled={isSubmitting}
        className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {isSubmitting ? (
          <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
        ) : null}
        {editingClass ? t('common.update') : t('common.create')}
      </button>
      <button
        type="button"
        onClick={() => setModalOpen(false)}
        disabled={isSubmitting}
        className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 dark:border-gray-600 shadow-sm px-4 py-2 bg-white dark:bg-gray-800 text-base font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {t('common.cancel')}
      </button>
    </>
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            {t('classes.title')}
          </h1>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-300">
            {t('classes.description')}
          </p>
        </div>
        <button
          onClick={handleCreate}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          <svg className="-ml-1 mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          {t('classes.add_class')}
        </button>
      </div>

      {/* Data Table */}
      <DataTable
        data={classes}
        columns={columns}
        loading={loading}
        onEdit={handleEdit}
        onDelete={handleDelete}
        pagination={{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: pagination.total,
          onChange: (page, pageSize) => {
            setPagination(prev => ({
              ...prev,
              current: page,
              pageSize,
            }));
          },
        }}
      />

      {/* Create/Edit Modal */}
      <Modal
        isOpen={modalOpen}
        onClose={() => setModalOpen(false)}
        title={editingClass ? t('classes.edit_class') : t('classes.add_class')}
        size="lg"
        footer={modalFooter}
      >
        <form id="class-form" onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              {t('classes.name')} *
            </label>
            <input
              type="text"
              {...register('name', { required: t('validation.required') })}
              className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder={t('classes.name_placeholder')}
            />
            {errors.name && (
              <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.name.message}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              {t('classes.homeroom_teacher')} *
            </label>
            <select
              {...register('homeroom_teacher', { required: t('validation.required') })}
              className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">{t('classes.select_teacher')}</option>
              {teachers.map((teacher) => (
                <option key={teacher.teacher_id} value={teacher.teacher_id}>
                  {teacher.first_name} {teacher.last_name} ({teacher.teacher_id})
                </option>
              ))}
            </select>
            {errors.homeroom_teacher && (
              <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.homeroom_teacher.message}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              {t('classes.description')}
            </label>
            <textarea
              {...register('description')}
              rows={3}
              className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder={t('classes.description_placeholder')}
            />
            {errors.description && (
              <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.description.message}</p>
            )}
          </div>
        </form>
      </Modal>

      {/* Delete Confirmation Modal */}
      <ConfirmModal
        isOpen={confirmModal.isOpen}
        onClose={() => setConfirmModal({ isOpen: false, classItem: null })}
        onConfirm={handleConfirmDelete}
        title={t('classes.delete_class')}
        message={t('classes.delete_confirm', { 
          name: confirmModal.classItem ? confirmModal.classItem.name : '' 
        })}
        confirmText={t('common.delete')}
        type="danger"
      />
    </div>
  );
};

export default ClassesPage;