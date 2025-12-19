import { useEffect, useState } from 'react'
import { Table, Button, Space, Modal, Form, Input, Select, InputNumber, Tag, message } from 'antd'
import { PlusOutlined, EditOutlined } from '@ant-design/icons'
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table'
import http from '@/utils/http'

interface Role {
  id: number
  code: string
  name: string
  description: string
  sort_order: number
  status: number
  status_name: string
  permissions?: string[]
}

const RoleManagement = () => {
  const [data, setData] = useState<Role[]>([])
  const [loading, setLoading] = useState(false)
  const [keyword, setKeyword] = useState('')
  const [status, setStatus] = useState<number | undefined>()
  const [modalVisible, setModalVisible] = useState(false)
  const [editing, setEditing] = useState<Role | null>(null)
  const [form] = Form.useForm()
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })

  const fetchData = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const response = await http.get('/admin/roles', {
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
      console.error('获取角色列表失败:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const columns: ColumnsType<Role> = [
    {
      title: '编码',
      dataIndex: 'code',
      key: 'code',
      width: 180,
    },
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      width: 160,
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
      render: (v: string) => v || '-',
    },
    {
      title: '排序',
      dataIndex: 'sort_order',
      key: 'sort_order',
      width: 90,
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
      title: '权限数',
      key: 'permissions',
      width: 90,
      render: (_, record) => (record.permissions ? record.permissions.length : 0),
    },
    {
      title: '操作',
      key: 'action',
      width: 140,
      render: (_, record) => (
        <Button
          type="link"
          icon={<EditOutlined />}
          onClick={() => openEdit(record)}
        >
          编辑
        </Button>
      ),
    },
  ]

  const openCreate = () => {
    setEditing(null)
    form.resetFields()
    form.setFieldsValue({ status: 1, sort_order: 0 })
    setModalVisible(true)
  }

  const openEdit = (record: Role) => {
    setEditing(record)
    form.resetFields()
    form.setFieldsValue({
      code: record.code,
      name: record.name,
      description: record.description,
      sort_order: record.sort_order,
      status: record.status,
    })
    setModalVisible(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editing) {
        const response = await http.put(`/admin/roles/${editing.id}`, values)
        if (response.data.code === 200000) {
          message.success('更新成功')
        }
      } else {
        const response = await http.post('/admin/roles', values)
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
        <h1 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>角色管理</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          新建角色
        </Button>
      </div>

      <div style={{ marginBottom: 12 }}>
        <Space wrap>
          <Input
            placeholder="关键词（编码/名称）"
            allowClear
            style={{ width: 240 }}
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
      />

      <Modal
        title={editing ? '编辑角色' : '新建角色'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={520}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="code"
            label="角色编码"
            rules={[{ required: true, message: '请输入角色编码' }]}
          >
            <Input placeholder="例如 super_admin" />
          </Form.Item>
          <Form.Item
            name="name"
            label="角色名称"
            rules={[{ required: true, message: '请输入角色名称' }]}
          >
            <Input placeholder="例如 超级管理员" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} placeholder="请输入描述" />
          </Form.Item>
          <Form.Item name="sort_order" label="排序">
            <InputNumber style={{ width: '100%' }} />
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
    </div>
  )
}

export default RoleManagement

