<script setup lang = "ts">

</script>

<template>
    <div class="login_view">
        <div class="login_input">
            <div class="input_box">
                <el-text class="text">账号</el-text>
                <el-input class="input" placeholder="Please input" />
            </div>
            <div class="input_box">
                <el-text class="text">密码</el-text>
                <el-input class="input"
                    type="password"
                    placeholder="Please input password"
                    show-password
                />
            </div>
        </div>

        <div class="login_btns">
            <el-button class="btn">登陆</el-button>
            <el-button class="btn">注册</el-button>
        </div>
    </div>
</template>

<style scroped>
.login_view {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    background-color: #1c1c1d;
    height: 100vh;
}
.login_input {
    width: 350px;
}
.input_box {
    display: flex;
    flex-direction: row;
    margin-top: 20px;
}
.input.el-input {
    height: 50px;
}
.text {
    font-size: 22px;
    white-space: nowrap;
    width: 100px;
}
.input .el-input__wrapper {
    background-color: #1c1c1d;
}
.login_btns {
    display: flex;
    align-items: center;
    align-content: center;
    flex-direction: column;
}
.btn.el-button {
    margin-top: 40px;
    width: 280px;
    height: 40px;
    margin-left: 0px;
    background-color: #1c1c1d;
}
</style>