<template>
  <view class="page">
    <view class="header">
      <view class="h1">门诊预约</view>
      <view class="h2">请选择科室开始预约</view>
    </view>

    <view class="panel">
      <view class="panel-title">搜索</view>
      <view class="search">
        <input class="input" v-model="keyword" placeholder="搜索科室或医生" @confirm="doSearch" />
        <button class="btn" @click="doSearch">搜索</button>
      </view>
      <view v-if="searching" class="muted">搜索中…</view>
      <view v-else-if="hasSearch">
        <view class="section-title">科室</view>
        <view v-if="deptResults.length === 0" class="muted">未找到匹配科室</view>
        <view v-for="d in deptResults" :key="'d'+d.id" class="row" @click="goDept(d)">
          <view class="row-main">
            <view class="row-title">{{ d.name }}</view>
            <view class="row-sub" v-if="d.description">{{ d.description }}</view>
          </view>
          <view class="row-arrow">›</view>
        </view>

        <view class="section-title">医生</view>
        <view v-if="doctorResults.length === 0" class="muted">未找到匹配医生</view>
        <view v-for="doc in doctorResults" :key="'u'+doc.id" class="row" @click="goDoctor(doc)">
          <view class="row-main">
            <view class="row-title">{{ doc.name }} <text class="tag" v-if="doc.title">{{ doc.title }}</text></view>
            <view class="row-sub">{{ doc.department_name || '' }}</view>
          </view>
          <view class="row-arrow">›</view>
        </view>
      </view>
    </view>

    <view class="panel">
      <view class="panel-title">快捷入口</view>
      <view class="grid">
        <view class="grid-item" @click="go('/pages/appointment/department')">
          <view class="grid-text">开始预约</view>
        </view>
        <view class="grid-item" @click="switchTab('/pages/appointment/list')">
          <view class="grid-text">我的预约</view>
        </view>
        <view class="grid-item" @click="switchTab('/pages/user/index')">
          <view class="grid-text">个人中心</view>
        </view>
      </view>
    </view>

    <view class="panel">
      <view class="panel-title">公告与须知</view>
      <view class="list">
        <view class="item" @click="go('/pages/legal/notice')">
          <text class="txt">就诊须知与预约规则</text><text class="arrow">›</text>
        </view>
        <view class="item" @click="go('/pages/user/subscribe')">
          <text class="txt">消息提醒授权</text><text class="arrow">›</text>
        </view>
        <view class="item" @click="go('/pages/legal/contact')">
          <text class="txt">联系方式</text><text class="arrow">›</text>
        </view>
      </view>
    </view>

    <view class="panel" v-if="recent.departments.length || recent.doctors.length">
      <view class="panel-title">最近访问</view>
      <view v-if="recent.departments.length" class="section-title">科室</view>
      <view class="chips" v-if="recent.departments.length">
        <view class="chip" v-for="d in recent.departments" :key="'rd'+d.id" @click="goDept(d)">
          {{ d.name }}
        </view>
      </view>
      <view v-if="recent.doctors.length" class="section-title">医生</view>
      <view class="chips" v-if="recent.doctors.length">
        <view class="chip" v-for="doc in recent.doctors" :key="'rdoc'+doc.id" @click="goDoctor(doc)">
          {{ doc.name }}
        </view>
      </view>
    </view>

    <view class="panel">
      <view class="panel-title">科室列表</view>
      <view v-if="loading" class="muted">加载中…</view>
      <view v-else>
        <view v-if="departments.length === 0" class="muted">暂无科室</view>
        <view v-for="d in departments" :key="d.id" class="row" @click="goDept(d)">
          <view class="row-main">
            <view class="row-title">{{ d.name }}</view>
            <view class="row-sub" v-if="d.description">{{ d.description }}</view>
          </view>
          <view class="row-arrow">›</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import { listDepartments } from '../../api/department'
import { listDoctors } from '../../api/doctor'
import { addRecentDepartment, getRecent } from '../../utils/recent'

const loading = ref(false)
const departments = ref([])
const keyword = ref('')
const searching = ref(false)
const hasSearch = ref(false)
const deptResults = ref([])
const doctorResults = ref([])
const recent = ref(getRecent())

function go(url) {
  uni.navigateTo({ url })
}
function switchTab(url) {
  uni.switchTab({ url })
}
function goDept(d) {
  addRecentDepartment(d)
  recent.value = getRecent()
  uni.navigateTo({ url: `/pages/appointment/doctor-list?department_id=${d.id}&department_name=${encodeURIComponent(d.name)}` })
}
function goDoctor(doc) {
  uni.navigateTo({ url: `/pages/appointment/doctor-detail?id=${doc.id}` })
}

async function doSearch() {
  const k = (keyword.value || '').trim()
  if (!k) {
    hasSearch.value = false
    deptResults.value = []
    doctorResults.value = []
    return
  }
  searching.value = true
  try {
    const allDepts = departments.value.length ? departments.value : ((await listDepartments()) || [])
    departments.value = allDepts
    deptResults.value = allDepts.filter((d) => (d.name || '').includes(k)).slice(0, 8)
    doctorResults.value = ((await listDoctors({ keyword: k })) || []).slice(0, 8)
    hasSearch.value = true
  } finally {
    searching.value = false
  }
}

async function load() {
  loading.value = true
  try {
    departments.value = (await listDepartments()) || []
    recent.value = getRecent()
  } finally {
    loading.value = false
  }
}

onShow(load)
</script>

<style scoped>
.page {
  min-height: 100vh;
  background: #f6f7f9;
  padding: 24rpx;
}
.header {
  padding: 18rpx 6rpx 12rpx;
}
.h1 {
  font-size: 40rpx;
  font-weight: 700;
  color: #111827;
}
.h2 {
  margin-top: 8rpx;
  font-size: 26rpx;
  color: #6b7280;
}
.panel {
  background: #fff;
  border: 1rpx solid #e5e7eb;
  border-radius: 16rpx;
  padding: 20rpx;
  margin-top: 18rpx;
}
.panel-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #111827;
  margin-bottom: 12rpx;
}
.section-title {
  margin-top: 14rpx;
  font-size: 26rpx;
  font-weight: 600;
  color: #111827;
}
.search {
  display: flex;
  gap: 12rpx;
  margin-bottom: 10rpx;
}
.input {
  flex: 1;
  height: 76rpx;
  border: 1rpx solid #e5e7eb;
  border-radius: 12rpx;
  padding: 0 18rpx;
  font-size: 28rpx;
  background: #fff;
}
.btn {
  height: 76rpx;
  line-height: 76rpx;
  border-radius: 12rpx;
  border: 1rpx solid #d1d5db;
  background: #fff;
  color: #111827;
  font-size: 28rpx;
  padding: 0 18rpx;
}
.grid {
  display: flex;
  gap: 12rpx;
}
.grid-item {
  flex: 1;
  border: 1rpx solid #e5e7eb;
  border-radius: 12rpx;
  padding: 22rpx 16rpx;
  background: #f9fafb;
}
.grid-text {
  font-size: 28rpx;
  color: #111827;
  font-weight: 600;
}
.row {
  display: flex;
  align-items: center;
  padding: 18rpx 8rpx;
  border-top: 1rpx solid #f3f4f6;
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
.row-arrow {
  font-size: 36rpx;
  color: #9ca3af;
  padding-left: 8rpx;
}
.tag {
  margin-left: 10rpx;
  font-size: 22rpx;
  color: #374151;
  border: 1rpx solid #e5e7eb;
  border-radius: 999rpx;
  padding: 2rpx 10rpx;
}
.list {
  border: 1rpx solid #e5e7eb;
  border-radius: 16rpx;
  overflow: hidden;
}
.item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 22rpx 18rpx;
  border-top: 1rpx solid #f3f4f6;
  background: #fff;
}
.item:first-of-type {
  border-top: none;
}
.txt {
  font-size: 28rpx;
  color: #111827;
}
.arrow {
  font-size: 34rpx;
  color: #9ca3af;
}
.chips {
  display: flex;
  flex-wrap: wrap;
  gap: 10rpx;
  margin-top: 10rpx;
}
.chip {
  padding: 10rpx 14rpx;
  border: 1rpx solid #e5e7eb;
  border-radius: 999rpx;
  background: #f9fafb;
  font-size: 24rpx;
  color: #111827;
}
.muted {
  color: #6b7280;
  font-size: 26rpx;
  padding: 12rpx 0;
}
</style>
