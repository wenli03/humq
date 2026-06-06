<template>
  <div>
    <div class="page-header">
      <h3>集群管理</h3>
      <el-button type="primary" @click="showCreate = true">注册集群</el-button>
    </div>
    <el-table :data="clusters" style="margin-top:20px" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="bootstrap_servers" label="连接地址" />
      <el-table-column prop="status" label="状态">
        <template #default="{ row }">
          <el-tag :type="row.status === 'online' ? 'success' : 'danger'">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" />
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button size="small" @click="viewInfo(row)">详情</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showCreate" title="注册集群" width="500px">
      <el-form :model="form">
        <el-form-item label="集群名称">
          <el-input v-model="form.name" placeholder="如：生产集群" />
        </el-form-item>
        <el-form-item label="连接地址">
          <el-input v-model="form.bootstrap_servers" placeholder="如：localhost:9092" />
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

const clusters = ref([])
const loading = ref(false)
const showCreate = ref(false)
const form = ref({ name: '', bootstrap_servers: '' })

async function fetchClusters() {
  loading.value = true
  try {
    const res = await request.get('/clusters')
    clusters.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  await request.post('/clusters', form.value)
  ElMessage.success('集群注册成功')
  showCreate.value = false
  form.value = { name: '', bootstrap_servers: '' }
  fetchClusters()
}

function viewInfo(row) {
  ElMessage.info(`集群 ${row.name} 详情`)
}

async function handleDelete(id) {
  await ElMessageBox.confirm('确认删除该集群？', '提示', { type: 'warning' })
  await request.delete(`/clusters/${id}`)
  ElMessage.success('删除成功')
  fetchClusters()
}

onMounted(fetchClusters)
</script>

<style scoped>
.page-header { display: flex; align-items: center; justify-content: space-between; }
</style>