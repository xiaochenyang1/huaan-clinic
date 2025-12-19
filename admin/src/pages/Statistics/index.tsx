import { useState, useEffect } from 'react'
import { Card, Row, Col, Statistic, Table, DatePicker, Spin } from 'antd'
import { UserOutlined, CheckCircleOutlined, CloseCircleOutlined, ClockCircleOutlined } from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import type { Dayjs } from 'dayjs'
import dayjs from 'dayjs'
import http from '@/utils/http'

const { RangePicker } = DatePicker

interface StatisticsData {
  appointment_stats: {
    total: number
    pending: number
    confirmed: number
    completed: number
    cancelled: number
    missed: number
  }
  doctor_stats: Array<{
    doctor_id: number
    doctor_name: string
    department_name: string
    total_appointments: number
    completed_appointments: number
    cancelled_appointments: number
  }>
  department_ranking: Array<{
    department_id: number
    department_name: string
    appointment_count: number
    completion_rate: number
  }>
  time_slot_distribution: Array<{
    period: string
    period_name: string
    count: number
    percentage: number
  }>
}

const Statistics = () => {
  const [data, setData] = useState<StatisticsData | null>(null)
  const [loading, setLoading] = useState(false)
  const [dateRange, setDateRange] = useState<[Dayjs, Dayjs]>([
    dayjs().subtract(30, 'day'),
    dayjs(),
  ])

  type DoctorStatRow = StatisticsData['doctor_stats'][number]
  type DepartmentRankingRow = StatisticsData['department_ranking'][number]

  const doctorColumns: ColumnsType<DoctorStatRow> = [
    {
      title: '医生姓名',
      dataIndex: 'doctor_name',
      key: 'doctor_name',
    },
    {
      title: '科室',
      dataIndex: 'department_name',
      key: 'department_name',
    },
    {
      title: '总预约数',
      dataIndex: 'total_appointments',
      key: 'total_appointments',
      sorter: (a, b) => a.total_appointments - b.total_appointments,
    },
    {
      title: '已完成',
      dataIndex: 'completed_appointments',
      key: 'completed_appointments',
      render: (val: number) => <span style={{ color: '#52c41a' }}>{val}</span>,
    },
    {
      title: '已取消',
      dataIndex: 'cancelled_appointments',
      key: 'cancelled_appointments',
      render: (val: number) => <span style={{ color: '#ff4d4f' }}>{val}</span>,
    },
  ]

  const departmentColumns: ColumnsType<DepartmentRankingRow> = [
    {
      title: '排名',
      key: 'rank',
      width: 80,
      render: (_, __, index) => index + 1,
    },
    {
      title: '科室名称',
      dataIndex: 'department_name',
      key: 'department_name',
    },
    {
      title: '预约数量',
      dataIndex: 'appointment_count',
      key: 'appointment_count',
      sorter: (a, b) => a.appointment_count - b.appointment_count,
    },
    {
      title: '完成率',
      dataIndex: 'completion_rate',
      key: 'completion_rate',
      render: (val: number) => `${(val * 100).toFixed(1)}%`,
    },
  ]

  useEffect(() => {
    fetchStatistics()
  }, [dateRange])

  const fetchStatistics = async () => {
    setLoading(true)
    try {
      const response = await http.get('/admin/statistics', {
        params: {
          start_date: dateRange[0].format('YYYY-MM-DD'),
          end_date: dateRange[1].format('YYYY-MM-DD'),
        },
      })
      if (response.data.code === 200000) {
        setData(response.data.data)
      }
    } catch (error) {
      console.error('获取统计数据失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleDateChange = (dates: [Dayjs | null, Dayjs | null] | null) => {
    if (dates && dates[0] && dates[1]) {
      setDateRange([dates[0], dates[1]])
    }
  }

  if (loading && !data) {
    return (
      <div style={{ textAlign: 'center', padding: '100px 0' }}>
        <Spin size="large" />
      </div>
    )
  }

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>数据统计</h1>
        <RangePicker
          value={dateRange}
          onChange={handleDateChange}
          format="YYYY-MM-DD"
          allowClear={false}
        />
      </div>

      {/* 预约统计卡片 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable>
            <Statistic
              title="总预约数"
              value={data?.appointment_stats.total || 0}
              prefix={<UserOutlined style={{ color: '#1890ff' }} />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable>
            <Statistic
              title="待确认"
              value={data?.appointment_stats.pending || 0}
              prefix={<ClockCircleOutlined style={{ color: '#faad14' }} />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable>
            <Statistic
              title="已完成"
              value={data?.appointment_stats.completed || 0}
              prefix={<CheckCircleOutlined style={{ color: '#52c41a' }} />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card hoverable>
            <Statistic
              title="已取消"
              value={data?.appointment_stats.cancelled || 0}
              prefix={<CloseCircleOutlined style={{ color: '#ff4d4f' }} />}
              valueStyle={{ color: '#ff4d4f' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 时段分布 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col span={24}>
          <Card title="预约时段分布" bordered={false}>
            <Row gutter={16}>
              {data?.time_slot_distribution.map((slot) => (
                <Col key={slot.period} xs={24} sm={12} lg={6}>
                  <Card>
                    <Statistic
                      title={slot.period_name}
                      value={slot.count}
                      suffix={`/ ${slot.percentage.toFixed(1)}%`}
                    />
                  </Card>
                </Col>
              ))}
            </Row>
          </Card>
        </Col>
      </Row>

      {/* 科室排行 */}
      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} lg={12}>
          <Card title="科室预约排行" bordered={false}>
            <Table
              columns={departmentColumns}
              dataSource={data?.department_ranking || []}
              rowKey="department_id"
              pagination={false}
              size="small"
            />
          </Card>
        </Col>

        {/* 医生统计 */}
        <Col xs={24} lg={12}>
          <Card title="医生预约统计 TOP10" bordered={false}>
            <Table
              columns={doctorColumns}
              dataSource={data?.doctor_stats.slice(0, 10) || []}
              rowKey="doctor_id"
              pagination={false}
              size="small"
            />
          </Card>
        </Col>
      </Row>

      {/* 完整医生统计表格 */}
      <Card title="医生详细统计" bordered={false}>
        <Table
          columns={doctorColumns}
          dataSource={data?.doctor_stats || []}
          rowKey="doctor_id"
          pagination={{
            pageSize: 10,
            showTotal: (total) => `共 ${total} 条`,
            showSizeChanger: true,
            showQuickJumper: true,
          }}
        />
      </Card>
    </div>
  )
}

export default Statistics
