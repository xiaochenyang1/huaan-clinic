import { useState, useEffect } from 'react'
import { Table, Button, Space, Modal, Descriptions } from 'antd'
import { EyeOutlined } from '@ant-design/icons'
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table'
import http from '@/utils/http'

interface Patient {
  id: number
  name: string
  phone: string
  id_card: string
  gender: string
  age: number
  address: string
  created_at: string
  last_login_ip: string
}

const PatientList = () => {
  const [data, setData] = useState<Patient[]>([])
  const [loading, setLoading] = useState(false)
  const [detailVisible, setDetailVisible] = useState(false)
  const [selectedPatient, setSelectedPatient] = useState<Patient | null>(null)
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  })

  const columns: ColumnsType<Patient> = [
    {
      title: '姓名',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '性别',
      dataIndex: 'gender',
      key: 'gender',
    },
    {
      title: '年龄',
      dataIndex: 'age',
      key: 'age',
      width: 80,
    },
    {
      title: '联系电话',
      dataIndex: 'phone',
      key: 'phone',
    },
    {
      title: '身份证号',
      dataIndex: 'id_card',
      key: 'id_card',
      render: (idCard: string) =>
        idCard ? idCard.replace(/(\d{6})\d{8}(\d{4})/, '$1********$2') : '-',
    },
    {
      title: '最后登录IP',
      dataIndex: 'last_login_ip',
      key: 'last_login_ip',
      width: 140,
      render: (ip: string) => ip || '-',
    },
    {
      title: '注册时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      render: (_, record) => (
        <Space size="middle">
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record)}
          >
            详情
          </Button>
        </Space>
      ),
    },
  ]

  const fetchData = async (page = 1, pageSize = 10) => {
    setLoading(true)
    try {
      const response = await http.get('/admin/patients', {
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
      console.error('获取患者列表失败:', error)
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

  const handleViewDetail = (patient: Patient) => {
    setSelectedPatient(patient)
    setDetailVisible(true)
  }

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <h1 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>患者管理</h1>
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
        title="患者详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailVisible(false)}>
            关闭
          </Button>,
        ]}
        width={600}
      >
        {selectedPatient && (
          <Descriptions column={2} bordered>
            <Descriptions.Item label="姓名">{selectedPatient.name}</Descriptions.Item>
            <Descriptions.Item label="性别">{selectedPatient.gender}</Descriptions.Item>
            <Descriptions.Item label="年龄">{selectedPatient.age}</Descriptions.Item>
            <Descriptions.Item label="联系电话">{selectedPatient.phone}</Descriptions.Item>
            <Descriptions.Item label="身份证号" span={2}>
              {selectedPatient.id_card}
            </Descriptions.Item>
            <Descriptions.Item label="地址" span={2}>
              {selectedPatient.address || '未填写'}
            </Descriptions.Item>
            <Descriptions.Item label="最后登录IP">
              {selectedPatient.last_login_ip || '未记录'}
            </Descriptions.Item>
            <Descriptions.Item label="注册时间">
              {selectedPatient.created_at}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </div>
  )
}

export default PatientList
