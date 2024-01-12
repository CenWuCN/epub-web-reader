import {createRouter, createWebHistory} from 'vue-router'

const reader = () => import("./components/Reader.vue")
const bookshelf = () => import("./components/BookShelf.vue")
// const app = () => import("./App.vue")

const router = createRouter({
    history: createWebHistory(),
    routes:[
        // { path:"/", component: app}
        { path:"/bookshelf", component: bookshelf},
        { path:"/reader", component: reader}
    ]
})

export default router