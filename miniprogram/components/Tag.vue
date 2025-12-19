<template>
  <text class="tag" :class="type">{{ text }}</text>
</template>

<script setup>
defineProps({
  text: { type: String, default: '' },
  type: { type: String, default: 'default' },
})
</script>

<style scoped>
.tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 4rpx 12rpx;
  border-radius: 999rpx;
  font-size: 22rpx;
  border: 1rpx solid #e5e7eb;
  background: #f9fafb;
  color: #374151;
}
.tag.success {
  border-color: #bbf7d0;
  background: #f0fdf4;
  color: #166534;
}
.tag.warning {
  border-color: #fde68a;
  background: #fffbeb;
  color: #92400e;
}
.tag.danger {
  border-color: #fecaca;
  background: #fef2f2;
  color: #991b1b;
}
.tag.info {
  border-color: #bfdbfe;
  background: #eff6ff;
  color: #1d4ed8;
}
</style>

