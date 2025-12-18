<template>
  <view class="page">
    <view class="panel">
      <view class="title">隐私政策</view>
      <view class="text">
        本隐私政策用于说明“华安医疗预约”小程序在提供门诊预约服务过程中，如何收集、使用、存储与保护你的个人信息。
      </view>

      <view class="h2">1. 我们收集的信息</view>
      <view class="text">
        - 账号与登录信息：用于登录与身份识别（如微信 OpenID、手机号、用户名等）。
      </view>
      <view class="text">
        - 就诊相关信息：就诊人姓名、身份证号、手机号、预约信息、就诊记录等，用于完成预约与就诊服务。
      </view>
      <view class="text">
        - 设备与日志信息：用于保障系统安全、故障排查与性能优化。
      </view>

      <view class="h2">2. 我们如何使用信息</view>
      <view class="text">
        - 提供预约挂号、取消、签到、查询等核心功能。
      </view>
      <view class="text">
        - 发送就诊提醒（需你授权订阅消息）。
      </view>

      <view class="h2">3. 信息共享与披露</view>
      <view class="text">
        我们不会向无关第三方出售你的个人信息。仅在法律法规要求或为提供服务所必需的情况下进行共享。
      </view>

      <view class="h2">4. 信息存储与保护</view>
      <view class="text">
        我们采用合理的安全措施保护你的信息，防止未授权访问、披露、篡改或毁损。
      </view>

      <view class="h2">5. 你的权利</view>
      <view class="text">
        你可在小程序内管理就诊人、查看预约与就诊记录，并可通过“设置”清理本地缓存或退出登录。
      </view>

      <view class="h2">6. 联系我们</view>
      <view class="text">
        如对本政策有疑问，请在“联系方式”页面查看客服信息。
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

