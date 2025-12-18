<template>
  <view class="page">
    <view class="panel">
      <view class="panel-title">科室列表</view>
      <view v-if="loading" class="muted">加载中…</view>
      <view v-else>
        <view v-if="departments.length === 0" class="muted">暂无科室</view>
        <view v-for="d in departments" :key="d.id" class="row" @click="goDept(d)">
          <view class="row-main">
            <view class="row-title">{{ d.name }}</view>
            <view class="row-sub" v-if="d.description">{{ d.description }}</view>
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
import { listDepartments } from '../../api/department'

const loading = ref(false)
const departments = ref([])

function goDept(d) {
  uni.navigateTo({ url: `/pages/appointment/doctor-list?department_id=${d.id}&department_name=${encodeURIComponent(d.name)}` })
}

async function load() {
  loading.value = true
  try {
    departments.value = (await listDepartments()) || []
  } finally {
    loading.value = false
  }
}

onLoad(load)
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

