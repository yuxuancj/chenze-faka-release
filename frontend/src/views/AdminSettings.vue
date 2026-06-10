<template>
  <div class="settings">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">系统设置</span>
        </div>
      </template>

      <el-tabs v-model="activeTab" type="border-card">
        <el-tab-pane label="基础设置" name="basic">
          <el-form
            ref="basicFormRef"
            :model="basicForm"
            :rules="basicRules"
            label-width="140px"
            class="settings-form"
          >
            <el-form-item label="网站名称" prop="siteName">
              <el-input v-model="basicForm.siteName" placeholder="请输入网站名称" />
            </el-form-item>
            <el-form-item label="网站标题" prop="siteTitle">
              <el-input v-model="basicForm.siteTitle" placeholder="请输入网站标题" />
            </el-form-item>
            <el-form-item label="网站LOGO">
              <el-upload
                class="logo-uploader"
                :show-file-list="false"
                :auto-upload="false"
                accept="image/*"
              >
                <img v-if="basicForm.logo" :src="basicForm.logo" class="logo-image" />
                <el-icon v-else class="logo-uploader-icon"><Plus /></el-icon>
              </el-upload>
              <div class="upload-tip">建议尺寸 200x200，支持 jpg/png 格式</div>
            </el-form-item>
            <el-form-item label="备案号">
              <el-input v-model="basicForm.icp" placeholder="请输入ICP备案号" />
            </el-form-item>
            <el-form-item label="客服邮箱">
              <el-input v-model="basicForm.email" placeholder="请输入客服邮箱" />
            </el-form-item>
            <el-form-item label="是否开启站点">
              <el-switch v-model="basicForm.siteOpen" active-text="开启" inactive-text="关闭" />
            </el-form-item>
            <el-form-item label="站点关闭提示">
              <el-input
                v-model="basicForm.closeTip"
                type="textarea"
                :rows="3"
                placeholder="站点关闭时显示的提示信息"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="saveBasic">保存设置</el-button>
              <el-button @click="resetBasic">重置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="支付设置" name="payment">
          <el-form
            ref="paymentFormRef"
            :model="paymentForm"
            label-width="140px"
            class="settings-form"
          >
            <el-form-item label="支付宝支付">
              <el-switch v-model="paymentForm.alipayEnabled" active-text="开启" inactive-text="关闭" />
            </el-form-item>
            <el-form-item label="支付宝AppID" v-show="paymentForm.alipayEnabled">
              <el-input v-model="paymentForm.alipayAppId" placeholder="请输入支付宝AppID" />
            </el-form-item>
            <el-form-item label="支付宝私钥" v-show="paymentForm.alipayEnabled">
              <el-input
                v-model="paymentForm.alipayPrivateKey"
                type="textarea"
                :rows="4"
                placeholder="请输入支付宝私钥"
                show-password
              />
            </el-form-item>
            <el-form-item label="微信支付">
              <el-switch v-model="paymentForm.wechatEnabled" active-text="开启" inactive-text="关闭" />
            </el-form-item>
            <el-form-item label="微信商户号" v-show="paymentForm.wechatEnabled">
              <el-input v-model="paymentForm.wechatMchId" placeholder="请输入微信商户号" />
            </el-form-item>
            <el-form-item label="微信API密钥" v-show="paymentForm.wechatEnabled">
              <el-input v-model="paymentForm.wechatApiKey" type="password" placeholder="请输入微信API密钥" show-password />
            </el-form-item>
            <el-form-item label="支付手续费(%)">
              <el-input-number v-model="paymentForm.fee" :min="0" :max="100" :step="0.1" />
            </el-form-item>
            <el-form-item label="最低充值金额">
              <el-input-number v-model="paymentForm.minRecharge" :min="1" :step="1" />
              <span style="margin-left: 10px; color: #909399;">元</span>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="savePayment">保存设置</el-button>
              <el-button @click="resetPayment">重置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="商品设置" name="product">
          <el-form
            ref="productFormRef"
            :model="productForm"
            label-width="140px"
            class="settings-form"
          >
            <el-form-item label="自动发货">
              <el-switch v-model="productForm.autoDeliver" active-text="开启" inactive-text="关闭" />
            </el-form-item>
            <el-form-item label="库存预警阈值">
              <el-input-number v-model="productForm.stockWarning" :min="0" :step="1" />
              <span style="margin-left: 10px; color: #909399;">件</span>
            </el-form-item>
            <el-form-item label="商品默认分类">
              <el-select v-model="productForm.defaultCategory" placeholder="请选择默认分类" style="width: 100%">
                <el-option label="影音会员" value="1" />
                <el-option label="游戏点卡" value="2" />
                <el-option label="软件激活" value="3" />
                <el-option label="其他商品" value="4" />
              </el-select>
            </el-form-item>
            <el-form-item label="允许用户评论">
              <el-switch v-model="productForm.allowComment" active-text="允许" inactive-text="禁止" />
            </el-form-item>
            <el-form-item label="评论需要审核">
              <el-switch v-model="productForm.commentReview" active-text="需要" inactive-text="不需要" />
            </el-form-item>
            <el-form-item label="商品展示方式">
              <el-radio-group v-model="productForm.displayMode">
                <el-radio-button label="grid">网格布局</el-radio-button>
                <el-radio-button label="list">列表布局</el-radio-button>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="每页商品数量">
              <el-input-number v-model="productForm.pageSize" :min="6" :max="100" :step="6" />
              <span style="margin-left: 10px; color: #909399;">件/页</span>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="saveProduct">保存设置</el-button>
              <el-button @click="resetProduct">重置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="管理员设置" name="admin">
          <el-form
            ref="adminFormRef"
            :model="adminForm"
            :rules="adminRules"
            label-width="140px"
            class="settings-form"
          >
            <el-form-item label="管理员账号">
              <el-input v-model="adminForm.username" disabled />
            </el-form-item>
            <el-form-item label="当前密码" prop="oldPassword">
              <el-input v-model="adminForm.oldPassword" type="password" show-password placeholder="请输入当前密码" />
            </el-form-item>
            <el-form-item label="新密码" prop="newPassword">
              <el-input v-model="adminForm.newPassword" type="password" show-password placeholder="请输入新密码" />
            </el-form-item>
            <el-form-item label="确认新密码" prop="confirmPassword">
              <el-input v-model="adminForm.confirmPassword" type="password" show-password placeholder="请再次输入新密码" />
            </el-form-item>
            <el-form-item label="登录IP限制">
              <el-switch v-model="adminForm.ipLimit" active-text="开启" inactive-text="关闭" />
            </el-form-item>
            <el-form-item label="允许登录IP">
              <el-input
                v-model="adminForm.allowedIps"
                type="textarea"
                :rows="3"
                placeholder="一行一个IP，支持IP段如 192.168.1.*，留空表示不限制"
              />
            </el-form-item>
            <el-form-item label="登录失败锁定">
              <el-input-number v-model="adminForm.loginFailLock" :min="3" :max="20" :step="1" />
              <span style="margin-left: 10px; color: #909399;">次后锁定账户</span>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="saveAdmin">保存修改</el-button>
              <el-button @click="resetAdmin">重置</el-button>
              <el-button type="danger" @click="clearCache">清除缓存</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

const activeTab = ref('basic')
const saving = ref(false)

const basicFormRef = ref(null)
const paymentFormRef = ref(null)
const productFormRef = ref(null)
const adminFormRef = ref(null)

const basicForm = reactive({
  siteName: '辰泽发卡',
  siteTitle: '辰泽发卡系统 - 专业的虚拟商品交易平台',
  logo: '',
  icp: '粤ICP备12345678号',
  email: 'support@chenze.com',
  siteOpen: true,
  closeTip: '系统维护中，请稍后访问',
})

const basicRules = {
  siteName: [{ required: true, message: '请输入网站名称', trigger: 'blur' }],
  siteTitle: [{ required: true, message: '请输入网站标题', trigger: 'blur' }],
}

const paymentForm = reactive({
  alipayEnabled: true,
  alipayAppId: '2021000000000000',
  alipayPrivateKey: '',
  wechatEnabled: true,
  wechatMchId: '1234567890',
  wechatApiKey: '',
  fee: 0.6,
  minRecharge: 10,
})

const productForm = reactive({
  autoDeliver: true,
  stockWarning: 10,
  defaultCategory: '1',
  allowComment: true,
  commentReview: false,
  displayMode: 'grid',
  pageSize: 12,
})

const adminForm = reactive({
  username: 'admin',
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
  ipLimit: false,
  allowedIps: '',
  loginFailLock: 5,
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入新密码'))
  } else if (value !== adminForm.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const adminRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
  ],
  confirmPassword: [{ required: true, validator: validateConfirmPassword, trigger: 'blur' }],
}

const saveBasic = async () => {
  if (!basicFormRef.value) return
  await basicFormRef.value.validate((valid) => {
    if (valid) {
      saving.value = true
      setTimeout(() => {
        ElMessage.success('基础设置已保存')
        saving.value = false
      }, 800)
    }
  })
}

const resetBasic = () => {
  basicForm.siteName = '辰泽发卡'
  basicForm.siteTitle = '辰泽发卡系统 - 专业的虚拟商品交易平台'
  basicForm.logo = ''
  basicForm.icp = ''
  basicForm.email = ''
  basicForm.siteOpen = true
  basicForm.closeTip = '系统维护中，请稍后访问'
  basicFormRef.value?.resetFields()
}

const savePayment = () => {
  saving.value = true
  setTimeout(() => {
    ElMessage.success('支付设置已保存')
    saving.value = false
  }, 800)
}

const resetPayment = () => {
  paymentForm.alipayEnabled = false
  paymentForm.alipayAppId = ''
  paymentForm.alipayPrivateKey = ''
  paymentForm.wechatEnabled = false
  paymentForm.wechatMchId = ''
  paymentForm.wechatApiKey = ''
  paymentForm.fee = 0
  paymentForm.minRecharge = 10
}

const saveProduct = () => {
  saving.value = true
  setTimeout(() => {
    ElMessage.success('商品设置已保存')
    saving.value = false
  }, 800)
}

const resetProduct = () => {
  productForm.autoDeliver = true
  productForm.stockWarning = 10
  productForm.defaultCategory = ''
  productForm.allowComment = true
  productForm.commentReview = false
  productForm.displayMode = 'grid'
  productForm.pageSize = 12
}

const saveAdmin = async () => {
  if (!adminFormRef.value) return
  await adminFormRef.value.validate((valid) => {
    if (valid) {
      saving.value = true
      setTimeout(() => {
        ElMessage.success('管理员设置已保存')
        adminForm.oldPassword = ''
        adminForm.newPassword = ''
        adminForm.confirmPassword = ''
        saving.value = false
      }, 800)
    }
  })
}

const resetAdmin = () => {
  adminForm.oldPassword = ''
  adminForm.newPassword = ''
  adminForm.confirmPassword = ''
  adminForm.ipLimit = false
  adminForm.allowedIps = ''
  adminForm.loginFailLock = 5
  adminFormRef.value?.resetFields()
}

const clearCache = () => {
  ElMessageBox.confirm('确定要清除所有系统缓存吗？此操作不可撤销', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(() => {
      ElMessage.success('缓存已清除')
    })
    .catch(() => {})
}
</script>

<style scoped>
.settings {
  padding: 0;
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

.settings-form {
  max-width: 800px;
  padding: 20px 0;
}

.logo-uploader {
  display: flex;
  align-items: center;
}

.logo-uploader :deep(.el-upload) {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: border-color 0.3s;
  width: 100px;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-uploader :deep(.el-upload:hover) {
  border-color: #409eff;
}

.logo-uploader-icon {
  font-size: 28px;
  color: #8c939d;
}

.logo-image {
  width: 100px;
  height: 100px;
  object-fit: contain;
  display: block;
}

.upload-tip {
  font-size: 12px;
  color: #909399;
  margin-left: 12px;
  display: inline-block;
}
</style>
