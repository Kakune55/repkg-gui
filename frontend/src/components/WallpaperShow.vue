<script setup>
import { ref, onMounted } from 'vue'
import { GetWallpapers, GetImgBase64 } from '../../wailsjs/go/main/App'

// 响应式数据
const wallpapers = ref([])
const loading = ref(true)
const error = ref(null)
// 存储壁纸封面的Map
const wallpaperCovers = ref(new Map())

// 定义emit事件
const emit = defineEmits(['select-wallpaper'])

// 封面请求队列与控制逻辑
const coverRequestQueue = []
let isProcessingQueue = false

// 加载壁纸列表
async function loadWallpapers() {
  try {
    loading.value = true
    error.value = null
    const result = await GetWallpapers()
    const parsedWallpapers = JSON.parse(result)
    wallpapers.value = parsedWallpapers
    loading.value = false
    
    // 初始化封面Map
    wallpaperCovers.value = new Map()
    
    // 将所有壁纸的封面请求加入队列
    parsedWallpapers.forEach(wallpaper => {
      if (wallpaper.coverPath) {
        enqueueCoverRequest(wallpaper)
      }
    })
  } catch (err) {
    console.error("加载壁纸失败:", err)
    error.value = "加载壁纸失败: " + (err.message || "未知错误")
    loading.value = false
  }
}

// 将封面请求添加到队列
function enqueueCoverRequest(wallpaper) {
  coverRequestQueue.push(wallpaper)
  processCoverQueue()
}

// 处理封面请求队列
async function processCoverQueue() {
  if (isProcessingQueue || coverRequestQueue.length === 0) return
  
  isProcessingQueue = true
  const wallpaper = coverRequestQueue.shift()
  
  try {
    const coverBase64 = await getWallpaperCover(wallpaper.coverPath)
    wallpaperCovers.value.set(wallpaper.path, coverBase64)
  } catch (err) {
    console.error(`加载壁纸封面失败 ${wallpaper.coverPath}:`, err)
  }
  
  // 延迟处理下一个请求
  setTimeout(() => {
    isProcessingQueue = false
    processCoverQueue()
  }, 0)
}

// 获取壁纸封面
async function getWallpaperCover(path) {
  try {
    const coverBase64 = await GetImgBase64(path)
    return coverBase64 || null
  } catch (err) {
    console.error(`加载壁纸封面失败 ${path}:`, err)
    return null
  }
}

// 选择壁纸，触发跳转
function selectWallpaper(wallpaper) {
  emit('select-wallpaper', wallpaper)
}

// 组件挂载时加载壁纸
onMounted(() => {
  loadWallpapers()
})
</script>

<template>
  <div class="wallpaper-container">
    <h2>壁纸列表</h2>
    
    <!-- 加载状态 -->
    <div v-if="loading" class="loading">
      加载中...
    </div>
    
    <!-- 错误信息 -->
    <div v-else-if="error" class="error">
      {{ error }}
    </div>
    
    <!-- 无壁纸提示 -->
    <div v-else-if="wallpapers.length === 0" class="no-wallpapers">
      未找到壁纸
    </div>
    
    <!-- 壁纸网格 -->
    <div v-else class="wallpaper-grid">
      <div 
        v-for="wallpaper in wallpapers" 
        :key="wallpaper.path" 
        class="wallpaper-card"
        @click="selectWallpaper(wallpaper)"
      >
        <div class="wallpaper-preview">
          <img 
            v-if="wallpaper.coverPath && wallpaperCovers.get(wallpaper.path)"
            :src="wallpaperCovers.get(wallpaper.path)"
            :alt="wallpaper.name"
            @error="$event.target.src = ''"
          />
          <div v-else class="no-preview">
            {{ wallpaper.coverPath ? '加载中...' : '无预览图' }}
          </div>
        </div>
        <div class="wallpaper-title">{{ wallpaper.name }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.wallpaper-container {
  padding: 20px;
  color: #333;
}

.loading, .error, .no-wallpapers {
  padding: 20px;
  text-align: center;
  color: #666;
}

.error {
  color: #f44336;
}

.wallpaper-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 15px;
  margin-top: 20px;
}

.wallpaper-card {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  transition: transform 0.2s, box-shadow 0.2s;
  display: flex;
  flex-direction: column;
  cursor: pointer;
  border: 2px solid transparent;
  height: 180px;
}

.wallpaper-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.12);
  border-color: #2196f3;
}

.wallpaper-preview {
  width: 100%;
  height: 124px;
  overflow: hidden;
  background-color: #eee;
  display: flex;
  align-items: center;
  justify-content: center;
}

.wallpaper-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-preview {
  color: #999;
  font-size: 14px;
}

.wallpaper-title {
  padding: 10px;
  font-weight: bold;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
