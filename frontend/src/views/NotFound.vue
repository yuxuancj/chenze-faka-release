<template>
    <div class="min-h-screen bg-gray-50 flex items-center justify-center px-4">
        <el-card class="w-full max-w-md" shadow="hover">
            <div class="text-center py-8 space-y-4">
                <div class="text-6xl font-bold text-gray-300">404</div>
                <h2 class="text-xl font-semibold text-gray-800">页面不存在</h2>
                <p class="text-gray-500">您访问的页面可能已被删除或链接错误。</p>
                <div class="pt-4 flex items-center justify-center space-x-3">
                    <router-link to="/">
                        <el-button type="primary">返回首页</el-button>
                    </router-link>
                    <router-link to="/products">
                        <el-button>浏览商品</el-button>
                    </router-link>
                </div>
            </div>
        </el-card>
    </div>
</template>

<script setup>
</script>
