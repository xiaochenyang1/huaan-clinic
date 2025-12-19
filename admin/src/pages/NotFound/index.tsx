import { Button, Result } from 'antd'
import { useLocation, useNavigate } from 'react-router-dom'

export default function NotFound() {
  const navigate = useNavigate()
  const location = useLocation()

  return (
    <Result
      status="404"
      title="404"
      subTitle={`页面不存在：${location.pathname}`}
      extra={[
        <Button type="primary" key="home" onClick={() => navigate('/dashboard', { replace: true })}>
          返回首页
        </Button>,
      ]}
    />
  )
}

