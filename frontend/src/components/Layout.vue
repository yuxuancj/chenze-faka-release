<template>
    <div class="min-h-screen bg-gray-50">
        <header class="bg-white border-b border-gray-200 sticky top-0 z-10">
            <div class="max-w-7xl mx-auto px-4 sm:px-6">
                <div class="flex items-center justify-between h-16">
                    <div class="flex items-center space-x-4">
                        <router-link to="/" class="text-xl font-bold text-blue-600">商城</router-link>
                        <nav class="hidden md:flex space-x-2">
                            <router-link to="/" class="nav-link">首页</router-link>
                            <router-link to="/products" class="nav-link">商品列表</router-link>
                            <router-link to="/cart" class="nav-link">
                                购物车
                                <span v-if="cartStore.totalCount > 0" class="ml-1 badge-red">({{ cartStore.totalCount }})</span>
                            </router-link>
                        </nav>
                    </div>
                    <div class="flex items-center space-x-2">
                        <template v-if="token">
                            <router-link to="/user/profile" class="nav-link">个人中心</router-link>
                            <router-link to="/user/orders" class="nav-link">我的订单</router-link>
                            <button @click="logout" class="btn-sm btn-secondary">退出</button>
                        </template>
                        <template v-else>
                            <router-link to="/user/login" class="btn-sm btn-primary">登录</router-link>
                            <router-link to="/user/register" class="btn-sm btn-secondary">注册</router-link>
                        </template>
                    </div>
                </div>
            </div>
        </header>
        <main class="max-w-7xl mx-auto px-4 sm:px-6 py-6">
            <slot />
        </main>
        <footer class="border-t border-gray-200 bg-white mt-12">
            <div class="max-w-7xl mx-auto px-4 sm:px-6 py-6 text-center text-sm text-gray-500">
                商城系统
            </div>
        </footer>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import { useCartStore } from '../stores/cart'

const cartStore = useCartStore()

const token = computed(() => localStorage.getItem('token'))

function logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('is_admin')
    window.location.href = '/'
}
</script>
