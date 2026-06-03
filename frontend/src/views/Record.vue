<template>
  <div class="record-page">
    <el-row :gutter="20">
      <el-col :span="14">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>记录碳排放</span>
              <el-date-picker v-model="recordDate" type="date" format="YYYY-MM-DD" value-format="YYYY-MM-DD" placeholder="选择日期" size="small" />
            </div>
          </template>

          <el-form :model="form" label-width="80px">
            <el-form-item label="分类">
              <el-radio-group v-model="form.category" @change="onCategoryChange">
                <el-radio-button value="transport">🚗 出行</el-radio-button>
                <el-radio-button value="electricity">⚡ 用电</el-radio-button>
                <el-radio-button value="diet">🍽️ 饮食</el-radio-button>
              </el-radio-group>
            </el-form-item>

            <el-form-item label="具体项目">
              <el-select v-model="form.item" placeholder="选择项目" style="width: 100%">
                <el-option
                  v-for="f in filteredFactors"
                  :key="f.item"
                  :label="`${f.item} (${f.factor} ${f.unit})`"
                  :value="f.item"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="数量">
              <el-input-number v-model="form.amount" :min="0.1" :step="1" :precision="1" style="width: 200px" />
              <span class="unit-hint">{{ currentUnit }}</span>
            </el-form-item>

            <el-form-item label="备注">
              <el-input v-model="form.note" placeholder="可选备注" />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="submitRecord" :loading="submitting">
                提交记录
              </el-button>
              <span class="emission-preview" v-if="previewEmission > 0">
                预估排放: <strong>{{ previewEmission.toFixed(2) }} kgCO₂</strong>
              </span>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <el-col :span="10">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>{{ recordDate }} 记录</span>
              <el-tag type="warning">总计: {{ totalEmission.toFixed(2) }} kgCO₂</el-tag>
            </div>
          </template>

          <el-table :data="records" style="width: 100%" empty-text="暂无记录">
            <el-table-column prop="category" label="分类" width="70">
              <template #default="{ row }">
                {{ getCategoryIcon(row.category) }}
              </template>
            </el-table-column>
            <el-table-column prop="item" label="项目" width="80" />
            <el-table-column prop="amount" label="数量" width="70" />
            <el-table-column prop="emission" label="排放(kg)" width="90">
              <template #default="{ row }">
                {{ row.emission.toFixed(2) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="60">
              <template #default="{ row }">
                <el-button type="danger" size="small" link @click="deleteRecord(row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'

const recordDate = ref(new Date().toISOString().slice(0, 10))
const factors = ref([])
const records = ref([])
const totalEmission = ref(0)
const submitting = ref(false)

const form = reactive({
  category: 'transport',
  item: '',
  amount: 1,
  note: ''
})

const filteredFactors = computed(() =>
  factors.value.filter(f => f.category === form.category)
)

const currentUnit = computed(() => {
  const f = factors.value.find(f => f.item === form.item)
  return f ? f.unit : ''
})

const previewEmission = computed(() => {
  const f = factors.value.find(f => f.item === form.item)
  return f ? form.amount * f.factor : 0
})

function getCategoryIcon(cat) {
  const icons = { transport: '🚗', electricity: '⚡', diet: '🍽️' }
  return icons[cat] || '📝'
}

function onCategoryChange() {
  form.item = ''
}

async function loadFactors() {
  const { data } = await api.get('/factors')
  factors.value = data
}

async function loadRecords() {
  const { data } = await api.get('/records', { params: { date: recordDate.value } })
  records.value = data.records || []
  totalEmission.value = data.total_emission || 0
}

async function submitRecord() {
  if (!form.item) {
    ElMessage.warning('请选择具体项目')
    return
  }

  submitting.value = true
  try {
    const { data } = await api.post('/records', {
      record_date: recordDate.value,
      category: form.category,
      item: form.item,
      amount: form.amount,
      note: form.note
    })
    ElMessage.success('记录成功')
    form.amount = 1
    form.note = ''
    await loadRecords()

    if (data.goal_achieved) {
      ElMessageBox.alert(data.encouragement, '🎉 减排目标达成！', {
        confirmButtonText: '太棒了',
        type: 'success'
      })
    }
  } catch (err) {
    ElMessage.error(err.response?.data?.error || '提交失败')
  } finally {
    submitting.value = false
  }
}

async function deleteRecord(id) {
  await api.delete(`/records/${id}`)
  ElMessage.success('已删除')
  await loadRecords()
}

watch(recordDate, loadRecords)

onMounted(async () => {
  await loadFactors()
  await loadRecords()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.unit-hint {
  margin-left: 12px;
  color: #999;
}

.emission-preview {
  margin-left: 16px;
  color: #e6a23c;
}
</style>
