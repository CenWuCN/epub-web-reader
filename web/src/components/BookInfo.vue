<script setup lang = "ts">
const props = defineProps({
    bookName: String,
    bookCoverPath: String,
    bookPath: String
})
</script>

<template>
    <div class = "book_info">
        <div class="cover_img_bg">
            <img class="cover_img" :src= "props.bookCoverPath"/>
        </div>
        <div class="book_name">{{bookName}}</div>
        
    </div>
</template>

<style>
.book_info {
    width: 128px;
    height: 235px;
    margin-left: 36px;
    margin-bottom: 44px;
}
.cover_img_bg {
    background-color: white;
}
.cover_img {
    width: 128px;
    height: 185px;
    vertical-align: top;
}
.book_name {
    display: -webkit-box;
    margin-top: 14px;
    max-height: 36px;
    font-size: 15px;
    text-align: left;
    text-overflow: ellipsis;
    line-height: 18px;
    overflow: hidden;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    color: #b2b4b8;

}
/* .bookname:hover{
    overflow: visible;
    text-overflow: unset;
} */
</style>