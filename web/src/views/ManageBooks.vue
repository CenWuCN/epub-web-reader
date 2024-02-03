<script setup lang="ts">
import HeaderComp from '../components/HeaderComp.vue';
import { ref } from 'vue'
import { ElTable } from 'element-plus'
import BookInfo from '../datastruct/BookInfo';
import {useStore} from '../stores/store'

const multipleTableRef = ref<InstanceType<typeof ElTable>>()
const tableData = ref<Array<BookInfo>>([])

const store = useStore()

fetch("/api/bookshelf", {
    headers: store.getHeaders(),
    method:"POST",
    
})
.then(response=>response.json())
.then(data=>{
    tableData.value = data
})

function DelBooks(){
    let rows = multipleTableRef.value?.getSelectionRows()
    console.log(rows)

    if (rows.len >0) {
        let stringArray:Array<string>
        console.log(multipleTableRef.getSelectionRows)
        fetch("/api/delbooks", {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
                "Authorization": "Bearer " + store.getToken()
            },
            body: jsonData
        })
    }

}

</script>

<template>
    <div class="middle-bg">
        <HeaderComp></HeaderComp>
        <el-table
            ref="multipleTableRef"
            :data="tableData"
            
        >
            <el-table-column type="selection" width="55"></el-table-column>
            <el-table-column label="书名" property="Name" width="400"></el-table-column>
            <el-table-column label="文件名" property="Path" width="500"></el-table-column>
            <el-table-column label="id" property="Id" width="120"></el-table-column>
        </el-table>
        <el-button @click="DelBooks">删除书籍</el-button>
    </div>
</template>

<style scoped>
@import "../css/middle-bg.css";
</style>