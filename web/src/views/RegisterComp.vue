<script setup lang="ts">
import { ref } from 'vue';
import router from '../router'
import { ElMessageBox } from 'element-plus'
import 'element-plus/theme-chalk/index.css';
import errorcode from '../errorcode'

const id = ref("")
const pw = ref("")
const invitecode = ref("")

function Register(){
    let data = new FormData()
    data.append("id", id.value)
    data.append("password", pw.value)
    data.append("invitecode", invitecode.value)

    fetch("/api/register", {
        method:"POST",
        body:data
    })
    .then(response=>response.json())
    .then(data=>{
        if(data.errorcode == 0){
            router.push("/login")
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
                <el-text size="large">账号</el-text>
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
            <div class="input_box">
                <el-text size="large">邀请码</el-text>
                <el-input v-model="invitecode" class="input" placeholder="Please input" size="large"/>
            </div>
        </div>
        <div class="login_btns">
            <el-button class="btn" size="large" type="primary" @click="Register()">注册</el-button>
            <div class="create_link">
                <el-link type="primary" href="/login">返回登陆界面</el-link>
            </div>
            
        </div>
    </div>
</template>

<style scoped>
@import "../css/uikit.css"
</style>