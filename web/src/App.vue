<script setup lang="ts">
// import HelloWorld from './components/HelloWorld.vue'
import JSZip from 'jszip'
import ePub, { NavItem } from 'epubjs'
import { ref, nextTick } from 'vue'

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

let book = ePub("/doc/11.epub");
let rendition = book.renderTo("area", { flow: "scrolled-doc", width: "800" });
let displayed = rendition.display();
rendition.themes.font("微软雅黑")

let toclist = ref(Array<TocItem>)


book.loaded.navigation.then((navi) => {
  console.log(navi.toc)
  toclist.value = GetTocList(navi.toc, 0)
  console.log(toclist)
})

let reg = /#(.*)/


async function JumpToToc(link: string): void {
  console.log("点击目录")
  rendition.display(link).then(()=>{
    let tocInHtmlId = reg.exec(link)
    console.log(tocInHtmlId[1])
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
    var offsetTop = el.offsetTop;

    // 滚动到目标元素的位置
      console.log(offsetTop)
      window.scrollTo(0, offsetTop);

  })

}

</script>

<template>
  <el-row>
    <el-col :span="4">
      <div class="left-space"></div>
    </el-col>
    <el-col :span="16">
      <div class="content">
        <div id='area'></div>
      </div>
    </el-col>
    <el-col :span="4">
      <div class="right-space">
        <div id="toc-container">
          <div id="toc">
            <ul>
              <li v-for="item in toclist" @click="JumpToToc(item.href)">{{ item.label }}</li>
            </ul>
          </div>
        </div>
      </div>
    </el-col>
  </el-row>


  <!-- <div>
    <a href="https://vitejs.dev" target="_blank">
      <img src="/vite.svg" class="logo" alt="Vite logo" />
    </a>
    <a href="https://vuejs.org/" target="_blank">
      <img src="./assets/vue.svg" class="logo vue" alt="Vue logo" />
    </a>
  </div>
  <HelloWorld msg="Vite + Vue" /> -->
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}

.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}

.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}

#toc-container {
  position: relative;
}

li {
  list-style: none;
  text-align: left;
  white-space: pre;
}

ul {
  height: 800px;
}
</style>
