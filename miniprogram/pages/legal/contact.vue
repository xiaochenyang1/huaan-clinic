<template>
  <view class="page">
    <view class="panel">
      <view class="title">联系方式</view>
      <view class="text">如需帮助，请通过以下方式联系我们：</view>

      <view class="card">
        <view class="row">
          <text class="k">客服热线</text>
          <text class="v selectable" @click="copy('400-000-0000')">400-000-0000</text>
        </view>
        <view class="row">
          <text class="k">服务时间</text>
          <text class="v">周一至周日 08:00-18:00</text>
        </view>
        <view class="row">
          <text class="k">邮箱</text>
          <text class="v selectable" @click="copy('support@huaan-medical.com')">support@huaan-medical.com</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
function copy(text) {
  uni.setClipboardData({
    data: text,
    success: () => uni.showToast({ title: '已复制', icon: 'none' }),
  })
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
  padding: 24rpx;
}
.title {
  font-size: 36rpx;
  font-weight: 700;
  color: #111827;
  margin-bottom: 12rpx;
}
.text {
  margin-top: 10rpx;
  font-size: 26rpx;
  color: #374151;
  line-height: 1.8;
}
.card {
  margin-top: 18rpx;
  border: 1rpx solid #e5e7eb;
  border-radius: 16rpx;
  padding: 16rpx 18rpx;
  background: #f9fafb;
}
.row {
  display: flex;
  padding: 12rpx 0;
}
.k {
  width: 160rpx;
  color: #6b7280;
  font-size: 26rpx;
}
.v {
  flex: 1;
  color: #111827;
  font-size: 26rpx;
}
.selectable {
  text-decoration: underline;
}
</style>

