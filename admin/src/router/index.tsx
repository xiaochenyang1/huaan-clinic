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
import RequirePermission from './RequirePermission'

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
          element: (
            <RequirePermission any={['statistics:view']}>
              <Dashboard />
            </RequirePermission>
          ),
        },
        {
          path: 'department',
          element: (
            <RequirePermission any={['department:view']}>
              <DepartmentList />
            </RequirePermission>
          ),
        },
        {
          path: 'doctor',
          element: (
            <RequirePermission any={['doctor:view']}>
              <DoctorList />
            </RequirePermission>
          ),
        },
        {
          path: 'schedule',
          element: (
            <RequirePermission any={['schedule:view']}>
              <ScheduleList />
            </RequirePermission>
          ),
        },
        {
          path: 'appointment',
          element: (
            <RequirePermission any={['appointment:view']}>
              <AppointmentList />
            </RequirePermission>
          ),
        },
        {
          path: 'patient',
          element: (
            <RequirePermission any={['patient:view']}>
              <PatientList />
            </RequirePermission>
          ),
        },
        {
          path: 'statistics',
          element: (
            <RequirePermission any={['statistics:view']}>
              <Statistics />
            </RequirePermission>
          ),
        },
        {
          path: 'system/admin',
          element: (
            <RequirePermission any={['admin:view']}>
              <AdminManagement />
            </RequirePermission>
          ),
        },
        {
          path: 'system/role',
          element: (
            <RequirePermission any={['role:view']}>
              <RoleManagement />
            </RequirePermission>
          ),
        },
        {
          path: 'logs',
          element: (
            <RequirePermission any={['log:view']}>
              <Logs />
            </RequirePermission>
          ),
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
