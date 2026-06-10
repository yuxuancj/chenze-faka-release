<template>
    <div class="min-h-screen bg-gray-50 flex">
        <aside class="w-64 bg-white border-r border-gray-200 min-h-screen">
            <div class="p-4 border-b border-gray-200">
                <h1 class="text-xl font-bold text-blue-600">管理后台</h1>
            </div>
            <nav class="p-2 space-y-1">
                <router-link to="/admin/" class="nav-link">仪表盘</router-link>
                <router-link to="/admin/products" class="nav-link">商品管理</router-link>
                <router-link to="/admin/cards" class="nav-link">卡密管理</router-link>
                <router-link to="/admin/categories" class="nav-link">分类管理</router-link>
                <router-link to="/admin/orders" class="nav-link">订单管理</router-link>
                <router-link to="/admin/users" class="nav-link">用户管理</router-link>
                <router-link to="/admin/settings" class="nav-link">系统设置</router-link>
            </nav>
            <div class="p-4 border-t border-gray-200 mt-4">
                <button @click="logout" class="btn-sm btn-secondary w-full">退出登录</button>
            </div>
        </aside>
        <div class="flex-1 flex flex-col">
            <header class="bg-white border-b border-gray-200 h-16 px-6 flex items-center justify-between">
                <h2 class="text-lg font-semibold text-gray-800">{{ pageTitle }}</h2>
                <div class="text-sm text-gray-500">管理员</div>
            </header>
            <main class="flex-1 p-6">
                <slot />
            </main>
        </div>
    </div>
</template>

<script setup>
defineProps({
    pageTitle: {
        type: String,
        default: ''
    }
})

function logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('is_admin')
    window.location.href = '/admin/login'
}
</script>
