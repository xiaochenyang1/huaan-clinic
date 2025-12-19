<template>
  <view class="item">
    <view class="label">
      <text v-if="required" class="req">*</text>
      <text>{{ label }}</text>
    </view>
    <view class="control">
      <slot />
    </view>
    <view v-if="hint" class="hint">{{ hint }}</view>
  </view>
</template>

<script setup>
defineProps({
  label: { type: String, default: '' },
  required: { type: Boolean, default: false },
  hint: { type: String, default: '' },
})
</script>

<style scoped>
.item {
  margin-bottom: 18rpx;
}
.label {
  display: flex;
  align-items: center;
  gap: 6rpx;
  color: #374151;
  font-size: 24rpx;
  margin-bottom: 8rpx;
}
.req {
  color: #dc2626;
}
.hint {
  margin-top: 8rpx;
  color: #6b7280;
  font-size: 22rpx;
  line-height: 1.6;
}
</style>

