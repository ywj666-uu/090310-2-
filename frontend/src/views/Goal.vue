<template>
  <div class="goal-page">
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>设定减排目标</template>
          <el-form :model="form" label-width="100px">
            <el-form-item label="日均目标">
              <el-input-number v-model="form.target_emission" :min="0.1" :step="0.5" :precision="1" />
              <span style="margin-left: 8px; color: #999">kgCO₂/天</span>
            </el-form-item>
            <el-form-item label="开始日期">
              <el-date-picker v-model="form.start_date" type="date" format="YYYY-MM-DD" value-format="YYYY-MM-DD" />
            </el-form-item>
            <el-form-item label="结束日期">
              <el-date-picker v-model="form.end_date" type="date" format="YYYY-MM-DD" value-format="YYYY-MM-DD" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="createGoal" :loading="submitting">
                设定目标
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>历史目标</template>
          <el-table :data="history" empty-text="暂无历史目标">
            <el-table-column prop="target_emission" label="目标(kg/天)" width="110" />
            <el-table-column prop="start_date" label="开始" width="110" />
            <el-table-column prop="end_date" label="结束" width="110" />
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card v-if="activeGoal" class="goal-status-card">
          <template #header>当前目标进度</template>
          <div class="goal-progress">
            <el-progress
              type="dashboard"
              :percentage="activeGoal.completion_rate"
              :color="progressColor"
              :width="180"
            >
              <template #default>
                <div class="progress-inner">
                  <div class="progress-value">{{ activeGoal.current_avg.toFixed(1) }}</div>
                  <div class="progress-label">kgCO₂/天</div>
                </div>
              </template>
            </el-progress>

            <div class="goal-info">
              <div class="info-row">
                <span class="info-label">目标排放:</span>
                <span class="info-value">{{ activeGoal.goal.target_emission }} kgCO₂/天</span>
              </div>
              <div class="info-row">
                <span class="info-label">当前日均:</span>
                <span class="info-value" :class="{ achieved: activeGoal.achieved }">
                  {{ activeGoal.current_avg.toFixed(2) }} kgCO₂/天
                </span>
              </div>
              <div class="info-row">
                <span class="info-label">剩余天数:</span>
                <span class="info-value">{{ activeGoal.days_remaining }} 天</span>
              </div>
              <div class="info-row">
                <span class="info-label">状态:</span>
                <el-tag :type="activeGoal.achieved ? 'success' : 'warning'" size="small">
                  {{ activeGoal.achieved ? '已达标 ✓' : '努力中...' }}
                </el-tag>
              </div>
            </div>
          </div>

          <el-alert
            v-if="activeGoal.encouragement"
            :title="activeGoal.encouragement"
            type="success"
            show-icon
            :closable="false"
            style="margin-top: 20px"
          />
        </el-card>

        <el-card v-else>
          <el-empty description="暂未设定减排目标，快来设一个吧！" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '../api'

const submitting = ref(false)
const activeGoal = ref(null)
const history = ref([])

const today = new Date().toISOString().slice(0, 10)
const form = reactive({
  target_emission: 5.0,
  start_date: today,
  end_date: ''
})

const progressColor = computed(() => {
  if (!activeGoal.value) return '#409eff'
  return activeGoal.value.achieved ? '#67c23a' : '#e6a23c'
})

function getStatusType(status) {
  const types = { active: 'primary', completed: 'success', cancelled: 'info', failed: 'danger' }
  return types[status] || 'info'
}

function getStatusText(status) {
  const texts = { active: '进行中', completed: '已完成', cancelled: '已取消', failed: '未达成' }
  return texts[status] || status
}

async function createGoal() {
  if (!form.start_date || !form.end_date) {
    ElMessage.warning('请选择开始和结束日期')
    return
  }

  submitting.value = true
  try {
    await api.post('/goals', form)
    ElMessage.success('目标设定成功！加油！')
    await loadGoalData()
  } catch (err) {
    ElMessage.error(err.response?.data?.error || '设定失败')
  } finally {
    submitting.value = false
  }
}

async function loadGoalData() {
  try {
    const { data } = await api.get('/goals/active')
    activeGoal.value = data
  } catch {
    activeGoal.value = null
  }

  try {
    const { data } = await api.get('/goals/history')
    history.value = data || []
  } catch {
    history.value = []
  }
}

onMounted(loadGoalData)
</script>

<style scoped>
.goal-progress {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
}

.progress-inner {
  text-align: center;
}

.progress-value {
  font-size: 24px;
  font-weight: bold;
  color: #2d8c4e;
}

.progress-label {
  font-size: 12px;
  color: #999;
}

.goal-info {
  margin-top: 20px;
  width: 100%;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.info-label {
  color: #666;
}

.info-value {
  font-weight: 500;
}

.info-value.achieved {
  color: #67c23a;
}
</style>
