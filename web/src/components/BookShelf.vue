<script setup lang="ts">
import BookInfo from './BookInfo.vue';
import { ref } from 'vue';
    
    interface bookinfo {
        Name: string
        Path: string
        CoverPath: string
    }

    const bookinfos = ref<Array<bookinfo>>([])
    fetch("/api/bookshelf")
        .then(response => response.json())
        .then(booktable => {
            console.log(booktable)
            bookinfos.value = booktable
            })
</script>

<template>
    <div id="BookShelf">
        <BookInfo v-for="bookinfo in bookinfos" 
            :bookName="bookinfo.Name"
            :bookCoverPath="bookinfo.CoverPath"
            :bookPath = "bookinfo.Path"
            />
    </div>
</template>

<style>
#BookShelf {
    display: flex;
    flex-wrap: wrap;
}
</style>