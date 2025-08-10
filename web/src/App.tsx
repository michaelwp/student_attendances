import React, { useEffect, useState } from 'react';
import { useThemeStore } from './stores/themeStore';
import { useAuthStore } from './stores/authStore';
import { apiService } from './services/api';
import LoginPage from './components/LoginPage';
import Dashboard from './pages/Dashboard';
import TeachersPage from './pages/TeachersPage';
import StudentsPage from './pages/StudentsPage';
import ClassesPage from './pages/ClassesPage';
import AttendancePage from './pages/AttendancePage';
import AdminsPage from './pages/AdminsPage';
import StudentHomepage from './pages/StudentHomepage';
import StudentPage from './pages/StudentPage';
import Layout from './components/Layout';
import { ToastContainer } from './components/Toast';
import './i18n';

const App: React.FC = () => {
  const { initializeTheme } = useThemeStore();
  const { initializeAuth, isAuthenticated, user } = useAuthStore();
  const [currentPath, setCurrentPath] = useState(window.location.hash || '#/');

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

  useEffect(() => {
    // Listen to hash changes for routing
    const handleHashChange = () => {
      setCurrentPath(window.location.hash || '#/');
    };

    window.addEventListener('hashchange', handleHashChange);
    return () => window.removeEventListener('hashchange', handleHashChange);
  }, []);

  // Check if we're on the student homepage (public route)
  const path = currentPath.replace('#', '');
  const isStudentHomepage = path === '/' || path === '/home' || path === '';
  
  // Show student homepage without authentication
  if (isStudentHomepage) {
    return (
      <>
        <StudentHomepage />
        <ToastContainer />
      </>
    );
  }
  
  // Show the login page if not authenticated for admin routes
  if (!isAuthenticated) {
    return <LoginPage />;
  }

  const MainContent = () => {
    // Show student-specific page for student users
    if (user?.type === 'student') {
      return <StudentPage />;
    }
    
    switch (path) {
      case '/teachers':
        return <TeachersPage />;
      case '/students':
        return <StudentsPage />;
      case '/classes':
        return <ClassesPage />;
      case '/attendance':
        return <AttendancePage />;
      case '/admins':
        return <AdminsPage />;
      case '/dashboard':
      default:
        return <Dashboard />;
    }
  };

  return (
    <>
      <Layout>
        <MainContent />
      </Layout>
      <ToastContainer />
    </>
  );
};

export default App;
