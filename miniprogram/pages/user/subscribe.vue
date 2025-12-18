<template>
  <view class="page">
    <view class="panel">
      <view class="panel-title">消息提醒</view>
      <view class="muted">
        就诊提醒需要你在微信端授权订阅。授权后，系统会在就诊前一天推送提醒（如后端已开启推送并配置模板）。
      </view>

      <view class="field">
        <text class="label">订阅状态</text>
        <view class="status">{{ statusText }}</view>
      </view>

      <button class="btn primary" @click="requestSubscribe" :disabled="loading">
        {{ loading ? '处理中…' : '授权订阅就诊提醒' }}
      </button>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { isLoggedIn, toLoginPage } from '../../utils/auth'
import { WECHAT_SUBSCRIBE_TEMPLATE_IDS, STORAGE_KEYS } from '../../utils/config'
import { getStorage, setStorage } from '../../utils/storage'

const loading = ref(false)
const state = ref(getStorage(STORAGE_KEYS.subscribe) || {})

const statusText = computed(() => {
  const v = state.value?.appointmentReminder
  if (v === 'accept') return '已授权'
  if (v === 'reject') return '已拒绝'
  if (v === 'ban') return '被禁用/关闭'
  return '未授权'
})

async function requestSubscribe() {
  const tpl = WECHAT_SUBSCRIBE_TEMPLATE_IDS.appointmentReminder
  if (!tpl) {
    uni.showToast({ title: '未配置模板ID', icon: 'none' })
    return
  }
  loading.value = true
  try {
    const res = await new Promise((resolve, reject) => {
      // #ifdef MP-WEIXIN
      uni.requestSubscribeMessage({
        tmplIds: [tpl],
        success: (r) => resolve(r),
        fail: (e) => reject(e),
      })
      // #endif
      // #ifndef MP-WEIXIN
      resolve({ [tpl]: 'unsupported' })
      // #endif
    })
    const v = res?.[tpl] || 'unknown'
    state.value = { ...(state.value || {}), appointmentReminder: v }
    setStorage(STORAGE_KEYS.subscribe, state.value)
    uni.showToast({ title: v === 'accept' ? '已授权' : '已记录', icon: 'none' })
  } finally {
    loading.value = false
  }
}

onLoad(() => {
  if (!isLoggedIn()) {
    toLoginPage()
    return
  }
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
.muted {
  color: #6b7280;
  font-size: 26rpx;
  line-height: 1.7;
}
.field {
  margin-top: 18rpx;
  margin-bottom: 18rpx;
}
.label {
  display: block;
  font-size: 24rpx;
  color: #374151;
  margin-bottom: 8rpx;
}
.status {
  height: 84rpx;
  border: 1rpx solid #e5e7eb;
  border-radius: 12rpx;
  padding: 0 20rpx;
  display: flex;
  align-items: center;
  font-size: 28rpx;
  color: #111827;
  background: #fff;
}
.btn {
  margin-top: 10rpx;
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

