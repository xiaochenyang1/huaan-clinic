<template>
  <view class="page">
    <view class="panel">
      <view class="title">个人中心</view>
      <view class="muted" v-if="!user">未登录</view>
      <view v-else class="kv">
        <text class="k">用户</text><text class="v">{{ user.nickname || user.name || user.open_id || '-' }}</text>
      </view>

      <view class="actions">
        <button class="btn" @click="goProfile">个人信息</button>
        <button class="btn" @click="goSubscribe">消息提醒</button>
        <button class="btn" @click="goPatient">就诊人管理</button>
        <button class="btn" @click="goAppointments">我的预约</button>
        <button class="btn" @click="goRecords">就诊记录</button>
        <button class="btn danger" @click="doLogout">退出登录</button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { getUserInfo } from '../../api/auth'
import { getUser, isLoggedIn, logout, setUser, toLoginPage } from '../../utils/auth'

const user = ref(getUser())

function goProfile() {
  uni.navigateTo({ url: '/pages/user/profile' })
}
function goSubscribe() {
  uni.navigateTo({ url: '/pages/user/subscribe' })
}
function goPatient() {
  uni.navigateTo({ url: '/pages/user/patient-list' })
}
function goAppointments() {
  uni.switchTab({ url: '/pages/appointment/list' })
}
function goRecords() {
  uni.navigateTo({ url: '/pages/record/list' })
}
async function doLogout() {
  await logout()
}

onShow(async () => {
  if (!isLoggedIn()) {
    user.value = null
    toLoginPage()
    return
  }
  user.value = getUser()
  try {
    const u = await getUserInfo()
    setUser(u)
    user.value = u
  } catch (e) {}
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background: #f6f7f9;
  padding: 24rpx;
}
.panel {
  background: #fff;
  border: 1rpx solid #e5e7eb;
  border-radius: 16rpx;
  padding: 24rpx;
}
.title {
  font-size: 36rpx;
  font-weight: 700;
  color: #111827;
}
.muted {
  margin-top: 10rpx;
  color: #6b7280;
  font-size: 26rpx;
}
.kv {
  margin-top: 14rpx;
  display: flex;
  padding: 10rpx 0;
}
.k {
  width: 140rpx;
  color: #6b7280;
  font-size: 26rpx;
}
.v {
  flex: 1;
  color: #111827;
  font-size: 26rpx;
}
.actions {
  margin-top: 18rpx;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
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
.btn.danger {
  border: 1rpx solid #991b1b;
  color: #991b1b;
}
</style>
