<template>
  <view v-if="modelValue" class="mask" @click="onMaskClick">
    <view class="dialog" @click.stop>
      <view class="header">
        <text class="title">{{ title }}</text>
      </view>
      <view class="body">
        <slot />
      </view>
      <view class="footer">
        <button v-if="showCancel" class="btn" @click="cancel">{{ cancelText }}</button>
        <button class="btn primary" @click="confirm">{{ okText }}</button>
      </view>
    </view>
  </view>
</template>

<script setup>
const props = defineProps({
  modelValue: { type: Boolean, default: false },
  title: { type: String, default: '' },
  showCancel: { type: Boolean, default: true },
  okText: { type: String, default: '确定' },
  cancelText: { type: String, default: '取消' },
  closeOnMask: { type: Boolean, default: true },
})

const emit = defineEmits(['update:modelValue', 'confirm', 'cancel'])

function close() {
  emit('update:modelValue', false)
}

function confirm() {
  emit('confirm')
  close()
}

function cancel() {
  emit('cancel')
  close()
}

function onMaskClick() {
  if (!props.closeOnMask) return
  cancel()
}
</script>

<style scoped>
.mask {
  position: fixed;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24rpx;
  z-index: 999;
}
.dialog {
  width: 100%;
  max-width: 640rpx;
  background: #fff;
  border-radius: 16rpx;
  overflow: hidden;
}
.header {
  padding: 20rpx;
  border-bottom: 1rpx solid #e5e7eb;
}
.title {
  font-size: 30rpx;
  font-weight: 600;
  color: #111827;
}
.body {
  padding: 20rpx;
}
.footer {
  padding: 16rpx 20rpx 20rpx;
  display: flex;
  gap: 12rpx;
  justify-content: flex-end;
}
.btn {
  height: 72rpx;
  line-height: 72rpx;
  padding: 0 22rpx;
  border-radius: 12rpx;
  border: 1rpx solid #d1d5db;
  background: #fff;
  color: #111827;
  font-size: 26rpx;
}
.btn.primary {
  border: 1rpx solid #111827;
  background: #111827;
  color: #fff;
}
</style>

