<script setup>
import { ref } from 'vue'
import WallpaperShow from './components/WallpaperShow.vue'
import WallpaperDetail from './components/WallpaperDetail.vue'

const currentView = ref('list')
const selectedWallpaper = ref(null)

function showWallpaperDetail(wallpaper) {
  selectedWallpaper.value = wallpaper
  currentView.value = 'detail'
}

function backToList() {
  currentView.value = 'list'
}
</script>

<template>
  <div>
    <WallpaperShow 
      v-if="currentView === 'list'" 
      @select-wallpaper="showWallpaperDetail" 
    />
    <WallpaperDetail 
      v-else-if="currentView === 'detail'" 
      :wallpaper="selectedWallpaper" 
      @back-to-list="backToList" 
    />
  </div>
</template>

<style>
body {
  margin: 0;
  font-family: "Nunito", -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto",
  "Oxygen", "Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
  sans-serif;
  /* background-color: #f5f5f5; */
}

#app {
  min-height: 100vh;
}
</style>
