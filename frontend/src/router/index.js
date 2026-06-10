import { createRouter, createWebHistory } from 'vue-router'

const routes = [
    { path: '/', name: 'home', component: () => import('../views/Home.vue') },
    { path: '/products', name: 'products', component: () => import('../views/ProductList.vue') },
    { path: '/product/:id', name: 'productDetail', component: () => import('../views/ProductDetail.vue') },
    { path: '/order/:order_no', name: 'orderDetail', component: () => import('../views/OrderDetail.vue') },
    { path: '/cart', name: 'cart', component: () => import('../views/Cart.vue') },
    { path: '/checkout', name: 'checkout', component: () => import('../views/Checkout.vue') },
    { path: '/seckill', name: 'seckill', component: () => import('../views/Seckill.vue') },
    { path: '/coupon/redeem', name: 'couponRedeem', component: () => import('../views/CouponRedeem.vue') },
    { path: '/user/login', name: 'login', component: () => import('../views/Login.vue'), alias: '/login' },
    { path: '/user/register', name: 'register', component: () => import('../views/Register.vue'), alias: '/register' },
    { path: '/user/profile', name: 'userCenter', component: () => import('../views/UserCenter.vue') },
    { path: '/user/orders', name: 'userOrders', component: () => import('../views/UserOrders.vue') },
    { path: '/user/cards', name: 'userCards', component: () => import('../views/UserCards.vue') },
    { path: '/user/coupons', name: 'userCoupons', component: () => import('../views/UserCoupons.vue') },
    { path: '/user/points', name: 'userPoints', component: () => import('../views/UserPoints.vue') },
    { path: '/user/distribution', name: 'distribution', component: () => import('../views/Distribution.vue') },
    { path: '/admin/login', name: 'adminLogin', component: () => import('../views/AdminLogin.vue') },
    { path: '/admin/', name: 'adminDashboard', component: () => import('../views/AdminDashboard.vue') },
    { path: '/admin/products', name: 'adminProducts', component: () => import('../views/AdminProducts.vue') },
    { path: '/admin/product/new', name: 'adminProductNew', component: () => import('../views/AdminProductEdit.vue') },
    { path: '/admin/product/:id', name: 'adminProductEdit', component: () => import('../views/AdminProductEdit.vue') },
    { path: '/admin/cards', name: 'adminCards', component: () => import('../views/AdminCards.vue') },
    { path: '/admin/categories', name: 'adminCategories', component: () => import('../views/AdminCategories.vue') },
    { path: '/admin/orders', name: 'adminOrders', component: () => import('../views/AdminOrders.vue') },
    { path: '/admin/users', name: 'adminUsers', component: () => import('../views/AdminUsers.vue') },
    { path: '/admin/coupons', name: 'adminCoupons', component: () => import('../views/AdminCoupons.vue') },
    { path: '/admin/seckills', name: 'adminSeckills', component: () => import('../views/AdminSeckills.vue') },
    { path: '/admin/distribution', name: 'adminDistribution', component: () => import('../views/AdminDistribution.vue') },
    { path: '/admin/points-settings', name: 'adminPointsSettings', component: () => import('../views/AdminPointsSettings.vue') },
    { path: '/admin/settings', name: 'adminSettings', component: () => import('../views/AdminSettings.vue') },
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

    if (token && (to.path === '/user/login' || to.path === '/login' || to.path === '/user/register' || to.path === '/register')) {
        next({ path: '/' })
        return
    }

    const authPaths = ['/cart', '/checkout', '/user/profile', '/user/orders', '/user/cards', '/user/coupons', '/user/points', '/user/distribution']
    if (authPaths.some((p) => to.path === p || to.path.startsWith(p + '/'))) {
        if (!token) {
            next({ path: '/user/login', query: { redirect: to.fullPath } })
        } else {
            next()
        }
        return
    }

    if (to.path.startsWith('/admin/') && to.path !== '/admin/login') {
        if (!token || !isAdmin) {
            next({ path: '/admin/login' })
        } else {
            next()
        }
        return
    }

    next()
})

export default router
