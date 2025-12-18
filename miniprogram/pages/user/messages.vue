<template>
  <view class="page">
    <view class="panel">
      <view class="panel-title">消息中心</view>
      <view class="muted">当前为本地版消息中心（后续可接入服务端通知列表）。</view>

      <view class="card">
        <view class="row">
          <text class="k">就诊提醒订阅</text>
          <text class="v">{{ subscribeText }}</text>
        </view>
        <view class="row">
          <text class="k">说明</text>
          <text class="v">授权后可收到就诊前一天提醒（需后端配置模板并开启推送）。</text>
        </view>
        <button class="btn" @click="goSubscribe">去授权</button>
      </view>

      <view class="card">
        <view class="row">
          <text class="k">常用入口</text>
          <text class="v">规则/联系/隐私</text>
        </view>
        <view class="links">
          <text class="link" @click="go('/pages/legal/notice')">就诊须知</text>
          <text class="link" @click="go('/pages/legal/contact')">联系方式</text>
          <text class="link" @click="go('/pages/legal/privacy')">隐私政策</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { STORAGE_KEYS } from '../../utils/config'
import { getStorage } from '../../utils/storage'

const state = ref(getStorage(STORAGE_KEYS.subscribe) || {})

const subscribeText = computed(() => {
  const v = state.value?.appointmentReminder
  if (v === 'accept') return '已授权'
  if (v === 'reject') return '已拒绝'
  if (v === 'ban') return '被禁用/关闭'
  return '未授权'
})

function go(url) {
  uni.navigateTo({ url })
}
function goSubscribe() {
  uni.navigateTo({ url: '/pages/user/subscribe' })
}

onShow(() => {
  state.value = getStorage(STORAGE_KEYS.subscribe) || {}
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
.muted {
  color: #6b7280;
  font-size: 26rpx;
  line-height: 1.7;
  margin-bottom: 14rpx;
}
.card {
  background: #fff;
  border: 1rpx solid #e5e7eb;
  border-radius: 16rpx;
  padding: 20rpx;
  margin-bottom: 14rpx;
}
.row {
  display: flex;
  padding: 10rpx 0;
}
.k {
  width: 180rpx;
  color: #6b7280;
  font-size: 26rpx;
}
.v {
  flex: 1;
  color: #111827;
  font-size: 26rpx;
  line-height: 1.6;
}
.btn {
  margin-top: 12rpx;
  height: 80rpx;
  line-height: 80rpx;
  border-radius: 12rpx;
  border: 1rpx solid #d1d5db;
  background: #fff;
  color: #111827;
  font-size: 28rpx;
}
.links {
  margin-top: 10rpx;
  display: flex;
  gap: 16rpx;
  flex-wrap: wrap;
}
.link {
  font-size: 26rpx;
  color: #111827;
  text-decoration: underline;
}
</style>

