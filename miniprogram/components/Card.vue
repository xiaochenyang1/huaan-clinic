<template>
  <view class="card" :style="{ padding }">
    <slot />
  </view>
</template>

<script setup>
defineProps({
  padding: {
    type: String,
    default: '20rpx',
  },
})
</script>

<style scoped>
.card {
  background: #fff;
  border: 1rpx solid #e5e7eb;
  border-radius: 16rpx;
}
</style>

