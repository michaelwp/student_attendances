import React from 'react';
import Dashboard from '../pages/Dashboard';
import TeachersPage from '../pages/TeachersPage';
import StudentsPage from '../pages/StudentsPage';
import ClassesPage from '../pages/ClassesPage';
import AttendancePage from '../pages/AttendancePage';
import AdminsPage from '../pages/AdminsPage';
import StudentHomepage from '../pages/StudentHomepage';

interface RouterProps {
  currentPath: string;
}

export const Router: React.FC<RouterProps> = ({ currentPath }) => {
  // Remove hash from path
  const path = currentPath.replace('#', '');

  switch (path) {
    case '/':
    case '/home':
      return <StudentHomepage />;
    case '/dashboard':
      return <Dashboard />;
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
    default:
      return <StudentHomepage />;
  }
};

export default Router;