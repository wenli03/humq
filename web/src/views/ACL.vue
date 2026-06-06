<template>
  <div>
    <h3>权限管理</h3>
    <el-tabs v-model="activeTab" style="margin-top:16px">
      <el-tab-pane label="用户管理" name="users">
        <div style="margin-bottom:16px">
          <el-button type="primary" @click="showCreateUser = true">添加用户</el-button>
        </div>
        <el-table :data="users" v-loading="loading">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="username" label="用户名" />
          <el-table-column prop="role" label="角色">
            <template #default="{ row }">
              <el-tag :type="row.role === 'admin' ? 'danger' : 'info'">{{ row.role === 'admin' ? '管理员' : '普通用户' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="创建时间" width="180" />
          <el-table-column label="操作" width="160">
            <template #default="{ row }">
              <el-button size="small" @click="editUser(row)">编辑</el-button>
              <el-button size="small" type="danger" @click="deleteUser(row.id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
      <el-tab-pane label="ACL规则" name="acls">
        <div style="margin-bottom:16px">
          <el-button type="primary" @click="showCreateACL = true">添加规则</el-button>
        </div>
        <el-table :data="acls" v-loading="loading">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="user_id" label="用户ID" width="100" />
          <el-table-column prop="resource_type" label="资源类型" width="120" />
          <el-table-column prop="resource_name" label="资源名称" />
          <el-table-column prop="operation" label="操作" width="100">
            <template #default="{ row }">
              <el-tag size="small">{{ row.operation }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button size="small" type="danger" @click="deleteACL(row.id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="showCreateUser" title="添加用户" width="400px">
      <el-form :model="userForm">
        <el-form-item label="用户名">
          <el-input v-model="userForm.username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="userForm.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="userForm.role" style="width:100%">
            <el-option label="管理员" value="admin" />
            <el-option label="普通用户" value="user" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateUser = false">取消</el-button>
        <el-button type="primary" @click="createUser">确认</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showCreateACL" title="添加ACL规则" width="400px">
      <el-form :model="aclForm">
        <el-form-item label="用户ID">
          <el-input-number v-model="aclForm.user_id" :min="1" style="width:100%" />
        </el-form-item>
        <el-form-item label="资源类型">
          <el-select v-model="aclForm.resource_type" style="width:100%">
            <el-option label="Topic" value="topic" />
            <el-option label="消费组" value="group" />
          </el-select>
        </el-form-item>
        <el-form-item label="资源名称">
          <el-input v-model="aclForm.resource_name" />
        </el-form-item>
        <el-form-item label="操作">
          <el-select v-model="aclForm.operation" style="width:100%">
            <el-option label="读取(read)" value="read" />
            <el-option label="写入(write)" value="write" />
            <el-option label="管理(admin)" value="admin" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateACL = false">取消</el-button>
        <el-button type="primary" @click="createACL">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '../api/request'
import { ElMessage, ElMessageBox } from 'element-plus'

const activeTab = ref('users')
const users = ref([])
const acls = ref([])
const loading = ref(false)
const showCreateUser = ref(false)
const showCreateACL = ref(false)
const userForm = ref({ username: '', password: '', role: 'user' })
const aclForm = ref({ user_id: 1, resource_type: 'topic', resource_name: '', operation: 'read' })

async function fetchData() {
  loading.value = true
  try {
    const [u, a] = await Promise.all([
      request.get('/users'),
      request.get('/acls')
    ])
    users.value = u.data || []
    acls.value = a.data || []
  } finally { loading.value = false }
}

async function createUser() {
  await request.post('/users', userForm.value)
  ElMessage.success('用户创建成功')
  showCreateUser.value = false
  fetchData()
}

function editUser(row) {
  ElMessage.info(`编辑用户: ${row.username}`)
}

async function deleteUser(id) {
  await ElMessageBox.confirm('确认删除该用户？', '提示', { type: 'warning' })
  await request.delete(`/users/${id}`)
  ElMessage.success('删除成功')
  fetchData()
}

async function createACL() {
  await request.post('/acls', aclForm.value)
  ElMessage.success('ACL规则创建成功')
  showCreateACL.value = false
  fetchData()
}

async function deleteACL(id) {
  await request.delete(`/acls/${id}`)
  ElMessage.success('删除成功')
  fetchData()
}

onMounted(fetchData)
</script>