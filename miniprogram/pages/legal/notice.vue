<template>
  <view class="page">
    <view class="panel">
      <view class="title">就诊须知与预约规则</view>
      <view class="text">
        为保障就诊秩序，请在预约前仔细阅读以下规则。具体以医院现场与系统提示为准。
      </view>

      <view class="h2">1. 预约规则</view>
      <view class="text">- 预约成功后可在“我的预约”中查看详情。</view>
      <view class="text">- 同一时段重复预约将被拦截。</view>
      <view class="text">- 取消预约需遵循系统限制（如就诊当天可能不可取消）。</view>

      <view class="h2">2. 签到规则</view>
      <view class="text">- 请在就诊时间范围内完成签到，逾期可能无法就诊。</view>

      <view class="h2">3. 爽约与限制</view>
      <view class="text">
        - 未按时就诊且未按规则取消的，可能被记为爽约；多次爽约账号可能被限制预约。
      </view>

      <view class="h2">4. 资料准备</view>
      <view class="text">- 请携带有效身份证件、医保卡（如有）等材料。</view>
      <view class="text">- 儿童或特殊人群就诊请遵循现场指引。</view>

      <view class="h2">5. 温馨提示</view>
      <view class="text">- 如遇系统问题，请在“联系方式”中获取客服渠道。</view>
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

