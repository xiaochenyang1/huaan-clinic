import { useState } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { Layout, Menu, Dropdown, Avatar } from 'antd'
import type { MenuProps } from 'antd'
import {
  DashboardOutlined,
  TeamOutlined,
  UserOutlined,
  CalendarOutlined,
  MedicineBoxOutlined,
  LogoutOutlined,
  UserSwitchOutlined,
  ScheduleOutlined,
  BarChartOutlined,
  FileTextOutlined,
  SettingOutlined,
  SafetyCertificateOutlined,
} from '@ant-design/icons'
import { hasAnyPermission } from '@/utils/permissions'

const { Header, Sider, Content } = Layout

const MainLayout = () => {
  const [collapsed, setCollapsed] = useState(false)
  const navigate = useNavigate()
  const location = useLocation()

  const rawMenuItems: Array<{ key: string; icon: JSX.Element; label: string; perms?: string[] }> = [
    { key: '/dashboard', icon: <DashboardOutlined />, label: '仪表盘', perms: ['statistics:view'] },
    { key: '/department', icon: <MedicineBoxOutlined />, label: '科室管理', perms: ['department:view'] },
    { key: '/doctor', icon: <TeamOutlined />, label: '医生管理', perms: ['doctor:view'] },
    { key: '/schedule', icon: <ScheduleOutlined />, label: '排班管理', perms: ['schedule:view'] },
    { key: '/appointment', icon: <CalendarOutlined />, label: '预约管理', perms: ['appointment:view'] },
    { key: '/patient', icon: <UserOutlined />, label: '患者管理', perms: ['patient:view'] },
    { key: '/statistics', icon: <BarChartOutlined />, label: '数据统计', perms: ['statistics:view'] },
    { key: '/system/admin', icon: <SettingOutlined />, label: '管理员管理', perms: ['admin:view'] },
    { key: '/system/role', icon: <SafetyCertificateOutlined />, label: '角色管理', perms: ['role:view'] },
    { key: '/logs', icon: <FileTextOutlined />, label: '系统日志', perms: ['log:view'] },
  ]

  const menuItems: MenuProps['items'] = rawMenuItems
    .filter((item) => !item.perms || hasAnyPermission(item.perms))
    .map(({ perms, ...rest }) => rest)

  const userMenuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserSwitchOutlined />,
      label: '个人信息',
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      danger: true,
    },
  ]

  const handleMenuClick: MenuProps['onClick'] = ({ key }) => {
    navigate(key)
  }

  const handleUserMenuClick: MenuProps['onClick'] = ({ key }) => {
    if (key === 'logout') {
      localStorage.removeItem('token')
      localStorage.removeItem('permissions')
      navigate('/login')
    } else if (key === 'profile') {
      navigate('/profile')
    }
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        collapsible
        collapsed={collapsed}
        onCollapse={(value) => setCollapsed(value)}
        style={{
          overflow: 'auto',
          height: '100vh',
          position: 'fixed',
          left: 0,
          top: 0,
          bottom: 0,
          background: '#ffffff',
          boxShadow: '2px 0 8px rgba(0,0,0,0.08)',
        }}
        theme="light"
      >
        <div style={{
          height: 64,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          gap: 8,
          borderBottom: '1px solid #e5e7eb',
          background: 'linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%)',
        }}>
          <MedicineBoxOutlined style={{ fontSize: 24, color: '#fff' }} />
          {!collapsed && (
            <span style={{
              fontSize: 18,
              fontWeight: 600,
              color: '#fff',
            }}>
              华安医疗
            </span>
          )}
        </div>
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={handleMenuClick}
          style={{
            border: 'none',
            background: '#fff',
          }}
        />
      </Sider>
      <Layout style={{ marginLeft: collapsed ? 80 : 200, transition: 'all 0.2s' }}>
        <Header style={{
          padding: '0 32px',
          background: '#ffffff',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          borderBottom: '1px solid #e5e7eb',
          boxShadow: '0 1px 4px rgba(0,0,0,0.04)',
        }}>
          <div style={{
            fontSize: 18,
            fontWeight: 600,
            color: '#1f2937',
          }}>
            华安医疗后台管理系统
          </div>
          <Dropdown
            menu={{ items: userMenuItems, onClick: handleUserMenuClick }}
            placement="bottomRight"
          >
            <div style={{
              cursor: 'pointer',
              display: 'flex',
              alignItems: 'center',
              gap: 12,
              padding: '8px 12px',
              borderRadius: 8,
              transition: 'background 0.3s',
            }}>
              <Avatar
                size={32}
                icon={<UserOutlined />}
                style={{
                  background: 'linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%)',
                }}
              />
              <span style={{
                color: '#1f2937',
                fontWeight: 500,
              }}>
                管理员
              </span>
            </div>
          </Dropdown>
        </Header>
        <Content style={{
          margin: 12,
          overflow: 'initial',
          background: '#f9fafb',
        }}>
          <div
            style={{
              padding: 24,
              background: '#ffffff',
              borderRadius: 12,
              minHeight: 'calc(100vh - 112px)',
              boxShadow: '0 1px 3px rgba(0,0,0,0.08)',
            }}
          >
            <Outlet />
          </div>
        </Content>
      </Layout>
    </Layout>
  )
}

export default MainLayout
