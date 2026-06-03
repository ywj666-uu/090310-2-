<template>
  <div class="login-container">
    <div class="login-card">
      <h2>🌍 注册账号</h2>
      <p class="subtitle">开始你的低碳之旅</p>
      <el-form :model="form" :rules="rules" ref="formRef">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="用户名（至少3位）" prefix-icon="User" size="large" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="密码（至少6位）" prefix-icon="Lock" size="large" show-password />
        </el-form-item>
        <el-form-item prop="email">
          <el-input v-model="form.email" placeholder="邮箱（选填）" prefix-icon="Message" size="large" />
        </el-form-item>
        <el-form-item prop="region">
          <el-select v-model="form.region" placeholder="选择地区" size="large" style="width: 100%">
            <el-option label="全国" value="全国" />
            <el-option label="北京" value="北京" />
            <el-option label="上海" value="上海" />
            <el-option label="广州" value="广州" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" style="width: 100%" @click="handleRegister" :loading="loading">
            注册
          </el-button>
        </el-form-item>
      </el-form>
      <p class="register-link">
        已有账号？<router-link to="/login">去登录</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import api from '../api'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  username: '',
  password: '',
  email: '',
  region: '全国'
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, message: '用户名至少3位', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ]
}

async function handleRegister() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await api.post('/register', form)
    ElMessage.success('注册成功，请登录')
    router.push('/login')
  } catch (err) {
    ElMessage.error(err.response?.data?.error || '注册失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #2d8c4e 100%);
}

.login-card {
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  width: 400px;
}

.login-card h2 {
  text-align: center;
  margin-bottom: 8px;
  color: #333;
}

.subtitle {
  text-align: center;
  color: #999;
  margin-bottom: 30px;
}

.register-link {
  text-align: center;
  color: #666;
  margin-top: 16px;
}
</style>
