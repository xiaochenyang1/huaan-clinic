<template>
  <view class="page">
    <view class="panel">
      <view class="title">预约已提交</view>
      <view class="muted">请按预约时间就诊，逾期将记为爽约。</view>

      <view v-if="appointment" class="block">
        <view class="kv"><text class="k">预约号</text><text class="v">{{ appointment.appointment_no || '-' }}</text></view>
        <view class="kv"><text class="k">日期</text><text class="v">{{ appointment.appointment_date || '-' }}</text></view>
        <view class="kv"><text class="k">时段</text><text class="v">{{ appointment.period_name || '-' }}</text></view>
        <view class="kv"><text class="k">医生</text><text class="v">{{ appointment.doctor_name || '-' }}</text></view>
        <view class="kv"><text class="k">科室</text><text class="v">{{ appointment.department_name || '-' }}</text></view>
        <view class="kv"><text class="k">就诊人</text><text class="v">{{ appointment.patient_name || '-' }}</text></view>
      </view>

      <button class="btn primary" @click="goList">查看我的预约</button>
      <button class="btn" @click="goHome">返回首页</button>
    </view>
  </view>
</template>

<script setup>
import { onLoad } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { getAppointment } from '../../api/appointment'

const appointmentId = ref('')
const appointment = ref(null)

function goList() {
  uni.switchTab({ url: '/pages/appointment/list' })
}
function goHome() {
  uni.switchTab({ url: '/pages/index/index' })
}

onLoad(async (q) => {
  appointmentId.value = q?.appointment_id || ''
  if (appointmentId.value) {
    try {
      appointment.value = await getAppointment(appointmentId.value)
    } catch (e) {}
  }
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background: #f6f7f9;
  padding: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}
.panel {
  width: 100%;
  background: #fff;
  border: 1rpx solid #e5e7eb;
  border-radius: 16rpx;
  padding: 28rpx;
}
.title {
  font-size: 38rpx;
  font-weight: 700;
  color: #111827;
}
.muted {
  margin-top: 10rpx;
  color: #6b7280;
  font-size: 26rpx;
}
.block {
  margin-top: 18rpx;
  border-top: 1rpx solid #f3f4f6;
  padding-top: 12rpx;
}
.kv {
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

