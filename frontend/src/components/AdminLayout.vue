<template>
    <div class="min-h-screen bg-gray-100 flex">
        <aside
            :class="[
                mobileSidebar ? 'fixed inset-y-0 left-0 z-20' : 'hidden md:flex',
                'w-56 bg-white border-r border-gray-200 flex-col'
            ]"
        >
            <div class="h-14 flex items-center justify-between px-4 border-b border-gray-200">
                <router-link to="/admin/" class="text-lg font-bold text-gray-800">管理后台</router-link>
                <button class="md:hidden text-gray-500 text-sm" @click="mobileSidebar = false">关闭</button>
            </div>
            <nav class="flex-1 py-3 overflow-y-auto">
                <router-link to="/admin/" class="nav-link block mx-2 rounded" @click="mobileSidebar = false">仪表盘</router-link>
                <router-link to="/admin/products" class="nav-link block mx-2 rounded" @click="mobileSidebar = false">商品管理</router-link>
                <router-link to="/admin/cards" class="nav-link block mx-2 rounded" @click="mobileSidebar = false">卡密管理</router-link>
                <router-link to="/admin/categories" class="nav-link block mx-2 rounded" @click="mobileSidebar = false">分类管理</router-link>
                <router-link to="/admin/orders" class="nav-link block mx-2 rounded" @click="mobileSidebar = false">订单管理</router-link>
                <router-link to="/admin/users" class="nav-link block mx-2 rounded" @click="mobileSidebar = false">用户管理</router-link>
                <router-link to="/admin/coupons" class="nav-link block mx-2 rounded" @click="mobileSidebar = false">优惠券管理</router-link>
                <router-link to="/admin/seckills" class="nav-link block mx-2 rounded" @click="mobileSidebar = false">秒杀活动</router-link>
                <router-link to="/admin/distribution" class="nav-link block mx-2 rounded" @click="mobileSidebar = false">分销设置</router-link>
                <router-link to="/admin/points-settings" class="nav-link block mx-2 rounded" @click="mobileSidebar = false">积分设置</router-link>
                <router-link to="/admin/settings" class="nav-link block mx-2 rounded" @click="mobileSidebar = false">系统设置</router-link>
            </nav>
            <div class="p-3 border-t border-gray-200 text-sm">
                <div class="flex items-center justify-between">
                    <span class="text-gray-600">管理员</span>
                    <button @click="logout" class="text-blue-600 hover:underline">退出</button>
                </div>
            </div>
        </aside>

        <div class="flex-1 flex flex-col min-w-0">
            <header class="h-14 bg-white border-b border-gray-200 flex items-center px-4 md:hidden">
                <button @click="mobileSidebar = !mobileSidebar" class="btn-secondary btn-sm">菜单</button>
                <span class="ml-3 font-semibold text-gray-800">{{ pageTitle }}</span>
            </header>

            <main class="flex-1 p-4 md:p-6 overflow-x-auto">
                <div class="hidden md:block mb-4">
                    <h1 class="text-2xl font-bold text-gray-800">{{ pageTitle }}</h1>
                </div>
                <slot></slot>
            </main>
        </div>

        <div
            v-if="mobileSidebar"
            class="fixed inset-0 bg-black bg-opacity-50 z-10 md:hidden"
            @click="mobileSidebar = false"
        ></div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

const props = defineProps({
    pageTitle: {
        type: String,
        default: '管理后台'
    }
})

const router = useRouter()
const userStore = useUserStore()
const mobileSidebar = ref(false)

function logout() {
    userStore.logout()
    router.push('/admin/login')
}
</script>
