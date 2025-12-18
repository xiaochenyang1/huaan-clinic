<template>
  <view class="page">
    <view class="panel">
      <view class="panel-title">{{ departmentName || '医生列表' }}</view>
      <view class="search">
        <input class="input" v-model="keyword" placeholder="搜索医生姓名/职称" @confirm="load" />
        <button class="btn" @click="load">搜索</button>
      </view>
      <view v-if="loading" class="muted">加载中…</view>
      <view v-else>
        <view v-if="doctors.length === 0" class="muted">暂无医生</view>
        <view v-for="d in doctors" :key="d.id" class="row" @click="goDoctor(d)">
          <view class="row-main">
            <view class="row-title">{{ d.name }} <text class="tag" v-if="d.title">{{ d.title }}</text></view>
            <view class="row-sub">{{ d.department_name || departmentName || '' }}</view>
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
import { listDoctors } from '../../api/doctor'

const departmentId = ref('')
const departmentName = ref('')
const keyword = ref('')
const loading = ref(false)
const doctors = ref([])

function goDoctor(d) {
  uni.navigateTo({ url: `/pages/appointment/doctor-detail?id=${d.id}` })
}

async function load() {
  loading.value = true
  try {
    doctors.value = (await listDoctors({ department_id: departmentId.value, keyword: keyword.value })) || []
  } finally {
    loading.value = false
  }
}

onLoad((q) => {
  departmentId.value = q?.department_id || ''
  departmentName.value = q?.department_name ? decodeURIComponent(q.department_name) : ''
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
.search {
  display: flex;
  gap: 12rpx;
  margin-bottom: 12rpx;
}
.input {
  flex: 1;
  height: 76rpx;
  border: 1rpx solid #e5e7eb;
  border-radius: 12rpx;
  padding: 0 18rpx;
  font-size: 28rpx;
  background: #fff;
}
.btn {
  height: 76rpx;
  line-height: 76rpx;
  border-radius: 12rpx;
  border: 1rpx solid #d1d5db;
  background: #fff;
  color: #111827;
  font-size: 28rpx;
  padding: 0 18rpx;
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
.tag {
  margin-left: 10rpx;
  font-size: 22rpx;
  color: #374151;
  border: 1rpx solid #e5e7eb;
  border-radius: 999rpx;
  padding: 2rpx 10rpx;
}
.muted {
  color: #6b7280;
  font-size: 26rpx;
  padding: 12rpx 0;
}
</style>

