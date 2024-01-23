import {createRouter, createWebHistory} from 'vue-router'

const reader = () => import("./components/Reader.vue")
const bookShelf = () => import("./components/BookShelf.vue")
const login = () => import("./components/Login.vue")
// const app = () => import("./App.vue")

const router = createRouter({
    history: createWebHistory(),
    routes:[
        // { path:"/", component: app}
        { name:"bookShelf", path:"/bookshelf", component: bookShelf},
        { name:"reader", path:"/reader", component: reader},
        { name:"login", path:"/login", component: login}
    ]
})

export default router