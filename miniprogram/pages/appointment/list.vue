<template>
  <view class="page">
    <view class="tabs">
      <view class="tab" :class="{ active: status === '' }" @click="setStatus('')">全部</view>
      <view class="tab" :class="{ active: status === 'pending' }" @click="setStatus('pending')">待就诊</view>
      <view class="tab" :class="{ active: status === 'checked_in' }" @click="setStatus('checked_in')">已签到</view>
      <view class="tab" :class="{ active: status === 'completed' }" @click="setStatus('completed')">已完成</view>
      <view class="tab" :class="{ active: status === 'cancelled' }" @click="setStatus('cancelled')">已取消</view>
    </view>

    <view class="panel">
      <view v-if="loading" class="muted">加载中…</view>
      <view v-else>
        <view v-if="items.length === 0" class="muted">暂无记录</view>
        <view v-for="a in items" :key="a.id" class="card" @click="goDetail(a)">
          <view class="card-title">
            <text class="strong">{{ a.department_name || '-' }}</text>
            <text class="badge">{{ a.status_name || a.status || '-' }}</text>
          </view>
          <view class="card-sub">{{ a.doctor_name || '-' }} · {{ a.patient_name || '-' }}</view>
          <view class="card-sub">{{ a.appointment_date || '-' }} {{ a.period_name || '' }}</view>
        </view>
      </view>
    </view>

    <view class="floating">
      <button class="btn primary" @click="goStart">开始预约</button>
    </view>
  </view>
</template>

<script setup>
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { isLoggedIn, toLoginPage } from '../../utils/auth'
import { listAppointments } from '../../api/appointment'

const status = ref('')
const loading = ref(false)
const items = ref([])

function goStart() {
  uni.navigateTo({ url: '/pages/appointment/department' })
}

function goDetail(a) {
  uni.navigateTo({ url: `/pages/appointment/detail?id=${a.id}` })
}

function setStatus(s) {
  status.value = s
  load()
}

async function load() {
  loading.value = true
  try {
    items.value = (await listAppointments({ status: status.value || undefined })) || []
  } finally {
    loading.value = false
  }
}

onShow(() => {
  if (!isLoggedIn()) {
    toLoginPage()
    return
  }
  load()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background: #f6f7f9;
  padding: 24rpx;
  padding-bottom: 120rpx;
}
.tabs {
  display: flex;
  gap: 10rpx;
  flex-wrap: wrap;
  margin-bottom: 14rpx;
}
.tab {
  padding: 12rpx 18rpx;
  border-radius: 999rpx;
  border: 1rpx solid #e5e7eb;
  background: #fff;
  color: #374151;
  font-size: 24rpx;
}
.tab.active {
  border-color: #111827;
  color: #111827;
  font-weight: 600;
}
.panel {
  background: transparent;
}
.card {
  background: #fff;
  border: 1rpx solid #e5e7eb;
  border-radius: 16rpx;
  padding: 20rpx;
  margin-bottom: 14rpx;
}
.card-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.strong {
  font-size: 30rpx;
  font-weight: 700;
  color: #111827;
}
.badge {
  font-size: 22rpx;
  color: #374151;
  border: 1rpx solid #e5e7eb;
  border-radius: 999rpx;
  padding: 4rpx 12rpx;
  background: #f9fafb;
}
.card-sub {
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #6b7280;
}
.muted {
  color: #6b7280;
  font-size: 26rpx;
  padding: 12rpx 0;
}
.floating {
  position: fixed;
  left: 24rpx;
  right: 24rpx;
  bottom: 24rpx;
}
.btn {
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 14rpx;
  border: 1rpx solid #d1d5db;
  background: #fff;
  color: #111827;
  font-size: 30rpx;
}
.btn.primary {
  border: 1rpx solid #111827;
  background: #111827;
  color: #fff;
}
</style>

