<template>
  <view class="page">
    <view class="panel">
      <view class="panel-title">设置</view>

      <view class="list">
        <view class="item" @click="go('/pages/legal/privacy')">
          <text class="txt">隐私政策</text><text class="arrow">›</text>
        </view>
        <view class="item" @click="go('/pages/legal/terms')">
          <text class="txt">用户协议</text><text class="arrow">›</text>
        </view>
        <view class="item" @click="go('/pages/legal/about')">
          <text class="txt">关于我们</text><text class="arrow">›</text>
        </view>
        <view class="item" @click="go('/pages/legal/contact')">
          <text class="txt">联系方式</text><text class="arrow">›</text>
        </view>
      </view>

      <view class="list">
        <view class="item" @click="clearCache">
          <text class="txt">清理缓存</text>
          <text class="muted">{{ cacheHint }}</text>
        </view>
      </view>

      <button class="btn danger" @click="confirmLogout">退出登录</button>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { logout } from '../../utils/auth'

const cacheHint = ref('仅清理本地存储')

function go(url) {
  uni.navigateTo({ url })
}

async function clearCache() {
  const ok = await new Promise((resolve) => {
    uni.showModal({
      title: '清理缓存',
      content: '将清理本地缓存与登录信息，需要重新登录。',
      success: (res) => resolve(res.confirm),
      fail: () => resolve(false),
    })
  })
  if (!ok) return
  uni.clearStorageSync()
  cacheHint.value = '已清理'
  uni.showToast({ title: '已清理', icon: 'none' })
}

async function confirmLogout() {
  const ok = await new Promise((resolve) => {
    uni.showModal({
      title: '退出登录',
      content: '确认退出当前账号？',
      success: (res) => resolve(res.confirm),
      fail: () => resolve(false),
    })
  })
  if (!ok) return
  await logout()
}
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
.list {
  border: 1rpx solid #e5e7eb;
  border-radius: 16rpx;
  overflow: hidden;
  margin-bottom: 18rpx;
}
.item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 22rpx 18rpx;
  border-top: 1rpx solid #f3f4f6;
  background: #fff;
}
.item:first-of-type {
  border-top: none;
}
.txt {
  font-size: 28rpx;
  color: #111827;
}
.muted {
  font-size: 24rpx;
  color: #6b7280;
}
.arrow {
  font-size: 34rpx;
  color: #9ca3af;
}
.btn {
  height: 84rpx;
  line-height: 84rpx;
  border-radius: 12rpx;
  border: 1rpx solid #d1d5db;
  background: #fff;
  color: #111827;
  font-size: 28rpx;
}
.btn.danger {
  border: 1rpx solid #991b1b;
  color: #991b1b;
}
</style>

