<script setup lang = "ts">
const props = defineProps({
    bookName: String,
    bookCoverPath: String,
    bookPath: String
})
</script>

<template>
    <div class = "bookinfo">
        <img class="coverImg" :src= "props.bookCoverPath"/>
        {{ bookName }}
        
    </div>
</template>

<style>
.bookinfo {
    width: 128px;
    height: 235px;
}
.coverImg {
    width: 128px;
    height: 185px;
}
</style>