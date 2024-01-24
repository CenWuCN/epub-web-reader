import {createRouter, createWebHistory} from 'vue-router'

const ReaderComp = () => import("./views/ReaderComp.vue")
const BookShelfComp = () => import("./views/BookShelfComp.vue")
const LoginComp = () => import("./views/LoginComp.vue")
// const app = () => import("./App.vue")

const router = createRouter({
    history: createWebHistory(),
    routes:[
        // { path:"/", component: app}
        { name:"bookShelf", path:"/bookshelf", component: BookShelfComp},
        { name:"reader", path:"/reader/:bookid", component: ReaderComp},
        { name:"login", path:"/login", component: LoginComp}
    ]
})

export default router