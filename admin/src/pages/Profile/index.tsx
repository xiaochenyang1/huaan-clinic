import { useState, useEffect } from 'react'
import { Card, Descriptions, Button, Modal, Form, Input, message } from 'antd'
import { LockOutlined } from '@ant-design/icons'
import http from '@/utils/http'

interface AdminInfo {
  id: number
  username: string
  nickname: string
  phone: string
  email: string
  status: number
  status_name: string
  created_at: string
  last_login_at: string
  last_login_ip: string
}

const Profile = () => {
  const [adminInfo, setAdminInfo] = useState<AdminInfo | null>(null)
  const [loading, setLoading] = useState(false)
  const [passwordModalVisible, setPasswordModalVisible] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    fetchAdminInfo()
  }, [])

  const fetchAdminInfo = async () => {
    setLoading(true)
    try {
      const response = await http.get('/admin/info')
      if (response.data.code === 200000) {
        setAdminInfo(response.data.data)
      }
    } catch (error) {
      console.error('获取管理员信息失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleChangePassword = () => {
    form.resetFields()
    setPasswordModalVisible(true)
  }

  const handlePasswordSubmit = async () => {
    try {
      const values = await form.validateFields()
      const response = await http.put('/admin/password', {
        old_password: values.old_password,
        new_password: values.new_password,
      })
      if (response.data.code === 200000) {
        message.success('密码修改成功，请重新登录')
        setPasswordModalVisible(false)
        // 3秒后自动跳转到登录页
        setTimeout(() => {
          localStorage.removeItem('token')
          window.location.href = '/login'
        }, 3000)
      }
    } catch (error) {
      console.error('修改密码失败:', error)
    }
  }

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>个人中心</h1>
        <Button
          type="primary"
          icon={<LockOutlined />}
          onClick={handleChangePassword}
        >
          修改密码
        </Button>
      </div>

      <Card
        title="基本信息"
        loading={loading}
        bordered={false}
        style={{ borderRadius: 8 }}
      >
        {adminInfo && (
          <Descriptions column={2} bordered>
            <Descriptions.Item label="用户名">{adminInfo.username}</Descriptions.Item>
            <Descriptions.Item label="昵称">{adminInfo.nickname}</Descriptions.Item>
            <Descriptions.Item label="手机号">{adminInfo.phone || '未填写'}</Descriptions.Item>
            <Descriptions.Item label="邮箱">{adminInfo.email || '未填写'}</Descriptions.Item>
            <Descriptions.Item label="状态">
              <span style={{ color: adminInfo.status === 1 ? '#52c41a' : '#ff4d4f' }}>
                {adminInfo.status_name}
              </span>
            </Descriptions.Item>
            <Descriptions.Item label="注册时间">{adminInfo.created_at}</Descriptions.Item>
            <Descriptions.Item label="最后登录时间">
              {adminInfo.last_login_at || '暂无记录'}
            </Descriptions.Item>
            <Descriptions.Item label="最后登录IP">
              {adminInfo.last_login_ip || '暂无记录'}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Card>

      <Modal
        title="修改密码"
        open={passwordModalVisible}
        onOk={handlePasswordSubmit}
        onCancel={() => setPasswordModalVisible(false)}
        width={500}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="old_password"
            label="当前密码"
            rules={[
              { required: true, message: '请输入当前密码' },
              { min: 6, message: '密码至少6位' },
            ]}
          >
            <Input.Password placeholder="请输入当前密码" />
          </Form.Item>
          <Form.Item
            name="new_password"
            label="新密码"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 6, message: '密码至少6位' },
              { max: 20, message: '密码最多20位' },
            ]}
          >
            <Input.Password placeholder="请输入新密码（6-20位）" />
          </Form.Item>
          <Form.Item
            name="confirm_password"
            label="确认密码"
            dependencies={['new_password']}
            rules={[
              { required: true, message: '请确认新密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('new_password') === value) {
                    return Promise.resolve()
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'))
                },
              }),
            ]}
          >
            <Input.Password placeholder="请再次输入新密码" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default Profile
