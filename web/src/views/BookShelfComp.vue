<script setup lang="ts">
import router from '../router';
import BookInfoComp from '../components/BookInfoComp.vue';
import { ref } from 'vue';
import {useStore} from '../stores/store'
import BookInfo from '../datastruct/BookInfo'
import HeaderComp from '../components/HeaderComp.vue';

const store = useStore()

const bookinfos = ref<Array<BookInfo>>([])
fetch("/api/bookshelf", {
        method:"POST",
        headers: store.getHeaders()
    })
    .then(response => {
        if (response.status == 401) {
            router.push("/login")
        }
        return response.json()
    })
    .then(booktable => {
        bookinfos.value = booktable
        })

function JumptoReader(bookid: string){
    // store.opfPath = bookid
    // localStorage.setItem("opf", bookid)
    router.push("/reader/"+ bookid)
}
</script>

<template>
    <div id = "container">
        <div id="app_content">
            <HeaderComp></HeaderComp>
            <div id="BookShelf">
                <BookInfoComp @click="JumptoReader(bookinfo.Id)" v-for="bookinfo in bookinfos" 
                    :bookName="bookinfo.Name"
                    :bookCoverPath="bookinfo.CoverPath"
                    :bookPath = "bookinfo.Path"
                    :opf = "bookinfo.Opf"
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
    /* padding-top: 72px; */
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