<template>
  <div class="trend-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>碳排放趋势分析</span>
          <div class="header-controls">
            <el-select v-model="region" placeholder="对比区域" size="small" style="width: 120px; margin-right: 12px" @change="loadTrend">
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
            <el-radio-group v-model="period" @change="loadTrend">
              <el-radio-button value="weekly">本周</el-radio-button>
              <el-radio-button value="monthly">本月</el-radio-button>
            </el-radio-group>
          </div>
        </div>
      </template>
      <v-chart :option="trendOption" style="height: 400px" autoresize />
    </el-card>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>分类碳排放对比</template>
          <v-chart :option="barOption" style="height: 300px" autoresize />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>与{{ region }}平均值对比</template>
          <div class="comparison">
            <div class="comp-item">
              <div class="comp-label">你的日均排放</div>
              <div class="comp-value user-value">{{ userAvg.toFixed(2) }} kgCO₂</div>
            </div>
            <div class="comp-vs">VS</div>
            <div class="comp-item">
              <div class="comp-label">区域平均排放</div>
              <div class="comp-value region-value">{{ regionalAvg.toFixed(2) }} kgCO₂</div>
            </div>
          </div>
          <div class="comp-result" :class="resultClass">
            {{ comparisonText }}
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent, MarkLineComponent } from 'echarts/components'
import api from '../api'

use([CanvasRenderer, LineChart, BarChart, GridComponent, TooltipComponent, LegendComponent, MarkLineComponent])

const period = ref('weekly')
const region = ref('全国')
const trendData = ref([])
const regionalAvg = ref(0)
const categoryData = ref([])
const userAvg = ref(0)

const trendOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  legend: { data: ['每日碳排放', '区域平均'] },
  xAxis: {
    type: 'category',
    data: trendData.value.map(d => d.date.slice(5))
  },
  yAxis: { type: 'value', name: 'kgCO₂' },
  series: [
    {
      name: '每日碳排放',
      type: 'line',
      data: trendData.value.map(d => d.total_emission.toFixed(2)),
      smooth: true,
      areaStyle: { opacity: 0.2 },
      itemStyle: { color: '#2d8c4e' },
      lineStyle: { width: 3 }
    },
    {
      name: '区域平均',
      type: 'line',
      data: trendData.value.map(() => regionalAvg.value),
      lineStyle: { type: 'dashed', color: '#ff6b6b' },
      itemStyle: { color: '#ff6b6b' },
      symbol: 'none'
    }
  ]
}))

const barOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  xAxis: {
    type: 'category',
    data: categoryData.value.map(d => getCategoryName(d.category))
  },
  yAxis: { type: 'value', name: 'kgCO₂' },
  series: [{
    type: 'bar',
    data: categoryData.value.map(d => ({
      value: parseFloat(d.total_emission.toFixed(2)),
      itemStyle: { color: getCategoryColor(d.category) }
    })),
    barWidth: '50%'
  }]
}))

const resultClass = computed(() => userAvg.value <= regionalAvg.value ? 'good' : 'bad')

const comparisonText = computed(() => {
  if (regionalAvg.value === 0) return '暂无区域对比数据'
  const diff = ((userAvg.value - regionalAvg.value) / regionalAvg.value * 100).toFixed(1)
  if (userAvg.value <= regionalAvg.value) {
    return `太棒了！你的碳排放比区域平均低 ${Math.abs(diff)}%`
  }
  return `你的碳排放比区域平均高 ${diff}%，继续努力减排！`
})

function getCategoryName(cat) {
  const names = { transport: '出行', electricity: '用电', diet: '饮食' }
  return names[cat] || cat
}

function getCategoryColor(cat) {
  const colors = { transport: '#409eff', electricity: '#e6a23c', diet: '#67c23a' }
  return colors[cat] || '#909399'
}

async function loadTrend() {
  try {
    const { data } = await api.get(`/trend/${period.value}`, { params: { region: region.value } })
    trendData.value = data.user_data || []
    regionalAvg.value = data.regional_average || 0

    if (trendData.value.length > 0) {
      const total = trendData.value.reduce((s, d) => s + d.total_emission, 0)
      userAvg.value = total / trendData.value.length
    }
  } catch (err) {
    console.error('加载趋势失败', err)
  }
}

async function loadCategory() {
  try {
    const { data } = await api.get('/summary/category')
    categoryData.value = data || []
  } catch (err) {
    console.error('加载分类失败', err)
  }
}

onMounted(() => {
  loadTrend()
  loadCategory()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-controls {
  display: flex;
  align-items: center;
}

.comparison {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30px 0;
  gap: 30px;
}

.comp-item {
  text-align: center;
}

.comp-label {
  color: #999;
  margin-bottom: 8px;
}

.comp-value {
  font-size: 24px;
  font-weight: bold;
}

.user-value { color: #2d8c4e; }
.region-value { color: #ff6b6b; }

.comp-vs {
  font-size: 20px;
  font-weight: bold;
  color: #ccc;
}

.comp-result {
  text-align: center;
  padding: 12px;
  border-radius: 8px;
  font-weight: 500;
}

.comp-result.good {
  background: #f0f9eb;
  color: #67c23a;
}

.comp-result.bad {
  background: #fef0f0;
  color: #f56c6c;
}
</style>
