<script setup>
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'
import { useUserStore } from './store'

onLaunch(async () => {
  const userStore = useUserStore()
  await userStore.bootstrap()
})

onShow(() => {})
onHide(() => {})
</script>

<style>
/* 每个页面公共 css */
</style>
