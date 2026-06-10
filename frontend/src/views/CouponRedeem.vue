<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">兑换优惠券</h2>

            <div class="card">
                <div class="card-body space-y-4">
                    <div>
                        <label class="form-label">兑换码</label>
                        <input v-model="code" type="text" class="form-input" placeholder="请输入兑换码">
                    </div>
                    <button @click="redeem" :disabled="loading" class="btn-primary">
                        {{ loading ? '兑换中...' : '立即兑换' }}
                    </button>
                </div>
            </div>

            <div class="card">
                <div class="card-header font-semibold">可领取优惠券</div>
                <div class="card-body">
                    <div v-if="listLoading" class="text-center text-gray-500 py-6">加载中...</div>
                    <div v-else-if="coupons.length === 0" class="text-center text-gray-500 py-6">暂无可领取优惠券</div>
                    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-3">
                        <div v-for="c in coupons" :key="c.id" class="border border-gray-200 rounded-md p-4">
                            <div class="text-2xl font-bold text-orange-600">{{ formatDiscount(c) }}</div>
                            <div class="text-xs text-gray-500 mt-1">满 {{ c.min_amount || 0 }} 元可用</div>
                            <div class="text-sm text-gray-700 mt-2">{{ c.name }}</div>
                            <div class="text-xs text-gray-500 mt-1">有效期至: {{ c.expire_time || '长期有效' }}</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { couponList, couponRedeem } from '../api/coupon'

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
        alert('请输入兑换码')
        return
    }
    loading.value = true
    couponRedeem(code.value).then(() => {
        alert('兑换成功')
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
