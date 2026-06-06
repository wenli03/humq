<template>
  <div>
    <h3>消费组监控</h3>
    <div style="margin-top:16px">
      <el-select v-model="clusterId" placeholder="选择集群" @change="fetchConsumers" style="width:240px">
        <el-option v-for="c in clusters" :key="c.id" :label="c.name" :value="c.id" />
      </el-select>
    </div>
    <el-table :data="consumers" style="margin-top:16px" v-loading="loading">
      <el-table-column prop="group_id" label="消费组" width="200" />
      <el-table-column prop="topics" label="订阅Topic">
        <template #default="{ row }">
          <el-tag v-for="t in row.topics" :key="t" size="small" style="margin-right:4px">{{ t }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="members" label="成员数" width="80" />
      <el-table-column prop="state" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.state === 'Stable' ? 'success' : 'warning'">{{ row.state }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="lag" label="积压量(Lag)" width="140" sortable>
        <template #default="{ row }">
          <span :style="{ color: row.lag > 1000 ? '#F56C6C' : '#67C23A', fontWeight: 'bold' }">
            {{ row.lag.toLocaleString() }}
          </span>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '../api/request'

const consumers = ref([])
const clusters = ref([])
const clusterId = ref(null)
const loading = ref(false)

async function fetchClusters() {
  const res = await request.get('/clusters')
  clusters.value = res.data || []
}

async function fetchConsumers() {
  if (!clusterId.value) return
  loading.value = true
  try {
    const res = await request.get('/consumers', { params: { cluster_id: clusterId.value } })
    consumers.value = res.data || []
  } finally { loading.value = false }
}

onMounted(fetchClusters)
</script>