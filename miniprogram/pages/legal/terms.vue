<template>
  <view class="page">
    <view class="panel">
      <view class="title">用户协议</view>
      <view class="text">
        欢迎使用“华安医疗预约”小程序。你在使用本服务前，应仔细阅读并同意本协议。
      </view>

      <view class="h2">1. 服务内容</view>
      <view class="text">本服务提供预约挂号、预约管理、就诊人管理、就诊记录查询等功能。</view>

      <view class="h2">2. 用户义务</view>
      <view class="text">- 提供真实、准确的就诊信息。</view>
      <view class="text">- 遵守预约规则，按时就诊；多次爽约可能导致账号被限制。</view>

      <view class="h2">3. 免责声明</view>
      <view class="text">
        因不可抗力、系统维护、网络原因等导致服务中断或延迟，我们将尽力恢复，但不承担由此产生的直接或间接损失。
      </view>

      <view class="h2">4. 协议变更</view>
      <view class="text">
        我们可能适时更新协议内容，并在小程序内公示。你继续使用服务视为同意更新后的协议。
      </view>
    </view>
  </view>
</template>

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
.h2 {
  margin-top: 18rpx;
  font-size: 28rpx;
  font-weight: 700;
  color: #111827;
}
.text {
  margin-top: 10rpx;
  font-size: 26rpx;
  color: #374151;
  line-height: 1.8;
}
</style>

