import { defineStore } from 'pinia'

export const useStore = defineStore("store", {
    state:()=>({token:""}),
    actions:{
        getToken(){
            if (this.token == ""){
                this.token = localStorage.getItem("token")
            }
            return this.token
        }
    }
})