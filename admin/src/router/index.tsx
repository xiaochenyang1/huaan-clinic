import { Navigate, useRoutes } from 'react-router-dom'
import MainLayout from '@/layouts/MainLayout'
import Login from '@/pages/Login'
import Dashboard from '@/pages/Dashboard'
import DepartmentList from '@/pages/Department/List'
import DoctorList from '@/pages/Doctor/List'
import ScheduleList from '@/pages/Schedule/List'
import AppointmentList from '@/pages/Appointment/List'
import PatientList from '@/pages/Patient/List'
import Statistics from '@/pages/Statistics'
import Profile from '@/pages/Profile'
import NotFound from '@/pages/NotFound'
import Logs from '@/pages/Logs'
import AdminManagement from '@/pages/System/Admin'
import RoleManagement from '@/pages/System/Role'
import RequireAuth from './RequireAuth'

const AppRouter = () => {
  const routes = useRoutes([
    {
      path: '/login',
      element: <Login />,
    },
    {
      path: '/',
      element: (
        <RequireAuth>
          <MainLayout />
        </RequireAuth>
      ),
      children: [
        {
          index: true,
          element: <Navigate to="/dashboard" replace />,
        },
        {
          path: 'dashboard',
          element: <Dashboard />,
        },
        {
          path: 'department',
          element: <DepartmentList />,
        },
        {
          path: 'doctor',
          element: <DoctorList />,
        },
        {
          path: 'schedule',
          element: <ScheduleList />,
        },
        {
          path: 'appointment',
          element: <AppointmentList />,
        },
        {
          path: 'patient',
          element: <PatientList />,
        },
        {
          path: 'statistics',
          element: <Statistics />,
        },
        {
          path: 'system/admin',
          element: <AdminManagement />,
        },
        {
          path: 'system/role',
          element: <RoleManagement />,
        },
        {
          path: 'logs',
          element: <Logs />,
        },
        {
          path: 'profile',
          element: <Profile />,
        },
        {
          path: '*',
          element: <NotFound />,
        },
      ],
    },
    {
      path: '*',
      element: <NotFound />,
    },
  ])

  return routes
}

export default AppRouter
