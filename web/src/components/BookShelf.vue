<script setup lang="ts">
import router from '../router';
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

    function JumptoReader(bookPath: string){
        router.push("/reader"+ bookPath)
    }
</script>

<template>
    <div id = "container">
        <div id="app_content">
            <div id="BookShelf">
                <BookInfo @click="JumptoReader(bookinfo.Path)" v-for="bookinfo in bookinfos" 
                    :bookName="bookinfo.Name"
                    :bookCoverPath="bookinfo.CoverPath"
                    :bookPath = "bookinfo.Path"
                    />
            </div>
        </div>
    </div>

</template>

<style scoped>
#container {
    background-color: #1c1c1d;
    min-width: 100%;
    min-height: 100vh;
}
#app_content {
    max-width: 1120px;
    padding-top: 72px;
    margin-left: auto;
    margin-right: auto;
    background-color: #1c1c1d;
}
#BookShelf {
    display: flex;
    flex-wrap: wrap;
    margin: auto;
}

</style>