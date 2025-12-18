<template>
  <view class="page">
    <view class="panel">
      <view class="panel-title">就诊人列表</view>
      <view v-if="loading" class="muted">加载中…</view>
      <view v-else>
        <view v-if="items.length === 0" class="muted">暂无就诊人</view>
        <view v-for="p in items" :key="p.id" class="row">
          <view class="row-main" @click="edit(p)">
            <view class="row-title">
              {{ p.name }}
              <text v-if="p.is_default === 1" class="tag">默认</text>
            </view>
            <view class="row-sub">{{ p.relation_name || p.relation || '' }} · {{ p.phone || '' }}</view>
          </view>
          <button class="mini danger" @click="remove(p)">删除</button>
        </view>
      </view>
      <button class="btn primary" @click="create">新增就诊人</button>
    </view>
  </view>
</template>

<script setup>
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { isLoggedIn, toLoginPage } from '../../utils/auth'
import { deletePatient, listPatients } from '../../api/patient'

const loading = ref(false)
const items = ref([])

function create() {
  uni.navigateTo({ url: '/pages/user/patient-edit' })
}
function edit(p) {
  uni.navigateTo({ url: `/pages/user/patient-edit?id=${p.id}` })
}

async function remove(p) {
  const ok = await new Promise((resolve) => {
    uni.showModal({
      title: '确认删除',
      content: `删除就诊人：${p.name}`,
      success: (res) => resolve(res.confirm),
      fail: () => resolve(false),
    })
  })
  if (!ok) return
  await deletePatient(p.id)
  await load()
}

async function load() {
  loading.value = true
  try {
    items.value = (await listPatients()) || []
  } finally {
    loading.value = false
  }
}

onShow(() => {
  if (!isLoggedIn()) {
    toLoginPage()
    return
  }
  load()
})
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
  padding: 20rpx;
}
.panel-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #111827;
  margin-bottom: 12rpx;
}
.row {
  display: flex;
  align-items: center;
  padding: 16rpx 8rpx;
  border-top: 1rpx solid #f3f4f6;
  gap: 12rpx;
}
.row:first-of-type {
  border-top: none;
}
.row-main {
  flex: 1;
}
.row-title {
  font-size: 30rpx;
  color: #111827;
}
.row-sub {
  margin-top: 6rpx;
  font-size: 24rpx;
  color: #6b7280;
}
.tag {
  margin-left: 10rpx;
  font-size: 22rpx;
  color: #111827;
  border: 1rpx solid #e5e7eb;
  border-radius: 999rpx;
  padding: 2rpx 10rpx;
  background: #f9fafb;
}
.mini {
  height: 64rpx;
  line-height: 64rpx;
  border-radius: 10rpx;
  border: 1rpx solid #d1d5db;
  background: #fff;
  color: #111827;
  font-size: 26rpx;
  padding: 0 16rpx;
}
.mini.danger {
  border-color: #991b1b;
  color: #991b1b;
}
.btn {
  margin-top: 16rpx;
  height: 84rpx;
  line-height: 84rpx;
  border-radius: 12rpx;
  border: 1rpx solid #d1d5db;
  background: #fff;
  color: #111827;
  font-size: 28rpx;
}
.btn.primary {
  border: 1rpx solid #111827;
  background: #111827;
  color: #fff;
}
.muted {
  color: #6b7280;
  font-size: 26rpx;
  padding: 12rpx 0;
}
</style>

