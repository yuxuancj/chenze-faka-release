<template>
    <div class="min-h-screen bg-gray-50">
        <nav class="bg-white border-b border-gray-200 sticky top-0 z-10">
            <div class="max-w-7xl mx-auto px-4">
                <div class="flex items-center justify-between h-14">
                    <router-link to="/" class="text-xl font-bold text-gray-800">
                        {{ siteName }}
                    </router-link>

                    <div class="hidden md:flex items-center space-x-6">
                        <router-link to="/" class="nav-link">首页</router-link>
                        <router-link to="/products" class="nav-link">商品</router-link>
                        <template v-if="token">
                            <router-link to="/user/orders" class="nav-link">我的订单</router-link>
                            <router-link to="/user/profile" class="nav-link">个人中心</router-link>
                            <button @click="logout" class="btn-secondary btn-sm">退出</button>
                        </template>
                        <template v-else>
                            <router-link to="/user/login" class="nav-link">登录</router-link>
                            <router-link to="/user/register" class="nav-link">注册</router-link>
                        </template>
                        <router-link to="/cart" class="btn-primary btn-sm">
                            购物车
                            <span v-if="cartCount > 0" class="ml-1 bg-white text-blue-600 rounded-full text-xs w-5 h-5 inline-flex items-center justify-center">{{ cartCount }}</span>
                        </router-link>
                    </div>

                    <button class="md:hidden btn-secondary btn-sm" @click="mobileMenuOpen = !mobileMenuOpen">
                        {{ mobileMenuOpen ? '关闭' : '菜单' }}
                    </button>
                </div>

                <div v-if="mobileMenuOpen" class="md:hidden border-t border-gray-200 py-3 space-y-2">
                    <router-link to="/" class="nav-link block">首页</router-link>
                    <router-link to="/products" class="nav-link block">商品</router-link>
                    <router-link to="/cart" class="nav-link block">购物车</router-link>
                    <template v-if="token">
                        <router-link to="/user/orders" class="nav-link block">我的订单</router-link>
                        <router-link to="/user/profile" class="nav-link block">个人中心</router-link>
                        <button @click="logout" class="nav-link block w-full text-left">退出登录</button>
                    </template>
                    <template v-else>
                        <router-link to="/user/login" class="nav-link block">登录</router-link>
                        <router-link to="/user/register" class="nav-link block">注册</router-link>
                    </template>
                </div>
            </div>
        </nav>

        <main class="max-w-7xl mx-auto px-4 py-6">
            <slot></slot>
        </main>

        <footer class="bg-white border-t border-gray-200 mt-8">
            <div class="max-w-7xl mx-auto px-4 py-6 text-center text-sm text-gray-500">
                版权所有 &copy; {{ new Date().getFullYear() }} {{ siteName }}
            </div>
        </footer>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useCartStore } from '../stores/cart'

const mobileMenuOpen = ref(false)
const cartStore = useCartStore()
const siteName = ref('发卡平台')
const token = computed(() => localStorage.getItem('token'))
const cartCount = computed(() => cartStore.totalCount)

function logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('is_admin')
    window.location.href = '/'
}
</script>
