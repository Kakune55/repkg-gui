<script setup>
import { defineProps, defineEmits, ref, onMounted } from 'vue'
import { GetImgBase64, OpenDirInExploer, GetRepkgVersion, ExtractPKG, GetWallpaperProjectInfo } from '../../wailsjs/go/main/App'

const props = defineProps({
    wallpaper: {
        type: Object,
        required: true,
    }
})

const emit = defineEmits(['backToList'])

function goBack() {
    emit('backToList')
}

// 封面base64
const coverBase64 = ref(null)
const logtext = ref("壁纸信息输出:\n")
const rePkgExtractStatus = ref("提取资源")

onMounted(async () => {
    if (props.wallpaper && props.wallpaper.coverPath) {
        try {
            coverBase64.value = await GetImgBase64(props.wallpaper.coverPath)
        } catch {
            coverBase64.value = null
        }
    }
})

async function openDirInExploer() {
    try {
        await OpenDirInExploer(props.wallpaper.path)
    } catch (err) {
        console.error(`打开文件夹失败 ${props.wallpaper.path}:`, err)
    }
}

function rePkgExtract() {
    rePkgExtractStatus.value = "提取中"
    ExtractPKG(props.wallpaper.pkgPath, props.wallpaper.path + "/output")
    rePkgExtractStatus.value = "提取完成"
    setTimeout(() => {
        rePkgExtractStatus.value = "提取资源"
    }, 3000)
}

async function getWallpaperProjectInfo() {
    // 获取壁纸信息
    logtext.value = await GetWallpaperProjectInfo(props.wallpaper.path)
}


</script>

<template>
    <div class="wallpaper-detail-container">
        <div class="list-container">
            <div class="wallpaper-detail-header">
                <!-- 使用 Sober 按钮 -->
                <s-button @click="goBack">
                    返回列表
                </s-button>
                <h2>{{ wallpaper.name }}</h2>
            </div>
        </div>
        <div class="list-container">
            <div class="wallpaper-detail-content">
                <s-card class="wallpaper-info">
                    <div class="wallpaper-cover">
                        <img v-if="coverBase64" :src="coverBase64" :alt="wallpaper.name" class="cover-img" />
                        <div v-else class="no-preview">无预览图</div>
                    </div>
                    <div class="wallpaper-properties">
                        <div class="property-item">
                            <div class="property-label">说明:</div>
                            <div class="property-value">{{ wallpaper.description }}</div>
                        </div>
                        <div class="property-item">
                            <div class="property-label">路径:</div>
                            <div class="property-value">{{ wallpaper.path }}</div>
                        </div>
                        <div class="property-item">
                            <div class="property-label">包文件路径:</div>
                            <div class="property-value">{{ wallpaper.pkgPath }}</div>
                        </div>
                        <div class="property-item">
                            <div class="property-label">分级:</div>
                            <div class="property-value">{{ wallpaper.contentrating }}</div>
                        </div>
                    </div>
                </s-card>
            </div>
            <div class="list-container">
                <s-card class="wallpaper-actions">
                    <s-button @click="rePkgExtract">
                        {{ rePkgExtractStatus }}
                    </s-button>
                    <s-button @click="openDirInExploer">
                        打开文件夹
                    </s-button>
                    <s-button @click="getWallpaperProjectInfo">
                        详细文件信息(package.json)
                    </s-button>
                    <s-button @click="GetRepkgVersion">
                        RePKG版本信息
                    </s-button>
                </s-card>
            </div>
            <div class="list-container">
                <s-card class="wallpaper-log">
                    <pre class="log-output">{{ logtext }}</pre>
                </s-card>
            </div>
        </div>
    </div>
</template>

<style scoped>
s-card {
    max-width: none;
}

.wallpaper-detail-container {
    padding: 10px 25px;
    max-width: 1000px;
    margin: 0 auto;
    gap: 20px;
}

.list-container {
    padding-top: 10px;
}

.wallpaper-detail-header {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
    gap: 15px;
}


.wallpaper-detail-content {
    display: block;
    gap: 20px;
}

.wallpaper-info {
    padding: 20px;
    display: flex;
    align-items: flex-start;
    gap: 20px;
}

.wallpaper-cover {
    width: 120px;
    height: 120px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 8px;
    overflow: hidden;
    margin-right: 20px;
}

.cover-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 8px;
}

.no-preview {
    color: #999;
    font-size: 14px;
    text-align: center;
}

.wallpaper-properties {
    display: flex;
    flex-direction: column;
    gap: 10px;
    flex: 1;
}

.property-item {
    display: flex;
    margin-bottom: 8px;
}

.property-label {
    min-width: 120px;
    font-weight: bold;
}

.property-value {
    flex: 1;
    word-break: break-all;
}

.wallpaper-actions {
    display: flex;
    justify-content: left;
    padding: 15px;
    gap: 10px;
}

.wallpaper-log {
    width: 100%;
    height: 250px;
    padding: 15px;
}

.log-output {
    width: 100%;
    height: 230px;
    margin: 0;
    padding: 10px;
    font-size: 15px;
    font-family: inherit;
    white-space: pre-wrap;
    word-break: break-all;
    overflow-y: auto;
    border: none;
    outline: none;
    box-shadow: none;
    resize: none;
    user-select: text;
    text-align: left;
    /* 隐藏滚动条 */
    scrollbar-width: none; /* Firefox */
    -ms-overflow-style: none; /* IE 10+ */
}
.log-output::-webkit-scrollbar {
    display: none; /* Chrome/Safari/Webkit */
}
</style>
