<template>
    <el-container class="layout-container">
        <el-header class="layout-header">
            <div class="header-inner">
                <div class="header-left">
                    <router-link to="/" class="site-title">
                        {{ siteName }}
                    </router-link>
                </div>

                <el-menu
                    mode="horizontal"
                    :router="true"
                    class="header-menu"
                    :default-active="activeMenu"
                >
                    <el-menu-item index="/">首页</el-menu-item>
                    <el-menu-item index="/products">商品</el-menu-item>
                </el-menu>

                <div class="header-right">
                    <el-badge
                        :value="cartCount"
                        :hidden="cartCount === 0"
                        class="cart-badge"
                    >
                        <router-link to="/cart">
                            <el-button type="primary" size="default">
                                购物车
                            </el-button>
                        </router-link>
                    </el-badge>

                    <template v-if="token">
                        <el-dropdown trigger="click" @command="handleUserCommand">
                            <span class="user-trigger">
                                {{ userNickname || '用户' }}
                                <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                            </span>
                            <template #dropdown>
                                <el-dropdown-menu>
                                    <el-dropdown-item command="/user/profile">个人资料</el-dropdown-item>
                                    <el-dropdown-item command="/user/orders">我的订单</el-dropdown-item>
                                    <el-dropdown-item command="/user/distribution">推广中心</el-dropdown-item>
                                    <el-dropdown-item command="/user/signin" divided>
                                        <span style="color: #67c23a; font-weight: 500;">每日签到</span>
                                    </el-dropdown-item>
                                    <el-dropdown-item command="/user/coupons">我的优惠券</el-dropdown-item>
                                    <el-dropdown-item command="/user/points">积分明细</el-dropdown-item>
                                    <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                                </el-dropdown-menu>
                            </template>
                        </el-dropdown>
                    </template>
                    <template v-else>
                        <router-link to="/user/login">
                            <el-button type="default" size="default">登录</el-button>
                        </router-link>
                        <router-link to="/user/register">
                            <el-button type="default" size="default">注册</el-button>
                        </router-link>
                    </template>

                    <el-button
                        class="mobile-menu-btn"
                        text
                        @click="mobileMenuVisible = true"
                    >
                        <el-icon><Menu /></el-icon>
                    </el-button>
                </div>
            </div>

            <el-drawer
                v-model="mobileMenuVisible"
                direction="rtl"
                size="240px"
                class="mobile-drawer"
            >
                <template #header>
                    <span>{{ siteName }}</span>
                </template>
                <el-menu
                    :default-active="activeMenu"
                    @select="handleMobileSelect"
                >
                    <el-menu-item index="/">
                        <el-icon><House /></el-icon>
                        <span>首页</span>
                    </el-menu-item>
                    <el-menu-item index="/products">
                        <el-icon><Goods /></el-icon>
                        <span>商品</span>
                    </el-menu-item>
                    <el-menu-item index="/cart">
                        <el-icon><ShoppingCart /></el-icon>
                        <span>购物车</span>
                    </el-menu-item>
                    <template v-if="token">
                        <el-menu-item index="/user/profile">
                            <el-icon><User /></el-icon>
                            <span>个人资料</span>
                        </el-menu-item>
                        <el-menu-item index="/user/orders">
                            <el-icon><Document /></el-icon>
                            <span>我的订单</span>
                        </el-menu-item>
                        <el-menu-item index="/user/distribution">
                            <el-icon><Share /></el-icon>
                            <span>推广中心</span>
                        </el-menu-item>
                        <el-menu-item index="/user/signin">
                            <el-icon><Calendar /></el-icon>
                            <span>每日签到</span>
                        </el-menu-item>
                        <el-menu-item index="/user/coupons">
                            <el-icon><Ticket /></el-icon>
                            <span>我的优惠券</span>
                        </el-menu-item>
                        <el-menu-item index="/user/points">
                            <el-icon><Coin /></el-icon>
                            <span>积分明细</span>
                        </el-menu-item>
                        <el-menu-item index="logout" @click="logout">
                            <el-icon><SwitchButton /></el-icon>
                            <span>退出登录</span>
                        </el-menu-item>
                    </template>
                    <template v-else>
                        <el-menu-item index="/user/login">
                            <el-icon><User /></el-icon>
                            <span>登录</span>
                        </el-menu-item>
                        <el-menu-item index="/user/register">
                            <el-icon><EditPen /></el-icon>
                            <span>注册</span>
                        </el-menu-item>
                    </template>
                </el-menu>
            </el-drawer>
        </el-header>

        <el-main class="layout-main">
            <div class="page-container">
                <slot></slot>
            </div>
        </el-main>

        <el-footer class="layout-footer">
            <div class="footer-inner">
                版权所有 &copy; {{ currentYear }} {{ siteName }}
            </div>
        </el-footer>
    </el-container>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { useUserStore } from '../stores/user'
import {
    ArrowDown,
    Menu,
    House,
    Goods,
    ShoppingCart,
    User,
    Document,
    Share,
    Calendar,
    Ticket,
    Coin,
    SwitchButton,
    EditPen
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const mobileMenuVisible = ref(false)
const cartStore = useCartStore()
const userStore = useUserStore()
const siteName = ref('发卡平台')
const currentYear = new Date().getFullYear()

const token = computed(() => userStore.token)
const cartCount = computed(() => cartStore.totalCount)
const userNickname = computed(() => userStore.userInfo?.nickname || userStore.userInfo?.username || '')
const activeMenu = computed(() => route.path)

function logout() {
    mobileMenuVisible.value = false
    userStore.logout()
    router.push('/')
}

function handleUserCommand(command) {
    if (command === 'logout') {
        logout()
    } else {
        router.push(command)
    }
}

function handleMobileSelect(index) {
    if (index === 'logout') {
        logout()
    } else {
        mobileMenuVisible.value = false
    }
}
</script>

<style scoped>
.layout-container {
    min-height: 100vh;
}

.layout-header {
    height: 56px;
    padding: 0;
    background-color: #ffffff;
    border-bottom: 1px solid #e4e7ed;
    position: sticky;
    top: 0;
    z-index: 100;
}

.header-inner {
    max-width: 1280px;
    margin: 0 auto;
    padding: 0 16px;
    height: 100%;
    display: flex;
    align-items: center;
    gap: 24px;
}

.site-title {
    font-size: 18px;
    font-weight: 700;
    color: #303133;
    white-space: nowrap;
}

.site-title:hover {
    color: #409eff;
}

.header-menu {
    border-bottom: none;
    flex: 1;
}

.header-right {
    display: flex;
    align-items: center;
    gap: 12px;
}

.user-trigger {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    cursor: pointer;
    color: #303133;
    padding: 8px 12px;
    border-radius: 4px;
    transition: color 0.2s;
}

.user-trigger:hover {
    color: #409eff;
}

.mobile-menu-btn {
    display: none;
}

.layout-main {
    padding: 0;
    background-color: #f5f7fa;
    flex: 1;
}

.layout-footer {
    background-color: #ffffff;
    border-top: 1px solid #e4e7ed;
    padding: 0;
    height: auto;
    margin-top: 24px;
}

.footer-inner {
    max-width: 1280px;
    margin: 0 auto;
    padding: 24px 16px;
    text-align: center;
    font-size: 13px;
    color: #909399;
}

@media (max-width: 768px) {
    .header-inner {
        gap: 12px;
    }

    .header-menu {
        display: none;
    }

    .header-right .cart-badge,
    .header-right .user-trigger,
    .header-right > router-link {
        display: none;
    }

    .mobile-menu-btn {
        display: inline-flex;
    }
}
</style>
