import { useState, useEffect, forwardRef, useImperativeHandle } from 'react';
import { useTranslation } from 'react-i18next';
import { useToast } from '../utils/toast-helpers';
import { teacherDashboardApi } from '../services/api';
import type { AbsentRequest } from '../types/models';

export interface TeacherAbsentRequestListHandle {
  refreshRequests: () => void;
}

interface TeacherAbsentRequestListProps {
  onRequestStatusChanged?: () => void;
}

export const TeacherAbsentRequestList = forwardRef<TeacherAbsentRequestListHandle, TeacherAbsentRequestListProps>(
  ({ onRequestStatusChanged }, ref) => {
    const { t } = useTranslation();
    const { showError, showSuccess } = useToast();
    const [requests, setRequests] = useState<AbsentRequest[]>([]);
    const [loading, setLoading] = useState(false);
    const [processingRequests, setProcessingRequests] = useState<Set<number>>(new Set());
    const [pagination, setPagination] = useState({
      current: 1,
      pageSize: 10,
      total: 0,
    });

    useEffect(() => {
      fetchRequests();
    }, [pagination.current, pagination.pageSize]);

    useImperativeHandle(ref, () => ({
      refreshRequests: fetchRequests,
    }));

    const fetchRequests = async () => {
      try {
        setLoading(true);
        const response = await teacherDashboardApi.getAbsentRequests({
          limit: pagination.pageSize,
          offset: (pagination.current - 1) * pagination.pageSize,
        });
        setRequests(response.data || []);
        setPagination(prev => ({
          ...prev,
          total: response.total || 0,
        }));
      } catch (error: unknown) {
        console.error('Failed to fetch absent requests:', error);
        showError(
          t('common.error'),
          (error as Error)?.message || t('teacher_page.fetch_requests_failed')
        );
      } finally {
        setLoading(false);
      }
    };

    const handleApprove = async (requestId: number) => {
      try {
        setProcessingRequests(prev => new Set(prev).add(requestId));
        await teacherDashboardApi.approveAbsentRequest(requestId);
        
        showSuccess(
          t('common.success'),
          t('teacher_page.request_approved_successfully')
        );
        
        // Refresh the list
        await fetchRequests();
        
        // Notify parent component
        if (onRequestStatusChanged) {
          onRequestStatusChanged();
        }
      } catch (error: unknown) {
        console.error('Failed to approve request:', error);
        showError(
          t('common.error'),
          (error as Error)?.message || t('teacher_page.approve_request_failed')
        );
      } finally {
        setProcessingRequests(prev => {
          const newSet = new Set(prev);
          newSet.delete(requestId);
          return newSet;
        });
      }
    };

    const handleReject = async (requestId: number) => {
      try {
        setProcessingRequests(prev => new Set(prev).add(requestId));
        await teacherDashboardApi.rejectAbsentRequest(requestId);
        
        showSuccess(
          t('common.success'),
          t('teacher_page.request_rejected_successfully')
        );
        
        // Refresh the list
        await fetchRequests();
        
        // Notify parent component
        if (onRequestStatusChanged) {
          onRequestStatusChanged();
        }
      } catch (error: unknown) {
        console.error('Failed to reject request:', error);
        showError(
          t('common.error'),
          (error as Error)?.message || t('teacher_page.reject_request_failed')
        );
      } finally {
        setProcessingRequests(prev => {
          const newSet = new Set(prev);
          newSet.delete(requestId);
          return newSet;
        });
      }
    };

    const getStatusBadge = (status: string) => {
      const statusMap = {
        pending: {
          color: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/50 dark:text-yellow-300',
          icon: '⏳',
          label: t('teacher_page.status_pending')
        },
        approved: {
          color: 'bg-green-100 text-green-800 dark:bg-green-900/50 dark:text-green-300',
          icon: '✅',
          label: t('teacher_page.status_approved')
        },
        rejected: {
          color: 'bg-red-100 text-red-800 dark:bg-red-900/50 dark:text-red-300',
          icon: '❌',
          label: t('teacher_page.status_rejected')
        }
      };

      const statusInfo = statusMap[status as keyof typeof statusMap] || statusMap.pending;

      return (
        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${statusInfo.color}`}>
          <span className="mr-1">{statusInfo.icon}</span>
          {statusInfo.label}
        </span>
      );
    };

    const formatDate = (dateString: string) => {
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      });
    };

    const formatDateTime = (dateString: string | null) => {
      if (!dateString) return 'N/A';
      return new Date(dateString).toLocaleString('en-US', {
        year: 'numeric',
        month: 'short',
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
                <div className="w-8 h-8 bg-purple-100 dark:bg-purple-900/50 rounded-lg flex items-center justify-center">
                  <span className="text-purple-600 dark:text-purple-400">📋</span>
                </div>
              </div>
              <div className="ml-3">
                <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white">
                  {t('teacher_page.absent_requests_list')}
                </h3>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  {t('teacher_page.absent_requests_list_desc')}
                </p>
              </div>
            </div>
            <button
              onClick={fetchRequests}
              className="inline-flex items-center px-3 py-2 border border-gray-300 dark:border-gray-600 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              {t('common.refresh')}
            </button>
          </div>

          {loading ? (
            <div className="flex items-center justify-center py-12">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>
          ) : requests.length === 0 ? (
            <div className="text-center py-12">
              <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <h3 className="mt-2 text-sm font-medium text-gray-900 dark:text-white">
                {t('teacher_page.no_requests')}
              </h3>
              <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                {t('teacher_page.no_requests_desc')}
              </p>
            </div>
          ) : (
            <div className="space-y-4">
              {requests.map((request) => {
                const isProcessing = processingRequests.has(request.id);
                const isPending = request.status === 'pending';
                
                return (
                  <div key={request.id} className="border border-gray-200 dark:border-gray-600 rounded-lg p-4">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="flex items-center space-x-4 mb-2">
                          <div className="text-sm font-medium text-gray-900 dark:text-white">
                            #{request.id}
                          </div>
                          <div className="text-sm font-medium text-blue-600 dark:text-blue-400">
                            {t('teacher_page.student')}: {request.student_id}
                          </div>
                          <div className="text-sm text-gray-500 dark:text-gray-400">
                            {formatDate(request.request_date)}
                          </div>
                          {getStatusBadge(request.status)}
                        </div>
                        
                        <div className="mb-3">
                          <h4 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                            {t('teacher_page.reason')}:
                          </h4>
                          <p className="text-sm text-gray-600 dark:text-gray-400 bg-gray-50 dark:bg-gray-700 rounded p-2">
                            {request.reason}
                          </p>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs text-gray-500 dark:text-gray-400">
                          <div>
                            <span className="font-medium">{t('teacher_page.submitted')}:</span>
                            <br />
                            {formatDateTime(request.created_at)}
                          </div>
                          
                          {request.status === 'approved' && (
                            <div>
                              <span className="font-medium">{t('teacher_page.approved_by')}:</span>
                              <br />
                              {request.approved_by ? `Teacher #${request.approved_by}` : 'N/A'}
                              <br />
                              <span className="font-medium">{t('teacher_page.approved_at')}:</span>
                              <br />
                              {formatDateTime(request.approved_at || null)}
                            </div>
                          )}
                          
                          {request.status === 'rejected' && (
                            <div>
                              <span className="font-medium">{t('teacher_page.rejected_by')}:</span>
                              <br />
                              {request.rejected_by ? `Teacher #${request.rejected_by}` : 'N/A'}
                              <br />
                              <span className="font-medium">{t('teacher_page.rejected_at')}:</span>
                              <br />
                              {formatDateTime(request.rejected_at || null)}
                            </div>
                          )}
                        </div>
                      </div>

                      {/* Action Buttons */}
                      {isPending && (
                        <div className="flex space-x-2 ml-4">
                          <button
                            onClick={() => handleApprove(request.id)}
                            disabled={isProcessing}
                            className="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
                          >
                            {isProcessing ? (
                              <svg className="animate-spin -ml-1 mr-1 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                              </svg>
                            ) : (
                              <svg className="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                              </svg>
                            )}
                            {t('teacher_page.approve')}
                          </button>
                          <button
                            onClick={() => handleReject(request.id)}
                            disabled={isProcessing}
                            className="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed"
                          >
                            {isProcessing ? (
                              <svg className="animate-spin -ml-1 mr-1 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                              </svg>
                            ) : (
                              <svg className="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                              </svg>
                            )}
                            {t('teacher_page.reject')}
                          </button>
                        </div>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          )}

          {/* Pagination */}
          {requests.length > 0 && (
            <div className="mt-6 flex items-center justify-between border-t border-gray-200 dark:border-gray-700 pt-4">
              <div className="text-sm text-gray-700 dark:text-gray-300">
                {t('common.showing')}{' '}
                <span className="font-medium">
                  {Math.min((pagination.current - 1) * pagination.pageSize + 1, pagination.total)}
                </span>{' '}
                {t('common.to')}{' '}
                <span className="font-medium">
                  {Math.min(pagination.current * pagination.pageSize, pagination.total)}
                </span>{' '}
                {t('common.of')}{' '}
                <span className="font-medium">{pagination.total}</span>{' '}
                {t('teacher_page.requests')}
              </div>
              <div className="flex space-x-2">
                <button
                  onClick={() => setPagination(prev => ({ ...prev, current: Math.max(1, prev.current - 1) }))}
                  disabled={pagination.current === 1}
                  className="px-3 py-1 text-sm border border-gray-300 dark:text-gray-300 rounded hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {t('common.previous')}
                </button>
                <span className="px-3 py-1 text-sm dark:text-gray-300">
                  {pagination.current}
                </span>
                <button
                  onClick={() => setPagination(prev => ({ ...prev, current: prev.current + 1 }))}
                  disabled={pagination.current * pagination.pageSize >= pagination.total}
                  className="px-3 py-1 text-sm border border-gray-300 dark:text-gray-300 rounded hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {t('common.next')}
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    );
  }
);

export default TeacherAbsentRequestList;