<template>
  <div>
    <h3>仪表盘</h3>
    <el-row :gutter="20" style="margin-top:20px">
      <el-col :span="6" v-for="card in cards" :key="card.title">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-info">
              <div class="stat-title">{{ card.title }}</div>
              <div class="stat-value">{{ card.value }}</div>
            </div>
            <el-icon :size="40" :color="card.color"><component :is="card.icon" /></el-icon>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top:20px">
      <el-col :span="12">
        <el-card>
          <template #header>吞吐量 (TPS)</template>
          <div ref="tpsChart" style="height:300px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>消费延迟 (Lag)</template>
          <div ref="lagChart" style="height:300px"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, markRaw } from 'vue'
import { Monitor, Connection, Coin, Warning } from '@element-plus/icons-vue'
import * as echarts from 'echarts'

const cards = [
  { title: '集群数', value: '0', icon: markRaw(Monitor), color: '#409EFF' },
  { title: 'Topic数', value: '0', icon: markRaw(Connection), color: '#67C23A' },
  { title: '消息总量', value: '0', icon: markRaw(Coin), color: '#E6A23C' },
  { title: '告警数', value: '0', icon: markRaw(Warning), color: '#F56C6C' }
]

const tpsChart = ref()
const lagChart = ref()

onMounted(() => {
  const tps = echarts.init(tpsChart.value)
  tps.setOption({
    xAxis: { type: 'category', data: [] },
    yAxis: { type: 'value', name: '条/秒' },
    series: [{ data: [], type: 'line', smooth: true, areaStyle: {} }]
  })

  const lag = echarts.init(lagChart.value)
  lag.setOption({
    xAxis: { type: 'category', data: [] },
    yAxis: { type: 'value', name: '条' },
    series: [{ data: [], type: 'bar' }]
  })
})
</script>

<style scoped>
.stat-card { display: flex; align-items: center; justify-content: space-between; }
.stat-title { font-size: 14px; color: #999; }
.stat-value { font-size: 28px; font-weight: bold; color: #333; margin-top: 8px; }
</style>