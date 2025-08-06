import { useThemeStore } from '../stores/themeStore';

export const useDarkMode = () => {
  const { isDark, toggleTheme } = useThemeStore();
  
  // Rename toggleTheme to toggleDarkMode to match the expected interface in Layout.tsx
  return {
    isDark,
    toggleDarkMode: toggleTheme
  };
};