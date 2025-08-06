import React from 'react';
import { useTranslation } from 'react-i18next';
import { Sun, Moon, Monitor } from 'lucide-react';
import { useThemeStore, type Theme } from '../stores/themeStore';

const ThemeToggle: React.FC = () => {
  const { t } = useTranslation();
  const { theme, setTheme } = useThemeStore();

  const themes: { value: Theme; label: string; icon: React.ReactNode }[] = [
    {
      value: 'light',
      label: t('themes.light'),
      icon: <Sun className="w-4 h-4" />,
    },
    {
      value: 'dark',
      label: t('themes.dark'),
      icon: <Moon className="w-4 h-4" />,
    },
    {
      value: 'system',
      label: 'System',
      icon: <Monitor className="w-4 h-4" />,
    },
  ];

  const currentTheme = themes.find(t => t.value === theme) || themes[0];

  return (
    <div className="relative group">
      <button
        className="flex items-center space-x-2 px-3 py-2 rounded-lg bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors duration-200 text-sm"
        title={t('themes.toggle')}
      >
        <span className="text-gray-600 dark:text-gray-400">
          {currentTheme.icon}
        </span>
        <span className="text-gray-700 dark:text-gray-300">
          {currentTheme.label}
        </span>
      </button>

      <div className="absolute top-full right-0 mt-2 w-36 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg shadow-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-10">
        <div className="py-1">
          {themes.map((themeOption) => (
            <button
              key={themeOption.value}
              onClick={() => setTheme(themeOption.value)}
              className={`w-full px-4 py-2 text-left hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-150 flex items-center space-x-3 text-sm ${
                theme === themeOption.value
                  ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400'
                  : 'text-gray-700 dark:text-gray-300'
              }`}
            >
              <span className="text-gray-500 dark:text-gray-400">
                {themeOption.icon}
              </span>
              <span>{themeOption.label}</span>
              {theme === themeOption.value && (
                <span className="ml-auto text-primary-500">✓</span>
              )}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ThemeToggle;