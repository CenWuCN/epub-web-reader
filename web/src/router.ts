import {createRouter, createWebHistory} from 'vue-router'

const reader = () => import("./components/Reader.vue")
const bookshelf = () => import("./components/BookShelf.vue")
// const app = () => import("./App.vue")

const router = createRouter({
    history: createWebHistory(),
    routes:[
        // { path:"/", component: app}
        { name:"bookShelf", path:"/bookshelf", component: bookshelf},
        { name:"Reader", path:"/reader", component: reader, props:{data:"data"}}
    ]
})

export default router