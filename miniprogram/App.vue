<script setup>
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'
import { bootstrapAuth } from './utils/auth'

onLaunch(async () => {
  await bootstrapAuth()
})

onShow(() => {})
onHide(() => {})
</script>

<style>
/* 每个页面公共 css */
</style>
