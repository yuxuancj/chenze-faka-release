<template>
    <AdminLayout :page-title="isEdit ? '编辑商品' : '新增商品'">
        <el-card v-loading="loading" shadow="never">
            <el-form :model="form" label-width="100px" class="space-y-4">
                <el-collapse v-model="activeSections">
                    <el-collapse-item title="基本信息" name="basic">
                        <el-form-item label="商品名称">
                            <el-input v-model="form.name" placeholder="请输入商品名称" maxlength="100" show-word-limit />
                        </el-form-item>
                        <el-form-item label="商品分类">
                            <el-select v-model="form.category_id" placeholder="请选择分类" style="width: 300px">
                                <el-option label="请选择分类" :value="0" />
                                <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
                            </el-select>
                        </el-form-item>
                        <el-row :gutter="20">
                            <el-col :span="8">
                                <el-form-item label="价格">
                                    <el-input-number v-model="form.price" :min="0" :precision="2" :step="1" style="width: 100%" />
                                </el-form-item>
                            </el-col>
                            <el-col :span="8">
                                <el-form-item label="库存">
                                    <el-input-number v-model="form.stock" :min="0" :step="1" style="width: 100%" />
                                </el-form-item>
                            </el-col>
                            <el-col :span="8">
                                <el-form-item label="排序">
                                    <el-input-number v-model="form.sort" :step="1" style="width: 100%" />
                                </el-form-item>
                            </el-col>
                        </el-row>
                        <el-form-item label="商品图片">
                            <el-input v-model="form.image" placeholder="图片URL" />
                        </el-form-item>
                        <el-form-item label="商品描述">
                            <el-input type="textarea" v-model="form.description" :rows="4" placeholder="请输入商品描述" />
                        </el-form-item>
                        <el-row :gutter="20">
                            <el-col :span="12">
                                <el-form-item label="状态">
                                    <el-radio-group v-model="form.status">
                                        <el-radio :label="1">上架</el-radio>
                                        <el-radio :label="0">下架</el-radio>
                                    </el-radio-group>
                                </el-form-item>
                            </el-col>
                            <el-col :span="12">
                                <el-form-item label="隐藏商品">
                                    <el-switch v-model="form.is_hidden" />
                                </el-form-item>
                            </el-col>
                        </el-row>
                    </el-collapse-item>

                    <el-collapse-item title="发货类型" name="delivery">
                        <el-form-item label="发货方式">
                            <el-tabs v-model="form.delivery_type" type="card" style="margin-bottom: 0">
                                <el-tab-pane label="卡密" name="card" />
                                <el-tab-pane label="文本" name="text" />
                                <el-tab-pane label="API授权" name="api" />
                                <el-tab-pane label="无限卡密" name="infinite" />
                                <el-tab-pane label="人工核销" name="manual" />
                            </el-tabs>
                        </el-form-item>

                        <div v-if="form.delivery_type === 'card'">
                            <el-form-item label="卡密批量导入">
                                <el-upload
                                    action="#"
                                    :auto-upload="false"
                                    :show-file-list="false"
                                    accept=".txt"
                                    :on-change="handleCardFileUpload"
                                >
                                    <el-button type="primary">
                                        <el-icon><Upload /></el-icon> 点击上传 .txt 文件
                                    </el-button>
                                    <span class="text-sm text-gray-500 ml-3">每行一个卡密，或使用 | 分隔卡号和密码</span>
                                </el-upload>
                            </el-form-item>
                            <div v-if="cardList.length > 0" class="mb-4">
                                <div class="flex items-center justify-between mb-2">
                                    <span class="text-sm text-gray-600">已导入 {{ cardList.length }} 条卡密</span>
                                    <el-button size="small" type="danger" @click="cardList = []">清空</el-button>
                                </div>
                                <el-input
                                    type="textarea"
                                    :rows="6"
                                    :model-value="cardList.slice(0, 20).join('\n')"
                                    readonly
                                    placeholder="已导入的卡密"
                                    v-if="cardList.length > 0"
                                />
                                <div v-if="cardList.length > 20" class="text-center text-gray-400 py-2 text-sm">
                                    ... 还有 {{ cardList.length - 20 }} 条
                                </div>
                            </div>
                        </div>

                        <div v-if="form.delivery_type === 'text'">
                            <el-form-item label="文本内容">
                                <el-input type="textarea" v-model="deliveryConfig.text_content" :rows="6" placeholder="购买后显示的固定文本内容" />
                            </el-form-item>
                        </div>

                        <div v-if="form.delivery_type === 'api'">
                            <el-form-item label="API URL">
                                <el-input v-model="deliveryConfig.api_url" placeholder="https://api.example.com/auth" />
                            </el-form-item>
                            <el-form-item label="请求方式">
                                <el-select v-model="deliveryConfig.api_method" style="width: 200px">
                                    <el-option label="GET" value="GET" />
                                    <el-option label="POST" value="POST" />
                                </el-select>
                            </el-form-item>
                            <el-form-item label="Header 参数">
                                <div v-for="(value, key) in deliveryConfig.api_headers" :key="key" class="flex gap-2 mb-2 items-center">
                                    <el-input :model-value="key" @update:model-value="updateApiHeaderKey($event, key)" placeholder="Key" style="flex: 1" />
                                    <el-input v-model="deliveryConfig.api_headers[key]" placeholder="Value" style="flex: 1" />
                                    <el-button type="danger" @click="deleteApiHeader(key)">删除</el-button>
                                </div>
                                <el-button @click="addApiHeader">添加 Header</el-button>
                            </el-form-item>
                            <el-form-item label="Body 参数">
                                <div v-for="(value, key) in deliveryConfig.api_body" :key="key" class="flex gap-2 mb-2 items-center">
                                    <el-input :model-value="key" @update:model-value="updateApiBodyKey($event, key)" placeholder="Key" style="flex: 1" />
                                    <el-input v-model="deliveryConfig.api_body[key]" placeholder="Value" style="flex: 1" />
                                    <el-button type="danger" @click="deleteApiBody(key)">删除</el-button>
                                </div>
                                <el-button @click="addApiBody">添加 Body 参数</el-button>
                            </el-form-item>
                            <el-row :gutter="20">
                                <el-col :span="12">
                                    <el-form-item label="成功标识字段">
                                        <el-input v-model="deliveryConfig.api_success_field" placeholder="code==0" />
                                    </el-form-item>
                                </el-col>
                                <el-col :span="12">
                                    <el-form-item label="授权码返回字段">
                                        <el-input v-model="deliveryConfig.api_auth_field" placeholder="data.auth_code" />
                                    </el-form-item>
                                </el-col>
                            </el-row>
                            <el-form-item label="到期天数">
                                <el-input-number v-model="deliveryConfig.api_expire_days" :min="1" :step="1" />
                            </el-form-item>
                        </div>

                        <div v-if="form.delivery_type === 'infinite'">
                            <el-form-item label="固定内容">
                                <el-input type="textarea" v-model="deliveryConfig.text_content" :rows="4" placeholder="用户购买后显示的固定内容" />
                            </el-form-item>
                        </div>

                        <div v-if="form.delivery_type === 'manual'">
                            <el-alert
                                title="人工核销模式：用户购买后生成工单，管理员需手动完成核销。"
                                type="info"
                                :closable="false"
                                show-icon
                            />
                        </div>
                    </el-collapse-item>

                    <el-collapse-item title="SKU 多规格管理" name="sku">
                        <el-form-item label="启用SKU">
                            <el-switch v-model="form.has_sku" />
                        </el-form-item>
                        <template v-if="form.has_sku">
                            <div class="mb-4">
                                <el-button type="primary" @click="addSku">添加规格</el-button>
                            </div>
                            <el-table :data="form.skus" border>
                                <el-table-column label="规格组合">
                                    <template #default="scope">
                                        <el-input v-model="scope.row.spec_names" placeholder="如：颜色:红色,尺寸:L" />
                                    </template>
                                </el-table-column>
                                <el-table-column label="SKU编码" width="180">
                                    <template #default="scope">
                                        <el-input v-model="scope.row.sku_code" placeholder="SKU编码" />
                                    </template>
                                </el-table-column>
                                <el-table-column label="价格" width="160">
                                    <template #default="scope">
                                        <el-input-number v-model="scope.row.price" :min="0" :precision="2" style="width: 100%" />
                                    </template>
                                </el-table-column>
                                <el-table-column label="库存" width="140">
                                    <template #default="scope">
                                        <el-input-number v-model="scope.row.stock" :min="0" style="width: 100%" />
                                    </template>
                                </el-table-column>
                                <el-table-column label="图片">
                                    <template #default="scope">
                                        <el-input v-model="scope.row.image" placeholder="图片URL" />
                                    </template>
                                </el-table-column>
                                <el-table-column label="操作" width="100" fixed="right">
                                    <template #default="scope">
                                        <el-button type="danger" size="small" @click="removeSku(scope.$index)">删除</el-button>
                                    </template>
                                </el-table-column>
                            </el-table>
                            <el-alert
                                v-if="form.skus.length === 0"
                                title="暂无 SKU，点击上方按钮添加"
                                type="info"
                                :closable="false"
                                class="mt-4"
                            />
                        </template>
                    </el-collapse-item>

                    <el-collapse-item title="会员等级价" name="member">
                        <el-form-item
                            v-for="level in memberLevels"
                            :key="level.id"
                            :label="level.name"
                        >
                            <el-input-number
                                :model-value="getMemberPrice(level.id)"
                                @update:model-value="(val) => updateMemberPrice(level.id, val)"
                                :min="0"
                                :precision="2"
                                :placeholder="'默认 ' + form.price"
                            />
                        </el-form-item>
                        <el-alert
                            title="留空则使用默认价格"
                            type="info"
                            :closable="false"
                            show-icon
                            class="mt-4"
                        />
                    </el-collapse-item>

                    <el-collapse-item title="SEO 设置" name="seo">
                        <el-form-item label="SEO 标题">
                            <el-input v-model="form.seo.title" placeholder="页面标题" />
                        </el-form-item>
                        <el-form-item label="SEO 关键词">
                            <el-input v-model="form.seo.keywords" placeholder="多个关键词用逗号分隔" />
                        </el-form-item>
                        <el-form-item label="SEO 描述">
                            <el-input type="textarea" v-model="form.seo.description" :rows="3" placeholder="页面描述" />
                        </el-form-item>
                    </el-collapse-item>
                </el-collapse>

                <div class="flex items-center gap-3 pt-6">
                    <el-button type="primary" size="large" :loading="saving" @click="save">保存</el-button>
                    <router-link to="/admin/products">
                        <el-button size="large">返回</el-button>
                    </router-link>
                </div>
            </el-form>
        </el-card>
    </AdminLayout>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import AdminLayout from '../components/AdminLayout.vue'
import { adminProductCreate, adminProductUpdate, adminProductDetail } from '../api/admin'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)
const loading = ref(false)
const saving = ref(false)
const categories = ref([])
const memberLevels = ref([])
const cardList = ref([])
const activeSections = ref(['basic', 'delivery', 'sku', 'member', 'seo'])

const deliveryConfig = reactive({
    card_list: [],
    text_content: '',
    api_url: '',
    api_method: 'POST',
    api_headers: {},
    api_body: {},
    api_success_field: '',
    api_auth_field: '',
    api_expire_days: 365
})

const form = reactive({
    name: '',
    category_id: 0,
    price: 0,
    stock: 0,
    description: '',
    image: '',
    images: [],
    status: 1,
    is_hidden: false,
    has_sku: false,
    sort: 0,
    delivery_type: 'card',
    delivery_config: deliveryConfig,
    skus: [],
    member_prices: [],
    seo: {
        title: '',
        keywords: '',
        description: ''
    }
})

function loadCategories() {
    fetch('/admin/api/categories').then(res => res.json()).then(data => {
        categories.value = data || []
    }).catch(() => { })
}

function loadMemberLevels() {
    fetch('/admin/api/settings/member_levels').then(res => res.json()).then(data => {
        memberLevels.value = data || []
    }).catch(() => {
        memberLevels.value = [
            { id: 1, name: '普通会员' },
            { id: 2, name: 'VIP会员' },
            { id: 3, name: '钻石会员' }
        ]
    })
}

function loadProduct() {
    if (!route.params.id) return
    loading.value = true
    adminProductDetail(route.params.id).then((data) => {
        if (data) {
            form.name = data.name || ''
            form.category_id = data.category_id || 0
            form.price = data.price || 0
            form.stock = data.stock || 0
            form.description = data.description || ''
            form.image = data.image || ''
            form.status = data.status === undefined ? 1 : data.status
            form.is_hidden = data.is_hidden || false
            form.has_sku = data.has_sku || false
            form.sort = data.sort || 0
            form.delivery_type = data.delivery_type || 'card'
            if (data.delivery_config) {
                Object.assign(deliveryConfig, data.delivery_config)
            }
            if (data.skus) {
                form.skus = data.skus.map(s => ({
                    spec_names: typeof s.spec_names === 'object' ? JSON.stringify(s.spec_names) : s.spec_names || '',
                    sku_code: s.sku_code || '',
                    price: s.price || 0,
                    stock: s.stock || 0,
                    image: s.image || ''
                }))
            }
            if (data.member_prices) {
                form.member_prices = data.member_prices
            }
            if (data.seo) {
                Object.assign(form.seo, data.seo)
            }
        }
    }).catch(() => { }).finally(() => {
        loading.value = false
    })
}

function getMemberPrice(levelId) {
    const mp = form.member_prices.find(p => p.level_id === levelId)
    return mp ? mp.price : ''
}

function updateMemberPrice(levelId, value) {
    const existing = form.member_prices.find(p => p.level_id === levelId)
    if (value && value > 0) {
        if (existing) {
            existing.price = value
        } else {
            form.member_prices.push({ level_id: levelId, price: value })
        }
    } else if (existing) {
        form.member_prices = form.member_prices.filter(p => p.level_id !== levelId)
    }
}

function handleCardFileUpload(file) {
    if (!file || !file.raw) return
    const reader = new FileReader()
    reader.onload = (e) => {
        const content = e.target.result
        cardList.value = content.split('\n').map(line => line.trim()).filter(line => line)
        deliveryConfig.card_list = cardList.value
    }
    reader.readAsText(file.raw)
}

function addApiHeader() {
    deliveryConfig.api_headers['new_key'] = ''
}

function deleteApiHeader(key) {
    delete deliveryConfig.api_headers[key]
}

function updateApiHeaderKey(newValue, oldKey) {
    if (newValue !== oldKey && newValue) {
        deliveryConfig.api_headers[newValue] = deliveryConfig.api_headers[oldKey]
        delete deliveryConfig.api_headers[oldKey]
    }
}

function addApiBody() {
    deliveryConfig.api_body['new_key'] = ''
}

function deleteApiBody(key) {
    delete deliveryConfig.api_body[key]
}

function updateApiBodyKey(newValue, oldKey) {
    if (newValue !== oldKey && newValue) {
        deliveryConfig.api_body[newValue] = deliveryConfig.api_body[oldKey]
        delete deliveryConfig.api_body[oldKey]
    }
}

function addSku() {
    form.skus.push({
        spec_names: '',
        sku_code: '',
        price: 0,
        stock: 0,
        image: ''
    })
}

function removeSku(index) {
    form.skus.splice(index, 1)
}

function save() {
    if (!form.name) {
        ElMessage.warning('请输入商品名称')
        return
    }
    if (form.price <= 0) {
        ElMessage.warning('请输入有效的价格')
        return
    }
    if (form.category_id === 0) {
        ElMessage.warning('请选择商品分类')
        return
    }
    form.delivery_config = deliveryConfig
    const validMemberPrices = []
    memberLevels.value.forEach(level => {
        const price = getMemberPrice(level.id)
        if (price && price > 0) {
            validMemberPrices.push({ level_id: level.id, price })
        }
    })
    form.member_prices = validMemberPrices

    saving.value = true
    const action = isEdit.value
        ? adminProductUpdate(route.params.id, form)
        : adminProductCreate(form)
    action.then(() => {
        ElMessage.success('保存成功')
        router.push('/admin/products')
    }).catch((err) => {
        ElMessage.error(err.message || '保存失败')
    }).finally(() => {
        saving.value = false
    })
}

onMounted(() => {
    loadCategories()
    loadMemberLevels()
    loadProduct()
})
</script>
