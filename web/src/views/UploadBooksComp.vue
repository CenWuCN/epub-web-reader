<script setup lang="ts">
import { UploadFilled } from '@element-plus/icons-vue'
import HeaderComp from '../components/HeaderComp.vue';
import { ref } from 'vue';
import {useStore} from '../stores/store'

const store = useStore()

const headers= ref({})
headers.value = store.getHeaders()

</script>

<template>
    <div class="middle-bg">
        <HeaderComp></HeaderComp>
        <el-upload
            class="upload-demo"
            drag
            action="/api/uploadbooks"
            :headers="headers"
            multiple
            accept=".epub"
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
@import "../css/middle-bg.css";

@media screen and (width < 768px) {
    .bg {
        width: 100dvw;
    }
}
</style>