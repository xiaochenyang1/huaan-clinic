<template>
  <view class="page">
    <view class="panel">
      <view class="panel-title">就诊记录</view>
      <view v-if="loading" class="muted">加载中…</view>
      <view v-else>
        <view v-if="items.length === 0" class="muted">暂无记录</view>
        <view v-for="r in items" :key="r.id" class="card" @click="goDetail(r)">
          <view class="card-title">
            <text class="strong">{{ r.department_name || '-' }}</text>
            <text class="badge">{{ r.created_at || r.visit_date || '-' }}</text>
          </view>
          <view class="card-sub">{{ r.doctor_name || '-' }} · {{ r.patient_name || '-' }}</view>
          <view class="card-sub" v-if="r.diagnosis || r.summary">{{ r.diagnosis || r.summary }}</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { isLoggedIn, toLoginPage } from '../../utils/auth'
import { listRecords } from '../../api/record'

const loading = ref(false)
const items = ref([])

function goDetail(r) {
  uni.navigateTo({ url: `/pages/record/detail?id=${r.id}` })
}

async function load() {
  loading.value = true
  try {
    items.value = (await listRecords()) || []
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
}
.panel {
  background: transparent;
}
.panel-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #111827;
  margin: 8rpx 0 12rpx;
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
</style>

