import {createRouter, createWebHistory} from 'vue-router'
import ManageBooksVue from './views/ManageBooks.vue'

const ReaderComp = () => import("./views/ReaderComp.vue")
const BookShelfComp = () => import("./views/BookShelfComp.vue")
const LoginComp = () => import("./views/LoginComp.vue")
const UploadBooksComp = () => import("./views/UploadBooksComp.vue")
const RegisterComp = () => import("./views/RegisterComp.vue")
const ManageBooksComp = () => import("./views/ManageBooks.vue")
// const app = () => import("./App.vue")

const router = createRouter({
    history: createWebHistory(),
    routes:[
        // { path:"/", component: app}
        { name:"bookshelf", path:"/bookshelf", component: BookShelfComp},
        { name:"reader", path:"/reader/:bookid", component: ReaderComp},
        { name:"login", path:"/login", component: LoginComp},
        { name:"uploadbooks", path:"/uploadbooks", component: UploadBooksComp},
        { name:"register", path:"/register", component:RegisterComp},
        { name:"managebooks", path:"/managebooks", component:ManageBooksComp}
    ]
})

export default router