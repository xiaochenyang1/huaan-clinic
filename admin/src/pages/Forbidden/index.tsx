import { Button, Result } from 'antd'
import { useNavigate } from 'react-router-dom'

export default function Forbidden() {
  const navigate = useNavigate()

  return (
    <Result
      status="403"
      title="403"
      subTitle="没有权限访问该页面"
      extra={[
        <Button type="primary" key="home" onClick={() => navigate('/dashboard', { replace: true })}>
          返回首页
        </Button>,
      ]}
    />
  )
}

