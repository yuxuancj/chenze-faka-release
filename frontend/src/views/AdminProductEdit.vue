<template>
    <AdminLayout :page-title="isEdit ? '编辑商品' : '新增商品'">
        <div class="card max-w-2xl">
            <div class="card-body space-y-4">
                <div>
                    <label class="form-label">商品名称</label>
                    <input v-model="form.name" type="text" class="form-input" placeholder="请输入商品名称">
                </div>
                <div>
                    <label class="form-label">分类</label>
                    <select v-model="form.category_id" class="form-input">
                        <option value="0">请选择分类</option>
                        <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                    </select>
                </div>
                <div class="grid grid-cols-2 gap-4">
                    <div>
                        <label class="form-label">价格</label>
                        <input v-model.number="form.price" type="number" step="0.01" class="form-input" placeholder="0.00">
                    </div>
                    <div>
                        <label class="form-label">库存</label>
                        <input v-model.number="form.stock" type="number" class="form-input" placeholder="0">
                    </div>
                </div>
                <div>
                    <label class="form-label">图片</label>
                    <input v-model="form.image" type="text" class="form-input" placeholder="图片URL">
                </div>
                <div>
                    <label class="form-label">商品描述</label>
                    <textarea v-model="form.description" class="form-input" rows="5" placeholder="请输入商品描述"></textarea>
                </div>
                <div>
                    <label class="form-label">状态</label>
                    <select v-model="form.status" class="form-input">
                        <option :value="1">上架</option>
                        <option :value="0">下架</option>
                    </select>
                </div>
                <div class="flex items-center space-x-2">
                    <button @click="save" :disabled="loading" class="btn-primary">
                        {{ loading ? '保存中...' : '保存' }}
                    </button>
                    <router-link to="/admin/products" class="btn-secondary">返回</router-link>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AdminLayout from '../components/AdminLayout.vue'
import { adminProductCreate, adminProductUpdate, adminCategoryList } from '../api/admin'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)
const loading = ref(false)
const categories = ref([])
const form = reactive({
    name: '',
    category_id: 0,
    price: 0,
    stock: 0,
    image: '',
    description: '',
    status: 1
})

function loadCategories() {
    adminCategoryList(1, 100).then((res) => {
        categories.value = res.data && res.data.list ? res.data.list : []
    }).catch(() => {})
}

function save() {
    if (!form.name) {
        alert('请输入商品名称')
        return
    }
    loading.value = true
    const action = isEdit.value
        ? adminProductUpdate(route.params.id, form)
        : adminProductCreate(form)
    action.then(() => {
        alert('保存成功')
        router.push('/admin/products')
    }).catch(() => {}).finally(() => {
        loading.value = false
    })
}

onMounted(loadCategories)
</script>
