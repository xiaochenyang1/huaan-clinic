import { useState, useEffect } from 'react'
import { Table, Button, Space, Tag, Select, message, Drawer, Descriptions, Spin } from 'antd'
import { CheckOutlined, CloseOutlined, DownloadOutlined, EyeOutlined } from '@ant-design/icons'
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table'
import http from '@/utils/http'

interface Appointment {
  id: number
  patient_name: string
  doctor_name: string
  department_name: string
  appointment_date: string
  appointment_time: string
  period_name: string
  status: string
  status_name: string
  created_at: string
}

interface AppointmentDetail {
  id: number
  appointment_no: string
  user_id: number
  patient_id: number
  patient_name: string
  doctor_id: number
  doctor_name: string
  doctor_title: string
  doctor_avatar?: string
  department_id: number
  department_name: string
  appointment_date: string
  period: string
  period_name: string
  appointment_time: string
  slot_number: number
  status: string
  status_name: string
  symptom?: string
  cancel_reason?: string
  cancelled_at?: string
  checked_in_at?: string
  completed_at?: string
  created_at: string
}

const statusMap: Record<string, { color: string }> = {
  pending: { color: 'default' },
  checked_in: { color: 'blue' },
  cancelled: { color: 'red' },
  completed: { color: 'green' },
  missed: { color: 'orange' },
}

const AppointmentList = () => {
  const [data, setData] = useState<Appointment[]>([])
  const [loading, setLoading] = useState(false)
  const [selectedStatus, setSelectedStatus] = useState<string>()
  const [detailOpen, setDetailOpen] = useState(false)
  const [detailLoading, setDetailLoading] = useState(false)
  const [detail, setDetail] = useState<AppointmentDetail | null>(null)
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })

  const columns: ColumnsType<Appointment> = [
    {
      title: '患者姓名',
      dataIndex: 'patient_name',
      key: 'patient_name',
    },
    {
      title: '医生',
      dataIndex: 'doctor_name',
      key: 'doctor_name',
    },
    {
      title: '科室',
      dataIndex: 'department_name',
      key: 'department_name',
    },
    {
      title: '预约日期',
      dataIndex: 'appointment_date',
      key: 'appointment_date',
    },
    {
      title: '预约时段',
      key: 'time_slot',
      render: (_, record) => (
        <span>
          {record.period_name || '-'} {record.appointment_time || ''}
        </span>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string, record) => (
        <Tag color={statusMap[status]?.color}>{record.status_name}</Tag>
      ),
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
      width: 220,
      render: (_, record) => (
        <Space size="middle">
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => openDetail(record.id)}
          >
            详情
          </Button>
          {record.status === 'pending' && (
            <>
              <Button
                type="link"
                size="small"
                icon={<CheckOutlined />}
                onClick={() => handleUpdateStatus(record.id, 'checked_in')}
              >
                签到
              </Button>
              <Button
                type="link"
                size="small"
                danger
                icon={<CloseOutlined />}
                onClick={() => handleUpdateStatus(record.id, 'cancelled')}
              >
                取消
              </Button>
            </>
          )}
          {record.status === 'checked_in' && (
            <>
              <Button
                type="link"
                size="small"
                icon={<CheckOutlined />}
                onClick={() => handleUpdateStatus(record.id, 'completed')}
              >
                完成
              </Button>
              <Button
                type="link"
                size="small"
                danger
                icon={<CloseOutlined />}
                onClick={() => handleUpdateStatus(record.id, 'cancelled')}
              >
                取消
              </Button>
            </>
          )}
        </Space>
      ),
    },
  ]

  const openDetail = async (id: number) => {
    setDetailOpen(true)
    setDetailLoading(true)
    setDetail(null)
    try {
      const response = await http.get(`/admin/appointments/${id}`)
      if (response.data.code === 200000) {
        setDetail(response.data.data)
      }
    } catch (error) {
      console.error('获取预约详情失败:', error)
    } finally {
      setDetailLoading(false)
    }
  }

  const fetchData = async (page = 1, pageSize = 10, status?: string) => {
    setLoading(true)
    try {
      const params: { page: number; page_size: number; status?: string } = {
        page,
        page_size: pageSize,
      }
      if (status) {
        params.status = status
      }
      const response = await http.get('/admin/appointments', { params })
      if (response.data.code === 200000) {
        setData(response.data.data.list)
        setPagination({
          current: page,
          pageSize: pageSize,
          total: response.data.data.total,
        })
      }
    } catch (error) {
      console.error('获取预约列表失败:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  const handleTableChange = (newPagination: TablePaginationConfig) => {
    fetchData(newPagination.current || 1, newPagination.pageSize || 10, selectedStatus)
  }

  const handleStatusChange = (value: string) => {
    setSelectedStatus(value)
    fetchData(1, pagination.pageSize, value)
  }

  const handleUpdateStatus = async (id: number, status: string) => {
    try {
      const response = await http.put(`/admin/appointments/${id}`, { status })
      if (response.data.code === 200000) {
        message.success('操作成功')
        fetchData(pagination.current, pagination.pageSize, selectedStatus)
      }
    } catch (error) {
      console.error('操作失败:', error)
    }
  }

  const handleExport = async () => {
    try {
      const params: { status?: string } = {}
      if (selectedStatus) {
        params.status = selectedStatus
      }

      message.loading({ content: '正在导出...', key: 'export' })

      const response = await http.get('/admin/appointments/export', {
        params,
        responseType: 'blob',
      })

      // 创建下载链接
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url

      // 从响应头获取文件名，如果没有则使用默认名称
      const contentDisposition = response.headers['content-disposition']
      let fileName = '预约数据.csv'
      if (contentDisposition) {
        const fileNameMatch = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/)
        if (fileNameMatch && fileNameMatch[1]) {
          fileName = decodeURIComponent(fileNameMatch[1].replace(/['"]/g, ''))
        }
      }

      link.setAttribute('download', fileName)
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)

      message.success({ content: '导出成功', key: 'export' })
    } catch (error) {
      console.error('导出失败:', error)
      message.error({ content: '导出失败', key: 'export' })
    }
  }

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>预约管理</h1>
        <Space>
          <Select
            placeholder="筛选状态"
            allowClear
            style={{ width: 150 }}
            onChange={handleStatusChange}
            options={[
              { label: '待就诊', value: 'pending' },
              { label: '已签到', value: 'checked_in' },
              { label: '已取消', value: 'cancelled' },
              { label: '已完成', value: 'completed' },
              { label: '爽约', value: 'missed' },
            ]}
          />
          <Button
            icon={<DownloadOutlined />}
            onClick={handleExport}
          >
            导出Excel
          </Button>
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

      <Drawer
        title="预约详情"
        open={detailOpen}
        onClose={() => setDetailOpen(false)}
        width={520}
      >
        {detailLoading ? (
          <div style={{ textAlign: 'center', padding: '48px 0' }}>
            <Spin />
          </div>
        ) : detail ? (
          <Descriptions bordered column={1} size="small">
            <Descriptions.Item label="预约编号">{detail.appointment_no}</Descriptions.Item>
            <Descriptions.Item label="患者">{detail.patient_name}</Descriptions.Item>
            <Descriptions.Item label="医生">{detail.doctor_name} {detail.doctor_title ? `(${detail.doctor_title})` : ''}</Descriptions.Item>
            <Descriptions.Item label="科室">{detail.department_name}</Descriptions.Item>
            <Descriptions.Item label="就诊时间">
              {detail.appointment_date} {detail.period_name} {detail.appointment_time}（号序 {detail.slot_number}）
            </Descriptions.Item>
            <Descriptions.Item label="状态">
              <Tag color={statusMap[detail.status]?.color}>{detail.status_name}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="症状描述">{detail.symptom || '-'}</Descriptions.Item>
            <Descriptions.Item label="取消原因">{detail.cancel_reason || '-'}</Descriptions.Item>
            <Descriptions.Item label="取消时间">{detail.cancelled_at || '-'}</Descriptions.Item>
            <Descriptions.Item label="签到时间">{detail.checked_in_at || '-'}</Descriptions.Item>
            <Descriptions.Item label="完成时间">{detail.completed_at || '-'}</Descriptions.Item>
            <Descriptions.Item label="创建时间">{detail.created_at}</Descriptions.Item>
          </Descriptions>
        ) : (
          <div style={{ color: '#6b7280' }}>暂无数据</div>
        )}
      </Drawer>
    </div>
  )
}

export default AppointmentList
