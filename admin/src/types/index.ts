export interface Department {
  id: string
  name: string
  description: string
  createdAt: string
  updatedAt: string
}

export interface Doctor {
  id: string
  name: string
  departmentId: string
  departmentName?: string
  title: string
  phone: string
  avatar?: string
  introduction?: string
  createdAt: string
  updatedAt: string
}

export interface Patient {
  id: string
  name: string
  phone: string
  idCard: string
  gender: string
  age: number
  address?: string
  createdAt: string
  updatedAt: string
}

export interface Appointment {
  id: string
  patientId: string
  patientName?: string
  doctorId: string
  doctorName?: string
  departmentId: string
  departmentName?: string
  appointmentDate: string
  appointmentTime: string
  status: 'pending' | 'confirmed' | 'cancelled' | 'completed'
  remark?: string
  createdAt: string
  updatedAt: string
}

export interface User {
  id: string
  username: string
  name: string
  role: string
  createdAt: string
}
