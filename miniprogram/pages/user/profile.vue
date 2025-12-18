<template>
  <view class="page">
    <view class="panel">
      <view class="panel-title">个人信息</view>

      <view class="field">
        <text class="label">昵称</text>
        <input class="input" v-model="form.nickname" placeholder="请输入昵称" maxlength="64" />
      </view>

      <view class="field">
        <text class="label">头像地址（可选）</text>
        <input class="input" v-model="form.avatar" placeholder="https://..." maxlength="512" />
      </view>

      <view class="field">
        <text class="label">手机号（可选）</text>
        <input class="input" v-model="form.phone" placeholder="11位手机号" maxlength="11" />
      </view>

      <view class="field">
        <text class="label">性别</text>
        <picker mode="selector" :range="genderLabels" :value="genderIndex" @change="onPickGender">
          <view class="picker">
            <text>{{ genderLabels[genderIndex] }}</text>
            <text class="arrow">›</text>
          </view>
        </picker>
      </view>

      <button class="btn primary" @click="submit" :disabled="loading">
        {{ loading ? '保存中…' : '保存' }}
      </button>
    </view>
  </view>
</template>

<script setup>
import { onLoad } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { getUserInfo, updateUserInfo } from '../../api/auth'
import { isLoggedIn, setUser, toLoginPage } from '../../utils/auth'

const loading = ref(false)

const genders = [
  { value: 0, label: '未知' },
  { value: 1, label: '男' },
  { value: 2, label: '女' },
]
const genderLabels = genders.map((g) => g.label)
const genderIndex = ref(0)

const form = ref({
  nickname: '',
  avatar: '',
  phone: '',
  gender: 0,
})

function onPickGender(e) {
  genderIndex.value = Number(e.detail.value || 0)
  form.value.gender = genders[genderIndex.value]?.value ?? 0
}

async function load() {
  loading.value = true
  try {
    const u = await getUserInfo()
    form.value = {
      nickname: u.nickname || '',
      avatar: u.avatar || '',
      phone: u.phone || '',
      gender: typeof u.gender === 'number' ? u.gender : 0,
    }
    const idx = genders.findIndex((g) => g.value === form.value.gender)
    genderIndex.value = idx >= 0 ? idx : 0
  } finally {
    loading.value = false
  }
}

async function submit() {
  loading.value = true
  try {
    const payload = {
      nickname: form.value.nickname || '',
      avatar: form.value.avatar || '',
      phone: form.value.phone || '',
      gender: form.value.gender ?? 0,
    }
    const u = await updateUserInfo(payload)
    setUser(u)
    uni.showToast({ title: '已保存', icon: 'none' })
    uni.navigateBack()
  } finally {
    loading.value = false
  }
}

onLoad(() => {
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
.field {
  margin-bottom: 18rpx;
}
.label {
  display: block;
  font-size: 24rpx;
  color: #374151;
  margin-bottom: 8rpx;
}
.input {
  height: 84rpx;
  border: 1rpx solid #e5e7eb;
  border-radius: 12rpx;
  padding: 0 20rpx;
  font-size: 28rpx;
  background: #fff;
}
.picker {
  height: 84rpx;
  border: 1rpx solid #e5e7eb;
  border-radius: 12rpx;
  padding: 0 20rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  font-size: 28rpx;
  color: #111827;
}
.arrow {
  color: #9ca3af;
  font-size: 34rpx;
}
.btn {
  margin-top: 10rpx;
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
</style>

