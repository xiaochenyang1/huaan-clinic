<template>
  <view class="page">
    <view class="panel">
      <view class="title">关于我们</view>
      <view class="text">
        “华安医疗预约”致力于提供便捷、可靠的门诊预约服务，帮助用户更高效地完成挂号、管理就诊人信息与查看就诊记录。
      </view>

      <view class="h2">版本信息</view>
      <view class="text">版本：1.0.0</view>

      <view class="h2">功能范围</view>
      <view class="text">- 预约挂号与预约管理</view>
      <view class="text">- 就诊人管理</view>
      <view class="text">- 就诊记录查询</view>
      <view class="text">- 就诊提醒（需订阅授权）</view>
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

