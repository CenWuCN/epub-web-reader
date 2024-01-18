<script setup lang="ts">
// import HelloWorld from './components/HelloWorld.vue'
import JSZip from 'jszip'
import ePub, { NavItem } from 'epubjs'
import { ref } from 'vue'
import { useRoute } from 'vue-router';
import router from "../router"

// // 创建一个 MutationObserver 实例
// let observer = new MutationObserver(function(mutations) {
//   mutations.forEach(function(mutation) {
//     // 在这里处理元素创建的逻辑
//     if (mutation.type === 'childList' && mutation.addedNodes.length > 0) {
//       // 处理新添加的元素
//       console.log('新元素已创建:', mutation.addedNodes);
//     }
//   });
// });

// // 选择要观察的目标节点
// let targetNode = document.body;

// // 配置观察选项（此处配置为监视子节点的变化）
// let config = { childList: true, subtree: true };

// // 传入目标节点和观察选项开始观察
// observer.observe(targetNode, config);

interface TocItem {
  id: string;
  href: string;
  label: string;
}

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

const epubwidth : string= "500px"
const route = useRoute()
console.log(route.params.bookpath)
// let bookpath = "/" + route.params.bookpath.join("/")
// console.log(bookpath)
let book = ePub("/epubs/11/OEBPS/content.opf");
// let rendition = book.renderTo("area", { flow: "scrolled-doc", width: epubwidth });
// let displayed = rendition.display();

// rendition.themes.font("微软雅黑")
// rendition.themes.fontSize("18px")
// rendition.themes.default({ ".p": { "line-height":"4 !important"}})
// rendition.themes.register("default", {
//   body:{color:"#d0d3d8", background:"#1c1c1d", lineHeight:"2"},
// })
// rendition.themes.select("default")
let area = document.getElementById("area")

let reg = /#(.*)/

function display(link: string){
  let section = book.spine.get(link)
  console.log("获取到的 section", section)

  console.log(link)
  let tocInHtmlId = reg.exec(link)
  console.log(tocInHtmlId)
  if (section) {
    section.render().then(function(html){
      console.log("需要渲染的 html", html)
      area = document.getElementById("area")
      if (area !== null){

        console.log("设置 html")
        area.innerHTML = html
        area.style.color = "#d0d3d8"
        area.style.backgroundColor = "#1c1c1d"
        area.style.fontSize = "18px"
        area.style.lineHeight="4em"

        if (tocInHtmlId !== null){
          let el = document.getElementById(tocInHtmlId[1])
          el?.scrollIntoView({behavior:'smooth'})
        }

        var next = document.getElementById("next_btn");
        next.addEventListener("click", function(e){
          window.scrollTo(0,0);
          let spineitem = section.next();
          console.log(spineitem)
          display(spineitem.href)
          e.preventDefault();
        }, false);
      }
    })
  }
}

let toclist = ref(Array<TocItem>)
const table = ref(false)
const booktitle = ref("")

book.loaded.navigation.then((navi) => {
  console.log(navi.toc)
  toclist.value = GetTocList(navi.toc, 0)
  console.log(toclist)

  display(toclist.value[0].href)

})

book.loaded.metadata.then((metadata)=>{
  booktitle.value = metadata.title
})



async function JumpToToc(link: string) {
  console.log("点击目录")
  display(link)

}

function JumpToBookShelf(){
  router.push("/bookshelf")
}

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

div {
  --epubwidth: 800px;
}

.app_content {
  display: flex;
  align-items: flex-start;
  flex-direction: column;
  height: auto;
  min-height: 100dvh;
  width: 100vw;
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
  background-color:"#1c1c1d",
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

#blank {
  width: var(--epubwidth);
  height: 300px;
  margin-left: auto;
  margin-right: auto;
  background-color: #1c1c1d;
}

#toc-container {
  position: relative;
}

.control_btns {
  position: fixed;
  left: 50%;
  margin-left: calc(var(--epubwidth) /2 + var(--epubwidth)/100);
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
  .drawer {
    width: 80dvw;
  }
}

</style>
