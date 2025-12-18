<template>
  <view class="page">
    <view class="panel">
      <view class="panel-title">预约信息</view>
      <view class="kv">
        <text class="k">日期</text><text class="v">{{ scheduleDate || '-' }}</text>
      </view>
      <view class="kv">
        <text class="k">时段</text><text class="v">{{ periodName || '-' }}</text>
      </view>
      <view class="kv">
        <text class="k">医生</text><text class="v">{{ doctorName || '-' }}</text>
      </view>
      <view class="kv">
        <text class="k">科室</text><text class="v">{{ departmentName || '-' }}</text>
      </view>
    </view>

    <view class="panel">
      <view class="panel-title">就诊人</view>
      <view v-if="patients.length === 0" class="muted">
        尚未添加就诊人
      </view>
      <view v-else class="picker-wrap">
        <picker mode="selector" :range="patientNames" :value="patientIndex" @change="onPickPatient">
          <view class="picker">
            <text>{{ patientNames[patientIndex] }}</text>
            <text class="arrow">›</text>
          </view>
        </picker>
      </view>
      <button class="btn" @click="goPatientList">管理就诊人</button>
    </view>

    <view class="panel">
      <view class="panel-title">症状描述（可选）</view>
      <textarea class="textarea" v-model="symptom" placeholder="请简要描述症状（选填）" maxlength="512" />
      <view class="agree" @click="agreed = !agreed">
        <view class="box" :class="{ on: agreed }"></view>
        <view class="agree-text">
          我已阅读并同意
          <text class="link" @click.stop="goNotice">《就诊须知与预约规则》</text>
        </view>
      </view>
      <button class="btn primary" @click="submit" :disabled="loading">
        {{ loading ? '提交中…' : '确认预约' }}
      </button>
    </view>
  </view>
</template>

<script setup>
import { onLoad, onShow } from '@dcloudio/uni-app'
import { computed, ref } from 'vue'
import { isLoggedIn, toLoginPage } from '../../utils/auth'
import { listPatients } from '../../api/patient'
import { getIdempotentToken } from '../../api/token'
import { createAppointment } from '../../api/appointment'

const scheduleId = ref('')
const scheduleDate = ref('')
const periodName = ref('')
const doctorName = ref('')
const departmentName = ref('')

const patients = ref([])
const patientIndex = ref(0)
const symptom = ref('')
const loading = ref(false)
const agreed = ref(false)

const patientNames = computed(() => patients.value.map((p) => `${p.name}（${p.relation_name || p.relation || ''}）`))

function goPatientList() {
  uni.navigateTo({ url: '/pages/user/patient-list' })
}

function onPickPatient(e) {
  patientIndex.value = Number(e.detail.value || 0)
}

function goNotice() {
  uni.navigateTo({ url: '/pages/legal/notice' })
}

async function loadPatients() {
  patients.value = (await listPatients()) || []
  const idx = patients.value.findIndex((p) => p.is_default === 1)
  patientIndex.value = idx >= 0 ? idx : 0
}

async function submit() {
  if (!scheduleId.value) {
    uni.showToast({ title: '排班信息缺失', icon: 'none' })
    return
  }
  if (patients.value.length === 0) {
    uni.showToast({ title: '请先添加就诊人', icon: 'none' })
    return
  }
  const patient = patients.value[patientIndex.value]
  if (!patient?.id) {
    uni.showToast({ title: '请选择就诊人', icon: 'none' })
    return
  }
  if (!agreed.value) {
    uni.showToast({ title: '请先阅读并同意就诊须知', icon: 'none' })
    return
  }

  loading.value = true
  try {
    const tokenData = await getIdempotentToken()
    const idempotentToken = tokenData?.token || tokenData?.idempotent_token || tokenData
    const apt = await createAppointment({
      idempotent_token: idempotentToken,
      schedule_id: Number(scheduleId.value),
      patient_id: Number(patient.id),
      symptom: symptom.value || '',
    })
    uni.redirectTo({ url: `/pages/appointment/success?appointment_id=${apt.id}` })
  } finally {
    loading.value = false
  }
}

onLoad((q) => {
  scheduleId.value = q?.schedule_id || ''
  scheduleDate.value = q?.schedule_date ? decodeURIComponent(q.schedule_date) : ''
  periodName.value = q?.period_name ? decodeURIComponent(q.period_name) : ''
  doctorName.value = q?.doctor_name ? decodeURIComponent(q.doctor_name) : ''
  departmentName.value = q?.department_name ? decodeURIComponent(q.department_name) : ''
})

onShow(async () => {
  if (!isLoggedIn()) {
    toLoginPage()
    return
  }
  await loadPatients()
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
  margin-bottom: 18rpx;
}
.panel-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #111827;
  margin-bottom: 12rpx;
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
.muted {
  color: #6b7280;
  font-size: 26rpx;
  margin-bottom: 12rpx;
}
.picker {
  height: 84rpx;
  border: 1rpx solid #e5e7eb;
  border-radius: 12rpx;
  padding: 0 18rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  font-size: 28rpx;
  color: #111827;
  margin-bottom: 12rpx;
}
.arrow {
  color: #9ca3af;
  font-size: 34rpx;
}
.btn {
  height: 84rpx;
  line-height: 84rpx;
  border-radius: 12rpx;
  border: 1rpx solid #d1d5db;
  background: #fff;
  color: #111827;
  font-size: 28rpx;
  margin-top: 10rpx;
}
.btn.primary {
  border: 1rpx solid #111827;
  background: #111827;
  color: #fff;
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
.agree {
  margin-top: 14rpx;
  display: flex;
  align-items: center;
  gap: 12rpx;
}
.box {
  width: 32rpx;
  height: 32rpx;
  border: 2rpx solid #9ca3af;
  border-radius: 6rpx;
  background: #fff;
}
.box.on {
  border-color: #111827;
  background: #111827;
}
.agree-text {
  font-size: 24rpx;
  color: #374151;
}
.link {
  color: #111827;
  text-decoration: underline;
}
</style>
