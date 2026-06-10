<template>
    <AdminLayout page-title="商品管理">
        <div class="mb-4 flex items-center justify-between flex-wrap gap-3">
            <div class="flex items-center gap-2">
                <el-input
                    v-model="keyword"
                    placeholder="搜索商品"
                    style="width: 240px"
                    clearable
                    @keyup.enter="search"
                >
                    <template #append>
                        <el-button @click="search">搜索</el-button>
                    </template>
                </el-input>
            </div>
            <router-link to="/admin/product/new">
                <el-button type="primary">新增商品</el-button>
            </router-link>
        </div>

        <el-card v-loading="loading" shadow="never">
            <el-table :data="products" style="width: 100%" empty-text="暂无数据" stripe>
                <el-table-column prop="id" label="ID" width="80" />
                <el-table-column prop="name" label="名称" min-width="200" />
                <el-table-column label="价格" width="120">
                    <template #default="scope">
                        <span class="text-blue-600 font-semibold">￥{{ scope.row.price }}</span>
                    </template>
                </el-table-column>
                <el-table-column prop="stock" label="库存" width="100" />
                <el-table-column prop="sales" label="已售" width="100" />
                <el-table-column label="状态" width="100">
                    <template #default="scope">
                        <el-tag v-if="scope.row.status === 1" type="success">上架</el-tag>
                        <el-tag v-else type="info">下架</el-tag>
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="160" fixed="right">
                    <template #default="scope">
                        <router-link :to="'/admin/product/' + scope.row.id">
                            <el-button size="small" type="primary">编辑</el-button>
                        </router-link>
                        <el-button size="small" type="danger" @click="del(scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>

            <div class="flex justify-center mt-4" v-if="pagination && pagination.total > 0">
                <el-pagination
                    v-model:current-page="page"
                    v-model:page-size="pagination.size"
                    :page-sizes="[10, 20, 50, 100]"
                    :total="pagination.total"
                    layout="prev, pager, next, total"
                    background
                    @current-change="load"
                    @size-change="load"
                />
            </div>
        </el-card>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import AdminLayout from '../components/AdminLayout.vue'
import { adminProductList, adminProductDelete } from '../api/admin'

const products = ref([])
const page = ref(1)
const pagination = reactive({ total: 0, size: 20 })
const keyword = ref('')
const loading = ref(false)

function load() {
    loading.value = true
    adminProductList(page.value, pagination.size, keyword.value).then((data) => {
        products.value = (data && data.list) ? data.list : []
        pagination.total = (data && data.total) ? data.total : 0
        pagination.size = (data && data.size) ? data.size : pagination.size
    }).catch(() => {
        products.value = []
        pagination.total = 0
    }).finally(() => {
        loading.value = false
    })
}

function search() {
    page.value = 1
    load()
}

function del(row) {
    ElMessageBox.confirm('确认删除该商品？', '提示', {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
        adminProductDelete(row.id).then(() => {
            ElMessage.success('删除成功')
            load()
        }).catch(() => {})
    }).catch(() => {})
}

onMounted(load)
</script>

<style scoped>
.text-blue-600 { color: #2563eb; }
.font-semibold { font-weight: 600; }
</style>
