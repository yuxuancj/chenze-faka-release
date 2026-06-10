<template>
    <AdminLayout :page-title="'仪表盘'">
        <div v-if="loading" class="text-center py-12 text-gray-500">
            <div class="animate-spin h-12 w-12 mx-auto"></div>
            <p class="mt-4">加载中...</p>
        </div>
        <div v-else class="space-y-6">
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div class="card p-4">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-gray-500 text-sm">今日订单</p>
                            <p class="text-3xl font-bold text-blue-600">{{ stats.today_orders || 0 }}</p>
                        </div>
                        <div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                            <span class="text-2xl">📋</span>
                        </div>
                    </div>
                    <p class="text-sm mt-2" :class="stats.yesterday_orders > 0 ? 'text-green-500' : 'text-gray-400'">
                        较昨日 {{ stats.yesterday_orders > 0 ? '+' + ((stats.today_orders / stats.yesterday_orders * 100 - 100).toFixed(1)) + '%' : '--' }}
                    </p>
                </div>

                <div class="card p-4">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-gray-500 text-sm">今日销售额</p>
                            <p class="text-3xl font-bold text-green-600">¥{{ (stats.today_sales || 0).toFixed(2) }}</p>
                        </div>
                        <div class="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
                            <span class="text-2xl">💰</span>
                        </div>
                    </div>
                    <p class="text-sm mt-2" :class="stats.yesterday_sales > 0 ? 'text-green-500' : 'text-gray-400'">
                        较昨日 {{ stats.yesterday_sales > 0 ? '+' + ((stats.today_sales / stats.yesterday_sales * 100 - 100).toFixed(1)) + '%' : '--' }}
                    </p>
                </div>

                <div class="card p-4">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-gray-500 text-sm">今日新增用户</p>
                            <p class="text-3xl font-bold text-purple-600">{{ stats.today_new_users || 0 }}</p>
                        </div>
                        <div class="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center">
                            <span class="text-2xl">👤</span>
                        </div>
                    </div>
                </div>

                <div class="card p-4">
                    <div class="flex items-center justify-between">
                        <div>
                            <p class="text-gray-500 text-sm">本周订单</p>
                            <p class="text-3xl font-bold text-orange-600">{{ stats.week_orders || 0 }}</p>
                        </div>
                        <div class="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center">
                            <span class="text-2xl">📊</span>
                        </div>
                    </div>
                </div>
            </div>

            <div class="card p-4">
                <p class="font-medium mb-4">累计数据</p>
                <div class="grid grid-cols-4 gap-4">
                    <div class="text-center">
                        <p class="text-2xl font-bold text-gray-800">{{ stats.total_users || 0 }}</p>
                        <p class="text-sm text-gray-500">累计用户</p>
                    </div>
                    <div class="text-center">
                        <p class="text-2xl font-bold text-gray-800">{{ stats.total_products || 0 }}</p>
                        <p class="text-sm text-gray-500">累计商品</p>
                    </div>
                    <div class="text-center">
                        <p class="text-2xl font-bold text-gray-800">{{ stats.total_orders || 0 }}</p>
                        <p class="text-sm text-gray-500">累计订单</p>
                    </div>
                    <div class="text-center">
                        <p class="text-2xl font-bold text-gray-800">¥{{ (stats.total_sales || 0).toFixed(2) }}</p>
                        <p class="text-sm text-gray-500">累计销售额</p>
                    </div>
                </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="card p-4">
                    <p class="font-medium mb-4">近7天订单趋势</p>
                    <div ref="orderChartRef" class="h-64"></div>
                </div>

                <div class="card p-4">
                    <p class="font-medium mb-4">近7天销售额趋势</p>
                    <div ref="salesChartRef" class="h-64"></div>
                </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="card p-4">
                    <p class="font-medium mb-4">商品销售排行 Top10</p>
                    <table class="w-full">
                        <thead>
                            <tr class="border-b">
                                <th class="text-left py-2">#</th>
                                <th class="text-left py-2">商品名称</th>
                                <th class="text-right py-2">销量</th>
                                <th class="text-right py-2">销售额</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(item, index) in stats.product_rank" :key="index" class="border-b last:border-0">
                                <td class="py-2">{{ index + 1 }}</td>
                                <td class="py-2 truncate">{{ item.name }}</td>
                                <td class="text-right py-2">{{ item.sales }}</td>
                                <td class="text-right py-2">¥{{ item.amount.toFixed(2) }}</td>
                            </tr>
                            <tr v-if="!stats.product_rank || stats.product_rank.length === 0">
                                <td colspan="4" class="text-center py-8 text-gray-400">暂无数据</td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <div class="card p-4">
                    <p class="font-medium mb-4">支付渠道占比</p>
                    <div ref="paymentChartRef" class="h-64"></div>
                </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div class="card p-4">
                    <p class="font-medium mb-4">系统信息</p>
                    <div class="space-y-2 text-sm">
                        <div class="flex justify-between">
                            <span class="text-gray-500">版本号</span>
                            <span class="font-medium">{{ stats.version || '-' }}</span>
                        </div>
                        <div class="flex justify-between">
                            <span class="text-gray-500">Go 版本</span>
                            <span class="font-medium">{{ stats.go_version || '-' }}</span>
                        </div>
                        <div class="flex justify-between">
                            <span class="text-gray-500">数据库</span>
                            <span class="font-medium">{{ stats.db_version || '-' }}</span>
                        </div>
                        <div class="flex justify-between">
                            <span class="text-gray-500">服务器时间</span>
                            <span class="font-medium">{{ stats.server_time || '-' }}</span>
                        </div>
                    </div>
                </div>

                <div class="card p-4">
                    <p class="font-medium mb-4">待处理事项</p>
                    <div class="space-y-3">
                        <div class="flex items-center justify-between p-3 bg-blue-50 rounded-lg">
                            <span class="text-gray-700">待审核工单</span>
                            <span class="bg-blue-500 text-white px-3 py-1 rounded-full text-sm">{{ stats.pending_orders || 0 }}</span>
                        </div>
                        <div class="flex items-center justify-between p-3 bg-green-50 rounded-lg">
                            <span class="text-gray-700">待处理提现</span>
                            <span class="bg-green-500 text-white px-3 py-1 rounded-full text-sm">{{ stats.pending_withdraws || 0 }}</span>
                        </div>
                        <div class="flex items-center justify-between p-3 bg-orange-50 rounded-lg">
                            <span class="text-gray-700">低库存商品</span>
                            <span class="bg-orange-500 text-white px-3 py-1 rounded-full text-sm">{{ stats.low_stock_count || 0 }}</span>
                        </div>
                    </div>
                </div>

                <div class="card p-4">
                    <p class="font-medium mb-4">快捷操作</p>
                    <div class="grid grid-cols-2 gap-3">
                        <button class="btn-secondary py-3">商品管理</button>
                        <button class="btn-secondary py-3">订单管理</button>
                        <button class="btn-secondary py-3">用户管理</button>
                        <button class="btn-secondary py-3">系统设置</button>
                    </div>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>import { ref, onMounted, watch, nextTick } from 'vue';
import AdminLayout from '../components/AdminLayout.vue';
import { request } from '../api/admin';
const loading = ref(true);
const stats = ref({});
const orderChartRef = ref(null);
const salesChartRef = ref(null);
const paymentChartRef = ref(null);
function loadData() {
 loading.value = true;
 request.get('/admin/api/dashboard').then(data => {
 stats.value = data;
 nextTick(() => {
 initCharts();
 });
 }).catch(() => {
 stats.value = {
 today_orders: 0,
 today_sales: 0,
 today_new_users: 0,
 yesterday_orders: 0,
 yesterday_sales: 0,
 week_orders: 0,
 week_sales: 0,
 total_users: 0,
 total_products: 0,
 total_orders: 0,
 total_sales: 0,
 last_7_days_orders: [0, 0, 0, 0, 0, 0, 0],
 last_7_days_sales: [0, 0, 0, 0, 0, 0, 0],
 product_rank: [],
 payment_ratio: [],
 version: 'v2.2.0',
 go_version: '',
 db_version: 'MySQL',
 server_time: new Date().toLocaleString(),
 pending_orders: 0,
 pending_withdraws: 0,
 low_stock_count: 0
 };
 nextTick(() => {
 initCharts();
 });
 }).finally(() => {
 loading.value = false;
 });
}
function initCharts() {
 const days = ['周一', '周二', '周三', '周四', '周五', '周六', '周日'];
 if (orderChartRef.value) {
 const canvas = document.createElement('canvas');
 orderChartRef.value.innerHTML = '';
 orderChartRef.value.appendChild(canvas);
 const ctx = canvas.getContext('2d');
 const orders = stats.value.last_7_days_orders || [0, 0, 0, 0, 0, 0, 0];
 const maxOrder = Math.max(...orders, 1);
 const padding = 40;
 const width = orderChartRef.value.clientWidth - padding * 2;
 const height = 200;
 ctx.clearRect(0, 0, width + padding * 2, height + padding);
 ctx.strokeStyle = '#e5e7eb';
 for (let i = 0; i <= 4; i++) {
 const y = padding + (height / 4) * i;
 ctx.beginPath();
 ctx.moveTo(padding, y);
 ctx.lineTo(width + padding, y);
 ctx.stroke();
 ctx.fillStyle = '#9ca3af';
 ctx.font = '10px sans-serif';
 ctx.textAlign = 'right';
 ctx.fillText(Math.round(maxOrder - (maxOrder / 4) * i), padding - 5, y + 3);
 }
 ctx.strokeStyle = '#3b82f6';
 ctx.lineWidth = 2;
 ctx.beginPath();
 orders.forEach((order, index) => {
 const x = padding + (width / 6) * index;
 const y = padding + height - (order / maxOrder) * height;
 if (index === 0) {
 ctx.moveTo(x, y);
 } else {
 ctx.lineTo(x, y);
 }
 });
 ctx.stroke();
 orders.forEach((order, index) => {
 const x = padding + (width / 6) * index;
 const y = padding + height - (order / maxOrder) * height;
 ctx.beginPath();
 ctx.arc(x, y, 4, 0, Math.PI * 2);
 ctx.fillStyle = '#3b82f6';
 ctx.fill();
 });
 ctx.fillStyle = '#6b7280';
 ctx.font = '12px sans-serif';
 ctx.textAlign = 'center';
 days.forEach((day, index) => {
 const x = padding + (width / 6) * index;
 ctx.fillText(day, x, height + padding - 5);
 });
 }
 if (salesChartRef.value) {
 const canvas = document.createElement('canvas');
 salesChartRef.value.innerHTML = '';
 salesChartRef.value.appendChild(canvas);
 const ctx = canvas.getContext('2d');
 const sales = stats.value.last_7_days_sales || [0, 0, 0, 0, 0, 0, 0];
 const maxSale = Math.max(...sales, 1);
 const padding = 40;
 const width = salesChartRef.value.clientWidth - padding * 2;
 const height = 200;
 ctx.clearRect(0, 0, width + padding * 2, height + padding);
 ctx.strokeStyle = '#e5e7eb';
 for (let i = 0; i <= 4; i++) {
 const y = padding + (height / 4) * i;
 ctx.beginPath();
 ctx.moveTo(padding, y);
 ctx.lineTo(width + padding, y);
 ctx.stroke();
 ctx.fillStyle = '#9ca3af';
 ctx.font = '10px sans-serif';
 ctx.textAlign = 'right';
 ctx.fillText('¥' + Math.round(maxSale - (maxSale / 4) * i), padding - 5, y + 3);
 }
 const colors = ['#22c55e', '#10b981', '#14b8a6', '#06b6d4', '#0ea5e9', '#3b82f6', '#6366f1'];
 sales.forEach((sale, index) => {
 const x = padding + (width / 6) * index - 15;
 const barWidth = 30;
 const heightRatio = sale / maxSale;
 const barHeight = heightRatio * height;
 ctx.fillStyle = colors[index % colors.length];
 ctx.fillRect(x, padding + height - barHeight, barWidth, barHeight);
 ctx.fillStyle = '#6b7280';
 ctx.font = '12px sans-serif';
 ctx.textAlign = 'center';
 ctx.fillText(days[index], x + barWidth / 2, height + padding - 5);
 });
 }
 if (paymentChartRef.value) {
 const canvas = document.createElement('canvas');
 paymentChartRef.value.innerHTML = '';
 paymentChartRef.value.appendChild(canvas);
 const ctx = canvas.getContext('2d');
 const payments = stats.value.payment_ratio || [];
 const centerX = 100;
 const centerY = 100;
 const radius = 80;
 let startAngle = -Math.PI / 2;
 const colors = ['#3b82f6', '#22c55e', '#f59e0b', '#ef4444', '#8b5cf6'];
 payments.forEach((payment, index) => {
 const sliceAngle = (payment.percent / 100) * Math.PI * 2;
 ctx.beginPath();
 ctx.moveTo(centerX, centerY);
 ctx.arc(centerX, centerY, radius, startAngle, startAngle + sliceAngle);
 ctx.closePath();
 ctx.fillStyle = colors[index % colors.length];
 ctx.fill();
 startAngle += sliceAngle;
 });
 ctx.fillStyle = '#ffffff';
 ctx.beginPath();
 ctx.arc(centerX, centerY, 50, 0, Math.PI * 2);
 ctx.fill();
 ctx.fillStyle = '#6b7280';
 ctx.font = '14px sans-serif';
 ctx.textAlign = 'center';
 ctx.fillText('支付', centerX, centerY - 5);
 ctx.fillText('占比', centerX, centerY + 15);
 ctx.font = '12px sans-serif';
 let legendY = 20;
 payments.forEach((payment, index) => {
 ctx.fillStyle = colors[index % colors.length];
 ctx.fillRect(160, legendY, 12, 12);
 ctx.fillStyle = '#374151';
 ctx.textAlign = 'left';
 ctx.fillText(payment.method + ': ' + payment.percent.toFixed(1) + '%', 180, legendY + 10);
 legendY += 25;
 });
 }
}
onMounted(() => {
 loadData();
});
</script>