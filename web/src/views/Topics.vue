<template>
  <div>
    <div class="page-header">
      <h3>Topic管理</h3>
      <el-button type="primary" @click="showCreate = true">创建Topic</el-button>
    </div>
    <div style="margin-top:16px">
      <el-select v-model="clusterId" placeholder="选择集群" @change="fetchTopics" style="width:240px">
        <el-option v-for="c in clusters" :key="c.id" :label="c.name" :value="c.id" />
      </el-select>
      <el-input v-model="keyword" placeholder="搜索Topic" style="width:240px;margin-left:12px" @change="fetchTopics" clearable />
    </div>
    <el-table :data="topics" style="margin-top:16px" v-loading="loading">
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="partitions" label="分区数" />
      <el-table-column prop="replication_factor" label="副本数" />
      <el-table-column prop="created_at" label="创建时间" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="viewDetail(row)">详情</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showCreate" title="创建Topic" width="500px">
      <el-form :model="form">
        <el-form-item label="集群">
          <el-select v-model="form.cluster_id" placeholder="选择集群" style="width:100%">
            <el-option v-for="c in clusters" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="Topic名称">
          <el-input v-model="form.name" placeholder="如：order-events" />
        </el-form-item>
        <el-form-item label="分区数">
          <el-input-number v-model="form.partitions" :min="1" :max="100" />
        </el-form-item>
        <el-form-item label="副本数">
          <el-input-number v-model="form.replication_factor" :min="1" :max="5" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreate = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '../api/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const topics = ref([])
const clusters = ref([])
const clusterId = ref(null)
const keyword = ref('')
const loading = ref(false)
const showCreate = ref(false)
const form = ref({ cluster_id: null, name: '', partitions: 3, replication_factor: 1 })

async function fetchClusters() {
  const res = await request.get('/clusters')
  clusters.value = res.data || []
}

async function fetchTopics() {
  if (!clusterId.value) return
  loading.value = true
  try {
    const res = await request.get('/topics', { params: { cluster_id: clusterId.value, keyword: keyword.value } })
    topics.value = res.data || []
  } finally { loading.value = false }
}

async function handleCreate() {
  await request.post('/topics', form.value)
  ElMessage.success('Topic创建成功')
  showCreate.value = false
  fetchTopics()
}

function viewDetail(row) {
  ElMessage.info(`查看Topic: ${row.name}`)
}

async function handleDelete(row) {
  await ElMessageBox.confirm(`确认删除Topic ${row.name}？`, '警告', { type: 'warning' })
  await request.delete(`/topics/${row.name}`, { params: { cluster_id: clusterId.value } })
  ElMessage.success('删除成功')
  fetchTopics()
}

onMounted(fetchClusters)
</script>

<style scoped>
.page-header { display: flex; align-items: center; justify-content: space-between; }
</style>