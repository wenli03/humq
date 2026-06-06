<template>
  <div>
    <h3>消息追踪</h3>
    <el-row :gutter="20" style="margin-top:20px">
      <el-col :span="5">
        <el-select v-model="clusterId" placeholder="选择集群" style="width:100%" @change="fetchClusters">
          <el-option v-for="c in clusters" :key="c.id" :label="c.name" :value="c.id" />
        </el-select>
      </el-col>
      <el-col :span="5">
        <el-input v-model="searchTopic" placeholder="Topic名称" />
      </el-col>
      <el-col :span="4">
        <el-input-number v-model="searchPartition" :min="0" placeholder="分区" style="width:100%" />
      </el-col>
      <el-col :span="4">
        <el-input-number v-model="searchOffset" :min="0" placeholder="Offset" style="width:100%" />
      </el-col>
      <el-col :span="6">
        <el-button type="primary" @click="traceMessages" :loading="tracing">查询消息</el-button>
      </el-col>
    </el-row>
    <el-table :data="messages" style="margin-top:16px" v-loading="tracing" empty-text="请输入查询条件后点击查询">
      <el-table-column prop="topic" label="Topic" width="160" />
      <el-table-column prop="partition" label="分区" width="80" />
      <el-table-column prop="offset" label="Offset" width="120" />
      <el-table-column prop="key" label="Key" width="200" />
      <el-table-column prop="value" label="Value" show-overflow-tooltip />
      <el-table-column prop="timestamp" label="时间" width="180" />
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '../api/request'

const clusters = ref([])
const clusterId = ref(null)
const searchTopic = ref('')
const searchPartition = ref(0)
const searchOffset = ref(0)
const messages = ref([])
const tracing = ref(false)

async function fetchClusters() {
  const res = await request.get('/clusters')
  clusters.value = res.data || []
}

onMounted(fetchClusters)

async function traceMessages() {
  tracing.value = true
  try {
    const res = await request.post('/messages/trace', {
      cluster_id: clusterId.value,
      topic: searchTopic.value,
      partition: searchPartition.value,
      offset: searchOffset.value
    })
    messages.value = res.data || []
  } finally { tracing.value = false }
}
</script>