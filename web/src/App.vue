<script setup lang="ts">
// import HelloWorld from './components/HelloWorld.vue'
import JSZip from 'jszip'
import ePub, { NavItem } from 'epubjs'
import {ref} from 'vue'

function GetTocList(toc:NavItem[], j:number):Array<string> {
  let spaces: string = ""
  let toclist: Array<string> =[]
  for(let i:number = 0; i < j; i++){
    spaces = spaces + "  "
  }

  toc.forEach((element:any) => {
    toclist.push(spaces + element.label.trim())
    if (element.subitems.length > 0){
      toclist = toclist.concat(GetTocList(element.subitems, j + 1))
    }
  });
  return toclist
}

let book = ePub("/doc/11.epub");
let rendition = book.renderTo("area", { flow: "scrolled-doc", width: "800"});
let displayed = rendition.display();
rendition.themes.font("微软雅黑")

let toclist = ref(Array<string>)
book.loaded.navigation.then((navi)=>{
  console.log(navi.toc)
  toclist.value = GetTocList(navi.toc, 0)
  console.log(toclist)
})


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
        <div id="toc">
          <ul>
            <li v-for="item in toclist">{{ item }}</li>
          </ul>
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

li {
  list-style: none;
  text-align: left;
  white-space: pre;
}

</style>
