<template>
  <view class="page">
    <view class="panel">
      <view class="panel-title">{{ id ? '编辑就诊人' : '新增就诊人' }}</view>

      <view class="field">
        <text class="label">姓名</text>
        <input class="input" v-model="form.name" placeholder="请输入姓名" />
      </view>

      <view class="field">
        <text class="label">身份证号</text>
        <input class="input" v-model="form.id_card" placeholder="请输入18位身份证号" />
      </view>

      <view class="field">
        <text class="label">手机号</text>
        <input class="input" v-model="form.phone" placeholder="请输入手机号" />
      </view>

      <view class="field">
        <text class="label">关系</text>
        <picker mode="selector" :range="relationLabels" :value="relationIndex" @change="onPickRelation">
          <view class="picker">
            <text>{{ relationLabels[relationIndex] }}</text>
            <text class="arrow">›</text>
          </view>
        </picker>
      </view>

      <view class="field row">
        <text class="label">设为默认</text>
        <switch :checked="form.is_default === 1" @change="(e) => (form.is_default = e.detail.value ? 1 : 0)" />
      </view>

      <button class="btn primary" @click="submit" :disabled="loading">
        {{ loading ? '保存中…' : '保存' }}
      </button>
    </view>
  </view>
</template>

<script setup>
import { onLoad } from '@dcloudio/uni-app'
import { computed, ref } from 'vue'
import { createPatient, getPatient, updatePatient } from '../../api/patient'
import { isLoggedIn, toLoginPage } from '../../utils/auth'

const id = ref('')
const loading = ref(false)

const relations = [
  { value: 'self', label: '本人' },
  { value: 'parent', label: '父母' },
  { value: 'child', label: '子女' },
  { value: 'spouse', label: '配偶' },
  { value: 'other', label: '其他' },
]

const relationLabels = relations.map((r) => r.label)
const relationIndex = ref(0)

const form = ref({
  name: '',
  id_card: '',
  phone: '',
  relation: 'self',
  is_default: 0,
})

const relationValue = computed(() => relations[relationIndex.value]?.value || 'self')

function onPickRelation(e) {
  relationIndex.value = Number(e.detail.value || 0)
  form.value.relation = relationValue.value
}

async function load() {
  if (!id.value) return
  const p = await getPatient(id.value)
  form.value = {
    name: p.name || '',
    id_card: p.id_card || '',
    phone: p.phone || '',
    relation: p.relation || 'self',
    is_default: p.is_default || 0,
  }
  const idx = relations.findIndex((r) => r.value === form.value.relation)
  relationIndex.value = idx >= 0 ? idx : 0
}

async function submit() {
  if (!form.value.name || form.value.name.length < 2) {
    uni.showToast({ title: '姓名至少2个字', icon: 'none' })
    return
  }
  if (!form.value.id_card || form.value.id_card.length !== 18) {
    uni.showToast({ title: '身份证号长度应为18位', icon: 'none' })
    return
  }
  if (!form.value.phone || form.value.phone.length !== 11) {
    uni.showToast({ title: '手机号长度应为11位', icon: 'none' })
    return
  }

  loading.value = true
  try {
    const payload = { ...form.value, relation: relationValue.value }
    if (id.value) await updatePatient(id.value, payload)
    else await createPatient(payload)
    uni.navigateBack()
  } finally {
    loading.value = false
  }
}

onLoad(async (q) => {
  if (!isLoggedIn()) {
    toLoginPage()
    return
  }
  id.value = q?.id || ''
  await load()
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
.row {
  display: flex;
  align-items: center;
  justify-content: space-between;
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

