<template>
  <div class="dashboard">
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-value">{{ todayEmission.toFixed(2) }}</div>
          <div class="stat-label">今日碳排放 (kgCO₂)</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-value">{{ weekAvg.toFixed(2) }}</div>
          <div class="stat-label">本周日均 (kgCO₂)</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-value">{{ regionalAvg.toFixed(2) }}</div>
          <div class="stat-label">区域平均 (kgCO₂)</div>
          <el-select v-model="region" size="small" style="margin-top: 8px; width: 100%" @change="loadData">
            <el-option label="全国" value="全国" />
            <el-option label="北京" value="北京" />
            <el-option label="上海" value="上海" />
            <el-option label="广州" value="广州" />
            <el-option label="深圳" value="深圳" />
            <el-option label="成都" value="成都" />
            <el-option label="杭州" value="杭州" />
            <el-option label="武汉" value="武汉" />
            <el-option label="南京" value="南京" />
            <el-option label="重庆" value="重庆" />
            <el-option label="西安" value="西安" />
          </el-select>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" :class="{ 'goal-achieved': goalAchieved }">
          <div class="stat-value">{{ goalStatus }}</div>
          <div class="stat-label">目标状态</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="14">
        <el-card>
          <template #header>近7日碳排放趋势</template>
          <v-chart :option="trendOption" style="height: 300px" autoresize />
        </el-card>
      </el-col>
      <el-col :span="10">
        <el-card>
          <template #header>碳排放分类构成</template>
          <v-chart :option="pieOption" style="height: 300px" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px" v-if="encouragement">
      <el-col :span="24">
        <el-alert :title="encouragement" type="success" show-icon :closable="false" />
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, PieChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent, MarkLineComponent } from 'echarts/components'
import api from '../api'

use([CanvasRenderer, LineChart, PieChart, GridComponent, TooltipComponent, LegendComponent, MarkLineComponent])

const todayEmission = ref(0)
const weekAvg = ref(0)
const regionalAvg = ref(0)
const goalAchieved = ref(false)
const goalStatus = ref('未设定')
const encouragement = ref('')
const weeklyData = ref([])
const categoryData = ref([])
const region = ref('全国')

const trendOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  xAxis: {
    type: 'category',
    data: weeklyData.value.map(d => d.date.slice(5))
  },
  yAxis: {
    type: 'value',
    name: 'kgCO₂'
  },
  series: [
    {
      name: '每日碳排放',
      type: 'line',
      data: weeklyData.value.map(d => d.total_emission),
      smooth: true,
      areaStyle: { opacity: 0.3 },
      itemStyle: { color: '#2d8c4e' },
      markLine: {
        data: [{ yAxis: regionalAvg.value, name: '区域平均' }],
        lineStyle: { color: '#ff6b6b', type: 'dashed' },
        label: { formatter: '区域平均: {c}' }
      }
    }
  ]
}))

const pieOption = computed(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c} kgCO₂ ({d}%)' },
  legend: { bottom: 0 },
  series: [{
    type: 'pie',
    radius: ['40%', '70%'],
    data: categoryData.value.map(d => ({
      name: getCategoryName(d.category),
      value: parseFloat(d.total_emission.toFixed(2))
    })),
    itemStyle: {
      borderRadius: 6,
      borderColor: '#fff',
      borderWidth: 2
    }
  }]
}))

function getCategoryName(cat) {
  const names = { transport: '出行', electricity: '用电', diet: '饮食' }
  return names[cat] || cat
}

async function loadData() {
  try {
    const [dailyRes, trendRes, categoryRes, goalRes] = await Promise.all([
      api.get('/records'),
      api.get('/trend/weekly', { params: { region: region.value } }),
      api.get('/summary/category'),
      api.get('/goals/active').catch(() => ({ data: null }))
    ])

    todayEmission.value = dailyRes.data.total_emission || 0
    weeklyData.value = trendRes.data.user_data || []
    regionalAvg.value = trendRes.data.regional_average || 0

    if (weeklyData.value.length > 0) {
      const total = weeklyData.value.reduce((s, d) => s + d.total_emission, 0)
      weekAvg.value = total / weeklyData.value.length
    }

    categoryData.value = categoryRes.data || []

    if (goalRes.data) {
      goalAchieved.value = goalRes.data.achieved
      goalStatus.value = goalRes.data.achieved ? '已达成 ✓' : '进行中'
      encouragement.value = goalRes.data.encouragement || ''
    }
  } catch (err) {
    console.error('加载数据失败', err)
  }
}

onMounted(loadData)
</script>

<style scoped>
.stats-row .stat-card {
  text-align: center;
  transition: transform 0.3s;
}

.stats-row .stat-card:hover {
  transform: translateY(-4px);
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #2d8c4e;
  margin-bottom: 8px;
}

.stat-label {
  color: #999;
  font-size: 14px;
}

.goal-achieved .stat-value {
  color: #67c23a;
}
</style>
