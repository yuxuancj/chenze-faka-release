<template>
    <AdminLayout page-title="仪表盘">
        <div v-if="loading" class="card p-8 text-center text-gray-500">
            加载中...
        </div>
        <div v-else>
            <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
                <div class="card">
                    <div class="card-body">
                        <div class="text-sm text-gray-500">用户总数</div>
                        <div class="text-2xl font-bold text-gray-800 mt-2">{{ stats.user_count || 0 }}</div>
                    </div>
                </div>
                <div class="card">
                    <div class="card-body">
                        <div class="text-sm text-gray-500">商品总数</div>
                        <div class="text-2xl font-bold text-blue-600 mt-2">{{ stats.product_count || 0 }}</div>
                    </div>
                </div>
                <div class="card">
                    <div class="card-body">
                        <div class="text-sm text-gray-500">订单总数</div>
                        <div class="text-2xl font-bold text-green-600 mt-2">{{ stats.order_count || 0 }}</div>
                    </div>
                </div>
                <div class="card">
                    <div class="card-body">
                        <div class="text-sm text-gray-500">总收入</div>
                        <div class="text-2xl font-bold text-yellow-600 mt-2">￥{{ stats.total_amount || 0 }}</div>
                    </div>
                </div>
            </div>

            <div class="card">
                <div class="card-header font-semibold">系统信息</div>
                <div class="card-body overflow-x-auto">
                    <table class="table w-full">
                        <tbody>
                            <tr><td class="w-32 text-gray-500">当前时间</td><td>{{ currentTime }}</td></tr>
                            <tr><td class="text-gray-500">版本</td><td>v1.0.0</td></tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AdminLayout from '../components/AdminLayout.vue'
import { adminDashboard } from '../api/admin'

const stats = ref({})
const loading = ref(true)
const currentTime = ref('')

function loadData() {
    const d = new Date()
    currentTime.value = d.toLocaleString('zh-CN')
    loading.value = true
    adminDashboard().then((data) => {
        stats.value = data || {}
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}

onMounted(loadData)
</script>
