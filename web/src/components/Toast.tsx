import React, { useEffect, useState, useCallback } from 'react';
import type { ToastData } from '../utils/toast-helpers';

export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface ToastProps {
  id: string;
  type: ToastType;
  title: string;
  message: string;
  duration?: number;
  onClose: (id: string) => void;
}

export const Toast: React.FC<ToastProps> = ({
  id,
  type,
  title,
  message,
  duration = 5000,
  onClose,
}) => {
  const [isVisible, setIsVisible] = useState(false);
  const [isClosing, setIsClosing] = useState(false);

  const handleClose = useCallback(() => {
    setIsClosing(true);
    setTimeout(() => {
      onClose(id);
    }, 300);
  }, [onClose, id]);

  useEffect(() => {
    // Show animation
    const showTimer = setTimeout(() => setIsVisible(true), 100);
    
    // Auto close after duration
    const closeTimer = setTimeout(() => {
      handleClose();
    }, duration);

    return () => {
      clearTimeout(showTimer);
      clearTimeout(closeTimer);
    };
  }, [duration, handleClose]);

  const getTypeStyles = () => {
    switch (type) {
      case 'success':
        return {
          bg: 'bg-green-500',
          border: 'border-green-400',
          icon: '✓',
          iconBg: 'bg-green-600',
        };
      case 'error':
        return {
          bg: 'bg-red-500',
          border: 'border-red-400',
          icon: '✕',
          iconBg: 'bg-red-600',
        };
      case 'warning':
        return {
          bg: 'bg-yellow-500',
          border: 'border-yellow-400',
          icon: '⚠',
          iconBg: 'bg-yellow-600',
        };
      case 'info':
        return {
          bg: 'bg-blue-500',
          border: 'border-blue-400',
          icon: 'ℹ',
          iconBg: 'bg-blue-600',
        };
      default:
        return {
          bg: 'bg-gray-500',
          border: 'border-gray-400',
          icon: 'ℹ',
          iconBg: 'bg-gray-600',
        };
    }
  };

  const styles = getTypeStyles();

  return (
    <div
      className={`
        fixed top-4 right-4 z-50 max-w-sm w-full
        transform transition-all duration-300 ease-in-out
        ${isVisible && !isClosing ? 'translate-x-0 opacity-100' : 'translate-x-full opacity-0'}
      `}
    >
      <div className={`rounded-lg shadow-lg overflow-hidden border ${styles.border}`}>
        <div className={`${styles.bg} p-4 text-white`}>
          <div className="flex items-start">
            <div className={`flex-shrink-0 w-6 h-6 rounded-full ${styles.iconBg} flex items-center justify-center mr-3`}>
              <span className="text-sm font-bold">{styles.icon}</span>
            </div>
            <div className="flex-1 min-w-0">
              <h4 className="text-sm font-medium mb-1">{title}</h4>
              <p className="text-sm opacity-90">{message}</p>
            </div>
            <button
              onClick={handleClose}
              className="flex-shrink-0 ml-2 text-white hover:text-gray-200 transition-colors"
            >
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};


export const ToastContainer: React.FC = () => {
  const [toasts, setToasts] = useState<ToastData[]>([]);

  useEffect(() => {
    // Direct assignment to the imported variable
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (window as any).__addToastCallback = (toast: ToastData) => {
      setToasts(prev => [...prev, toast]);
    };

    return () => {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (window as any).__addToastCallback = null;
    };
  }, []);

  const removeToast = (id: string) => {
    setToasts(prev => prev.filter(toast => toast.id !== id));
  };

  return (
    <>
      {toasts.map((toast, index) => (
        <div
          key={toast.id}
          style={{ top: `${20 + index * 80}px` }}
          className="fixed right-4 z-50"
        >
          <Toast
            {...toast}
            onClose={removeToast}
          />
        </div>
      ))}
    </>
  );
};

export default Toast;