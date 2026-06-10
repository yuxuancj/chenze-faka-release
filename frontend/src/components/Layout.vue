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
                        <router-link to="/cart" class="btn-primary btn-sm">
                            购物车
                            <span v-if="cartCount > 0" class="ml-1 bg-white text-blue-600 rounded-full text-xs w-5 h-5 inline-flex items-center justify-center">{{ cartCount }}</span>
                        </router-link>
                        <template v-if="token">
                            <div class="relative" ref="userMenuRef">
                                <button @click="userMenuOpen = !userMenuOpen" class="nav-link flex items-center gap-1">
                                    <span>{{ userNickname || '用户' }}</span>
                                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path></svg>
                                </button>
                                <div v-if="userMenuOpen" class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg border border-gray-200 py-1">
                                    <router-link to="/user/profile" @click="userMenuOpen = false" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">个人资料</router-link>
                                    <router-link to="/user/orders" @click="userMenuOpen = false" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">我的订单</router-link>
                                    <router-link to="/user/distribution" @click="userMenuOpen = false" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">推广中心</router-link>
                                    <router-link to="/user/signin" @click="userMenuOpen = false" class="block px-4 py-2 text-sm text-green-600 hover:bg-gray-100 font-medium">每日签到</router-link>
                                    <router-link to="/user/coupons" @click="userMenuOpen = false" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">我的优惠券</router-link>
                                    <router-link to="/user/points" @click="userMenuOpen = false" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">积分明细</router-link>
                                    <div class="border-t border-gray-100 my-1"></div>
                                    <button @click="logout" class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">退出登录</button>
                                </div>
                            </div>
                        </template>
                        <template v-else>
                            <router-link to="/user/login" class="nav-link">登录</router-link>
                            <router-link to="/user/register" class="nav-link">注册</router-link>
                        </template>
                    </div>

                    <button class="md:hidden btn-secondary btn-sm" @click="mobileMenuOpen = !mobileMenuOpen">
                        {{ mobileMenuOpen ? '关闭' : '菜单' }}
                    </button>
                </div>

                <div v-if="mobileMenuOpen" class="md:hidden border-t border-gray-200 py-3 space-y-1">
                    <router-link to="/" class="nav-link block">首页</router-link>
                    <router-link to="/products" class="nav-link block">商品</router-link>
                    <router-link to="/cart" class="nav-link block">购物车</router-link>
                    <template v-if="token">
                        <div class="border-t border-gray-200 pt-2 mt-2">
                            <router-link to="/user/profile" class="nav-link block">个人资料</router-link>
                            <router-link to="/user/orders" class="nav-link block">我的订单</router-link>
                            <router-link to="/user/distribution" class="nav-link block">推广中心</router-link>
                            <router-link to="/user/signin" class="nav-link block text-green-600">每日签到</router-link>
                            <router-link to="/user/coupons" class="nav-link block">我的优惠券</router-link>
                            <router-link to="/user/points" class="nav-link block">积分明细</router-link>
                            <button @click="logout" class="nav-link block w-full text-left">退出登录</button>
                        </div>
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
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { useUserStore } from '../stores/user'

const router = useRouter()
const mobileMenuOpen = ref(false)
const userMenuOpen = ref(false)
const userMenuRef = ref(null)
const cartStore = useCartStore()
const userStore = useUserStore()
const siteName = ref('发卡平台')

function logout() {
    userMenuOpen.value = false
    userStore.logout()
    router.push('/')
}

const token = computed(() => userStore.token)
const cartCount = computed(() => cartStore.totalCount)
const userNickname = computed(() => userStore.userInfo?.nickname || userStore.userInfo?.username || '')

function handleClickOutside(event) {
    if (userMenuRef.value && !userMenuRef.value.contains(event.target)) {
        userMenuOpen.value = false
    }
}

onMounted(() => {
    document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside)
})
</script>
