<script setup lang="ts">
// import HelloWorld from './components/HelloWorld.vue'
import JSZip from 'jszip'
import ePub, { NavItem } from 'epubjs'
import { ref } from 'vue'
import { routerKey, useRoute } from 'vue-router';
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

const route = useRoute()
console.log(route.params.bookpath)
let bookpath = "/" + route.params.bookpath.join("/")
console.log(bookpath)
let book = ePub(bookpath);
let rendition = book.renderTo("area", { flow: "scrolled-doc", width: "800" });
let displayed = rendition.display();

rendition.themes.font("微软雅黑")
rendition.themes.fontSize("18px")
rendition.themes.default({ ".p": { "line-height":"4 !important"}})
rendition.themes.register("default", {
  body:{color:"#d0d3d8", background:"#1c1c1d", lineHeight:"2"},
})
rendition.themes.select("default")


let toclist = ref(Array<TocItem>)
const table = ref(false)
const booktitle = ref("")

book.loaded.navigation.then((navi) => {
  console.log(navi.toc)
  toclist.value = GetTocList(navi.toc, 0)
  console.log(toclist)
})

book.loaded.metadata.then((metadata)=>{
  booktitle.value = metadata.title
})

let reg = /#(.*)/


async function JumpToToc(link: string): void {
  console.log("点击目录")
  rendition.display(link).then(()=>{
    console.log(link)
    let tocInHtmlId = reg.exec(link)
    if (tocInHtmlId !== null) {
      console.log(tocInHtmlId)
      console.log(typeof(tocInHtmlId[1]))
      
      // await nextTick()
      // let el = document.getElementById(tocInHtmlId[1])
      let iframes = document.getElementsByTagName("iframe")
      console.log(iframes[0].id)
      let innerDoc = iframes[0].contentDocument || iframes[0].contentWindow?.document
      let el = innerDoc?.getElementById(tocInHtmlId[1])
      // let el = document.querySelector(".chapterCaption")

      // console.log(el)
      // el.scrollIntoView({behavior:'smooth'})
      var offsetTop = el?.offsetTop;

      // 滚动到目标元素的位置
        console.log(offsetTop)
        window.scrollTo(0, offsetTop);
    }


  })

}

function JumpToBookShelf(){
  router.push("/bookshelf")
}

</script>

<template>
    <div class="app_content">
      <div id='area'></div>
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
        size="30%"
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
  justify-content: center;
  min-height: 100vh;
  min-width: 1000px;
  margin-left: auto;
  margin-right: auto;
  background-color: #262628;
}
#area {
  margin-left: auto;
  margin-right: auto;
}

#toc-container {
  position: relative;
}

.control_btns {
  position: fixed;
  left: 50%;
  margin-left: 548px;
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

</style>
