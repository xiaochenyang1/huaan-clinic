import { useState, useEffect } from 'react'
import { Table, Button, Space, Modal, Form, Input, Select, DatePicker, message, InputNumber, Row, Col, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, AppstoreAddOutlined } from '@ant-design/icons'
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table'
import dayjs, { Dayjs } from 'dayjs'
import http from '@/utils/http'

interface Schedule {
  id: number
  doctor_id: number
  doctor_name: string
  department_name: string
  schedule_date: string
  period: string
  period_name: string
  start_time: string
  end_time: string
  total_slots: number
  booked_slots: number
  available_slots: number
  status: number
  status_name: string
}

interface Doctor {
  id: number
  name: string
  department_id: number
  department_name: string
}

interface Department {
  id: number
  name: string
}

const periodMap: Record<string, { color: string }> = {
  morning: { color: 'blue' },
  afternoon: { color: 'orange' },
}

const statusMap: Record<number, { color: string }> = {
  0: { color: 'red' },
  1: { color: 'green' },
}

const ScheduleList = () => {
  const [data, setData] = useState<Schedule[]>([])
  const [doctors, setDoctors] = useState<Doctor[]>([])
  const [departments, setDepartments] = useState<Department[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [batchModalVisible, setBatchModalVisible] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const [batchForm] = Form.useForm()
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })
  const [filters, setFilters] = useState<{
    doctor_id?: number
    department_id?: number
    start_date?: string
    end_date?: string
    status?: number
  }>({})

  const columns: ColumnsType<Schedule> = [
    {
      title: '医生姓名',
      dataIndex: 'doctor_name',
      key: 'doctor_name',
      width: 120,
    },
    {
      title: '科室',
      dataIndex: 'department_name',
      key: 'department_name',
      width: 120,
    },
    {
      title: '排班日期',
      dataIndex: 'schedule_date',
      key: 'schedule_date',
      width: 120,
    },
    {
      title: '时段',
      dataIndex: 'period',
      key: 'period',
      width: 100,
      render: (period: string, record) => (
        <Tag color={periodMap[period]?.color}>{record.period_name}</Tag>
      ),
    },
    {
      title: '时间',
      key: 'time',
      width: 130,
      render: (_, record) => `${record.start_time} - ${record.end_time}`,
    },
    {
      title: '总号源',
      dataIndex: 'total_slots',
      key: 'total_slots',
      width: 80,
      align: 'center',
    },
    {
      title: '已约',
      dataIndex: 'booked_slots',
      key: 'booked_slots',
      width: 70,
      align: 'center',
    },
    {
      title: '剩余',
      dataIndex: 'available_slots',
      key: 'available_slots',
      width: 70,
      align: 'center',
      render: (slots: number) => (
        <span style={{ color: slots === 0 ? '#ff4d4f' : slots < 5 ? '#faad14' : '#52c41a' }}>
          {slots}
        </span>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 80,
      render: (status: number, record) => (
        <Tag color={statusMap[status]?.color}>{record.status_name}</Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            size="small"
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record.id)}
            disabled={record.booked_slots > 0}
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
      const params: any = {
        page,
        page_size: pageSize,
        ...filters,
      }
      const response = await http.get('/admin/schedules', { params })
      if (response.data.code === 200000) {
        setData(response.data.data.list)
        setPagination({
          current: page,
          pageSize: pageSize,
          total: response.data.data.total,
        })
      }
    } catch (error) {
      console.error('获取排班列表失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchDoctors = async () => {
    try {
      const response = await http.get('/admin/doctors', {
        params: {
          page: 1,
          page_size: 100,
        },
      })
      if (response.data.code === 200000) {
        setDoctors(response.data.data.list)
      }
    } catch (error) {
      console.error('获取医生列表失败:', error)
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
    fetchDoctors()
    fetchDepartments()
  }, [])

  useEffect(() => {
    fetchData(1, pagination.pageSize)
  }, [filters])

  const handleTableChange = (newPagination: TablePaginationConfig) => {
    fetchData(newPagination.current || 1, newPagination.pageSize || 10)
  }

  const handleFilterChange = (key: string, value: any) => {
    setFilters(prev => ({
      ...prev,
      [key]: value,
    }))
  }

  const handleAdd = () => {
    setEditingId(null)
    form.resetFields()
    form.setFieldsValue({
      schedule_date: dayjs(),
      period: 'morning',
      start_time: '08:00',
      end_time: '12:00',
      total_slots: 20,
      status: 1,
    })
    setModalVisible(true)
  }

  const handleEdit = (record: Schedule) => {
    setEditingId(record.id)
    form.setFieldsValue({
      doctor_id: record.doctor_id,
      schedule_date: dayjs(record.schedule_date),
      period: record.period,
      start_time: record.start_time,
      end_time: record.end_time,
      total_slots: record.total_slots,
      status: record.status,
    })
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个排班吗？',
      onOk: async () => {
        try {
          const response = await http.delete(`/admin/schedules/${id}`)
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
      const submitData = {
        ...values,
        schedule_date: values.schedule_date.format('YYYY-MM-DD'),
      }

      if (editingId) {
        const response = await http.put(`/admin/schedules/${editingId}`, {
          start_time: submitData.start_time,
          end_time: submitData.end_time,
          total_slots: submitData.total_slots,
          status: submitData.status,
        })
        if (response.data.code === 200000) {
          message.success('更新成功')
        }
      } else {
        const response = await http.post('/admin/schedules', submitData)
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

  const handleBatchAdd = () => {
    batchForm.resetFields()
    batchForm.setFieldsValue({
      start_date: dayjs(),
      end_date: dayjs().add(7, 'day'),
      periods: ['morning', 'afternoon'],
      week_days: [1, 2, 3, 4, 5],
      start_times: ['08:00', '14:00'],
      end_times: ['12:00', '17:00'],
      total_slots: 20,
    })
    setBatchModalVisible(true)
  }

  const handleBatchSubmit = async () => {
    try {
      const values = await batchForm.validateFields()
      const submitData = {
        ...values,
        start_date: values.start_date.format('YYYY-MM-DD'),
        end_date: values.end_date.format('YYYY-MM-DD'),
      }

      const response = await http.post('/admin/schedules/batch', submitData)
      if (response.data.code === 200000) {
        message.success(`批量创建成功，共创建 ${response.data.data.created_count} 条排班`)
        setBatchModalVisible(false)
        fetchData(pagination.current, pagination.pageSize)
      }
    } catch (error) {
      console.error('批量创建失败:', error)
    }
  }

  const handlePeriodChange = (value: string) => {
    if (value === 'morning') {
      form.setFieldsValue({
        start_time: '08:00',
        end_time: '12:00',
      })
    } else if (value === 'afternoon') {
      form.setFieldsValue({
        start_time: '14:00',
        end_time: '17:00',
      })
    }
  }

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>排班管理</h1>
        <Space>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
            添加排班
          </Button>
          <Button icon={<AppstoreAddOutlined />} onClick={handleBatchAdd}>
            批量排班
          </Button>
        </Space>
      </div>

      <div style={{ marginBottom: 16, padding: 16, background: '#fff', borderRadius: 8 }}>
        <Row gutter={16}>
          <Col span={6}>
            <Select
              placeholder="筛选科室"
              allowClear
              style={{ width: '100%' }}
              onChange={(value) => handleFilterChange('department_id', value)}
              options={departments.map(dept => ({
                label: dept.name,
                value: dept.id,
              }))}
            />
          </Col>
          <Col span={6}>
            <Select
              placeholder="筛选医生"
              allowClear
              style={{ width: '100%' }}
              onChange={(value) => handleFilterChange('doctor_id', value)}
              options={doctors.map(doc => ({
                label: `${doc.name} - ${doc.department_name}`,
                value: doc.id,
              }))}
            />
          </Col>
          <Col span={6}>
            <DatePicker
              placeholder="开始日期"
              style={{ width: '100%' }}
              onChange={(date) => handleFilterChange('start_date', date ? date.format('YYYY-MM-DD') : undefined)}
            />
          </Col>
          <Col span={6}>
            <DatePicker
              placeholder="结束日期"
              style={{ width: '100%' }}
              onChange={(date) => handleFilterChange('end_date', date ? date.format('YYYY-MM-DD') : undefined)}
            />
          </Col>
        </Row>
      </div>

      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        scroll={{ x: 1200 }}
        pagination={{
          ...pagination,
          showTotal: (total) => `共 ${total} 条`,
          showSizeChanger: true,
          showQuickJumper: true,
        }}
        onChange={handleTableChange}
      />

      <Modal
        title={editingId ? '编辑排班' : '添加排班'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="doctor_id"
            label="选择医生"
            rules={[{ required: true, message: '请选择医生' }]}
          >
            <Select
              placeholder="请选择医生"
              showSearch
              optionFilterProp="label"
              disabled={!!editingId}
              options={doctors.map(doc => ({
                label: `${doc.name} - ${doc.department_name}`,
                value: doc.id,
              }))}
            />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="schedule_date"
                label="排班日期"
                rules={[{ required: true, message: '请选择排班日期' }]}
              >
                <DatePicker
                  style={{ width: '100%' }}
                  disabledDate={(current) => current && current < dayjs().startOf('day')}
                  disabled={!!editingId}
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="period"
                label="时段"
                rules={[{ required: true, message: '请选择时段' }]}
              >
                <Select
                  placeholder="请选择时段"
                  disabled={!!editingId}
                  onChange={handlePeriodChange}
                >
                  <Select.Option value="morning">上午</Select.Option>
                  <Select.Option value="afternoon">下午</Select.Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="start_time"
                label="开始时间"
                rules={[{ required: true, message: '请输入开始时间' }]}
              >
                <Input placeholder="HH:mm 如 08:00" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="end_time"
                label="结束时间"
                rules={[{ required: true, message: '请输入结束时间' }]}
              >
                <Input placeholder="HH:mm 如 12:00" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="total_slots"
                label="总号源数"
                rules={[{ required: true, message: '请输入总号源数' }]}
              >
                <InputNumber
                  min={1}
                  max={999}
                  style={{ width: '100%' }}
                  placeholder="如 20"
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="status"
                label="状态"
                rules={[{ required: true, message: '请选择状态' }]}
              >
                <Select placeholder="请选择状态">
                  <Select.Option value={1}>正常</Select.Option>
                  <Select.Option value={0}>停诊</Select.Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>

      <Modal
        title="批量排班"
        open={batchModalVisible}
        onOk={handleBatchSubmit}
        onCancel={() => setBatchModalVisible(false)}
        width={700}
      >
        <Form form={batchForm} layout="vertical">
          <Form.Item
            name="doctor_id"
            label="选择医生"
            rules={[{ required: true, message: '请选择医生' }]}
          >
            <Select
              placeholder="请选择医生"
              showSearch
              optionFilterProp="label"
              options={doctors.map(doc => ({
                label: `${doc.name} - ${doc.department_name}`,
                value: doc.id,
              }))}
            />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="start_date"
                label="开始日期"
                rules={[{ required: true, message: '请选择开始日期' }]}
              >
                <DatePicker
                  style={{ width: '100%' }}
                  disabledDate={(current) => current && current < dayjs().startOf('day')}
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="end_date"
                label="结束日期"
                rules={[{ required: true, message: '请选择结束日期' }]}
              >
                <DatePicker
                  style={{ width: '100%' }}
                  disabledDate={(current) => current && current < dayjs().startOf('day')}
                />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item
            name="week_days"
            label="工作日"
            rules={[{ required: true, message: '请选择工作日' }]}
          >
            <Select
              mode="multiple"
              placeholder="请选择工作日"
              options={[
                { label: '周一', value: 1 },
                { label: '周二', value: 2 },
                { label: '周三', value: 3 },
                { label: '周四', value: 4 },
                { label: '周五', value: 5 },
                { label: '周六', value: 6 },
                { label: '周日', value: 0 },
              ]}
            />
          </Form.Item>
          <Form.Item
            name="periods"
            label="时段"
            rules={[{ required: true, message: '请选择时段' }]}
          >
            <Select
              mode="multiple"
              placeholder="请选择时段"
              options={[
                { label: '上午', value: 'morning' },
                { label: '下午', value: 'afternoon' },
              ]}
            />
          </Form.Item>
          <Form.Item
            name="start_times"
            label="开始时间"
            extra="格式: [上午开始时间, 下午开始时间]"
            rules={[{ required: true, message: '请输入开始时间' }]}
          >
            <Select
              mode="tags"
              placeholder="请输入开始时间，如 08:00 和 14:00"
              maxCount={2}
            />
          </Form.Item>
          <Form.Item
            name="end_times"
            label="结束时间"
            extra="格式: [上午结束时间, 下午结束时间]"
            rules={[{ required: true, message: '请输入结束时间' }]}
          >
            <Select
              mode="tags"
              placeholder="请输入结束时间，如 12:00 和 17:00"
              maxCount={2}
            />
          </Form.Item>
          <Form.Item
            name="total_slots"
            label="每时段号源数"
            rules={[{ required: true, message: '请输入号源数' }]}
          >
            <InputNumber
              min={1}
              max={999}
              style={{ width: '100%' }}
              placeholder="如 20"
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default ScheduleList
