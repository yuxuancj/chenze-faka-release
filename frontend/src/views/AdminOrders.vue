<template>
    <AdminLayout page-title="订单管理">
        <div class="space-y-4">
            <el-card shadow="never">
                <div class="flex items-center gap-2">
                    <el-input
                        v-model="keyword"
                        placeholder="搜索订单"
                        style="width: 240px"
                        @keyup.enter="search"
                    />
                    <el-button type="primary" @click="search">搜索</el-button>
                </div>
            </el-card>

            <el-card shadow="never">
                <el-table :data="orders" stripe v-loading="loading">
                    <el-table-column prop="order_no" label="订单号" min-width="160" />
                    <el-table-column label="商品" min-width="200">
                        <template #default="scope">
                            {{ scope.row.product_snapshot || scope.row.product_name || '-' }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="quantity" label="数量" width="80" />
                    <el-table-column label="金额" width="120">
                        <template #default="scope">￥{{ scope.row.amount }}</template>
                    </el-table-column>
                    <el-table-column label="状态" width="100">
                        <template #default="scope">
                            <el-tag :type="statusType(scope.row.status)" effect="light">
                                {{ statusText(scope.row.status) }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="created_at" label="创建时间" min-width="160" />
                    <el-table-column label="操作" width="100" fixed="right">
                        <template #default="scope">
                            <el-button type="primary" link @click="showDetail(scope.row)">查看</el-button>
                        </template>
                    </el-table-column>
                </el-table>

                <div class="mt-4 flex items-center justify-end">
                    <el-pagination
                        v-if="pagination.total > 0"
                        v-model:current-page="page"
                        v-model:page-size="pagination.size"
                        :total="pagination.total"
                        layout="prev, pager, next"
                        @current-change="load"
                    />
                </div>
            </el-card>

            <el-dialog
                v-model="detailVisible"
                title="订单详情"
                width="600px"
            >
                <template v-if="detail">
                    <el-descriptions :column="1" border>
                        <el-descriptions-item label="订单号">{{ detail.order_no }}</el-descriptions-item>
                        <el-descriptions-item label="商品">
                            {{ detail.product_snapshot || detail.product_name || '-' }}
                        </el-descriptions-item>
                        <el-descriptions-item label="数量">{{ detail.quantity }}</el-descriptions-item>
                        <el-descriptions-item label="金额">￥{{ detail.amount }}</el-descriptions-item>
                        <el-descriptions-item label="状态">
                            <el-tag :type="statusType(detail.status)" effect="light">
                                {{ statusText(detail.status) }}
                            </el-tag>
                        </el-descriptions-item>
                        <el-descriptions-item label="邮箱">{{ detail.email }}</el-descriptions-item>
                        <el-descriptions-item label="备注">{{ detail.remark || '-' }}</el-descriptions-item>
                        <el-descriptions-item label="创建时间">{{ detail.created_at }}</el-descriptions-item>
                    </el-descriptions>
                    <div v-if="detail.cards && detail.cards.length > 0" class="mt-4">
                        <h4 class="font-semibold mb-2">卡密</h4>
                        <el-table :data="detail.cards" stripe size="small">
                            <el-table-column label="卡号">
                                <template #default="scope">
                                    <span class="font-mono text-sm">{{ scope.row.card_data }}</span>
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>
                </template>
            </el-dialog>
        </div>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import AdminLayout from '../components/AdminLayout.vue'
import { adminOrderList, adminOrderDetail } from '../api/admin'

const orders = ref([])
const detail = ref(null)
const detailVisible = ref(false)
const page = ref(1)
const pageSize = 20
const pagination = reactive({ total: 0, size: 20 })
const keyword = ref('')
const loading = ref(false)

function statusType(status) {
    if (status === 0) return 'warning'
    if (status === 1) return 'success'
    if (status === 2) return 'primary'
    if (status === 3) return 'danger'
    return 'info'
}

function statusText(status) {
    if (status === 0) return '待支付'
    if (status === 1) return '已支付'
    if (status === 2) return '已完成'
    if (status === 3) return '已关闭'
    return '未知'
}

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
        detailVisible.value = true
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
        detailVisible.value = true
    }).catch(() => {
        detail.value = order
        detailVisible.value = true
    })
}

onMounted(load)
</script>
