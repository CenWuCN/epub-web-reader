<script setup lang = "ts">
import { ref } from 'vue';
import { ElMessageBox } from 'element-plus'
import 'element-plus/theme-chalk/index.css';
import router from '../router'
import errorcode from '../errorcode'

const id = ref("")
const pw = ref("")

function Login(){
    const formData = new FormData()
    formData.append("username", id.value)
    formData.append("password", pw.value)
    fetch("/api/login", {
        method:"POST",
        body: formData
    })
        .then(response=> response.json())
        .then(data=>{
            console.log(data)
            if (data.errorcode == ""){
                router.push("/bookshelf")
            }else{
                ElMessageBox.confirm(errorcode[data.errorcode])
                    .then(() => {
                    })
                    .catch(() => {

                    })
            }
        })
}

</script>

<template>
    <div class="login_view">
        <div class="login_input">
            <div class="input_box">
                <el-text class="text" size="large">账号</el-text>
                <el-input v-model="id" class="input" placeholder="Please input" size="large"/>
            </div>
            <div class="input_box">
                <el-text class="text" size="large">密码</el-text>
                <el-input v-model="pw" class="input"
                    type="password"
                    placeholder="Please input password"
                    show-password
                    size="large"
                />
            </div>
        </div>

        <div class="login_btns">
            <el-button class="btn" size="large" type="primary" @click="Login()">登陆</el-button>
            <div class="create_link">
                <el-text>如果你还没有账号</el-text>
                <el-link type="primary" href="/register">创建一个</el-link>
            </div>

        </div>
    </div>
</template>

<style scroped>
@import "../css/uikit.css"
</style>