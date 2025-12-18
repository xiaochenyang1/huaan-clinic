<template>
  <view class="page">
    <view class="panel" v-if="a">
      <view class="title">{{ a.department_name || '-' }}</view>
      <view class="sub">{{ a.doctor_name || '-' }} · {{ a.patient_name || '-' }}</view>

      <view class="block">
        <view class="kv"><text class="k">预约号</text><text class="v">{{ a.appointment_no || '-' }}</text></view>
        <view class="kv"><text class="k">日期</text><text class="v">{{ a.appointment_date || '-' }}</text></view>
        <view class="kv"><text class="k">时段</text><text class="v">{{ a.period_name || '-' }}</text></view>
        <view class="kv"><text class="k">号序</text><text class="v">{{ a.slot_number ?? '-' }}</text></view>
        <view class="kv"><text class="k">状态</text><text class="v">{{ a.status_name || a.status || '-' }}</text></view>
      </view>

      <view class="actions">
        <button v-if="a.status === 'pending'" class="btn primary" @click="checkin">签到</button>
        <button v-if="a.status === 'pending'" class="btn danger" @click="openCancel">取消预约</button>
      </view>

      <view class="more">
        <text class="more-link" @click="goNotice">查看就诊须知与预约规则</text>
      </view>
    </view>

    <view v-else class="panel">
      <view class="muted">{{ loading ? '加载中…' : '未找到预约' }}</view>
    </view>

    <view v-if="cancelVisible" class="modal-mask" @click="cancelVisible = false">
      <view class="modal" @click.stop>
        <view class="modal-title">取消预约</view>
        <textarea class="textarea" v-model="cancelReason" placeholder="请输入取消原因（必填）" maxlength="256" />
        <view class="modal-actions">
          <button class="btn" @click="cancelVisible = false">返回</button>
          <button class="btn danger" @click="doCancel" :disabled="submitting">确认取消</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { onLoad, onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { isLoggedIn, toLoginPage } from '../../utils/auth'
import { cancelAppointment, checkinAppointment, getAppointment } from '../../api/appointment'

const id = ref('')
const loading = ref(false)
const a = ref(null)

const cancelVisible = ref(false)
const cancelReason = ref('')
const submitting = ref(false)

function openCancel() {
  cancelReason.value = ''
  cancelVisible.value = true
}

async function load() {
  if (!id.value) return
  loading.value = true
  try {
    a.value = await getAppointment(id.value)
  } finally {
    loading.value = false
  }
}

async function checkin() {
  if (!id.value) return
  await checkinAppointment(id.value)
  await load()
}

async function doCancel() {
  if (!cancelReason.value || cancelReason.value.length < 2) {
    uni.showToast({ title: '请输入取消原因', icon: 'none' })
    return
  }
  submitting.value = true
  try {
    await cancelAppointment(id.value, { reason: cancelReason.value })
    cancelVisible.value = false
    await load()
  } finally {
    submitting.value = false
  }
}

function goNotice() {
  uni.navigateTo({ url: '/pages/legal/notice' })
}

onLoad((q) => {
  id.value = q?.id || ''
})

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
.sub {
  margin-top: 8rpx;
  font-size: 26rpx;
  color: #6b7280;
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
.actions {
  margin-top: 18rpx;
  display: flex;
  gap: 12rpx;
}
.btn {
  flex: 1;
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
.btn.danger {
  border: 1rpx solid #991b1b;
  background: #fff;
  color: #991b1b;
}
.muted {
  color: #6b7280;
  font-size: 26rpx;
}
.modal-mask {
  position: fixed;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24rpx;
}
.modal {
  width: 100%;
  background: #fff;
  border-radius: 16rpx;
  padding: 20rpx;
  border: 1rpx solid #e5e7eb;
}
.modal-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #111827;
  margin-bottom: 12rpx;
}
.textarea {
  width: 100%;
  min-height: 160rpx;
  border: 1rpx solid #e5e7eb;
  border-radius: 12rpx;
  padding: 16rpx;
  font-size: 28rpx;
  box-sizing: border-box;
  background: #fff;
}
.modal-actions {
  margin-top: 14rpx;
  display: flex;
  gap: 12rpx;
}
.more {
  margin-top: 16rpx;
}
.more-link {
  font-size: 24rpx;
  color: #111827;
  text-decoration: underline;
}
</style>
