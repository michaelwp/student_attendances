import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';

export type Theme = 'light' | 'dark' | 'system';

interface ThemeState {
  theme: Theme;
  isDark: boolean;
  setTheme: (theme: Theme) => void;
  toggleTheme: () => void;
  initializeTheme: () => void;
}

const getSystemTheme = (): boolean => {
  return window.matchMedia('(prefers-color-scheme: dark)').matches;
};

const applyTheme = (isDark: boolean) => {
  const root = document.documentElement;
  if (isDark) {
    root.classList.add('dark');
  } else {
    root.classList.remove('dark');
  }
};

export const useThemeStore = create<ThemeState>()(
  persist(
    (set, get) => ({
      theme: 'system',
      isDark: false,

      initializeTheme: () => {
        const { theme } = get();
        let isDark = false;

        if (theme === 'system') {
          isDark = getSystemTheme();
        } else {
          isDark = theme === 'dark';
        }

        set({ isDark });
        applyTheme(isDark);

        // Listen for system theme changes
        const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
        const handleChange = (e: MediaQueryListEvent) => {
          const { theme: currentTheme } = get();
          if (currentTheme === 'system') {
            set({ isDark: e.matches });
            applyTheme(e.matches);
          }
        };

        mediaQuery.addEventListener('change', handleChange);
        
        // Return cleanup function
        return () => mediaQuery.removeEventListener('change', handleChange);
      },

      setTheme: (theme: Theme) => {
        let isDark = false;

        if (theme === 'system') {
          isDark = getSystemTheme();
        } else {
          isDark = theme === 'dark';
        }

        set({ theme, isDark });
        applyTheme(isDark);
      },

      toggleTheme: () => {
        const { theme } = get();
        const newTheme = theme === 'dark' ? 'light' : 'dark';
        get().setTheme(newTheme);
      },
    }),
    {
      name: 'theme-storage',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({ theme: state.theme }),
      onRehydrateStorage: () => (state) => {
        if (state) {
          state.initializeTheme();
        }
      },
    }
  )
);