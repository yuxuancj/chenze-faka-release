<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6" v-for="(stat, index) in stats" :key="index">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon" :style="{ background: stat.bg }">
              <el-icon class="icon"><component :is="stat.icon" /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-label">{{ stat.label }}</div>
              <div class="stat-value">{{ stat.value }}</div>
              <div class="stat-change" :class="stat.trend === 'up' ? 'positive' : 'negative'">
                <el-icon><TrendCharts v-if="stat.trend === 'up'" /><TrendCharts v-else style="transform: rotate(180deg)" /></el-icon>
                {{ stat.change }}
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="16">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="card-title">近期订单</span>
              <el-button type="primary" link @click="goToOrders">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentOrders" style="width: 100%" stripe>
            <el-table-column prop="orderNo" label="订单号" width="180" />
            <el-table-column prop="product" label="商品" />
            <el-table-column prop="buyer" label="买家" width="120" />
            <el-table-column prop="amount" label="金额" width="100">
              <template #default="scope">
                <span class="amount">¥{{ scope.row.amount }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="statusTagType(scope.row.status)">{{ scope.row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="time" label="时间" width="170" />
          </el-table>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="card-title">热门商品</span>
            </div>
          </template>
          <el-table :data="topProducts" style="width: 100%">
            <el-table-column prop="rank" label="排名" width="60">
              <template #default="scope">
                <el-tag
                  :type="scope.row.rank === 1 ? 'danger' : scope.row.rank === 2 ? 'warning' : scope.row.rank === 3 ? 'success' : 'info'"
                  size="small"
                >
                  {{ scope.row.rank }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="name" label="商品" />
            <el-table-column prop="sales" label="销量" width="80" align="right" />
            <el-table-column prop="revenue" label="收入" width="100" align="right">
              <template #default="scope">¥{{ scope.row.revenue }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="card-title">系统信息</span>
            </div>
          </template>
          <el-descriptions :column="3" border>
            <el-descriptions-item label="系统名称">辰泽发卡系统</el-descriptions-item>
            <el-descriptions-item label="系统版本">v2.3.0</el-descriptions-item>
            <el-descriptions-item label="运行环境">Production</el-descriptions-item>
            <el-descriptions-item label="服务器时间">{{ serverTime }}</el-descriptions-item>
            <el-descriptions-item label="PHP版本">8.1.0</el-descriptions-item>
            <el-descriptions-item label="数据库">MySQL 8.0</el-descriptions-item>
            <el-descriptions-item label="注册用户">1,284</el-descriptions-item>
            <el-descriptions-item label="总订单数">5,672</el-descriptions-item>
            <el-descriptions-item label="总收入">¥89,420</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  ShoppingCart,
  Money,
  User,
  Goods,
  TrendCharts,
} from '@element-plus/icons-vue'

const router = useRouter()

const stats = ref([
  { label: '今日订单', value: '128', change: '+12.5%', trend: 'up', icon: ShoppingCart, bg: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
  { label: '今日收入', value: '¥3,840', change: '+8.3%', trend: 'up', icon: Money, bg: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' },
  { label: '新增用户', value: '24', change: '+5.2%', trend: 'up', icon: User, bg: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
  { label: '商品总数', value: '156', change: '-2.1%', trend: 'down', icon: Goods, bg: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' },
])

const recentOrders = ref([
  { orderNo: 'CZ202412100001', product: 'Netflix 30天会员', buyer: 'user123', amount: '29.90', status: '已完成', time: '2024-12-10 14:23' },
  { orderNo: 'CZ202412100002', product: 'Spotify Premium', buyer: 'musicfan', amount: '15.00', status: '已完成', time: '2024-12-10 13:45' },
  { orderNo: 'CZ202412100003', product: 'YouTube Premium', buyer: 'ytviewer', amount: '18.50', status: '处理中', time: '2024-12-10 12:30' },
  { orderNo: 'CZ202412100004', product: 'Netflix 30天会员', buyer: 'movielover', amount: '29.90', status: '已完成', time: '2024-12-10 11:15' },
  { orderNo: 'CZ202412100005', product: 'Steam 充值卡', buyer: 'gamer88', amount: '100.00', status: '已取消', time: '2024-12-10 10:00' },
])

const topProducts = ref([
  { rank: 1, name: 'Netflix 30天会员', sales: 234, revenue: '6,996' },
  { rank: 2, name: 'Spotify Premium', sales: 189, revenue: '2,835' },
  { rank: 3, name: 'YouTube Premium', sales: 156, revenue: '2,886' },
  { rank: 4, name: 'Steam 50元充值', sales: 98, revenue: '4,900' },
  { rank: 5, name: 'Apple ID 充值', sales: 87, revenue: '4,350' },
])

const serverTime = ref('')
let timer = null

const updateTime = () => {
  const now = new Date()
  serverTime.value = now.toLocaleString('zh-CN', { hour12: false })
}

onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

const statusTagType = (status) => {
  const map = {
    已完成: 'success',
    处理中: 'warning',
    已取消: 'info',
    待付款: 'danger',
  }
  return map[status] || ''
}

const goToOrders = () => {
  router.push('/admin/orders')
}
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.stat-card {
  border-radius: 8px;
  border: none;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon .icon {
  font-size: 28px;
  color: #fff;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 4px;
}

.stat-change {
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 2px;
}

.stat-change.positive {
  color: #67c23a;
}

.stat-change.negative {
  color: #f56c6c;
}

.mt-20 {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.amount {
  color: #f56c6c;
  font-weight: 600;
}
</style>
