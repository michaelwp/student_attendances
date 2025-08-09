import React, { useState, useRef, useCallback } from 'react';
import { useTranslation } from 'react-i18next';

interface PhotoUploadProps {
  currentPhotoUrl?: string;
  onPhotoSelect: (file: File) => void;
  onPhotoRemove?: () => void;
  className?: string;
  size?: 'sm' | 'md' | 'lg';
}

export const PhotoUpload: React.FC<PhotoUploadProps> = ({
  currentPhotoUrl,
  onPhotoSelect,
  onPhotoRemove,
  className = '',
  size = 'md',
}) => {
  const { t } = useTranslation();
  const [dragOver, setDragOver] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const sizeClasses = {
    sm: 'w-20 h-20',
    md: 'w-32 h-32',
    lg: 'w-48 h-48',
  };

  const validateFile = useCallback((file: File): boolean => {
    setError(null);

    // Check file type
    if (!file.type.startsWith('image/')) {
      setError(t('validation.invalid_file_type'));
      return false;
    }

    // Check file size (5MB limit)
    const maxSize = 5 * 1024 * 1024; // 5MB
    if (file.size > maxSize) {
      setError(t('validation.file_too_large', { size: 5 }));
      return false;
    }

    return true;
  }, [t]);

  const handleFileSelect = useCallback((file: File) => {
    if (validateFile(file)) {
      onPhotoSelect(file);
    }
  }, [validateFile, onPhotoSelect]);

  const handleFileInput = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      handleFileSelect(file);
    }
  }, [handleFileSelect]);

  const handleDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault();
    setDragOver(true);
  }, []);

  const handleDragLeave = useCallback((event: React.DragEvent) => {
    event.preventDefault();
    setDragOver(false);
  }, []);

  const handleDrop = useCallback((event: React.DragEvent) => {
    event.preventDefault();
    setDragOver(false);

    const file = event.dataTransfer.files[0];
    if (file) {
      handleFileSelect(file);
    }
  }, [handleFileSelect]);

  const handleClick = () => {
    fileInputRef.current?.click();
  };

  const handleRemove = () => {
    if (onPhotoRemove) {
      onPhotoRemove();
    }
    // Reset file input
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
    setError(null);
  };

  return (
    <div className={`space-y-2 ${className}`}>
      <div className="relative">
        <div
          className={`
            ${sizeClasses[size]}
            border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg
            flex flex-col items-center justify-center cursor-pointer
            transition-colors duration-200
            ${dragOver ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20' : 'hover:border-blue-400 hover:bg-gray-50 dark:hover:bg-gray-700'}
            ${currentPhotoUrl ? 'border-solid' : ''}
          `}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          onClick={handleClick}
        >
          {currentPhotoUrl ? (
            <>
              <img
                src={currentPhotoUrl}
                alt="Current photo"
                className={`${sizeClasses[size]} object-cover rounded-lg`}
              />
              {onPhotoRemove && (
                <button
                  type="button"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleRemove();
                  }}
                  className="absolute -top-2 -right-2 p-1 bg-red-500 text-white rounded-full hover:bg-red-600 transition-colors"
                  title={t('common.remove_photo')}
                >
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              )}
              <div className="absolute inset-0 bg-black bg-opacity-0 hover:bg-opacity-20 rounded-lg transition-opacity flex items-center justify-center">
                <span className="text-white opacity-0 hover:opacity-100 text-sm font-medium">
                  {t('common.change_photo')}
                </span>
              </div>
            </>
          ) : (
            <>
              <svg className="w-8 h-8 text-gray-400 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <p className="text-sm text-gray-500 dark:text-gray-400 text-center px-2">
                {t('common.drag_drop_photo')}
              </p>
            </>
          )}
        </div>

        <input
          ref={fileInputRef}
          type="file"
          accept="image/*"
          onChange={handleFileInput}
          className="hidden"
        />
      </div>

      {/* Action Buttons */}
      <div className="flex space-x-2">
        <button
          type="button"
          onClick={handleClick}
          className="px-3 py-1 text-xs bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
        >
          {currentPhotoUrl ? t('common.change_photo') : t('common.upload_photo')}
        </button>
        {currentPhotoUrl && onPhotoRemove && (
          <button
            type="button"
            onClick={handleRemove}
            className="px-3 py-1 text-xs bg-red-600 text-white rounded hover:bg-red-700 transition-colors"
          >
            {t('common.remove_photo')}
          </button>
        )}
      </div>

      {/* Requirements */}
      <p className="text-xs text-gray-500 dark:text-gray-400">
        {t('common.photo_requirements')}
      </p>

      {/* Error Message */}
      {error && (
        <p className="text-xs text-red-600 dark:text-red-400">
          {error}
        </p>
      )}
    </div>
  );
};

export default PhotoUpload;