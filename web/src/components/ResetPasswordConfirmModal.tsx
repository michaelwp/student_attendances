import React from 'react';
import { useTranslation } from 'react-i18next';
import { ConfirmModal } from './Modal';

interface ResetPasswordConfirmModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => Promise<void>;
  title: string;
  message: string;
  loading?: boolean;
}

export const ResetPasswordConfirmModal: React.FC<ResetPasswordConfirmModalProps> = ({
  isOpen,
  onClose,
  onConfirm,
  title,
  message,
  loading = false,
}) => {
  const { t } = useTranslation();

  const handleConfirm = async () => {
    await onConfirm();
  };

  return (
    <ConfirmModal
      isOpen={isOpen}
      onClose={onClose}
      onConfirm={handleConfirm}
      title={title}
      message={message}
      confirmText={t('common.reset_password')}
      type="warning"
      loading={loading}
    />
  );
};

export default ResetPasswordConfirmModal;