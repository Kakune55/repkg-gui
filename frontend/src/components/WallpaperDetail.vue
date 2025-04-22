<script setup>
import { defineProps, defineEmits, ref, onMounted } from 'vue'
import { GetImgBase64, OpenDirInExploer, GetRepkgVersion, ExtractPKG } from '../../wailsjs/go/main/App'

const props = defineProps({
    wallpaper: {
        type: Object,
        required: true
    }
})

const emit = defineEmits(['backToList'])

function goBack() {
    emit('backToList')
}

// 封面base64
const coverBase64 = ref(null)

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


</script>

<template>
    <div class="wallpaper-detail-container">
        <div class="wallpaper-detail-header">
            <button class="back-button" @click="goBack">
                返回列表
            </button>
            <h2 style="color: black;">{{ wallpaper.name }}</h2>
        </div>

        <div class="wallpaper-detail-content">
            <div class="wallpaper-info">
                <div class="wallpaper-cover">
                    <img v-if="coverBase64" :src="coverBase64" :alt="wallpaper.name" class="cover-img" />
                    <div v-else class="no-preview">无预览图</div>
                </div>
                <div class="wallpaper-properties">
                    <div class="property-item">
                        <div class="property-label">壁纸名称:</div>
                        <div class="property-value">{{ wallpaper.name }}</div>
                    </div>
                    <div class="property-item">
                        <div class="property-label">路径:</div>
                        <div class="property-value">{{ wallpaper.path }}</div>
                    </div>
                    <div class="property-item">
                        <div class="property-label">包文件路径:</div>
                        <div class="property-value">{{ wallpaper.pkgPath }}</div>
                    </div>
                </div>
            </div>

            <div class="wallpaper-actions">
                <div class="action-buttons">
                    <button class="action-button" @click="rePkgExtract">
                        {{ rePkgExtractStatus }}
                    </button>
                    <button class="action-button" @click="openDirInExploer">
                        打开文件夹
                    </button>
                    <button class="action-button" @click="rePkgExtract">
                        详细文件信息(package.json)
                    </button>
                    <button class="action-button" @click="GetRepkgVersion">
                        RePKG版本信息
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.wallpaper-detail-container {
    padding: 20px;
    max-width: 1000px;
    margin: 0 auto;
}

.wallpaper-detail-header {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
    gap: 15px;
}

.back-button {
    padding: 8px 16px;
    background-color: #2196f3;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    display: flex;
    align-items: center;
}

.back-button:hover {
    background-color: #0d8aee;
}

.wallpaper-detail-content {
    display: grid;
    grid-template-columns: 1fr;
    gap: 20px;
}

.wallpaper-info {
    background-color: rgba(201, 201, 201, 1);
    border-radius: 8px;
    padding: 20px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
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
    background: #eee;
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
    color: #121212;
}

.property-value {
    color: #121212;
    flex: 1;
    word-break: break-all;
}

.wallpaper-actions {
    background-color: rgba(201, 201, 201, 1);
    border-radius: 8px;
    padding: 10px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.action-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    margin: 15px;
}

.action-button {
    padding: 10px 20px;
    background-color: #2196f3;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    min-width: 150px;
}

.action-button:hover {
    background-color: #0d8aee;
}

@media (max-width: 600px) {
    .wallpaper-info {
        flex-direction: column;
        align-items: stretch;
    }

    .wallpaper-cover {
        margin-right: 0;
        margin-bottom: 10px;
        width: 100%;
        height: 120px;
    }

    .property-item {
        flex-direction: column;
    }

    .property-label {
        margin-bottom: 4px;
    }
}
</style>
