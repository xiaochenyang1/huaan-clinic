<template>
  <view class="page">
    <view class="panel" v-if="doctor">
      <view class="name">{{ doctor.name }}</view>
      <view class="meta">
        <text v-if="doctor.title" class="pill">{{ doctor.title }}</text>
        <text v-if="doctor.department_name" class="pill">{{ doctor.department_name }}</text>
      </view>
      <view class="desc" v-if="doctor.introduction">{{ doctor.introduction }}</view>

      <button class="btn primary" @click="goSchedule">选择时段预约</button>
    </view>

    <view v-else class="panel">
      <view class="muted">{{ loading ? '加载中…' : '未找到医生信息' }}</view>
    </view>
  </view>
</template>

<script setup>
import { onLoad } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { getDoctor } from '../../api/doctor'
import { addRecentDoctor } from '../../utils/recent'

const id = ref('')
const loading = ref(false)
const doctor = ref(null)

function goSchedule() {
  if (!id.value) return
  uni.navigateTo({ url: `/pages/appointment/schedule?doctor_id=${id.value}` })
}

async function load() {
  if (!id.value) return
  loading.value = true
  try {
    doctor.value = await getDoctor(id.value)
    if (doctor.value?.id) addRecentDoctor(doctor.value)
  } finally {
    loading.value = false
  }
}

onLoad((q) => {
  id.value = q?.id || ''
  load()
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
.name {
  font-size: 38rpx;
  font-weight: 700;
  color: #111827;
}
.meta {
  margin-top: 10rpx;
  display: flex;
  gap: 12rpx;
  flex-wrap: wrap;
}
.pill {
  font-size: 24rpx;
  color: #374151;
  border: 1rpx solid #e5e7eb;
  border-radius: 999rpx;
  padding: 6rpx 14rpx;
  background: #f9fafb;
}
.desc {
  margin-top: 18rpx;
  font-size: 26rpx;
  color: #374151;
  line-height: 1.7;
}
.btn {
  margin-top: 22rpx;
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
.muted {
  color: #6b7280;
  font-size: 26rpx;
}
</style>
