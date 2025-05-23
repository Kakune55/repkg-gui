<script setup>
import { ref, onMounted, computed, onBeforeUnmount, watch } from 'vue'
import { GetWallpapers, SelectBaseDir } from '../../wailsjs/go/main/App'

// 响应式数据
const wallpapers = ref([])
const loading = ref(true)
const error = ref(null)
const searchQuery = ref('') // 添加搜索关键词变量

// 筛选开关：是否只显示pkgPath不为空的壁纸
const showOnlyPkg = ref(false)
// contentrating多选筛选
const contentratingFilter = ref(['Everyone', 'Questionable'])


// 计算属性：根据开关和contentrating筛选壁纸
const filteredWallpapers = computed(() => {
  let result = wallpapers.value
  
  // 首先按搜索关键词过滤
  if (searchQuery.value.trim() !== '') {
    const query = searchQuery.value.toLowerCase().trim()
    result = result.filter(w => 
      w.name.toLowerCase().includes(query) || 
      (w.title && w.title.toLowerCase().includes(query))
    )
  }
  
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

// 选择壁纸，触发跳转
function selectWallpaper(wallpaper) {
  emit('select-wallpaper', wallpaper)
}

// 组件挂载时加载壁纸
onMounted(() => {
  loadWallpapers()
})

function selectBaseDir() {
  SelectBaseDir()
  loadWallpapers()
}

// 分页相关
const currentPage = ref(1)
const pageSize = ref(32) // 每页壁纸数
const pageCount = computed(() => Math.ceil(filteredWallpapers.value.length / pageSize.value))

// 当前页壁纸
const pagedWallpapers = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredWallpapers.value.slice(start, start + pageSize.value)
})

// 监听筛选变化时重置页码
watch([filteredWallpapers, pageSize], () => {
  currentPage.value = 1
})

// 加载壁纸列表
async function loadWallpapers() {
  try {
    loading.value = true
    error.value = null
    const result = await GetWallpapers()
    const parsedWallpapers = JSON.parse(result)
    wallpapers.value = parsedWallpapers
    loading.value = false
  } catch (err) {
    console.error("加载壁纸失败:", err)
    error.value = "加载壁纸失败: " + (err.message || "未知错误")
    loading.value = false
  }
}

// 获取封面图片src（取coverPath最后两段）
function getCoverSrc(coverPath) {
  if (!coverPath) return '';
  const parts = coverPath.replace(/\\/g, '/').split('/');
  if (parts.length < 2) return coverPath;
  return parts.slice(-2).join('/');
}

function toggleContentrating(rating, checked) {
  if (checked) {
    contentratingFilter.value.push(rating)
  } else {
    const index = contentratingFilter.value.indexOf(rating)
    if (index > -1) {
      contentratingFilter.value.splice(index, 1)
    }
  }
}

// 清除搜索框
function clearSearch() {
  searchQuery.value = ''
}

</script>

<template>
  <s-drawer id="drawer">
    <div slot="start" class="setting-bar">
      <div>
        <s-button @click="loadWallpapers">重载壁纸</s-button>
        <s-button @click="selectBaseDir">修改目录</s-button>
      </div>
      <div>
        <s-text-field style="display: grid; width: auto" label="搜索壁纸" v-model="searchQuery">
          <s-icon slot="start" name="search"></s-icon>
          <s-icon-button slot="end" @click="clearSearch" v-if="searchQuery">
            <s-icon name="close"></s-icon>
          </s-icon-button>
        </s-text-field>
      </div>

      <s-fold>
        <s-button slot="trigger">分级过滤</s-button>
        <div class="fold-content">
          <div class="fold-checkbox-row">
            <s-checkbox :checked="contentratingFilter.includes('Everyone')"
              @change="e => toggleContentrating('Everyone', e.target.checked)" type="checkbox"></s-checkbox>
            <span>Everyone</span>
          </div>
          <div class="fold-checkbox-row">
            <s-checkbox :checked="contentratingFilter.includes('Questionable')"
              @change="e => toggleContentrating('Questionable', e.target.checked)" type="checkbox"></s-checkbox>
            <span>Questionable</span>
          </div>
          <div class="fold-checkbox-row">
            <s-checkbox :checked="contentratingFilter.includes('Mature')"
              @change="e => toggleContentrating('Mature', e.target.checked)" type="checkbox"></s-checkbox>
            <span>Mature</span>
          </div>
        </div>
      </s-fold>

      <div class="switch-label">
        <s-switch v-model.lazy="showOnlyPkg" type="checkbox"></s-switch>
        <span>只显示可解包</span>
      </div>

    </div>
    <s-appbar>
      <s-icon-button slot="navigation" onclick="document.querySelector('#drawer').toggle()">
        <s-icon name="menu"></s-icon>
      </s-icon-button>
      <div slot="headline"> 壁纸列表 </div>
    </s-appbar>
    <s-scroll-view class="wallpaper-container">


      <!-- 加载状态 -->
      <div v-if="loading" class="loading">
        加载中...
      </div>

      <!-- 错误信息 -->
      <div v-else-if="error" class="error">
        {{ error }}
      </div>

      <!-- 无壁纸提示 -->
      <div v-else-if="pagedWallpapers.length === 0" class="no-wallpapers">
        <s-empty style="height: 540px;">{{ wallpapers.length === 0 ? '暂时没有数据' : '没有符合筛选条件的壁纸' }}</s-empty>
      </div>

      <!-- 壁纸网格 -->
      <div v-else class="wallpaper-grid">
        <s-card v-for="wallpaper in pagedWallpapers" :key="wallpaper.path" class="wallpaper-card"
          @click="selectWallpaper(wallpaper)">
          <div class="wallpaper-preview">
            <img v-if="wallpaper.coverPath" :src="getCoverSrc(wallpaper.coverPath)" :alt="wallpaper.name"
              @error="$event.target.src = ''" />
            <div v-else class="no-preview">
              {{ wallpaper.coverPath ? '加载中...' : '无预览图' }}
            </div>
          </div>
          <div class="wallpaper-title">{{ wallpaper.name }}</div>
        </s-card>
      </div>

      <!-- 分页控件 -->
      <div v-if="pageCount > 1" class="pagination-bar">
        <s-button :disabled="currentPage === 1" @click="currentPage--">上一页</s-button>
        <span>第 {{ currentPage }} / {{ pageCount }} 页</span>
        <s-button :disabled="currentPage === pageCount" @click="currentPage++">下一页</s-button>
        <s-text-field label="Page" type="number" min="1" :max="pageCount" v-model.number="currentPage"
          style="width: 50px; margin-left: 10px;" ></s-text-field>
      </div>
      <!-- 其他内容 -->
      <h4>RePKG-GUI v1.2.1 | RePKG v0.4.0 | CopyRight Kakune55</h4>
    </s-scroll-view>
  </s-drawer>
</template>

<style scoped>
.wallpaper-container {
  padding: 20px;
  height: 100%;
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
  border-radius: 8px;
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
  gap: 10px;
  align-items: center;
  user-select: none;
}


.setting-bar {
  display: block;
  padding: 20px;
}

.setting-bar > div {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  gap: 15px;
}

.fold-content {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 20px;
}

.fold-checkbox-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 2px 0;
}

/* 分页样式 */
.pagination-bar {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 18px 0 0 0;
  gap: 12px;
}
</style>
