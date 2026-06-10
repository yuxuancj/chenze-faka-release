<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">我的订单</h2>
            <el-card shadow="never" v-loading="loading">
                <el-table :data="orders" style="width: 100%" stripe empty-text="暂无订单">
                    <el-table-column prop="order_no" label="订单号" min-width="180">
                        <template #default="scope">
                            {{ scope.row.order_no || '-' }}
                        </template>
                    </el-table-column>
                    <el-table-column label="商品" min-width="200">
                        <template #default="scope">
                            {{ scope.row.product_snapshot || scope.row.product_name || '-' }}
                        </template>
                    </el-table-column>
                    <el-table-column label="数量" width="80">
                        <template #default="scope">
                            {{ scope.row.quantity || 1 }}
                        </template>
                    </el-table-column>
                    <el-table-column label="金额" width="100">
                        <template #default="scope">
                            ￥{{ scope.row.amount }}
                        </template>
                    </el-table-column>
                    <el-table-column label="状态" width="100">
                        <template #default="scope">
                            <el-tag v-if="scope.row.status === 0" type="warning">待支付</el-tag>
                            <el-tag v-else-if="scope.row.status === 1" type="success">已支付</el-tag>
                            <el-tag v-else-if="scope.row.status === 2" type="primary">已完成</el-tag>
                            <el-tag v-else-if="scope.row.status === 3" type="danger">已关闭</el-tag>
                            <el-tag v-else type="info">未知</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="created_at" label="创建时间" min-width="180">
                        <template #default="scope">
                            {{ scope.row.created_at }}
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="100" fixed="right">
                        <template #default="scope">
                            <el-button type="primary" link @click="viewOrder(scope.row.order_no)">
                                查看
                            </el-button>
                        </template>
                    </el-table-column>
                </el-table>
                <div class="mt-4 flex justify-end">
                    <el-pagination
                        v-model:current-page="page"
                        v-model:page-size="pageSize"
                        :page-sizes="[10, 20, 50, 100]"
                        :total="pagination.total"
                        layout="total, sizes, prev, pager, next, jumper"
                        @size-change="handleSizeChange"
                        @current-change="handleCurrentChange"
                    />
                </div>
            </el-card>
        </div>
    </Layout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Layout from '../components/Layout.vue'
import { orderList } from '../api/order'

const router = useRouter()
const orders = ref([])
const page = ref(1)
const pageSize = ref(20)
const pagination = reactive({ total: 0, size: 20 })
const loading = ref(false)

function viewOrder(orderNo) {
    router.push('/order/' + orderNo)
}

function loadOrders() {
    loading.value = true
    orderList(page.value, pageSize.value).then((data) => {
        orders.value = (data && data.list) ? data.list : []
        pagination.total = (data && data.total) ? data.total : 0
        pagination.size = (data && data.size) ? data.size : pageSize.value
    }).catch(() => {
        orders.value = []
        pagination.total = 0
    }).finally(() => {
        loading.value = false
    })
}

function handleSizeChange(val) {
    pageSize.value = val
    page.value = 1
    loadOrders()
}

function handleCurrentChange(val) {
    page.value = val
    loadOrders()
}

onMounted(loadOrders)
</script>
