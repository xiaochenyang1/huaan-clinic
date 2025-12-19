<template>
  <view class="page">
    <view class="card">
      <view class="title">完善个人信息</view>
      <view class="muted">首次登录建议完善昵称与性别，用于预约信息展示。</view>

      <view class="field">
        <text class="label">昵称</text>
        <input class="input" v-model="form.nickname" placeholder="请输入昵称（最多64字）" maxlength="64" />
      </view>

      <view class="field">
        <text class="label">性别</text>
        <view class="seg">
          <view class="seg-item" :class="{ active: form.gender === 0 }" @click="form.gender = 0">未知</view>
          <view class="seg-item" :class="{ active: form.gender === 1 }" @click="form.gender = 1">男</view>
          <view class="seg-item" :class="{ active: form.gender === 2 }" @click="form.gender = 2">女</view>
        </view>
      </view>

      <view class="hint">
        手机号绑定：当前后端 `PUT /user/info` 仅做手机号格式与唯一性校验，未接入短信校验绑定流程。
        如需手机号登录/绑定，请在登录页选择“短信（测试）”方式。
      </view>

      <button class="btn primary" @click="save" :disabled="loading">
        {{ loading ? '保存中…' : '保存并继续' }}
      </button>
      <button class="btn" @click="skip" :disabled="loading">暂不完善</button>
    </view>
  </view>
</template>

<script setup>
import { onLoad } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import { getUserInfo, updateUserInfo } from '../../api/auth'
import { getUser, isLoggedIn, setUser, toLoginPage } from '../../utils/auth'

const loading = ref(false)
const form = reactive({
  nickname: '',
  gender: 0,
})

async function hydrate() {
  const u = getUser()
  if (u?.nickname) form.nickname = u.nickname
  if (typeof u?.gender === 'number') form.gender = u.gender
  try {
    const fresh = await getUserInfo()
    setUser(fresh)
    if (fresh?.nickname) form.nickname = fresh.nickname
    if (typeof fresh?.gender === 'number') form.gender = fresh.gender
  } catch (e) {}
}

async function save() {
  loading.value = true
  try {
    const payload = {
      nickname: (form.nickname || '').trim(),
      gender: form.gender,
    }
    const u = await updateUserInfo(payload)
    setUser(u)
    uni.switchTab({ url: '/pages/index/index' })
  } finally {
    loading.value = false
  }
}

function skip() {
  uni.switchTab({ url: '/pages/index/index' })
}

onLoad(() => {
  if (!isLoggedIn()) {
    toLoginPage()
    return
  }
  hydrate()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 24rpx;
  background: #f6f7f9;
  display: flex;
  align-items: center;
  justify-content: center;
}
.card {
  width: 100%;
  max-width: 680rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 32rpx;
  border: 1rpx solid #e5e7eb;
}
.title {
  font-size: 40rpx;
  font-weight: 700;
  color: #111827;
}
.muted {
  margin-top: 10rpx;
  font-size: 26rpx;
  color: #6b7280;
  line-height: 1.7;
}
.field {
  margin-top: 22rpx;
}
.label {
  display: block;
  font-size: 24rpx;
  color: #374151;
  margin-bottom: 8rpx;
}
.input {
  height: 84rpx;
  border: 1rpx solid #e5e7eb;
  border-radius: 12rpx;
  padding: 0 20rpx;
  font-size: 28rpx;
  background: #fff;
}
.seg {
  display: flex;
  border: 1rpx solid #e5e7eb;
  border-radius: 12rpx;
  overflow: hidden;
}
.seg-item {
  flex: 1;
  text-align: center;
  padding: 18rpx 0;
  font-size: 26rpx;
  color: #374151;
  background: #f9fafb;
}
.seg-item.active {
  background: #111827;
  color: #fff;
  font-weight: 600;
}
.hint {
  margin-top: 18rpx;
  font-size: 24rpx;
  color: #6b7280;
  line-height: 1.7;
}
.btn {
  margin-top: 14rpx;
  height: 84rpx;
  line-height: 84rpx;
  border-radius: 12rpx;
  border: 1rpx solid #d1d5db;
  background: #fff;
  color: #111827;
  font-size: 28rpx;
}
.btn.primary {
  border: 1rpx solid #111827;
  background: #111827;
  color: #fff;
}
</style>

