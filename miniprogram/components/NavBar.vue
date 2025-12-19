<template>
  <view class="nav">
    <view class="left">
      <view v-if="back" class="back" @click="goBack">‹</view>
      <view class="title">{{ title }}</view>
    </view>
    <view class="right">
      <slot name="right" />
    </view>
  </view>
</template>

<script setup>
defineProps({
  title: { type: String, default: '' },
  back: { type: Boolean, default: true },
})

function goBack() {
  const pages = getCurrentPages()
  if (pages && pages.length > 1) {
    uni.navigateBack()
    return
  }
  uni.switchTab({ url: '/pages/index/index' })
}
</script>

<style scoped>
.nav {
  height: 96rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20rpx;
  border-bottom: 1rpx solid #e5e7eb;
  background: #fff;
}
.left {
  display: flex;
  align-items: center;
  gap: 12rpx;
}
.back {
  width: 64rpx;
  height: 64rpx;
  border-radius: 999rpx;
  border: 1rpx solid #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40rpx;
  color: #111827;
  background: #f9fafb;
}
.title {
  font-size: 32rpx;
  font-weight: 600;
  color: #111827;
}
.right {
  display: flex;
  align-items: center;
}
</style>

