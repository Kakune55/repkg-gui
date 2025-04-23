<script setup>
import { ref, onMounted, computed, onBeforeUnmount } from 'vue'
import { GetWallpapers, GetImgBase64 , SelectBaseDir} from '../../wailsjs/go/main/App'

// 响应式数据
const wallpapers = ref([])
const loading = ref(true)
const error = ref(null)
// 存储壁纸封面的Map
const wallpaperCovers = ref(new Map())

// 筛选开关：是否只显示pkgPath不为空的壁纸
const showOnlyPkg = ref(false)
// contentrating多选筛选
const contentratingFilter = ref(['Everyone', 'Questionable'])
const showContentratingDropdown = ref(false)
const contentratingOptions = ['Everyone', 'Questionable', 'Mature']

// 点击外部关闭下拉
function handleClickOutside(event) {
  const dropdown = document.getElementById('contentrating-dropdown')
  if (dropdown && !dropdown.contains(event.target)) {
    showContentratingDropdown.value = false
  }
}
onMounted(() => {
  document.addEventListener('mousedown', handleClickOutside)
})
onBeforeUnmount(() => {
  document.removeEventListener('mousedown', handleClickOutside)
})

// 计算属性：根据开关和contentrating筛选壁纸
const filteredWallpapers = computed(() => {
  let result = wallpapers.value
  if (showOnlyPkg.value) {
    result = result.filter(w => w.pkgPath && w.pkgPath !== '')
  }
  if (contentratingFilter.value.length < 3) {
    result = result.filter(w => contentratingFilter.value.includes(w.contentrating))
  }
  return result
})

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

  // 处理下一个请求
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
    <div class="filter-bar">
      <button @click="loadWallpapers" class="button">重载壁纸</button>
      <button @click="SelectBaseDir" class="button">修改目录</button>
      
      <!-- 下拉多选内容分级 -->
      <div class="dropdown-multiselect" id="contentrating-dropdown">
        <div class="dropdown-multiselect-label" @click="showContentratingDropdown = !showContentratingDropdown">
          <span class="desc-filter-text">内容分级：</span>
          <span class="dropdown-multiselect-selected">
            <template v-if="contentratingFilter.length === 0">无</template>
            <template v-else-if="contentratingFilter.length === contentratingOptions.length">全部</template>
            <template v-else>{{ contentratingFilter.join(', ') }}</template>
          </span>
          <span class="dropdown-arrow" :class="{ open: showContentratingDropdown }">&#9662;</span>
        </div>
        <div v-if="showContentratingDropdown" class="dropdown-multiselect-list">
          <label v-for="desc in contentratingOptions" :key="desc" class="desc-checkbox-label">
            <input
              type="checkbox"
              :value="desc"
              v-model="contentratingFilter"
              class="desc-checkbox-input"
            />
            <span class="desc-checkbox-custom"></span>
            <span class="desc-checkbox-text">{{ desc }}</span>
          </label>
        </div>
      </div>

      <label class="switch-label">
        <input type="checkbox" v-model="showOnlyPkg" class="switch-input" />
        <span class="switch-slider"></span>
        <span class="switch-text">只显示可解包</span>
      </label>

    </div>

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
      <div v-for="wallpaper in filteredWallpapers" :key="wallpaper.path" class="wallpaper-card"
        @click="selectWallpaper(wallpaper)">
        <div class="wallpaper-preview">
          <img v-if="wallpaper.coverPath && wallpaperCovers.get(wallpaper.path)"
            :src="wallpaperCovers.get(wallpaper.path)" :alt="wallpaper.name" @error="$event.target.src = ''" />
          <div v-else class="no-preview">
            {{ wallpaper.coverPath ? '加载中...' : '无预览图' }}
          </div>
        </div>
        <div class="wallpaper-title">{{ wallpaper.name }}</div>
      </div>
    </div>
    <!-- 其他内容 -->
  </div>
  <h4 style="color: #333;">RePKG-GUI v1.1.0 | RePKG v0.4.0 | CopyRight Kakune55</h4>
</template>

<style scoped>
.wallpaper-container {
  padding: 20px;
  color: #333;
}

.loading,
.error,
.no-wallpapers {
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

.switch-label {
  display: flex;
  align-items: center;
  user-select: none;
  cursor: pointer;
  margin-bottom: 0;
  margin-left: 6px;
  margin-right: 6px;
}

.switch-input {
  display: none;
}

.switch-slider {
  width: 46px;
  height: 26px;
  background: #b3c0d1;
  border-radius: 26px;
  position: relative;
  transition: background 0.2s;
  margin-right: 10px;
  flex-shrink: 0;
  box-shadow: 0 1px 3px rgba(33,150,243,0.04);
}

.switch-slider::before {
  content: "";
  position: absolute;
  left: 3px;
  top: 3px;
  width: 20px;
  height: 20px;
  background: #fff;
  border-radius: 50%;
  transition: transform 0.2s;
  box-shadow: 0 1px 3px rgba(33,150,243,0.10);
}

.switch-input:checked + .switch-slider {
  background: #2196f3;
}

.switch-input:checked + .switch-slider::before {
  transform: translateX(20px);
}

.switch-text {
  font-size: 15px;
  color: #2196f3;
  font-weight: 500;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 18px;
  margin-bottom: 18px;
  flex-wrap: wrap;
}

/* 优化按钮样式 */
.button {
  padding: 7px 20px;
  background: linear-gradient(90deg, #2196f3 60%, #42a5f5 100%);
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(33,150,243,0.08);
  transition: background 0.2s, box-shadow 0.2s, transform 0.1s;
  margin: 0;
  outline: none;
  min-width: 90px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.button:hover, .button:focus {
  background: linear-gradient(90deg, #1976d2 60%, #2196f3 100%);
  box-shadow: 0 4px 16px rgba(33,150,243,0.15);
  transform: translateY(-2px) scale(1.03);
}

/* 下拉多选样式 */
.dropdown-multiselect {
  position: relative;
  min-width: 170px;
  user-select: none;
}
.dropdown-multiselect-label {
  display: flex;
  align-items: center;
  border: 1.5px solid #b3c0d1;
  border-radius: 7px;
  background: #f8fbfd;
  padding: 5px 14px;
  cursor: pointer;
  min-height: 34px;
  transition: border 0.2s, box-shadow 0.2s;
  gap: 7px;
  box-shadow: 0 1px 3px rgba(33,150,243,0.04);
}
.dropdown-multiselect-label:hover, .dropdown-multiselect-label:focus {
  border-color: #2196f3;
  box-shadow: 0 2px 8px rgba(33,150,243,0.10);
}
.dropdown-multiselect-selected {
  flex: 1;
  font-size: 15px;
  color: #1976d2;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-weight: 500;
}
.dropdown-arrow {
  margin-left: 8px;
  font-size: 13px;
  color: #2196f3;
  transition: transform 0.2s;
}
.dropdown-arrow.open {
  transform: rotate(180deg);
}
.dropdown-multiselect-list {
  position: absolute;
  top: 110%;
  left: 0;
  z-index: 10;
  background: #fff;
  border: 1.5px solid #b3c0d1;
  border-radius: 7px;
  box-shadow: 0 4px 16px rgba(33,150,243,0.10);
  padding: 10px 0;
  min-width: 170px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  animation: dropdownIn 0.18s;
}
@keyframes dropdownIn {
  from { opacity: 0; transform: translateY(-8px);}
  to { opacity: 1; transform: translateY(0);}
}
.desc-checkbox-label {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 3px 18px;
  user-select: none;
  border-radius: 4px;
  transition: background 0.15s;
}
.desc-checkbox-label:hover {
  background: #f0f7ff;
}
.desc-checkbox-input {
  display: none;
}
.desc-checkbox-custom {
  width: 17px;
  height: 17px;
  border: 2px solid #b3c0d1;
  border-radius: 4px;
  background: #fff;
  margin-right: 6px;
  position: relative;
  transition: border 0.2s, background 0.2s;
}
.desc-checkbox-input:checked + .desc-checkbox-custom {
  border-color: #2196f3;
  background: #2196f3;
}
.desc-checkbox-input:checked + .desc-checkbox-custom::after {
  content: "";
  position: absolute;
  left: 4px;
  top: 1px;
  width: 4px;
  height: 8px;
  border: solid #fff;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}
.desc-checkbox-custom::after {
  content: "";
  display: none;
}
.desc-checkbox-input:checked + .desc-checkbox-custom::after {
  display: block;
}
.desc-checkbox-text {
  font-size: 15px;
  color: #333;
  font-weight: 400;
}
.desc-filter-text {
  font-size: 15px;
  color: #2196f3;
  margin-right: 2px;
  font-weight: 500;
}
</style>
