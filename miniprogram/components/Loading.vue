<template>
  <view class="wrap">
    <view class="spinner" />
    <text class="text">{{ text }}</text>
  </view>
</template>

<script setup>
defineProps({
  text: { type: String, default: '加载中…' },
})
</script>

<style scoped>
.wrap {
  padding: 32rpx 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  color: #6b7280;
  font-size: 26rpx;
}
.spinner {
  width: 28rpx;
  height: 28rpx;
  border-radius: 50%;
  border: 3rpx solid #e5e7eb;
  border-top-color: #111827;
  animation: spin 0.9s linear infinite;
}
@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
.text {
  color: #6b7280;
}
</style>

