<template>
  <view class="page">
    <view class="panel">
      <view class="panel-title">可预约时段</view>
      <view class="muted" v-if="loading">加载中…</view>
      <view v-else>
        <view v-if="schedules.length === 0" class="muted">暂无可预约号源</view>
        <view v-for="s in schedules" :key="s.id" class="row" @click="selectSchedule(s)">
          <view class="row-main">
            <view class="row-title">{{ s.schedule_date }} {{ s.period_name }}</view>
            <view class="row-sub">
              <text v-if="s.doctor_name">{{ s.doctor_name }}</text>
              <text v-if="s.department_name"> · {{ s.department_name }}</text>
              <text v-if="typeof s.remaining_slots === 'number'"> · 余号 {{ s.remaining_slots }}</text>
            </view>
          </view>
          <view class="row-arrow">›</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { onLoad } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { listAvailableSchedules } from '../../api/schedule'
import { nextDaysRange } from '../../utils/date'

const doctorId = ref('')
const departmentId = ref('')
const schedules = ref([])
const loading = ref(false)

function selectSchedule(s) {
  const q = [
    `schedule_id=${s.id}`,
    `schedule_date=${encodeURIComponent(s.schedule_date || '')}`,
    `period_name=${encodeURIComponent(s.period_name || '')}`,
    `doctor_name=${encodeURIComponent(s.doctor_name || '')}`,
    `department_name=${encodeURIComponent(s.department_name || '')}`,
  ].join('&')
  uni.navigateTo({ url: `/pages/appointment/confirm?${q}` })
}

async function load() {
  const { startDate, endDate } = nextDaysRange(7)
  loading.value = true
  try {
    schedules.value = await listAvailableSchedules({
      doctor_id: doctorId.value || undefined,
      department_id: departmentId.value || undefined,
      start_date: startDate,
      end_date: endDate,
    })
  } finally {
    loading.value = false
  }
}

onLoad((q) => {
  doctorId.value = q?.doctor_id || ''
  departmentId.value = q?.department_id || ''
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
  padding: 20rpx;
}
.panel-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #111827;
  margin-bottom: 12rpx;
}
.row {
  display: flex;
  align-items: center;
  padding: 18rpx 8rpx;
  border-top: 1rpx solid #f3f4f6;
}
.row:first-of-type {
  border-top: none;
}
.row-main {
  flex: 1;
}
.row-title {
  font-size: 30rpx;
  color: #111827;
}
.row-sub {
  margin-top: 6rpx;
  font-size: 24rpx;
  color: #6b7280;
}
.row-arrow {
  font-size: 36rpx;
  color: #9ca3af;
  padding-left: 8rpx;
}
.muted {
  color: #6b7280;
  font-size: 26rpx;
  padding: 12rpx 0;
}
</style>
