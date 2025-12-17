import { useState, useEffect } from 'react'
import { Table, Button, Space, Modal, Form, Input, Select, message } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table'
import http from '@/utils/http'
import AvatarUpload from '@/components/Upload/AvatarUpload'

interface Doctor {
  id: number
  name: string
  department_id: number
  department_name: string
  title: string
  title_name: string
  phone: string
  specialty: string
  avatar?: string
}

interface Department {
  id: number
  name: string
}

const DoctorList = () => {
  const [data, setData] = useState<Doctor[]>([])
  const [departments, setDepartments] = useState<Department[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })

  const columns: ColumnsType<Doctor> = [
    {
      title: '医生姓名',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '所属科室',
      dataIndex: 'department_name',
      key: 'department_name',
    },
    {
      title: '职称',
      dataIndex: 'title_name',
      key: 'title_name',
    },
    {
      title: '擅长',
      dataIndex: 'specialty',
      key: 'specialty',
      ellipsis: true,
    },
    {
      title: '联系电话',
      dataIndex: 'phone',
      key: 'phone',
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
      const response = await http.get('/admin/doctors', {
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
      console.error('获取医生列表失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchDepartments = async () => {
    try {
      const response = await http.get('/admin/departments', {
        params: {
          page: 1,
          page_size: 100,
        },
      })
      if (response.data.code === 200000) {
        setDepartments(response.data.data.list)
      }
    } catch (error) {
      console.error('获取科室列表失败:', error)
    }
  }

  useEffect(() => {
    fetchData()
    fetchDepartments()
  }, [])

  const handleTableChange = (newPagination: TablePaginationConfig) => {
    fetchData(newPagination.current || 1, newPagination.pageSize || 10)
  }

  const handleAdd = () => {
    setEditingId(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (record: Doctor) => {
    setEditingId(record.id)
    form.setFieldsValue({
      avatar: record.avatar,
      name: record.name,
      department_id: record.department_id,
      title: record.title,
      specialty: record.specialty,
      phone: record.phone,
    })
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个医生吗？',
      onOk: async () => {
        try {
          const response = await http.delete(`/admin/doctors/${id}`)
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
        const response = await http.put(`/admin/doctors/${editingId}`, values)
        if (response.data.code === 200000) {
          message.success('更新成功')
        }
      } else {
        const response = await http.post('/admin/doctors', values)
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
        <h1 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>医生管理</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          添加医生
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
        title={editingId ? '编辑医生' : '添加医生'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="avatar"
            label="医生头像"
          >
            <AvatarUpload />
          </Form.Item>
          <Form.Item
            name="name"
            label="医生姓名"
            rules={[{ required: true, message: '请输入医生姓名' }]}
          >
            <Input placeholder="请输入医生姓名" />
          </Form.Item>
          <Form.Item
            name="department_id"
            label="所属科室"
            rules={[{ required: true, message: '请选择所属科室' }]}
          >
            <Select
              placeholder="请选择科室"
              options={departments.map(dept => ({
                label: dept.name,
                value: dept.id,
              }))}
            />
          </Form.Item>
          <Form.Item
            name="title"
            label="职称"
            rules={[{ required: true, message: '请选择职称' }]}
          >
            <Select placeholder="请选择职称">
              <Select.Option value="chief_physician">主任医师</Select.Option>
              <Select.Option value="associate_chief_physician">副主任医师</Select.Option>
              <Select.Option value="attending_physician">主治医师</Select.Option>
              <Select.Option value="resident_physician">住院医师</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="specialty"
            label="擅长"
            rules={[{ required: true, message: '请输入擅长领域' }]}
          >
            <Input.TextArea rows={3} placeholder="请输入擅长领域" />
          </Form.Item>
          <Form.Item
            name="phone"
            label="联系电话"
            rules={[
              { required: true, message: '请输入联系电话' },
              { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号' },
            ]}
          >
            <Input placeholder="请输入联系电话" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default DoctorList
