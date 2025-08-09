import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useForm } from 'react-hook-form';
import DataTable from '../components/DataTable';
import type { Column } from '../components/DataTable';
import Modal, { ConfirmModal } from '../components/Modal';
import { useToast } from '../components/Toast';
import { attendanceApi, studentsApi, classesApi } from '../services/api';
import type { Attendance, AttendanceFormData, Student, Class } from '../types/models';

export const AttendancePage: React.FC = () => {
  const { t } = useTranslation();
  const { showSuccess, showError } = useToast();
  const [attendance, setAttendance] = useState<Attendance[]>([]);
  const [students, setStudents] = useState<Student[]>([]);
  const [classes, setClasses] = useState<Class[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [confirmModal, setConfirmModal] = useState<{
    isOpen: boolean;
    attendanceRecord: Attendance | null;
  }>({ isOpen: false, attendanceRecord: null });
  const [editingAttendance, setEditingAttendance] = useState<Attendance | null>(null);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [filters, setFilters] = useState({
    class_id: '',
    date: '',
    status: '',
  });

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
    setValue,
  } = useForm<AttendanceFormData>();

  useEffect(() => {
    fetchAttendance();
    fetchStudents();
    fetchClasses();
  }, [pagination.current, pagination.pageSize, filters]);

  const fetchAttendance = async () => {
    try {
      setLoading(true);

      const pageLimit = pagination.pageSize;
      const pageOffset = (pagination.current - 1) * pagination.pageSize;

      const hasClass = !!filters.class_id;
      const hasDate = !!filters.date;
      const hasStatus = !!filters.status;

      // Helper to apply client-side filters (status and optional second filter), then paginate locally
      const applyClientFiltersAndPaginate = (items: Attendance[]) => {
        let filtered = items;
        if (hasClass) {
          filtered = filtered.filter(a => String(a.class_id) === String(filters.class_id));
        }
        if (hasDate) {
          filtered = filtered.filter(a => new Date(a.date).toISOString().slice(0,10) === filters.date);
        }
        if (hasStatus) {
          filtered = filtered.filter(a => a.status === filters.status);
        }
        const total = filtered.length;
        const start = pageOffset;
        const end = start + pageLimit;
        const paged = filtered.slice(start, end);
        setAttendance(paged);
        setPagination(prev => ({ ...prev, total }));
      };

      // Decide which API to use based on filters
      if (!hasClass && !hasDate && !hasStatus) {
        // No filters: server-side pagination
        const res = await attendanceApi.getAll({ limit: pageLimit, offset: pageOffset });
        setAttendance(res.data || []);
        setPagination(prev => ({ ...prev, total: res.total || 0 }));
        return;
      }

      // At least one filter is applied.
      // Use the most restrictive server filter available, then apply remaining filters client-side.
      // Fetch a larger batch to allow accurate client-side filtering and pagination.
      const BULK_LIMIT = 1000;

      if (hasClass && !hasDate && !hasStatus) {
        // Only class filter -> server-side pagination okay
        const res = await attendanceApi.getByClassId(Number(filters.class_id), { limit: pageLimit, offset: pageOffset });
        const items = (res.data as unknown as Attendance[]) || [];
        // Backend endpoint returns ApiResponse with data array but no total; estimate by length when page size reached
        setAttendance(items);
        setPagination(prev => ({ ...prev, total: (res as any).total ?? prev.total }));
        return;
      }

      if (hasDate && !hasClass && !hasStatus) {
        // Only date filter -> server-side pagination okay using date-range with same start/end
        const res = await attendanceApi.getByDateRange(filters.date, filters.date, { limit: pageLimit, offset: pageOffset });
        const items = (res.data as unknown as Attendance[]) || [];
        setAttendance(items);
        setPagination(prev => ({ ...prev, total: (res as any).total ?? prev.total }));
        return;
      }

      // For combinations or status filter, fetch bulk with a primary server filter to reduce records
      let bulkItems: Attendance[] = [];
      if (hasClass) {
        const res = await attendanceApi.getByClassId(Number(filters.class_id), { limit: BULK_LIMIT, offset: 0 });
        bulkItems = ((res as any).data || []) as Attendance[];
      } else if (hasDate) {
        const res = await attendanceApi.getByDateRange(filters.date, filters.date, { limit: BULK_LIMIT, offset: 0 });
        bulkItems = ((res as any).data || []) as Attendance[];
      } else {
        const res = await attendanceApi.getAll({ limit: BULK_LIMIT, offset: 0 });
        bulkItems = (res.data || []) as Attendance[];
      }

      applyClientFiltersAndPaginate(bulkItems);

    } catch (error: any) {
      console.error('Failed to fetch attendance:', error);
      showError(
        t('common.error'),
        error?.message || t('error.failed_to_load')
      );
    } finally {
      setLoading(false);
    }
  };

  const fetchStudents = async () => {
    try {
      const response = await studentsApi.getAll({ limit: 1000, offset: 0 });
      setStudents(response.data || []);
    } catch (error: any) {
      console.error('Failed to fetch students:', error);
      showError(
        t('common.error'),
        error?.message || 'Failed to fetch students'
      );
    }
  };

  const fetchClasses = async () => {
    try {
      const response = await classesApi.getAll({ limit: 100, offset: 0 });
      setClasses(response.data || []);
    } catch (error: any) {
      console.error('Failed to fetch classes:', error);
      showError(
        t('common.error'),
        error?.message || 'Failed to fetch classes'
      );
    }
  };

  const handleCreate = () => {
    setEditingAttendance(null);
    reset();
    setModalOpen(true);
  };

  const handleEdit = (attendanceRecord: Attendance) => {
    setEditingAttendance(attendanceRecord);
    setValue('student_id', attendanceRecord.student_id);
    setValue('class_id', attendanceRecord.class_id);
    setValue('date', attendanceRecord.date.split('T')[0]);
    setValue('status', attendanceRecord.status);
    setValue('description', attendanceRecord.description || '');
    setModalOpen(true);
  };

  const handleDelete = (attendanceRecord: Attendance) => {
    setConfirmModal({ isOpen: true, attendanceRecord });
  };

  const handleConfirmDelete = async () => {
    if (!confirmModal.attendanceRecord) return;
    
    try {
      await attendanceApi.delete(confirmModal.attendanceRecord.id);
      showSuccess(
        t('attendance.deleted_successfully'),
        t('attendance.attendance_deleted_message')
      );
      await fetchAttendance();
      setConfirmModal({ isOpen: false, attendanceRecord: null });
    } catch (error: any) {
      console.error('Failed to delete attendance:', error);
      showError(
        t('common.error'),
        error?.message || t('attendance.delete_failed')
      );
    }
  };

  const onSubmit = async (data: AttendanceFormData) => {
    try {
      const formattedData = {
        ...data,
        date: new Date(data.date).toISOString(),
      };
      
      if (editingAttendance) {
        await attendanceApi.update(editingAttendance.id, formattedData);
        showSuccess(
          t('attendance.updated_successfully'),
          t('attendance.attendance_updated_message')
        );
      } else {
        await attendanceApi.create(formattedData);
        showSuccess(
          t('attendance.created_successfully'),
          t('attendance.attendance_created_message')
        );
      }
      await fetchAttendance();
      setModalOpen(false);
      setEditingAttendance(null);
      reset();
    } catch (error: any) {
      console.error('Failed to save attendance:', error);
      showError(
        t('common.error'),
        error?.message || t('attendance.save_failed')
      );
    }
  };

  const getStudentName = (studentId: string) => {
    const student = students.find(s => s.student_id === studentId);
    return student ? `${student.first_name} ${student.last_name}` : studentId;
  };

  const getClassName = (classId: number) => {
    const classItem = classes.find(c => c.id === classId);
    return classItem ? classItem.name : `Class ${classId}`;
  };

  const getStatusBadge = (status: string) => {
    const statusClasses = {
      present: 'bg-green-100 text-green-800 dark:bg-green-900/50 dark:text-green-300',
      absent: 'bg-red-100 text-red-800 dark:bg-red-900/50 dark:text-red-300',
      late: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/50 dark:text-yellow-300',
      excused: 'bg-blue-100 text-blue-800 dark:bg-blue-900/50 dark:text-blue-300',
    };
    
    return (
      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${statusClasses[status as keyof typeof statusClasses] || 'bg-gray-100 text-gray-800'}`}>
        {t(`attendance.status.${status}`)}
      </span>
    );
  };

  const columns: Column<Attendance>[] = [
    {
      key: 'date',
      title: t('attendance.date'),
      render: (value) => new Date(value).toLocaleDateString(),
      width: '32',
    },
    {
      key: 'student_id',
      title: t('attendance.student'),
      render: (value) => getStudentName(value),
    },
    {
      key: 'class_id',
      title: t('attendance.class'),
      render: (value) => getClassName(value),
    },
    {
      key: 'status',
      title: t('attendance.status.label'),
      render: (value) => getStatusBadge(value),
    },
    {
      key: 'description',
      title: t('attendance.notes'),
      render: (value) => value || '-',
    },
    {
      key: 'created_at',
      title: t('common.created_at'),
      render: (value) => new Date(value).toLocaleDateString(),
    },
  ];

  const modalFooter = (
    <>
      <button
        type="submit"
        form="attendance-form"
        disabled={isSubmitting}
        className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {isSubmitting ? (
          <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
          </svg>
        ) : null}
        {editingAttendance ? t('common.update') : t('common.create')}
      </button>
      <button
        type="button"
        onClick={() => {
          setModalOpen(false);
          setEditingAttendance(null);
          reset();
        }}
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
            {t('attendance.title')}
          </h1>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-300">
            {t('attendance.description')}
          </p>
        </div>
        <button
          onClick={handleCreate}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          <svg className="-ml-1 mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          {t('attendance.add_record')}
        </button>
      </div>

      {/* Filters */}
      <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('attendance.filter_class')}
            </label>
            <select
              value={filters.class_id}
              onChange={(e) => setFilters(prev => ({ ...prev, class_id: e.target.value }))}
              className="block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">{t('common.all')}</option>
              {classes.map((classItem) => (
                <option key={classItem.id} value={classItem.id}>
                  {classItem.name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('attendance.filter_date')}
            </label>
            <input
              type="date"
              value={filters.date}
              onChange={(e) => setFilters(prev => ({ ...prev, date: e.target.value }))}
              className="block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              {t('attendance.filter_status')}
            </label>
            <select
              value={filters.status}
              onChange={(e) => setFilters(prev => ({ ...prev, status: e.target.value }))}
              className="block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">{t('common.all')}</option>
              <option value="present">{t('attendance.status.present')}</option>
              <option value="absent">{t('attendance.status.absent')}</option>
              <option value="late">{t('attendance.status.late')}</option>
              <option value="excused">{t('attendance.status.excused')}</option>
            </select>
          </div>

          <div className="flex items-end">
            <button
              onClick={() => setFilters({ class_id: '', date: '', status: '' })}
              className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              {t('common.clear_filters')}
            </button>
          </div>
        </div>
      </div>

      {/* Data Table */}
      <DataTable
        data={attendance}
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
        onClose={() => {
          setModalOpen(false);
          setEditingAttendance(null);
          reset();
        }}
        title={editingAttendance ? t('attendance.edit_record') : t('attendance.add_record')}
        size="lg"
        footer={modalFooter}
      >
        <form id="attendance-form" onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('attendance.student')} *
              </label>
              <select
                {...register('student_id', { required: t('validation.required') })}
                className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="">{t('attendance.select_student')}</option>
                {students.map((student) => (
                  <option key={student.student_id} value={student.student_id}>
                    {student.first_name} {student.last_name} ({student.student_id})
                  </option>
                ))}
              </select>
              {errors.student_id && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.student_id.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('attendance.class')} *
              </label>
              <select
                {...register('class_id', { 
                  required: t('validation.required'),
                  valueAsNumber: true,
                })}
                className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="">{t('attendance.select_class')}</option>
                {classes.map((classItem) => (
                  <option key={classItem.id} value={classItem.id}>
                    {classItem.name}
                  </option>
                ))}
              </select>
              {errors.class_id && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.class_id.message}</p>
              )}
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('attendance.date')} *
              </label>
              <input
                type="date"
                {...register('date', { required: t('validation.required') })}
                className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              />
              {errors.date && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.date.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {t('attendance.status.label')} *
              </label>
              <select
                {...register('status', { required: t('validation.required') })}
                className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="">{t('attendance.select_status')}</option>
                <option value="present">{t('attendance.status.present')}</option>
                <option value="absent">{t('attendance.status.absent')}</option>
                <option value="late">{t('attendance.status.late')}</option>
                <option value="excused">{t('attendance.status.excused')}</option>
              </select>
              {errors.status && (
                <p className="mt-1 text-sm text-red-600 dark:text-red-400">{errors.status.message}</p>
              )}
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
              {t('attendance.notes')}
            </label>
            <textarea
              {...register('description')}
              rows={3}
              className="mt-1 block w-full border border-gray-300 dark:border-gray-600 rounded-md px-3 py-2 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
              placeholder={t('attendance.notes_placeholder')}
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
        onClose={() => setConfirmModal({ isOpen: false, attendanceRecord: null })}
        onConfirm={handleConfirmDelete}
        title={t('attendance.delete_record')}
        message={t('attendance.delete_confirm', { 
          student: confirmModal.attendanceRecord ? getStudentName(confirmModal.attendanceRecord.student_id) : '',
          date: confirmModal.attendanceRecord ? new Date(confirmModal.attendanceRecord.date).toLocaleDateString() : ''
        })}
        confirmText={t('common.delete')}
        type="danger"
      />
    </div>
  );
};

export default AttendancePage;