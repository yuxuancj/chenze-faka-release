<template>
    <AdminLayout page-title="卡密管理">
        <el-card shadow="never" class="mb-4">
            <div class="flex items-center justify-between flex-wrap gap-3">
                <div class="flex items-center gap-2">
                    <el-select
                        v-model="filter_product_id"
                        placeholder="选择商品"
                        style="width: 240px"
                        clearable
                    >
                        <el-option label="全部商品" :value="0" />
                        <el-option v-for="p in products" :key="p.id" :label="p.name" :value="p.id" />
                    </el-select>
                    <el-button type="primary" @click="loadCards">加载卡密</el-button>
                </div>
                <div class="flex items-center gap-2">
                    <el-upload
                        ref="fileUploadRef"
                        action="#"
                        :auto-upload="false"
                        :show-file-list="false"
                        accept=".txt"
                        :on-change="onFileChange"
                    >
                        <el-button type="primary">批量导入卡密</el-button>
                    </el-upload>
                    <el-button @click="showTextModal = true">文本导入</el-button>
                </div>
            </div>
        </el-card>

        <el-card v-loading="loading" shadow="never">
            <el-table :data="cards" style="width: 100%" empty-text="请选择商品加载卡密" stripe>
                <el-table-column prop="id" label="ID" width="80" />
                <el-table-column label="卡密">
                    <template #default="scope">
                        <span class="font-mono text-sm">{{ scope.row.card_data }}</span>
                    </template>
                </el-table-column>
                <el-table-column label="状态" width="120">
                    <template #default="scope">
                        <el-tag v-if="scope.row.status === 1" type="success">已使用</el-tag>
                        <el-tag v-else type="info">未使用</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="created_at" label="创建时间" width="180" />
            </el-table>

            <div class="flex justify-center mt-4" v-if="pagination && pagination.total > 0">
                <el-pagination
                    v-model:current-page="page"
                    v-model:page-size="pagination.size"
                    :page-sizes="[10, 20, 50, 100]"
                    :total="pagination.total"
                    layout="prev, pager, next, total"
                    background
                    @current-change="loadCards"
                    @size-change="loadCards"
                />
            </div>
        </el-card>

        <el-dialog
            v-model="showTextModal"
            title="文本导入卡密（每行一条）"
            width="600px"
        >
            <el-input
                type="textarea"
                v-model="importText"
                :rows="10"
                placeholder="每行一个卡密"
                style="margin-bottom: 12px"
            />
            <div class="flex items-center justify-between">
                <span class="text-sm text-gray-500">选择商品ID: {{ filter_product_id }}</span>
                <el-button type="primary" :loading="importing" @click="importTextCards">
                    导入
                </el-button>
            </div>
        </el-dialog>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import AdminLayout from '../components/AdminLayout.vue'
import { adminCardList, adminCardImport, adminProductList } from '../api/admin'

const cards = ref([])
const products = ref([])
const filter_product_id = ref(0)
const page = ref(1)
const pagination = reactive({ total: 0, size: 20 })
const loading = ref(false)
const showTextModal = ref(false)
const importText = ref('')
const importing = ref(false)

function loadProducts() {
    adminProductList(1, 100, '').then((data) => {
        products.value = (data && data.list) ? data.list : []
    }).catch(() => {})
}

function loadCards() {
    loading.value = true
    adminCardList(page.value, pagination.size, filter_product_id.value).then((data) => {
        cards.value = (data && data.list) ? data.list : []
        pagination.total = (data && data.total) ? data.total : 0
        pagination.size = (data && data.size) ? data.size : pagination.size
    }).catch(() => {
        cards.value = []
        pagination.total = 0
    }).finally(() => {
        loading.value = false
    })
}

function onFileChange(file) {
    if (!file || !file.raw) return
    if (!filter_product_id.value) {
        ElMessage.warning('请先选择商品')
        return
    }
    adminCardImport(filter_product_id.value, file.raw).then(() => {
        ElMessage.success('导入成功')
        loadCards()
    }).catch(() => {})
}

function importTextCards() {
    if (!importText.value.trim()) {
        ElMessage.warning('请输入卡密内容')
        return
    }
    if (!filter_product_id.value) {
        ElMessage.warning('请先选择商品')
        return
    }
    importing.value = true
    const blob = new Blob([importText.value], { type: 'text/plain' })
    const file = new File([blob], 'cards.txt', { type: 'text/plain' })
    adminCardImport(filter_product_id.value, file).then(() => {
        ElMessage.success('导入成功')
        showTextModal.value = false
        importText.value = ''
        loadCards()
    }).catch(() => {}).finally(() => {
        importing.value = false
    })
}

onMounted(() => {
    loadProducts()
    loadCards()
})
</script>
