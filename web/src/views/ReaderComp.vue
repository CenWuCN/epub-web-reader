<script setup lang="ts">

import JSZip from 'jszip'
import ePub, { NavItem, Book } from 'epubjs'
import { onMounted, ref } from 'vue'
import router from "../router"
import {useStore} from '../stores/store'
import BookInfo from '../datastruct/BookInfo'
import { onBeforeRouteLeave, useRoute } from 'vue-router'

const store = useStore()

interface TocItem {
  id: string;
  href: string;
  label: string;
}

let reg = /#(.*)/
let currentSectionLink = ""
let readingPos = 0
let currentSection
let toclist = ref(Array<TocItem>)
const table = ref(false)
const booktitle = ref("")
const drawerwidth = ref("")

const route = useRoute()
let bookid = route.params.bookid

let data = new FormData()
data.append("bookid", bookid)

let book: Book
let area = document.getElementById("area")

fetch("/api/bookinfo", {
  method:"POST",
  body:data,
  headers: store.getHeaders()
})
.then(response=>{
        if (response.status == 401) {
            router.push("/login")
        }
        return response.json()
    })
.then((bookinfo)=>{
  console.log(bookinfo)
  bookinfo as BookInfo
  book = ePub(bookinfo.Opf)
  book.loaded.navigation.then((navi) => {
    console.log(navi.toc)
    toclist.value = GetTocList(navi.toc, 0)
    console.log(toclist)

    currentSectionLink = bookinfo.ReadingPos.Link
    readingPos = bookinfo.ReadingPos.Percentage
    if (currentSectionLink == ""){
      display(toclist.value[0].href, 0)
    }
    else{
      console.log("重新刷新", currentSectionLink, readingPos)
      display(currentSectionLink, readingPos)
    }
})

  book.loaded.metadata.then((metadata)=>{
    booktitle.value = metadata.title.substring(0, 30) + "..."
  })
})



function GetTocList(toc: NavItem[], j: number): Array<TocItem> {
  let spaces: string = ""
  let toclist: Array<TocItem> = []
  for (let i: number = 0; i < j; i++) {
    spaces = spaces + "    "
  }

  toc.forEach((element: NavItem) => {
    toclist.push({
      id: element.id,
      href: element.href,
      label: spaces + element.label.trim()
    })
    if (element.subitems.length > 0) {
      toclist = toclist.concat(GetTocList(element.subitems, j + 1))
    }
  });
  return toclist
}

function display(link: string, readingPos:number){
  let section = book.spine.get(link)
  currentSection = section
  currentSectionLink = link
  console.log("获取到的 section", section)

  console.log(link)
  let tocInHtmlId = reg.exec(link)
  console.log(tocInHtmlId)
  if (section) {
    section.render().then(function(html){
      // console.log("需要渲染的 html", html)
      area = document.getElementById("area")
      if (area !== null){

        console.log("设置 html")
        area.innerHTML = html
        area.style.color = "#d0d3d8"
        area.style.backgroundColor = "#1c1c1d"
        area.style.fontSize = "18px"
        // area.style.lineHeight="4em"
        // area.style.lineHeight="important"
        
        if (readingPos == 0){
          if (tocInHtmlId !== null){
            let el = document.getElementById(tocInHtmlId[1])
            el?.scrollIntoView({behavior:'smooth'})
          }
        }else{
          setTimeout(function(){
            let y = document.body.scrollHeight * Number(readingPos)
              window.scrollTo({
                top: y,
                behavior:"smooth"
              })
          }, 500)
        }
      }
    })
  }
}

async function JumpToToc(link: string) {
  console.log("点击目录")
  display(link, 0)

}

function JumpToBookShelf(){
  router.push("/bookshelf")
}

function ScrollendEnvent(){
  let readingPos = window.scrollY/document.body.scrollHeight

  let data = new FormData()
  data.append("id", bookid)
  data.append("link", currentSectionLink)
  data.append("percentage", String(readingPos))
  fetch("/api/readingposset", {
    method:"POST",
    body: data,
    headers: store.getHeaders()
  })
}

function ResizeEvent(){
  if (window.innerWidth < 768){
    drawerwidth.value = "80%"
  }
  else{
    drawerwidth.value = "30%"
  }
}

onMounted(()=>{
  var next = document.getElementById("next_btn");
  next.addEventListener("click", function(e){
    window.scrollTo(0,0);
    let spineitem = currentSection.next();
    console.log(spineitem)
    display(spineitem.href, 0)
    e.preventDefault();
  }, false);
})
onBeforeRouteLeave(()=>{
  window.removeEventListener("scrollend", ScrollendEnvent)
  window.removeEventListener("resize", ResizeEvent)
})

window.addEventListener("scrollend", ScrollendEnvent)

window.addEventListener("resize", ResizeEvent)

</script>

<template>
    <div class="app_content">
      <div id="reader">
        <div id='area'></div>
        <div class="next_chapter">
          <el-button id="next_btn" class="control_btn" href="#next">下一章</el-button>
        </div>
      </div>

    </div>
    
    <div class="control_btns">
      <el-button class="control_btn" @click="table = true">目录</el-button>
      <el-button class="control_btn" @click="JumpToBookShelf">书架</el-button>
    </div>
    <div class="drawer">
      <el-drawer
        :title=booktitle
        v-model="table"
        direction="rtl"
        :size = "drawerwidth"
        :style="{ background: '#3a3a3c',color: '#eef0f4' }"
        >
        <div id="toc">
          <ul>
            <li v-for="item in toclist" @click="JumpToToc(item.href)">{{ item.label }}</li>
          </ul>
        </div>
      </el-drawer>
    </div>

</template>

<style scoped>
.app_content {
  display: flex;
  align-items: flex-start;
  flex-direction: column;
  height: auto;
  min-height: 100dvh;
  width: 100dvw;
  margin-left: auto;
  margin-right: auto;
  background-color: #262628;
}
#reader {
  display: flex;
  align-items: flex-start;
  flex-direction: column;
  width: 1000px;
  height: auto;
  min-height: 100dvh;
  margin-left: auto;
  margin-right: auto;
  background-color: #1c1c1d;
}
#area {
  margin-left: auto;
  margin-right: auto;
  width: 798px;
  height: auto;
  color:"#d0d3d8";
  background-color:"#1c1c1d";
}
.p{
  color: #d0d3d8;
}

.next_chapter {
  margin-left: auto;
  margin-right: auto;
  margin-bottom: 0;
  height: 260px;
}

#next_btn {
  margin-top: 100px;
  width: 400px;
  height: 60px;
  font-size: 16px;
  background-color: #333334;
  color: #eef0f4;
}

#next_btn:hover {
  background-color: #282829;
}

#toc-container {
  position: relative;
}

.control_btns {
  position: fixed;
  left: 80dvw;
  bottom: 48px;
  width: 48px;
  display: flex;
  align-items: flex-start;
  flex-direction: column;
}

.control_btn{
  margin-top: 24px;
  margin-left: 0px;
  background-color: #1c1c1d;
  border: 0;
}

.el-drawer__title {
  color: #eef0f4;
}

li {
  list-style: none;
  text-align: left;
  white-space: pre;
}

li:hover {
  background-color: #3474b4;
}

li:active,
li:visited,
li:focus {
  color: #0097ff;
}

ul {
  height: 800px;
}

@media screen and (width < 768px) {
  .app_content{
    width: 100%;
    height: auto;
    min-height: 100dvh;
  }
  #reader {
    width: 100%;
    height: auto;
    min-height: 100dvh;
  }
  #area {
    width: 90%;
    height: auto;
  }
  #next_btn {
    width: 80dvw;
    height: 4rem;
  }
  .control_btns {
    position: fixed;
    margin-left: 80dvw;
    left: 0dvw;
  }
}

</style>
