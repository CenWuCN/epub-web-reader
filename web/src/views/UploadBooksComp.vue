<script setup lang="ts">
import { UploadFilled } from '@element-plus/icons-vue'
import HeaderComp from '../components/HeaderComp.vue';
</script>

<template>
    <div class="bg">
        <HeaderComp></HeaderComp>
        <el-upload
            class="upload-demo"
            drag
            action="https://run.mocky.io/v3/9d059bf9-4660-45f2-925d-ce80ad6c4d15"
            multiple
            >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
            Drop file here or <em>click to upload</em>
            </div>
            <template #tip>
                <div class="el-upload__tip">
                    only epub files can be uploaded
                </div>
            </template>
        </el-upload>
    </div>

</template>

<style scoped>
.bg {
    display: flex;
    flex-direction: column;
    width: 1000px;
    height: 100vh;
    margin-left: auto;
    margin-right: auto;
}

.upload-demo {
    margin-top: 400px;
}

@media screen and (width < 768px) {
    .bg {
        width: 100dvw;
    }
}
</style>