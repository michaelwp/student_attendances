import React, { useEffect } from 'react';
import { useThemeStore } from './stores/themeStore';
import { useAuthStore } from './stores/authStore';
import { apiService } from './services/api';
import LoginPage from './components/LoginPage';
import Dashboard from './pages/Dashboard';
import Layout from './components/Layout';
import './i18n';

const App: React.FC = () => {
  const { initializeTheme } = useThemeStore();
  const { initializeAuth, isAuthenticated, user } = useAuthStore();

  useEffect(() => {
    // Initialize theme on an app load
    initializeTheme();
    
    // Initialize auth state from localStorage
    initializeAuth();
  }, [initializeTheme, initializeAuth]);

  useEffect(() => {
    // Set the token in the API service if it exists
    if (user?.token) {
      apiService.setAuthToken(user.token);
    }
  }, [user?.token]);

  // Show the login page if not authenticated, otherwise show a dashboard
  if (!isAuthenticated) {
    return <LoginPage />;
  }

  return (
    <Layout>
      <Dashboard />
    </Layout>
  );
};

export default App;
