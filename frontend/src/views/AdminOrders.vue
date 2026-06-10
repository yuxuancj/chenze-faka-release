<template>
    <AdminLayout page-title="订单管理">
        <div class="space-y-4">
            <div class="card">
                <div class="card-body">
                    <table class="table">
                        <thead>
                            <tr>
                                <th>订单号</th>
                                <th>用户</th>
                                <th>商品</th>
                                <th>数量</th>
                                <th>金额</th>
                                <th>邮箱</th>
                                <th>状态</th>
                                <th>创建时间</th>
                                <th>操作</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-if="orders.length === 0">
                                <td colspan="9" class="text-center text-gray-500 py-8">暂无订单</td>
                            </tr>
                            <tr v-for="order in orders" :key="order.order_no">
                                <td>{{ order.order_no }}</td>
                                <td>{{ order.user_email || '-' }}</td>
                                <td>{{ order.product_name || '-' }}</td>
                                <td>{{ order.quantity }}</td>
                                <td>￥{{ order.amount }}</td>
                                <td>{{ order.email }}</td>
                                <td>
                                    <span v-if="order.status === 1" class="badge-yellow">待支付</span>
                                    <span v-else-if="order.status === 2" class="badge-green">已完成</span>
                                    <span v-else-if="order.status === 3" class="badge-red">已取消</span>
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
            <div v-if="pagination && pagination.total > pagination.size" class="flex items-center justify-center space-x-2">
                <button @click="prevPage" :disabled="page <= 1" class="btn-sm btn-secondary">上一页</button>
                <span class="text-sm text-gray-600">第 {{ page }} / {{ totalPages }} 页</span>
                <button @click="nextPage" :disabled="page >= totalPages" class="btn-sm btn-secondary">下一页</button>
            </div>

            <div v-if="detail" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-20 p-4">
                <div class="card w-full max-w-2xl">
                    <div class="card-header flex items-center justify-between">
                        <span class="font-semibold">订单详情</span>
                        <button @click="detail = null" class="btn-sm btn-secondary">关闭</button>
                    </div>
                    <div class="card-body">
                        <table class="table">
                            <tbody>
                                <tr><td class="w-32 text-gray-500">订单号</td><td>{{ detail.order_no }}</td></tr>
                                <tr><td class="text-gray-500">用户邮箱</td><td>{{ detail.user_email || '-' }}</td></tr>
                                <tr><td class="text-gray-500">商品名称</td><td>{{ detail.product_name }}</td></tr>
                                <tr><td class="text-gray-500">数量</td><td>{{ detail.quantity }}</td></tr>
                                <tr><td class="text-gray-500">金额</td><td>￥{{ detail.amount }}</td></tr>
                                <tr><td class="text-gray-500">支付方式</td><td>{{ detail.pay_type }}</td></tr>
                                <tr><td class="text-gray-500">状态</td>
                                    <td>
                                        <span v-if="detail.status === 1" class="badge-yellow">待支付</span>
                                        <span v-else-if="detail.status === 2" class="badge-green">已完成</span>
                                        <span v-else-if="detail.status === 3" class="badge-red">已取消</span>
                                    </td>
                                </tr>
                                <tr><td class="text-gray-500">邮箱</td><td>{{ detail.email }}</td></tr>
                                <tr><td class="text-gray-500">备注</td><td>{{ detail.remark || '-' }}</td></tr>
                                <tr><td class="text-gray-500">创建时间</td><td>{{ detail.created_at }}</td></tr>
                            </tbody>
                        </table>
                        <div v-if="detail.cards && detail.cards.length > 0" class="mt-4">
                            <h4 class="font-semibold mb-2">卡密</h4>
                            <table class="table">
                                <thead>
                                    <tr>
                                        <th>卡号</th>
                                        <th>密码</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr v-for="(c, i) in detail.cards" :key="i">
                                        <td>{{ c.card_no }}</td>
                                        <td>{{ c.card_pwd }}</td>
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
import { ref, computed, onMounted } from 'vue'
import AdminLayout from '../components/AdminLayout.vue'
import { adminOrderList } from '../api/admin'

const orders = ref([])
const detail = ref(null)
const page = ref(1)
const pageSize = 20
const pagination = ref({ total: 0, size: 20 })

const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.size) || 1)

function load() {
    adminOrderList(page.value, pageSize).then((res) => {
        orders.value = res.data && res.data.list ? res.data.list : []
        pagination.value.total = (res.data && res.data.total) ? res.data.total : 0
    }).catch(() => {})
}

function showDetail(order) {
    detail.value = order
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
