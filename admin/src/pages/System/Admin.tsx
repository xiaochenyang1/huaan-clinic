import { useEffect, useMemo, useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, Select, Tag, message } from 'antd'
import { PlusOutlined, EditOutlined, KeyOutlined } from '@ant-design/icons'
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table'
import http from '@/utils/http'

interface Role {
  id: number
  code: string
  name: string
  status: number
}

interface Admin {
  id: number
  username: string
  nickname: string
  phone: string
  email: string
  status: number
  status_name: string
  roles?: string[]
  role_names?: string[]
  last_login_at?: string
  created_at: string
}

const AdminManagement = () => {
  const [data, setData] = useState<Admin[]>([])
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(false)
  const [keyword, setKeyword] = useState('')
  const [status, setStatus] = useState<number | undefined>()
  const [modalVisible, setModalVisible] = useState(false)
  const [resetVisible, setResetVisible] = useState(false)
  const [editing, setEditing] = useState<Admin | null>(null)
  const [resetting, setResetting] = useState<Admin | null>(null)
  const [form] = Form.useForm()
  const [resetForm] = Form.useForm()
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })

  const codeToRoleId = useMemo(() => {
    const map: Record<string, number> = {}
    roles.forEach((r) => { map[r.code] = r.id })
    return map
  }, [roles])

  const fetchRoles = async () => {
    try {
      const response = await http.get('/admin/roles', {
        params: {
          page: 1,
          page_size: 100,
        },
      })
      if (response.data.code === 200000) {
        setRoles(response.data.data.list || [])
      }
    } catch (error) {
      console.error('获取角色列表失败:', error)
    }
  }

  const fetchData = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const response = await http.get('/admin/admins', {
        params: {
          page,
          page_size: pageSize,
          keyword: keyword || undefined,
          status: typeof status === 'number' ? status : undefined,
        },
      })
      if (response.data.code === 200000) {
        setData(response.data.data.list || [])
        setPagination({
          current: page,
          pageSize,
          total: response.data.data.total || 0,
        })
      }
    } catch (error) {
      console.error('获取管理员列表失败:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchRoles()
    fetchData()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const columns: ColumnsType<Admin> = [
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
      width: 160,
    },
    {
      title: '昵称',
      dataIndex: 'nickname',
      key: 'nickname',
      width: 160,
      render: (v: string) => v || '-',
    },
    {
      title: '角色',
      key: 'roles',
      render: (_, record) => (
        <Space size={[4, 4]} wrap>
          {(record.role_names || []).map((name) => (
            <Tag key={name}>{name}</Tag>
          ))}
          {(record.role_names || []).length === 0 ? <span style={{ color: '#6b7280' }}>-</span> : null}
        </Space>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 110,
      render: (_: number, record) => (
        <Tag color={record.status === 1 ? 'success' : 'default'}>{record.status_name}</Tag>
      ),
    },
    {
      title: '手机号',
      dataIndex: 'phone',
      key: 'phone',
      width: 140,
      render: (v: string) => v || '-',
    },
    {
      title: '邮箱',
      dataIndex: 'email',
      key: 'email',
      width: 200,
      ellipsis: true,
      render: (v: string) => v || '-',
    },
    {
      title: '最后登录',
      dataIndex: 'last_login_at',
      key: 'last_login_at',
      width: 180,
      render: (v?: string) => v || '-',
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
    },
    {
      title: '操作',
      key: 'action',
      width: 260,
      render: (_, record) => (
        <Space size="middle">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => openEdit(record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            size="small"
            icon={<KeyOutlined />}
            onClick={() => openResetPassword(record)}
          >
            重置密码
          </Button>
          <Button
            type="link"
            size="small"
            danger={record.status === 1}
            onClick={() => toggleStatus(record)}
          >
            {record.status === 1 ? '禁用' : '启用'}
          </Button>
        </Space>
      ),
    },
  ]

  const openCreate = () => {
    setEditing(null)
    form.resetFields()
    form.setFieldsValue({ status: 1 })
    setModalVisible(true)
  }

  const openEdit = (record: Admin) => {
    setEditing(record)
    form.resetFields()
    form.setFieldsValue({
      username: record.username,
      nickname: record.nickname,
      phone: record.phone,
      email: record.email,
      status: record.status,
      role_ids: (record.roles || []).map((code) => codeToRoleId[code]).filter(Boolean),
    })
    setModalVisible(true)
  }

  const openResetPassword = (record: Admin) => {
    setResetting(record)
    resetForm.resetFields()
    setResetVisible(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editing) {
        const response = await http.put(`/admin/admins/${editing.id}`, {
          nickname: values.nickname,
          phone: values.phone,
          email: values.email,
          status: values.status,
          role_ids: values.role_ids,
        })
        if (response.data.code === 200000) {
          message.success('更新成功')
        }
      } else {
        const response = await http.post('/admin/admins', {
          username: values.username,
          password: values.password,
          nickname: values.nickname,
          phone: values.phone,
          email: values.email,
          status: values.status,
          role_ids: values.role_ids,
        })
        if (response.data.code === 200000) {
          message.success('创建成功')
        }
      }
      setModalVisible(false)
      fetchData(pagination.current, pagination.pageSize)
    } catch (error) {
      console.error('操作失败:', error)
    }
  }

  const handleResetPassword = async () => {
    try {
      const values = await resetForm.validateFields()
      if (!resetting) return
      const response = await http.put(`/admin/admins/${resetting.id}/password`, {
        password: values.password,
      })
      if (response.data.code === 200000) {
        message.success('密码重置成功')
        setResetVisible(false)
      }
    } catch (error) {
      console.error('重置密码失败:', error)
    }
  }

  const toggleStatus = async (record: Admin) => {
    try {
      const nextStatus = record.status === 1 ? 0 : 1
      const response = await http.put(`/admin/admins/${record.id}`, { status: nextStatus })
      if (response.data.code === 200000) {
        message.success('操作成功')
        fetchData(pagination.current, pagination.pageSize)
      }
    } catch (error) {
      console.error('更新状态失败:', error)
    }
  }

  const handleTableChange = (newPagination: TablePaginationConfig) => {
    fetchData(newPagination.current || 1, newPagination.pageSize || 10)
  }

  const handleSearch = () => {
    fetchData(1, pagination.pageSize)
  }

  const handleReset = () => {
    setKeyword('')
    setStatus(undefined)
    fetchData(1, pagination.pageSize)
  }

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>管理员管理</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建管理员
        </Button>
      </div>

      <div style={{ marginBottom: 12 }}>
        <Space wrap>
          <Input
            placeholder="关键词（用户名/昵称/手机号/邮箱）"
            allowClear
            style={{ width: 260 }}
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            onPressEnter={handleSearch}
          />
          <Select
            placeholder="状态"
            allowClear
            style={{ width: 140 }}
            value={typeof status === 'number' ? status : undefined}
            onChange={(v) => setStatus(v)}
            options={[
              { label: '启用', value: 1 },
              { label: '禁用', value: 0 },
            ]}
          />
          <Button type="primary" onClick={handleSearch}>查询</Button>
          <Button onClick={handleReset}>重置</Button>
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        pagination={{
          ...pagination,
          showTotal: (total) => `共 ${total} 条`,
          showSizeChanger: true,
          showQuickJumper: true,
        }}
        onChange={handleTableChange}
        scroll={{ x: 1200 }}
      />

      <Modal
        title={editing ? '编辑管理员' : '新建管理员'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={520}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="username"
            label="用户名"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input placeholder="4-20位字母/数字/下划线，字母开头" disabled={!!editing} />
          </Form.Item>

          {!editing && (
            <Form.Item
              name="password"
              label="密码"
              rules={[{ required: true, message: '请输入密码' }, { min: 6, message: '至少6位' }]}
            >
              <Input.Password placeholder="请输入密码" />
            </Form.Item>
          )}

          <Form.Item name="nickname" label="昵称">
            <Input placeholder="请输入昵称" />
          </Form.Item>
          <Form.Item name="phone" label="手机号">
            <Input placeholder="请输入手机号" />
          </Form.Item>
          <Form.Item name="email" label="邮箱">
            <Input placeholder="请输入邮箱" />
          </Form.Item>

          <Form.Item name="role_ids" label="角色">
            <Select
              mode="multiple"
              placeholder="请选择角色"
              options={roles.map((r) => ({ label: r.name, value: r.id, disabled: r.status !== 1 }))}
            />
          </Form.Item>

          <Form.Item
            name="status"
            label="状态"
            rules={[{ required: true, message: '请选择状态' }]}
          >
            <Select
              options={[
                { label: '启用', value: 1 },
                { label: '禁用', value: 0 },
              ]}
            />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={`重置密码${resetting ? ` - ${resetting.username}` : ''}`}
        open={resetVisible}
        onOk={handleResetPassword}
        onCancel={() => setResetVisible(false)}
        width={480}
      >
        <Form form={resetForm} layout="vertical">
          <Form.Item
            name="password"
            label="新密码"
            rules={[{ required: true, message: '请输入新密码' }, { min: 6, message: '至少6位' }]}
          >
            <Input.Password placeholder="请输入新密码" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default AdminManagement

