<script setup lang="ts">

</script>

<template>
    <div class="header">
        <RouterLink to="/uploadbooks"><el-button size = "large">上传书籍</el-button></RouterLink>
        <RouterLink to="/managebooks"><el-button size = "large">管理书籍</el-button></RouterLink> 
        <RouterLink to="/bookshelf"><el-button size = "large">书架</el-button></RouterLink>
        <el-text size="large">example</el-text>
    </div>
</template>

<style scoped>
.header {
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
}
</style>