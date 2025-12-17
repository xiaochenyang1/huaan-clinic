import { useState } from 'react'
import { Upload, message } from 'antd'
import { LoadingOutlined, PlusOutlined } from '@ant-design/icons'
import type { UploadChangeParam } from 'antd/es/upload'
import type { RcFile, UploadFile } from 'antd/es/upload/interface'
import http from '@/utils/http'

interface AvatarUploadProps {
  value?: string
  onChange?: (url: string) => void
}

const getBase64 = (img: RcFile, callback: (url: string) => void) => {
  const reader = new FileReader()
  reader.addEventListener('load', () => callback(reader.result as string))
  reader.readAsDataURL(img)
}

const beforeUpload = (file: RcFile) => {
  const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'
  if (!isJpgOrPng) {
    message.error('只能上传 JPG/PNG 格式的图片!')
    return false
  }
  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isLt2M) {
    message.error('图片大小不能超过 2MB!')
    return false
  }
  return true
}

const AvatarUpload: React.FC<AvatarUploadProps> = ({ value, onChange }) => {
  const [loading, setLoading] = useState(false)
  const [imageUrl, setImageUrl] = useState<string | undefined>(value)

  const handleChange = (info: UploadChangeParam<UploadFile>) => {
    if (info.file.status === 'uploading') {
      setLoading(true)
      return
    }
    if (info.file.status === 'done') {
      if (info.file.response?.code === 200000) {
        const url = info.file.response.data.url
        setImageUrl(url)
        setLoading(false)
        onChange?.(url)
        message.success('上传成功')
      } else {
        setLoading(false)
        message.error(info.file.response?.message || '上传失败')
      }
    }
    if (info.file.status === 'error') {
      setLoading(false)
      message.error('上传失败')
    }
  }

  const uploadButton = (
    <div>
      {loading ? <LoadingOutlined /> : <PlusOutlined />}
      <div style={{ marginTop: 8 }}>上传头像</div>
    </div>
  )

  // 获取token
  const getToken = () => {
    return localStorage.getItem('token') || ''
  }

  return (
    <Upload
      name="file"
      listType="picture-card"
      className="avatar-uploader"
      showUploadList={false}
      action={`${import.meta.env.VITE_API_BASE_URL || ''}/api/admin/upload/avatar`}
      headers={{
        Authorization: `Bearer ${getToken()}`,
      }}
      beforeUpload={beforeUpload}
      onChange={handleChange}
    >
      {imageUrl ? (
        <img src={imageUrl} alt="avatar" style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
      ) : (
        uploadButton
      )}
    </Upload>
  )
}

export default AvatarUpload
