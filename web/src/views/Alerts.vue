<template>
  <div>
    <div class="page-header">
      <h3>告警中心</h3>
      <el-button type="primary" @click="showCreateRule = true">创建规则</el-button>
    </div>
    <el-select v-model="clusterId" placeholder="选择集群" @change="fetchData" style="width:240px;margin-top:16px">
      <el-option v-for="c in clusters" :key="c.id" :label="c.name" :value="c.id" />
    </el-select>

    <el-tabs v-model="activeTab" style="margin-top:16px">
      <el-tab-pane label="告警规则" name="rules">
        <el-table :data="rules" v-loading="loading">
          <el-table-column prop="name" label="规则名称" />
          <el-table-column prop="metric" label="监控指标" width="140">
            <template #default="{ row }">
              <el-tag size="small">{{ metricLabel(row.metric) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="条件" width="120">
            <template #default="{ row }">{{ row.operator }} {{ row.threshold }}</template>
          </el-table-column>
          <el-table-column prop="enabled" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '停用' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="160">
            <template #default="{ row }">
              <el-button size="small" @click="toggleRule(row)">{{ row.enabled ? '停用' : '启用' }}</el-button>
              <el-button size="small" type="danger" @click="deleteRule(row.id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
      <el-tab-pane label="告警事件" name="events">
        <el-table :data="events" v-loading="loading">
          <el-table-column prop="level" label="级别" width="100">
            <template #default="{ row }">
              <el-tag :type="row.level === 'critical' ? 'danger' : 'warning'">{{ row.level }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="message" label="事件描述" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'triggered' ? 'danger' : 'success'">{{ row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="triggered_at" label="触发时间" width="180" />
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="showCreateRule" title="创建告警规则" width="500px">
      <el-form :model="ruleForm">
        <el-form-item label="集群">
          <el-select v-model="ruleForm.cluster_id" style="width:100%">
            <el-option v-for="c in clusters" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="规则名称">
          <el-input v-model="ruleForm.name" placeholder="如：消费积压告警" />
        </el-form-item>
        <el-form-item label="监控指标">
          <el-select v-model="ruleForm.metric" style="width:100%">
            <el-option label="消费积压(Lag)" value="lag" />
            <el-option label="TPS异常下降" value="tps_drop" />
            <el-option label="Broker离线" value="broker_down" />
            <el-option label="磁盘使用率" value="disk_usage" />
          </el-select>
        </el-form-item>
        <el-form-item label="条件">
          <el-select v-model="ruleForm.operator" style="width:100px">
            <el-option label=">" value=">" />
            <el-option label=">=" value=">=" />
            <el-option label="=" value="=" />
          </el-select>
        </el-form-item>
        <el-form-item label="阈值">
          <el-input-number v-model="ruleForm.threshold" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateRule = false">取消</el-button>
        <el-button type="primary" @click="createRule">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '../api/request'
import { ElMessage } from 'element-plus'

const clusters = ref([])
const clusterId = ref(null)
const rules = ref([])
const events = ref([])
const loading = ref(false)
const activeTab = ref('rules')
const showCreateRule = ref(false)
const ruleForm = ref({ cluster_id: null, name: '', metric: 'lag', operator: '>', threshold: 1000 })

function metricLabel(m) {
  const map = { lag: '消费积压', tps_drop: 'TPS下降', broker_down: 'Broker离线', disk_usage: '磁盘使用率' }
  return map[m] || m
}

async function fetchClusters() {
  const res = await request.get('/clusters')
  clusters.value = res.data || []
}

async function fetchData() {
  if (!clusterId.value) return
  loading.value = true
  try {
    const [r, e] = await Promise.all([
      request.get('/alerts/rules', { params: { cluster_id: clusterId.value } }),
      request.get('/alerts/events', { params: { cluster_id: clusterId.value } })
    ])
    rules.value = r.data || []
    events.value = e.data || []
  } finally { loading.value = false }
}

async function createRule() {
  await request.post('/alerts/rules', ruleForm.value)
  ElMessage.success('规则创建成功')
  showCreateRule.value = false
  fetchData()
}

async function toggleRule(row) {
  await request.put(`/alerts/rules/${row.id}`, { enabled: !row.enabled })
  ElMessage.success('操作成功')
  fetchData()
}

async function deleteRule(id) {
  await request.delete(`/alerts/rules/${id}`)
  ElMessage.success('删除成功')
  fetchData()
}

onMounted(fetchClusters)
</script>

<style scoped>
.page-header { display: flex; align-items: center; justify-content: space-between; }
</style>