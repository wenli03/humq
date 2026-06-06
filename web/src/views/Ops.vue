<template>
  <div>
    <h3>运维工具</h3>
    <el-row :gutter="20" style="margin-top:20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>消息积压</span>
            <el-button text type="primary" @click="fetchBacklog" style="float:right">刷新</el-button>
          </template>
          <el-table :data="backlogList" v-loading="backlogLoading" size="small">
            <el-table-column prop="consumer_group" label="消费组" width="120" />
            <el-table-column prop="topic" label="Topic" width="120" />
            <el-table-column prop="lag" label="积压量" width="100">
              <template #default="{ row }">
                <span :style="{ color: row.severity === 'critical' ? '#F56C6C' : '#E6A23C', fontWeight: 'bold' }">
                  {{ row.lag.toLocaleString() }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="操作">
              <template #default="{ row }">
                <el-button size="small" type="warning" @click="resolve(row, 'scale')">扩容</el-button>
                <el-button size="small" type="danger" @click="resolve(row, 'skip')">跳过</el-button>
                <el-button size="small" @click="resolve(row, 'throttle')">暂停生产</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <el-card style="margin-top:20px">
          <template #header>处理方案说明</template>
          <div style="font-size:13px;color:#666;line-height:2">
            <p><b>扩容消费者</b> — 增加消费者实例数至与分区数一致，提升并行消费能力</p>
            <p><b>跳过旧消息</b> — 将消费位置重置到最新，丢弃已无价值的积压消息</p>
            <p><b>暂停生产者</b> — 暂停低优先级数据源，为消费者赢得消化时间</p>
          </div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header>操作日志</template>
          <div v-if="logs.length === 0" style="text-align:center;padding:30px;color:#999">暂无记录</div>
          <el-timeline v-else>
            <el-timeline-item v-for="l in logs" :key="l.time" :timestamp="l.time" :color="l.action === 'skip' ? '#F56C6C' : '#409EFF'">
              {{ l.msg }}
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '../api/request'
import { ElMessage } from 'element-plus'

const backlogList = ref([])
const backlogLoading = ref(false)
const logs = ref([])

async function fetchBacklog() {
  backlogLoading.value = true
  try {
    const res = await request.get('/ops/backlog')
    backlogList.value = res.data || []
  } finally { backlogLoading.value = false }
}

async function resolve(row, action) {
  const names = { scale: '扩容消费者', skip: '跳过消息', throttle: '暂停生产者' }
  try {
    const res = await request.post('/ops/resolve-backlog', { consumer_group: row.consumer_group, action })
    ElMessage.success(res.data?.msg || '操作成功')
    logs.value.unshift({ time: new Date().toLocaleTimeString('zh-CN'), action, msg: `${names[action]}: ${row.consumer_group} (${row.topic})` })
    fetchBacklog()
  } catch { ElMessage.error('操作失败') }
}

onMounted(() => { fetchBacklog(); request.get('/ops/logs').then(r => logs.value = r.data || []) })
</script>
