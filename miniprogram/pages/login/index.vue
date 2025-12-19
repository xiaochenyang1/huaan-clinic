<template>
  <view class="page">
    <view class="card">
      <view class="title">华安医疗预约</view>
      <view class="sub">请登录后继续</view>

      <view class="seg">
        <view class="seg-item" :class="{ active: mode === 'wechat' }" @click="mode = 'wechat'">微信</view>
        <view class="seg-item" :class="{ active: mode === 'password' }" @click="mode = 'password'">密码</view>
        <view class="seg-item" :class="{ active: mode === 'phone' }" @click="mode = 'phone'">短信（测试）</view>
      </view>

      <view v-if="mode === 'wechat'" class="section">
        <button class="btn primary" @click="handleWeChatLogin" :disabled="loading">
          {{ loading ? '登录中…' : '微信一键登录' }}
        </button>
        <view class="hint">如在非微信环境运行，请使用密码/短信登录。</view>
      </view>

      <view v-else-if="mode === 'password'" class="section">
        <view class="field">
          <text class="label">用户名</text>
          <input class="input" v-model="username" placeholder="请输入用户名" />
        </view>
        <view class="field">
          <text class="label">密码</text>
          <input class="input" v-model="password" placeholder="请输入密码" password />
        </view>
        <button class="btn primary" @click="handlePasswordLogin" :disabled="loading">
          {{ loading ? '登录中…' : '登录' }}
        </button>
      </view>

      <view v-else class="section">
        <view class="hint strong">提示：需后端开启 `sms.enabled` 并配置 `sms.provider`；验证码不会回传到客户端（请看后端日志/真实短信）。</view>
        <view class="field">
          <text class="label">手机号</text>
          <input class="input" v-model="phone" placeholder="请输入手机号" />
        </view>
        <view class="field row">
          <view class="col">
            <text class="label">验证码</text>
            <input class="input" v-model="smsCode" placeholder="请输入验证码" />
          </view>
          <button class="btn" @click="handleSendCode" :disabled="smsSending || !phone">
            {{ smsSending ? `${smsCountdown}s` : '发送' }}
          </button>
        </view>
        <button class="btn primary" @click="handlePhoneLogin" :disabled="loading">
          {{ loading ? '登录中…' : '登录' }}
        </button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onUnload } from '@dcloudio/uni-app'
import { wechatLogin, passwordLogin, phoneLogin, sendSmsCode, getUserInfo } from '../../api/auth'
import { useUserStore } from '../../store'

const mode = ref('wechat')
const loading = ref(false)

const username = ref('')
const password = ref('')

const phone = ref('')
const smsCode = ref('')
const smsSending = ref(false)
const smsCountdown = ref(60)
let smsTimer = null
const userStore = useUserStore()

function stopTimer() {
  if (smsTimer) clearInterval(smsTimer)
  smsTimer = null
}

async function afterLogin(data) {
  const accessToken = data?.token || data?.access_token
  const refreshToken = data?.refresh_token
  userStore.setTokens({ accessToken, refreshToken })

  if (data?.user) {
    userStore.setUser(data.user)
    if (data.user.has_phone === false) {
      uni.redirectTo({ url: '/pages/login/complete' })
      return
    }
  }

  try {
    const user = await getUserInfo()
    userStore.setUser(user)
    if (user && user.has_phone === false) {
      uni.redirectTo({ url: '/pages/login/complete' })
      return
    }
  } catch (e) {}
  uni.switchTab({ url: '/pages/index/index' })
}

async function handleWeChatLogin() {
  loading.value = true
  try {
    const loginRes = await new Promise((resolve, reject) => {
      uni.login({
        provider: 'weixin',
        success: (res) => resolve(res),
        fail: (err) => reject(err),
      })
    })
    const data = await wechatLogin({ code: loginRes.code })
    await afterLogin(data)
  } finally {
    loading.value = false
  }
}

async function handlePasswordLogin() {
  if (!username.value || !password.value) {
    uni.showToast({ title: '请输入用户名和密码', icon: 'none' })
    return
  }
  loading.value = true
  try {
    const data = await passwordLogin({ username: username.value, password: password.value })
    await afterLogin(data)
  } finally {
    loading.value = false
  }
}

async function handleSendCode() {
  if (!phone.value) return
  smsSending.value = true
  smsCountdown.value = 60
  stopTimer()
  try {
    await sendSmsCode({ phone: phone.value })
    uni.showToast({ title: '已发送', icon: 'none' })
    smsTimer = setInterval(() => {
      smsCountdown.value -= 1
      if (smsCountdown.value <= 0) {
        smsSending.value = false
        stopTimer()
      }
    }, 1000)
  } catch (e) {
    smsSending.value = false
    stopTimer()
  }
}

async function handlePhoneLogin() {
  if (!phone.value || !smsCode.value) {
    uni.showToast({ title: '请输入手机号和验证码', icon: 'none' })
    return
  }
  loading.value = true
  try {
    const data = await phoneLogin({ phone: phone.value, code: smsCode.value })
    await afterLogin(data)
  } finally {
    loading.value = false
  }
}

onUnload(() => {
  stopTimer()
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
  font-weight: 600;
  color: #111827;
}
.sub {
  margin-top: 8rpx;
  font-size: 26rpx;
  color: #6b7280;
}
.seg {
  margin-top: 24rpx;
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
  background: #ffffff;
  color: #111827;
  font-weight: 600;
}
.section {
  margin-top: 24rpx;
}
.field {
  margin-bottom: 18rpx;
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
.row {
  display: flex;
  gap: 16rpx;
  align-items: flex-end;
}
.col {
  flex: 1;
}
.btn {
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
.hint {
  margin-top: 12rpx;
  font-size: 24rpx;
  color: #6b7280;
}
.hint.strong {
  margin-top: 0;
  margin-bottom: 12rpx;
  color: #991b1b;
}
</style>
