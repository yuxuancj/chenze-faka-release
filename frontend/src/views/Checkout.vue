<template>
    <Layout>
        <div class="space-y-4">
            <h2 class="text-xl font-bold text-gray-800">结算</h2>

            <el-card v-if="loadingProduct && fromProduct" shadow="never">
                <el-skeleton active :rows="3" />
            </el-card>
            <el-card v-else-if="!hasItems" class="text-center" shadow="never">
                <el-empty description="没有可以结算的商品。">
                    <el-button type="primary" @click="$router.push('/products')">去购物</el-button>
                </el-empty>
            </el-card>

            <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div class="md:col-span-2 space-y-4">
                    <el-card shadow="never">
                        <template #header>
                            <span class="font-semibold">订单信息</span>
                        </template>
                        <el-form :model="form" label-width="100px">
                            <el-form-item label="联系邮箱" :error="emailError">
                                <el-input v-model="form.email" placeholder="请输入接收卡密的邮箱" />
                            </el-form-item>
                            <el-form-item label="支付方式">
                                <el-select v-model="form.pay_type" style="width: 100%">
                                    <el-option label="在线支付" value="epay" />
                                    <el-option label="余额支付" value="balance" />
                                </el-select>
                            </el-form-item>
                            <el-form-item label="备注">
                                <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="选填" />
                            </el-form-item>
                        </el-form>
                    </el-card>

                    <el-card shadow="never">
                        <template #header>
                            <span class="font-semibold">选择优惠券</span>
                        </template>
                        <div v-if="couponsLoading" class="text-sm text-gray-500">加载中...</div>
                        <el-empty v-else-if="availableCoupons.length === 0" description="暂无可使用的优惠券" :image-size="80" />
                        <div v-else class="space-y-2">
                            <el-radio-group v-model="selectedCouponId" class="w-full">
                                <div
                                    v-for="c in availableCoupons"
                                    :key="c.id"
                                    class="flex items-start gap-3 border border-gray-200 rounded-md p-3 cursor-pointer hover:bg-gray-50"
                                    @click="selectedCouponId = c.id"
                                >
                                    <el-radio :value="c.id" />
                                    <div class="flex-1">
                                        <div class="text-orange-600 font-bold">
                                            {{ formatDiscount(c) }}
                                            <span class="text-xs text-gray-500 font-normal ml-2">满 {{ c.min_amount || 0 }} 元可用</span>
                                        </div>
                                        <div class="text-xs text-gray-600">{{ c.name }} · 有效期至 {{ c.expire_time || '长期有效' }}</div>
                                    </div>
                                </div>
                                <div
                                    class="flex items-center gap-3 border border-gray-200 rounded-md p-3 cursor-pointer hover:bg-gray-50"
                                    @click="selectedCouponId = 0"
                                >
                                    <el-radio :value="0" />
                                    <span class="text-sm text-gray-600">不使用优惠券</span>
                                </div>
                            </el-radio-group>
                        </div>
                    </el-card>

                    <el-card shadow="never">
                        <template #header>
                            <span class="font-semibold">商品清单</span>
                        </template>
                        <el-table :data="items" stripe>
                            <el-table-column prop="name" label="商品" />
                            <el-table-column label="单价" width="120">
                                <template #default="scope">￥{{ scope.row.price }}</template>
                            </el-table-column>
                            <el-table-column prop="quantity" label="数量" width="100" />
                            <el-table-column label="小计" width="150">
                                <template #default="scope">￥{{ (scope.row.price * scope.row.quantity).toFixed(2) }}</template>
                            </el-table-column>
                        </el-table>
                    </el-card>
                </div>

                <el-card class="h-fit" shadow="never">
                    <template #header>
                        <span class="font-semibold">结算</span>
                    </template>
                    <div class="space-y-3">
                        <div class="flex justify-between text-sm">
                            <span class="text-gray-600">商品总数</span>
                            <span>{{ totalCount }} 件</span>
                        </div>
                        <div class="flex justify-between text-sm">
                            <span class="text-gray-600">商品总额</span>
                            <span>￥{{ totalPrice.toFixed(2) }}</span>
                        </div>
                        <div v-if="couponDiscount > 0" class="flex justify-between text-sm text-orange-600">
                            <span>优惠券</span>
                            <span>-￥{{ couponDiscount.toFixed(2) }}</span>
                        </div>
                        <div class="flex justify-between text-lg font-bold pt-2 border-t border-gray-200">
                            <span class="text-gray-800">应付</span>
                            <span class="text-blue-600">￥{{ finalPrice.toFixed(2) }}</span>
                        </div>
                        <el-button type="primary" class="w-full" :loading="loading" @click="submitOrder">
                            提交订单
                        </el-button>
                    </div>
                </el-card>
            </div>
        </div>
    </Layout>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Layout from '../components/Layout.vue'
import { useCartStore } from '../stores/cart'
import { createOrder, payOrder } from '../api/order'
import { productDetail } from '../api/product'
import { userCoupons } from '../api/coupon'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()
const loading = ref(false)
const loadingProduct = ref(false)
const couponsLoading = ref(false)
const emailError = ref('')
const form = reactive({ email: '', pay_type: 'epay', remark: '' })
const productItem = ref(null)
const availableCoupons = ref([])
const selectedCouponId = ref(0)

const fromProduct = computed(() => !!route.query.product_id)

const items = computed(() => {
    if (productItem.value) {
        return [productItem.value]
    }
    return cartStore.items
})

const totalCount = computed(() => items.value.reduce((sum, item) => sum + item.quantity, 0))
const totalPrice = computed(() => items.value.reduce((sum, item) => sum + item.price * item.quantity, 0))
const hasItems = computed(() => items.value.length > 0)

const couponDiscount = computed(() => {
    if (!selectedCouponId.value) return 0
    const c = availableCoupons.value.find((x) => x.id === selectedCouponId.value)
    if (!c) return 0
    if (c.type === 2 || c.discount_rate) {
        const rate = (c.discount_rate || c.value) / 10
        return Math.max(0, totalPrice.value - totalPrice.value * rate)
    }
    return parseFloat(c.discount_amount || c.value || 0)
})

const finalPrice = computed(() => Math.max(0, totalPrice.value - couponDiscount.value))

function formatDiscount(c) {
    if (c.type === 2 || c.discount_rate) {
        return (c.discount_rate || c.value) + '折'
    }
    return '￥' + (c.discount_amount || c.value || 0)
}

function loadProduct() {
    const productId = route.query.product_id
    const quantity = parseInt(route.query.quantity) || 1
    if (!productId) return
    loadingProduct.value = true
    productDetail(productId).then((data) => {
        if (data && data.id) {
            productItem.value = {
                product_id: data.id,
                name: data.name,
                price: data.price,
                quantity: quantity
            }
        } else {
            ElMessage.error('商品不存在')
            router.push('/products')
        }
    }).catch(() => {
        ElMessage.error('商品不存在')
        router.push('/products')
    }).finally(() => {
        loadingProduct.value = false
    })
}

function loadCoupons() {
    couponsLoading.value = true
    userCoupons().then((d) => {
        const list = (d && d.list) ? d.list : (Array.isArray(d) ? d : [])
        availableCoupons.value = list.filter((c) => {
            return (c.status === 0 || c.status === undefined) && (totalPrice.value >= (c.min_amount || 0))
        })
    }).catch(() => {}).finally(() => {
        couponsLoading.value = false
    })
}

function submitOrder() {
    emailError.value = ''
    if (!form.email) {
        emailError.value = '请输入邮箱'
        ElMessage.warning('请输入邮箱')
        return
    }
    const emailRe = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRe.test(form.email)) {
        emailError.value = '邮箱格式不正确'
        ElMessage.warning('邮箱格式不正确')
        return
    }
    if (items.value.length === 0) {
        ElMessage.warning('没有可结算的商品')
        return
    }
    loading.value = true

    const orderItems = items.value.map((it) => ({
        product_id: it.product_id || it.productId || 0,
        sku_id: it.sku_id || it.skuId || 0,
        quantity: it.quantity
    }))

    const payload = {
        items: orderItems,
        email: form.email,
        remark: form.remark,
        pay_type: form.pay_type
    }
    if (selectedCouponId.value) {
        payload.coupon_id = selectedCouponId.value
    }
    if (orderItems.length === 1 && !productItem.value) {
        payload.product_id = orderItems[0].product_id
        payload.quantity = orderItems[0].quantity
    }

    createOrder(payload).then((data) => {
        if (!productItem.value) {
            cartStore.clearCart()
        }
        const orderNo = (data && data.order_no) ? data.order_no : ((data && data.order && data.order.order_no) ? data.order.order_no : null)
        if (!orderNo) {
            router.push('/user/orders')
            return
        }
        return payOrder(orderNo).then((payData) => {
            if (payData && payData.pay_url) {
                window.location.href = payData.pay_url
            } else {
                router.push('/order/' + orderNo)
            }
        }).catch(() => {
            router.push('/order/' + orderNo)
        })
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}

onMounted(() => {
    if (fromProduct.value) {
        loadProduct()
    }
    loadCoupons()
})
</script>
