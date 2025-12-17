import { useState, useEffect } from 'react'
import { Table, Button, Space, Modal, Form, Input, message } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table'
import http from '@/utils/http'

interface Department {
  id: number
  name: string
  description: string
  icon: string
  sort_order: number
  status: number
  status_name: string
}

const DepartmentList = () => {
  const [data, setData] = useState<Department[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })

  const columns: ColumnsType<Department> = [
    {
      title: '科室名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '科室描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '排序',
      dataIndex: 'sort_order',
      key: 'sort_order',
      width: 100,
    },
    {
      title: '状态',
      dataIndex: 'status_name',
      key: 'status_name',
      width: 100,
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_, record) => (
        <Space size="middle">
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.id)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ]

  const fetchData = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const response = await http.get('/admin/departments', {
        params: {
          page,
          page_size: pageSize,
        },
      })
      if (response.data.code === 200000) {
        setData(response.data.data.list)
        setPagination({
          current: page,
          pageSize: pageSize,
          total: response.data.data.total,
        })
      }
    } catch (error) {
      console.error('获取科室列表失败:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  const handleTableChange = (newPagination: TablePaginationConfig) => {
    fetchData(newPagination.current || 1, newPagination.pageSize || 10)
  }

  const handleAdd = () => {
    setEditingId(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (record: Department) => {
    setEditingId(record.id)
    form.setFieldsValue(record)
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个科室吗？',
      onOk: async () => {
        try {
          const response = await http.delete(`/admin/departments/${id}`)
          if (response.data.code === 200000) {
            message.success('删除成功')
            fetchData(pagination.current, pagination.pageSize)
          }
        } catch (error) {
          console.error('删除失败:', error)
        }
      },
    })
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingId) {
        const response = await http.put(`/admin/departments/${editingId}`, values)
        if (response.data.code === 200000) {
          message.success('更新成功')
        }
      } else {
        const response = await http.post('/admin/departments', values)
        if (response.data.code === 200000) {
          message.success('添加成功')
        }
      }
      setModalVisible(false)
      fetchData(pagination.current, pagination.pageSize)
    } catch (error) {
      console.error('操作失败:', error)
    }
  }

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>科室管理</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          添加科室
        </Button>
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
        title={editingId ? '编辑科室' : '添加科室'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="科室名称"
            rules={[{ required: true, message: '请输入科室名称' }]}
          >
            <Input placeholder="请输入科室名称" />
          </Form.Item>
          <Form.Item
            name="description"
            label="科室描述"
            rules={[{ required: true, message: '请输入科室描述' }]}
          >
            <Input.TextArea rows={4} placeholder="请输入科室描述" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default DepartmentList
