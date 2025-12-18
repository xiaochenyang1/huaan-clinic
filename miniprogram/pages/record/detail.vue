<template>
  <view class="page">
    <view class="panel" v-if="r">
      <view class="title">{{ r.department_name || '-' }}</view>
      <view class="sub">{{ r.doctor_name || '-' }} · {{ r.patient_name || '-' }}</view>

      <view class="block">
        <view class="kv"><text class="k">记录时间</text><text class="v">{{ r.created_at || r.visit_date || '-' }}</text></view>
        <view class="kv" v-if="r.appointment_no"><text class="k">预约号</text><text class="v">{{ r.appointment_no }}</text></view>
        <view class="kv" v-if="r.symptom"><text class="k">主诉</text><text class="v">{{ r.symptom }}</text></view>
        <view class="kv" v-if="r.diagnosis"><text class="k">诊断</text><text class="v">{{ r.diagnosis }}</text></view>
        <view class="kv" v-if="r.treatment"><text class="k">处理</text><text class="v">{{ r.treatment }}</text></view>
        <view class="kv" v-if="r.prescription"><text class="k">处方</text><text class="v">{{ r.prescription }}</text></view>
        <view class="kv" v-if="r.remark"><text class="k">备注</text><text class="v">{{ r.remark }}</text></view>
      </view>
    </view>

    <view v-else class="panel">
      <view class="muted">{{ loading ? '加载中…' : '未找到记录' }}</view>
    </view>
  </view>
</template>

<script setup>
import { onLoad, onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { isLoggedIn, toLoginPage } from '../../utils/auth'
import { getRecord } from '../../api/record'

const id = ref('')
const loading = ref(false)
const r = ref(null)

async function load() {
  if (!id.value) return
  loading.value = true
  try {
    r.value = await getRecord(id.value)
  } finally {
    loading.value = false
  }
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
  line-height: 1.6;
}
.muted {
  color: #6b7280;
  font-size: 26rpx;
}
</style>

