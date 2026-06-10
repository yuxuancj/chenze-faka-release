<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">购物车</h2>
            <el-card v-if="cartStore.items.length === 0" class="text-center" shadow="never">
                <template #default>
                    <el-empty description="购物车是空的，去选购商品吧。">
                        <el-button type="primary" @click="$router.push('/products')">去购物</el-button>
                    </el-empty>
                </template>
            </el-card>
            <template v-else>
                <el-card shadow="never">
                    <el-table :data="cartStore.items" stripe>
                        <el-table-column label="商品">
                            <template #default="scope">
                                <router-link :to="'/product/' + scope.row.product_id" class="text-blue-600 hover:underline">
                                    {{ scope.row.name }}
                                </router-link>
                            </template>
                        </el-table-column>
                        <el-table-column label="单价" width="120">
                            <template #default="scope">￥{{ scope.row.price }}</template>
                        </el-table-column>
                        <el-table-column label="数量" width="200">
                            <template #default="scope">
                                <el-input-number
                                    v-model="scope.row.quantity"
                                    :min="1"
                                    size="small"
                                    @change="cartStore.updateQuantity(scope.row.product_id, scope.row.quantity)"
                                />
                            </template>
                        </el-table-column>
                        <el-table-column label="小计" width="120">
                            <template #default="scope">￥{{ (scope.row.price * scope.row.quantity).toFixed(2) }}</template>
                        </el-table-column>
                        <el-table-column label="操作" width="120">
                            <template #default="scope">
                                <el-button type="danger" size="small" @click="cartStore.removeItem(scope.row.product_id)">删除</el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-card>
                <el-card shadow="never">
                    <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
                        <div class="text-lg">
                            共 <span class="font-semibold">{{ cartStore.totalCount }}</span> 件商品，
                            合计: <span class="text-blue-600 font-bold text-xl">￥{{ cartStore.totalPrice.toFixed(2) }}</span>
                        </div>
                        <div class="flex items-center gap-3">
                            <el-button @click="handleClearCart">清空购物车</el-button>
                            <el-button type="primary" @click="$router.push('/checkout')">去结算</el-button>
                        </div>
                    </div>
                </el-card>
            </template>
        </div>
    </Layout>
</template>

<script setup>
import Layout from '../components/Layout.vue'
import { useCartStore } from '../stores/cart'
import { ElMessage, ElMessageBox } from 'element-plus'

const cartStore = useCartStore()

function handleClearCart() {
    ElMessageBox.confirm('确定要清空购物车吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
        cartStore.clearCart()
        ElMessage.success('购物车已清空')
    }).catch(() => {})
}
</script>
