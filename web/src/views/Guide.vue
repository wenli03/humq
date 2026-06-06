<template>
  <div>
    <h3>使用指引</h3>
    <el-row :gutter="20" style="margin-top:20px">
      <el-col :span="6" v-for="item in guides" :key="item.title">
        <el-card shadow="hover" @click="activeGuide = item.key" class="guide-card">
          <div class="guide-item">
            <el-icon :size="32" color="#409EFF"><component :is="item.icon" /></el-icon>
            <h4>{{ item.title }}</h4>
            <p>{{ item.desc }}</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-card style="margin-top:20px" v-if="activeGuide">
      <template #header>{{ guideTitle }}</template>
      <div class="guide-content">
        <template v-if="activeGuide === 'quick'">
          <h4>快速开始</h4>
          <p>1. 安装 Docker 和 Docker Compose</p>
          <p>2. 克隆 HU MQ 项目：<code>git clone https://github.com/wenli03/humq.git</code></p>
          <p>3. 进入项目目录，启动服务：<code>docker compose -f deploy/docker-compose.yml up -d</code></p>
          <p>4. 访问监控后台：<code>http://localhost:8080</code></p>
          <p>5. 默认账号：admin / admin</p>
        </template>
        <template v-if="activeGuide === 'sdk'">
          <h4>SDK 接入示例</h4>
          <p><strong>Go：</strong>使用 Sarama 客户端连接 HU MQ（兼容 Kafka 协议）</p>
          <p><strong>Java：</strong>使用 Kafka Client 连接，bootstrap.servers 配置为 HU MQ 地址</p>
          <p><strong>Python：</strong>使用 kafka-python 或 confluent-kafka 连接</p>
          <p>所有标准 Kafka 客户端均可直接接入 HU MQ。</p>
        </template>
        <template v-if="activeGuide === 'best'">
          <h4>最佳实践</h4>
          <p><strong>Topic命名规范：</strong>使用 {业务线}.{子系统}.{事件类型} 格式</p>
          <p><strong>分区规划：</strong>根据预估吞吐量设置分区数，建议 TPS <= 10000 时单分区即可</p>
          <p><strong>告警阈值建议：</strong>Lag > 10000 告警，Broker 离线立即告警</p>
          <p><strong>副本数：</strong>生产环境建议 3 副本</p>
        </template>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Guide, Connection, Star } from '@element-plus/icons-vue'

const activeGuide = ref('quick')
const guides = [
  { key: 'quick', title: '快速开始', desc: '5分钟部署HU MQ', icon: 'Guide' },
  { key: 'sdk', title: 'SDK指南', desc: '多语言SDK接入', icon: 'Connection' },
  { key: 'best', title: '最佳实践', desc: 'Topic规划与调优', icon: 'Star' }
]
const guideTitle = computed(() => {
  const map = { quick: '快速开始', sdk: 'SDK指南', best: '最佳实践' }
  return map[activeGuide.value] || ''
})
</script>

<style scoped>
.guide-card { cursor: pointer; }
.guide-card:hover { box-shadow: 0 4px 16px rgba(0,0,0,0.12); }
.guide-item { text-align: center; padding: 12px 0; }
.guide-item h4 { margin: 12px 0 8px; }
.guide-item p { color: #999; font-size: 13px; margin: 0; }
.guide-content p { margin: 8px 0; line-height: 1.8; }
.guide-content code { background: #f5f5f5; padding: 2px 6px; border-radius: 3px; }
</style>