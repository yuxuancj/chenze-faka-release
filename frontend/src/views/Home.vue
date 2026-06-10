<template>
    <Layout>
        <el-card
            shadow="never"
            class="hero-card"
            body-style="padding: 0"
        >
            <div class="hero-inner">
                <h1 class="hero-title">商品精选</h1>
                <p class="hero-desc">浏览各类优质商品，快速下单，即刻到账。</p>
                <div class="hero-actions">
                    <router-link to="/products">
                        <el-button type="primary" size="large" plain>
                            浏览全部商品
                        </el-button>
                    </router-link>
                    <router-link to="/user/login">
                        <el-button size="large" class="hero-btn-outline">
                            登录账户
                        </el-button>
                    </router-link>
                </div>
                <div class="hero-tag">
                    <span class="hero-tag-label">优质商品</span>
                    <span class="hero-tag-value">快速发货</span>
                </div>
            </div>
        </el-card>

        <el-card
            shadow="never"
            class="section-card"
        >
            <template #header>
                <div class="section-header">
                    <span class="section-title">热门商品</span>
                    <router-link to="/products" class="section-link">
                        查看全部
                    </router-link>
                </div>
            </template>

            <el-skeleton
                v-if="loading"
                :rows="3"
                animated
                class="loading-skeleton"
            />

            <el-empty
                v-else-if="products.length === 0"
                description="暂无商品"
                class="empty-state"
            />

            <el-row
                v-else
                :gutter="16"
            >
                <el-col
                    v-for="item in products"
                    :key="item.id"
                    :xs="24"
                    :sm="12"
                    :md="8"
                >
                    <router-link :to="'/product/' + item.id" class="product-link">
                        <el-card
                            shadow="hover"
                            class="product-card"
                            body-style="padding: 0"
                        >
                            <div class="product-thumb">{{ item.name }}</div>
                            <div class="product-info">
                                <div class="product-name">{{ item.name }}</div>
                                <div class="product-price">￥{{ item.price }}</div>
                                <div class="product-stock">库存: {{ item.stock || 0 }}</div>
                            </div>
                        </el-card>
                    </router-link>
                </el-col>
            </el-row>
        </el-card>
    </Layout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Layout from '../components/Layout.vue'
import { productList } from '../api/product'

const products = ref([])
const loading = ref(true)

function load() {
    loading.value = true
    productList(1, 6, '').then((data) => {
        products.value = (data && data.list) ? data.list : []
    }).catch(() => {
        products.value = []
    }).finally(() => {
        loading.value = false
    })
}

onMounted(load)
</script>

<style scoped>
.hero-card {
    margin-bottom: 24px;
    border: none;
    border-radius: 12px;
    overflow: hidden;
}

.hero-inner {
    background: linear-gradient(120deg, #409eff 0%, #2c7be5 100%);
    padding: 48px 32px;
    color: #ffffff;
    position: relative;
}

.hero-title {
    font-size: 28px;
    font-weight: 700;
    margin: 0 0 12px 0;
}

.hero-desc {
    font-size: 15px;
    opacity: 0.9;
    margin: 0 0 24px 0;
    max-width: 420px;
}

.hero-actions {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
}

.hero-btn-outline {
    background-color: transparent;
    color: #ffffff;
    border-color: #ffffff;
}

.hero-btn-outline:hover {
    background-color: rgba(255, 255, 255, 0.1);
    color: #ffffff;
    border-color: #ffffff;
}

.hero-tag {
    position: absolute;
    right: 32px;
    bottom: 24px;
    text-align: right;
}

.hero-tag-label {
    display: block;
    font-size: 12px;
    opacity: 0.8;
}

.hero-tag-value {
    display: block;
    font-size: 22px;
    font-weight: 700;
    margin-top: 4px;
}

.section-card {
    border: none;
    border-radius: 8px;
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.section-title {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
}

.section-link {
    font-size: 13px;
    color: #409eff;
}

.section-link:hover {
    color: #66b1ff;
    text-decoration: underline;
}

.loading-skeleton {
    padding: 16px;
}

.empty-state {
    padding: 48px 0;
}

.product-link {
    text-decoration: none;
    color: inherit;
    display: block;
    margin-bottom: 16px;
}

.product-card {
    border-radius: 8px;
    overflow: hidden;
}

.product-thumb {
    height: 160px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #f0f2f5 0%, #e4e7ed 100%);
    color: #909399;
    font-size: 14px;
    text-align: center;
    padding: 0 16px;
    overflow: hidden;
}

.product-info {
    padding: 12px 16px;
}

.product-name {
    font-size: 14px;
    font-weight: 500;
    color: #303133;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.product-price {
    font-size: 18px;
    font-weight: 700;
    color: #f56c6c;
    margin-top: 8px;
}

.product-stock {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
}

@media (max-width: 768px) {
    .hero-inner {
        padding: 32px 20px;
    }

    .hero-title {
        font-size: 22px;
    }

    .hero-tag {
        position: static;
        text-align: left;
        margin-top: 24px;
    }
}
</style>
