import { defineStore } from 'pinia'

export const useStore = defineStore("store", {
    state:()=>({token:""}),
    actions:{
        getToken(){
            this.token = localStorage.getItem("token")
            return this.token
        },
        getHeaders(){
            return {"Authorization": "Bearer " + this.getToken() }
        }
    }
})