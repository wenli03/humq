<template>
  <div>
    <h3>运维工具</h3>
    <el-row :gutter="20" style="margin-top:20px">
      <el-col :span="12">
        <el-card>
          <template #header>分区重分配</template>
          <el-form>
            <el-form-item label="集群">
              <el-select v-model="opsClusterId" placeholder="选择集群" style="width:100%">
                <el-option v-for="c in clusters" :key="c.id" :label="c.name" :value="c.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="Topic (JSON数组)">
              <el-input v-model="reassignTopics" type="textarea" :rows="3" placeholder='["topic1", "topic2"]' />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleReassign">执行重分配</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>操作日志</template>
          <el-empty description="暂无操作记录" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '../api/request'
import { ElMessage } from 'element-plus'

const clusters = ref([])
const opsClusterId = ref(null)
const reassignTopics = ref('')

async function fetchClusters() {
  const res = await request.get('/clusters')
  clusters.value = res.data || []
}

async function handleReassign() {
  ElMessage.info('分区重分配功能已触发')
}

onMounted(fetchClusters)
</script>