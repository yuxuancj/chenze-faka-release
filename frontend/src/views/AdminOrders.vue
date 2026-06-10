<template>
    <AdminLayout page-title="订单管理">
        <div class="space-y-4">
            <div class="flex flex-col md:flex-row md:items-center gap-3">
                <div class="flex items-center gap-2">
                    <input
                        v-model="keyword"
                        type="text"
                        placeholder="搜索订单"
                        class="form-input w-48"
                        @keyup.enter="search"
                    >
                    <button @click="search" class="btn-primary btn-sm">搜索</button>
                </div>
            </div>

            <div v-if="loading" class="card p-8 text-center text-gray-500">
                加载中...
            </div>
            <div v-else class="card">
                <div class="card-body overflow-x-auto">
                    <table class="table w-full">
                        <thead>
                            <tr>
                                <th>订单号</th>
                                <th>商品</th>
                                <th>数量</th>
                                <th>金额</th>
                                <th>状态</th>
                                <th>创建时间</th>
                                <th>操作</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-if="orders.length === 0">
                                <td colspan="7" class="text-center text-gray-500 py-8">暂无订单</td>
                            </tr>
                            <tr v-for="order in orders" :key="order.id">
                                <td>{{ order.order_no }}</td>
                                <td>{{ order.product_snapshot || order.product_name || '-' }}</td>
                                <td>{{ order.quantity }}</td>
                                <td>￥{{ order.amount }}</td>
                                <td>
                                    <span v-if="order.status === 0" class="badge-yellow">待支付</span>
                                    <span v-else-if="order.status === 1" class="badge-green">已支付</span>
                                    <span v-else-if="order.status === 2" class="badge-blue">已完成</span>
                                    <span v-else-if="order.status === 3" class="badge-red">已关闭</span>
                                    <span v-else class="badge-gray">未知</span>
                                </td>
                                <td>{{ order.created_at }}</td>
                                <td>
                                    <button @click="showDetail(order)" class="btn-sm btn-primary">查看</button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>

            <div v-if="pagination && pagination.total > pagination.size" class="flex items-center justify-center gap-2">
                <button @click="prevPage" :disabled="page <= 1" class="btn-sm btn-secondary">上一页</button>
                <span class="text-sm text-gray-600">第 {{ page }} / {{ totalPages }} 页</span>
                <button @click="nextPage" :disabled="page >= totalPages" class="btn-sm btn-secondary">下一页</button>
            </div>

            <div v-if="detail" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-20 p-4">
                <div class="card w-full max-w-2xl max-h-[80vh] overflow-y-auto">
                    <div class="card-header flex items-center justify-between">
                        <span class="font-semibold">订单详情</span>
                        <button @click="detail = null" class="btn-sm btn-secondary">关闭</button>
                    </div>
                    <div class="card-body">
                        <table class="table w-full">
                            <tbody>
                                <tr><td class="w-32 text-gray-500">订单号</td><td>{{ detail.order_no }}</td></tr>
                                <tr><td class="text-gray-500">商品</td><td>{{ detail.product_snapshot || detail.product_name || '-' }}</td></tr>
                                <tr><td class="text-gray-500">数量</td><td>{{ detail.quantity }}</td></tr>
                                <tr><td class="text-gray-500">金额</td><td>￥{{ detail.amount }}</td></tr>
                                <tr><td class="text-gray-500">状态</td>
                                    <td>
                                        <span v-if="detail.status === 0" class="badge-yellow">待支付</span>
                                        <span v-else-if="detail.status === 1" class="badge-green">已支付</span>
                                        <span v-else-if="detail.status === 2" class="badge-blue">已完成</span>
                                        <span v-else-if="detail.status === 3" class="badge-red">已关闭</span>
                                        <span v-else class="badge-gray">未知</span>
                                    </td>
                                </tr>
                                <tr><td class="text-gray-500">邮箱</td><td>{{ detail.email }}</td></tr>
                                <tr><td class="text-gray-500">备注</td><td>{{ detail.remark || '-' }}</td></tr>
                                <tr><td class="text-gray-500">创建时间</td><td>{{ detail.created_at }}</td></tr>
                            </tbody>
                        </table>
                        <div v-if="detail.cards && detail.cards.length > 0" class="mt-4">
                            <h4 class="font-semibold mb-2">卡密</h4>
                            <table class="table w-full">
                                <thead>
                                    <tr><th>卡密</th></tr>
                                </thead>
                                <tbody>
                                    <tr v-for="(c, i) in detail.cards" :key="i">
                                        <td class="font-mono text-sm">{{ c.card_data }}</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import AdminLayout from '../components/AdminLayout.vue'
import { adminOrderList, adminOrderDetail } from '../api/admin'

const orders = ref([])
const detail = ref(null)
const page = ref(1)
const pageSize = 20
const pagination = reactive({ total: 0, size: 20 })
const keyword = ref('')
const loading = ref(false)

const totalPages = computed(() => Math.ceil(pagination.total / pagination.size) || 1)

function load() {
    loading.value = true
    adminOrderList(page.value, pageSize, keyword.value).then((data) => {
        orders.value = (data && data.list) ? data.list : []
        pagination.total = (data && data.total) ? data.total : 0
        pagination.size = (data && data.size) ? data.size : pageSize
    }).catch(() => {
        orders.value = []
        pagination.total = 0
    }).finally(() => {
        loading.value = false
    })
}

function search() {
    page.value = 1
    load()
}

function showDetail(order) {
    if (!order.id) {
        detail.value = order
        return
    }
    adminOrderDetail(order.id).then((data) => {
        if (data && data.order) {
            detail.value = { ...order, ...data.order, cards: data.cards || [] }
        } else if (data) {
            detail.value = data
        } else {
            detail.value = order
        }
    }).catch(() => {
        detail.value = order
    })
}

function prevPage() {
    if (page.value > 1) {
        page.value--
        load()
    }
}

function nextPage() {
    if (page.value < totalPages.value) {
        page.value++
        load()
    }
}

onMounted(load)
</script>
