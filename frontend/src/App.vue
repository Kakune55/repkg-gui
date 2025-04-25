<script setup>
import { ref } from 'vue'
import WallpaperShow from './components/WallpaperShow.vue'
import WallpaperDetail from './components/WallpaperDetail.vue'

const currentView = ref('list')
const selectedWallpaper = ref(null)

// 新增主题状态
const theme = ref('auto')
const themeText = {
  auto: '自动',
  light: '明亮',
  dark: '黑暗'
}
function toggleTheme() {
  if (theme.value === 'auto') theme.value = 'light'
  else if (theme.value === 'light') theme.value = 'dark'
  else theme.value = 'auto'
}

function showWallpaperDetail(wallpaper) {
  selectedWallpaper.value = wallpaper
  currentView.value = 'detail'
}

function backToList() {
  currentView.value = 'list'
}
</script>

<template>
  <s-page :theme="theme">
    <WallpaperShow 
      v-if="currentView === 'list'" 
      @select-wallpaper="showWallpaperDetail" 
    />
    <WallpaperDetail 
      v-else-if="currentView === 'detail'" 
      :wallpaper="selectedWallpaper" 
      @back-to-list="backToList" 
    />
    <!-- 悬浮切换主题按钮 -->
    <s-fab class="theme-switch-btn" @click="toggleTheme">
      切换主题：{{ themeText[theme] }}
    </s-fab>
  </s-page>
</template>

<style>
body {
  margin: 0;
  font-family: "Nunito", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto",
  "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
  sans-serif;
}

#app {
  min-height: 100vh;
}

.theme-switch-btn {
  position: fixed;
  right: 30px;
  bottom: 30px;
  z-index: 9999;
  padding: 10px 22px;
  cursor: pointer;
  font-size: 15px;
  user-select: none;
  transition: background 0.2s, color 0.2s;
}
</style>
