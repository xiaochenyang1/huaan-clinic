import { useState, useEffect } from 'react'
import { Card, Row, Col, Statistic, Spin } from 'antd'
import {
  UserOutlined,
  TeamOutlined,
  CalendarOutlined,
  MedicineBoxOutlined,
  CheckCircleOutlined,
} from '@ant-design/icons'
import http from '@/utils/http'

interface DashboardData {
  today_appointments: number
  today_checkins: number
  total_appointments: number
  total_users: number
  total_doctors: number
  total_departments: number
}

const Dashboard = () => {
  const [data, setData] = useState<DashboardData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchDashboardData()
  }, [])

  const fetchDashboardData = async () => {
    try {
      const response = await http.get('/admin/dashboard')
      if (response.data.code === 200000) {
        setData(response.data.data)
      }
    } catch (error) {
      console.error('获取仪表盘数据失败:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '100px 0' }}>
        <Spin size="large" />
      </div>
    )
  }

  return (
    <div>
      <h1 style={{ margin: 0, marginBottom: 24, fontSize: 20, fontWeight: 600 }}>数据概览</h1>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable style={{ borderRadius: 8 }}>
            <Statistic
              title="今日预约"
              value={data?.today_appointments || 0}
              prefix={<CalendarOutlined style={{ color: '#52c41a' }} />}
              valueStyle={{ color: '#52c41a', fontWeight: 600 }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable style={{ borderRadius: 8 }}>
            <Statistic
              title="今日签到"
              value={data?.today_checkins || 0}
              prefix={<CheckCircleOutlined style={{ color: '#1890ff' }} />}
              valueStyle={{ color: '#1890ff', fontWeight: 600 }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable style={{ borderRadius: 8 }}>
            <Statistic
              title="医生总数"
              value={data?.total_doctors || 0}
              prefix={<TeamOutlined style={{ color: '#722ed1' }} />}
              valueStyle={{ color: '#722ed1', fontWeight: 600 }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable style={{ borderRadius: 8 }}>
            <Statistic
              title="科室总数"
              value={data?.total_departments || 0}
              prefix={<MedicineBoxOutlined style={{ color: '#fa8c16' }} />}
              valueStyle={{ color: '#fa8c16', fontWeight: 600 }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} sm={12} lg={8}>
          <Card hoverable style={{ borderRadius: 8 }}>
            <Statistic
              title="总预约数"
              value={data?.total_appointments || 0}
              prefix={<CalendarOutlined style={{ color: '#13c2c2' }} />}
              valueStyle={{ color: '#13c2c2', fontWeight: 600 }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={8}>
          <Card hoverable style={{ borderRadius: 8 }}>
            <Statistic
              title="患者总数"
              value={data?.total_users || 0}
              prefix={<UserOutlined style={{ color: '#eb2f96' }} />}
              valueStyle={{ color: '#eb2f96', fontWeight: 600 }}
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Dashboard
