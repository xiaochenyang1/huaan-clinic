import { useState, useEffect } from 'react'
import { Table, Button, Space, Tag, Select, message } from 'antd'
import { CheckOutlined, CloseOutlined, DownloadOutlined } from '@ant-design/icons'
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table'
import http from '@/utils/http'

interface Appointment {
  id: number
  patient_name: string
  doctor_name: string
  department_name: string
  appointment_date: string
  appointment_time: string
  time_slot: string
  status: string
  status_name: string
  created_at: string
}

const statusMap: Record<string, { color: string }> = {
  pending: { color: 'default' },
  confirmed: { color: 'blue' },
  cancelled: { color: 'red' },
  completed: { color: 'green' },
  missed: { color: 'orange' },
}

const AppointmentList = () => {
  const [data, setData] = useState<Appointment[]>([])
  const [loading, setLoading] = useState(false)
  const [selectedStatus, setSelectedStatus] = useState<string>()
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
      dataIndex: 'time_slot',
      key: 'time_slot',
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
      width: 180,
      render: (_, record) => (
        <Space size="middle">
          {record.status === 'pending' && (
            <>
              <Button
                type="link"
                size="small"
                icon={<CheckOutlined />}
                onClick={() => handleUpdateStatus(record.id, 'confirmed')}
              >
                确认
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

  const fetchData = async (page = 1, pageSize = 10, status?: string) => {
    setLoading(true)
    try {
      const params: any = {
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
      const params: any = {}
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
      let fileName = '预约数据.xlsx'
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
              { label: '待确认', value: 'pending' },
              { label: '已确认', value: 'confirmed' },
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
    </div>
  )
}

export default AppointmentList
