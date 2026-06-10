import { createRouter, createWebHistory } from 'vue-router'

const routes = [
    { path: '/', name: 'home', component: () => import('../views/Home.vue') },
    { path: '/products', name: 'products', component: () => import('../views/ProductList.vue') },
    { path: '/product/:id', name: 'productDetail', component: () => import('../views/ProductDetail.vue') },
    { path: '/order/:order_no', name: 'orderDetail', component: () => import('../views/OrderDetail.vue') },
    { path: '/cart', name: 'cart', component: () => import('../views/Cart.vue'), meta: { auth: true } },
    { path: '/checkout', name: 'checkout', component: () => import('../views/Checkout.vue'), meta: { auth: true } },
    { path: '/user/login', name: 'login', component: () => import('../views/Login.vue'), alias: '/login' },
    { path: '/user/register', name: 'register', component: () => import('../views/Register.vue'), alias: '/register' },
    { path: '/user/profile', name: 'userCenter', component: () => import('../views/UserCenter.vue'), meta: { auth: true } },
    { path: '/user/orders', name: 'userOrders', component: () => import('../views/UserOrders.vue'), meta: { auth: true } },
    { path: '/admin/login', name: 'adminLogin', component: () => import('../views/AdminLogin.vue') },
    { path: '/admin/', name: 'adminDashboard', component: () => import('../views/AdminDashboard.vue'), meta: { admin: true } },
    { path: '/admin/products', name: 'adminProducts', component: () => import('../views/AdminProducts.vue'), meta: { admin: true } },
    { path: '/admin/product/new', name: 'adminProductNew', component: () => import('../views/AdminProductEdit.vue'), meta: { admin: true } },
    { path: '/admin/product/:id', name: 'adminProductEdit', component: () => import('../views/AdminProductEdit.vue'), meta: { admin: true } },
    { path: '/admin/cards', name: 'adminCards', component: () => import('../views/AdminCards.vue'), meta: { admin: true } },
    { path: '/admin/categories', name: 'adminCategories', component: () => import('../views/AdminCategories.vue'), meta: { admin: true } },
    { path: '/admin/orders', name: 'adminOrders', component: () => import('../views/AdminOrders.vue'), meta: { admin: true } },
    { path: '/admin/users', name: 'adminUsers', component: () => import('../views/AdminUsers.vue'), meta: { admin: true } },
    { path: '/admin/settings', name: 'adminSettings', component: () => import('../views/AdminSettings.vue'), meta: { admin: true } },
    { path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('../views/NotFound.vue') }
]

const router = createRouter({
    history: createWebHistory(),
    routes,
    scrollBehavior() {
        return { top: 0 }
    }
})

router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('token')
    const isAdmin = localStorage.getItem('is_admin') === 'true'

    // If logged in and trying to access login/register, redirect to home
    if (token && (to.path === '/user/login' || to.path === '/login' || to.path === '/user/register' || to.path === '/register')) {
        next({ path: '/' })
        return
    }

    if (to.meta.auth) {
        if (!token) {
            next({ path: '/user/login', query: { redirect: to.fullPath } })
        } else {
            next()
        }
    } else if (to.meta.admin) {
        if (!token || !isAdmin) {
            next({ path: '/admin/login' })
        } else {
            next()
        }
    } else {
        next()
    }
})

export default router
