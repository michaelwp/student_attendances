import type { ToastType } from '../components/Toast';

export interface ToastData {
  id: string;
  type: ToastType;
  title: string;
  message: string;
  duration?: number;
}

export interface UseToastReturn {
  showToast: (toast: Omit<ToastData, 'id'>) => void;
  showSuccess: (title: string, message: string) => void;
  showError: (title: string, message: string) => void;
  showWarning: (title: string, message: string) => void;
  showInfo: (title: string, message: string) => void;
}

export let toastId = 0;

export const generateToastId = () => `toast-${++toastId}`;

export const useToast = (): UseToastReturn => {
  const showToast = (toast: Omit<ToastData, 'id'>) => {
    const id = generateToastId();
    const newToast: ToastData = {
      ...toast,
      id,
    };
    
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const callback = (window as any).__addToastCallback;
    if (callback) {
      callback(newToast);
    }
  };

  const showSuccess = (title: string, message: string) => {
    showToast({ type: 'success', title, message });
  };

  const showError = (title: string, message: string) => {
    showToast({ type: 'error', title, message });
  };

  const showWarning = (title: string, message: string) => {
    showToast({ type: 'warning', title, message });
  };

  const showInfo = (title: string, message: string) => {
    showToast({ type: 'info', title, message });
  };

  return { showToast, showSuccess, showError, showWarning, showInfo };
};