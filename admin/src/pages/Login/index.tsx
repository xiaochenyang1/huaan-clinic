import { useState } from 'react'
import { Form, Input, Button, Card, message } from 'antd'
import { UserOutlined, LockOutlined, MedicineBoxOutlined } from '@ant-design/icons'
import { type Location, useLocation, useNavigate } from 'react-router-dom'
import http from '@/utils/http'

interface LoginForm {
  username: string
  password: string
}

const Login = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)

  const onFinish = async (values: LoginForm) => {
    setLoading(true)
    try {
      const response = await http.post('/admin/login', values)
      if (response.data.code === 200000) {
        localStorage.setItem('token', response.data.data.token)
        message.success('登录成功')
        const from = (location.state as { from?: Location } | null)?.from?.pathname || '/dashboard'
        navigate(from, { replace: true })
      } else {
        message.error(response.data.message || '登录失败')
      }
    } catch (error) {
      console.error('登录错误:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{
      width: '100vw',
      height: '100vh',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)',
    }}>
      <Card
        style={{
          width: 450,
          boxShadow: '0 8px 24px rgba(0,0,0,0.12)',
          borderRadius: 12,
          background: '#ffffff',
        }}
        bordered={false}
      >
        <div style={{
          textAlign: 'center',
          marginBottom: 32,
        }}>
          <div style={{
            display: 'inline-flex',
            alignItems: 'center',
            justifyContent: 'center',
            width: 64,
            height: 64,
            borderRadius: '50%',
            background: 'linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%)',
            marginBottom: 16,
          }}>
            <MedicineBoxOutlined style={{ fontSize: 32, color: '#fff' }} />
          </div>
          <h2 style={{
            fontSize: 24,
            fontWeight: 600,
            color: '#1f2937',
            margin: 0,
            marginBottom: 8,
          }}>
            华安医疗后台管理系统
          </h2>
          <p style={{
            fontSize: 14,
            color: '#6b7280',
            margin: 0,
          }}>
            Hospital Management System
          </p>
        </div>
        <Form
          form={form}
          name="login"
          onFinish={onFinish}
          autoComplete="off"
          size="large"
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input
              prefix={<UserOutlined style={{ color: '#9ca3af' }} />}
              placeholder="请输入用户名"
              style={{ height: 48 }}
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password
              prefix={<LockOutlined style={{ color: '#9ca3af' }} />}
              placeholder="请输入密码"
              style={{ height: 48 }}
            />
          </Form.Item>

          <Form.Item style={{ marginBottom: 0 }}>
            <Button
              type="primary"
              htmlType="submit"
              block
              loading={loading}
              style={{
                height: 48,
                fontSize: 16,
                fontWeight: 500,
                background: 'linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%)',
                border: 'none',
                boxShadow: '0 4px 12px rgba(59, 130, 246, 0.3)',
              }}
            >
              登录系统
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

export default Login
