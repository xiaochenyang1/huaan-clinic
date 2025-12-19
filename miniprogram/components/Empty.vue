<template>
  <view class="wrap">
    <view class="icon">∅</view>
    <text class="text">{{ text }}</text>
  </view>
</template>

<script setup>
defineProps({
  text: { type: String, default: '暂无数据' },
})
</script>

<style scoped>
.wrap {
  padding: 48rpx 0;
  text-align: center;
}
.icon {
  width: 88rpx;
  height: 88rpx;
  margin: 0 auto 14rpx;
  border-radius: 999rpx;
  border: 1rpx solid #e5e7eb;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 44rpx;
  color: #9ca3af;
}
.text {
  color: #6b7280;
  font-size: 26rpx;
}
</style>

