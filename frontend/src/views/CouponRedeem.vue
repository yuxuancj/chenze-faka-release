<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">兑换优惠券</h2>

            <el-card shadow="hover">
                <div class="space-y-4">
                    <el-input v-model="code" placeholder="请输入兑换码" />
                    <el-button type="primary" @click="redeem" :loading="loading">
                        立即兑换
                    </el-button>
                </div>
            </el-card>

            <el-card shadow="hover">
                <template #header>
                    <span class="font-semibold">可领取优惠券</span>
                </template>
                <el-empty v-if="listLoading" description="加载中..." />
                <el-empty v-else-if="coupons.length === 0" description="暂无可领取优惠券" />
                <el-row v-else :gutter="16">
                    <el-col v-for="c in coupons" :key="c.id" :xs="24" :sm="12" :md="8">
                        <el-card shadow="hover">
                            <div class="text-2xl font-bold text-orange-600">{{ formatDiscount(c) }}</div>
                            <div class="text-xs text-gray-500 mt-1">满 {{ c.min_amount || 0 }} 元可用</div>
                            <div class="text-sm text-gray-700 mt-2">{{ c.name }}</div>
                            <div class="text-xs text-gray-500 mt-1">有效期至: {{ c.expire_time || '长期有效' }}</div>
                        </el-card>
                    </el-col>
                </el-row>
            </el-card>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { couponList, couponRedeem } from '../api/coupon'
import { ElMessage } from 'element-plus'

const code = ref('')
const loading = ref(false)
const coupons = ref([])
const listLoading = ref(true)

function formatDiscount(c) {
    if (c.type === 2 || c.discount_rate) {
        return (c.discount_rate || c.value) + '折'
    }
    return '￥' + (c.discount_amount || c.value || 0)
}

function redeem() {
    if (!code.value) {
        ElMessage.warning('请输入兑换码')
        return
    }
    loading.value = true
    couponRedeem(code.value).then(() => {
        ElMessage.success('兑换成功')
        code.value = ''
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}

function load() {
    listLoading.value = true
    couponList().then((d) => {
        coupons.value = (d && d.list) ? d.list : (Array.isArray(d) ? d : [])
    }).catch(() => {}).finally(() => {
        listLoading.value = false
    })
}

onMounted(load)
</script>
