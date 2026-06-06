<template>
  <div>
    <h3>运维工具</h3>
    <el-row :gutter="20" style="margin-top:20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>消息积压处理</span>
            <el-button text type="primary" @click="fetchBacklog" style="float:right">刷新</el-button>
          </template>
          <el-table :data="backlogList" style="width:100%" v-loading="backlogLoading">
            <el-table-column prop="consumer_group" label="消费组" width="140" />
            <el-table-column prop="topic" label="Topic" width="140" />
            <el-table-column prop="lag" label="积压量" width="120">
              <template #default="{ row }">
                <span :style="{ color: row.severity === 'critical' ? '#F56C6C' : '#E6A23C', fontWeight: 'bold' }">
                  {{ row.lag.toLocaleString() }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="severity" label="严重程度" width="100">
              <template #default="{ row }">
                <el-tag :type="row.severity === 'critical' ? 'danger' : 'warning'">
                  {{ row.severity === 'critical' ? '严重' : '警告' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="suggestion" label="建议" min-width="160" />
            <el-table-column label="操作" width="260">
              <template #default="{ row }">
                <el-button size="small" type="warning" @click="resolveBacklog(row, 'scale')">扩容消费者</el-button>
                <el-button size="small" type="danger" @click="resolveBacklog(row, 'skip')">跳过旧消息</el-button>
                <el-button size="small" @click="resolveBacklog(row, 'throttle')">暂停生产者</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>积压处理方案说明</template>
          <div class="guide-content">
            <h4>方案1：扩容消费者</h4>
            <p>增加消费者实例数量，将单个消费组的消费者数扩展至分区数。建议消费者数 = 分区数。</p>
            <h4>方案2：增加分区数</h4>
            <p>通过增加 Topic 分区数提升并行消费能力。需在 Topic 管理页面操作。</p>
            <h4>方案3：跳过旧消息</h4>
            <p>当历史消息已无价值时，可将消费位置前移，跳过积压的旧消息。</p>
            <h4>方案4：暂停非关键生产者</h4>
            <p>临时暂停低优先级生产者，待积压消除后恢复。</p>
          </div>
        </el-card>

        <el-card style="margin-top:20px">
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
              <el-button type="primary" @click="handleRebalance">执行重分配</el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <el-card style="margin-top:20px">
          <template #header>操作日志</template>
          <div v-if="actionLogs.length === 0" style="text-align:center;padding:20px;color:#999">暂无操作记录</div>
          <el-timeline v-else>
            <el-timeline-item
              v-for="log in actionLogs"
              :key="log.time"
              :timestamp="log.time"
              :type="log.type"
            >
              {{ log.msg }}
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

const clusters = ref([])
const opsClusterId = ref(null)
const reassignTopics = ref('')

const backlogList = ref([])
const backlogLoading = ref(false)
const actionLogs = ref([])

async function fetchClusters() {
  try {
    const res = await request.get('/clusters')
    clusters.value = res.data || []
  } catch (e) { /* demo mode */ }
}

async function fetchBacklog() {
  backlogLoading.value = true
  try {
    const res = await request.get('/ops/backlog')
    backlogList.value = res.data || []
  } catch (e) {
    ElMessage.warning('获取积压数据失败')
  } finally {
    backlogLoading.value = false
  }
}

async function resolveBacklog(row, action) {
  const actionName = { scale: '扩容消费者', skip: '跳过旧消息', throttle: '暂停生产者' }
  try {
    const res = await request.post('/ops/resolve-backlog', {
      consumer_group: row.consumer_group,
      action: action
    })
    ElMessage.success(res.data?.msg || '操作成功')
    actionLogs.value.unshift({
      time: new Date().toLocaleString('zh-CN'),
      msg: `${actionName[action]}: ${row.consumer_group}`,
      type: action === 'skip' ? 'danger' : 'primary'
    })
    fetchBacklog()
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

async function handleRebalance() {
  try {
    const res = await request.post('/ops/rebalance', { topics: reassignTopics.value })
    ElMessage.success(res.data?.msg || '分区重分配已提交')
    actionLogs.value.unshift({
      time: new Date().toLocaleString('zh-CN'),
      msg: `分区重分配: ${reassignTopics.value}`,
      type: 'success'
    })
  } catch (e) {
    ElMessage.error('重分配失败')
  }
}

onMounted(() => {
  fetchClusters()
  fetchBacklog()
})
</script>

<style scoped>
.guide-content h4 { margin: 12px 0 6px; color: #333; }
.guide-content p { margin: 4px 0 12px; color: #666; font-size: 13px; line-height: 1.6; }
</style>
