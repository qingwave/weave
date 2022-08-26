<template>
  <div class="flex w-full justify-center">
    <div class="flex flex-col w-full px-[4rem] py-[2rem] h-max justify-center">
      <el-card class="w-full">
      <template #header>
        <div class="flex w-full items-center">
          <Peoples class="ml-[1rem]" theme="filled" size="42" fill="#94A3B8" />
          <span class="m-[0.75rem] text-2xl font-600">Group Info</span>
        </div>
      </template>
      <el-descriptions :model="group" size="large" :column=1>
        <el-descriptions-item label="Name">{{ group.name }}</el-descriptions-item>
        <el-descriptions-item label="Describe">{{ group.describe }}</el-descriptions-item>
        <el-descriptions-item label="CreatorId">{{ group.creatorId }}</el-descriptions-item>
        <el-descriptions-item label="CreateAt">{{ group.createdAt }}</el-descriptions-item>
        <el-descriptions-item label="Remarks" v-show="tag">
          <el-tag size="small">{{ tag }}</el-tag>
        </el-descriptions-item>
      </el-descriptions>

    </el-card>

    <el-card class="mt-[1rem]">
        <div class="flex my-[0.5rem] justify-between">
          <span class="mr-[10rem] text-bold">Group Users</span>
          <el-button plain :icon="User" @Click="showCreate = true">AddUser</el-button>
        </div>
        
        <el-dialog v-model="showCreate" top="5vh" title="Add User" width="50%">
          <el-form ref="createFormRef" :model="newUser" label-position="top" label-width="auto">
            <el-form-item label="Name" prop="name" required>
              <el-input v-model="newUser.name" />
              <span class="text-gray-400">The user name</span>
            </el-form-item>
            <el-form-item label="Role" prop="role" required>
              <el-select  class="w-full" v-model="newUser.role">
                <el-option v-for="role in roleOptions" :label="role.label" :value="role.value" />
              </el-select>
              <span class="text-gray-400">The user role in this group</span>
            </el-form-item>
          </el-form>
          <template #footer>
            <span class="dialog-footer">
              <el-button type="primary" @click="addUser">Confirm</el-button>
              <el-button @click="showCreate = false">Cancel</el-button>
            </span>
          </template>
        </el-dialog>

        <el-table :data="users">
          <el-table-column prop="name" label="Name" />
          <el-table-column prop="role" label="Role" />
          <el-table-column label="Operation">
            <template #default="scope">
              <el-popover :visible="showDelete == scope.$index" placement="top" :width="180">
                <template #reference>
                  <el-button size="small" type="danger" @click="showDelete = scope.$index" :icon="Delete" circle
                    class="wl-[1rem]" />
                </template>
                <p>Are you sure to delete this user from group?</p>
                <div class="my-[0.5rem]">
                  <el-button size="small" text @click="showDelete = -1">cancel</el-button>
                  <el-button size="small" type="danger" @click="deleteUser(scope.row)">confirm</el-button>
                </div>
              </el-popover>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>
  </div>
</template>

<style scoped>
</style>

<script setup>
import { Delete, User, Peoples } from '@icon-park/vue-next';
import { ref, unref, onMounted } from 'vue';
import { ElMessage } from "element-plus";
import request from '@/axios';
import { useRoute } from 'vue-router';

const route = useRoute();

const id = route.params.id;

const roleOptions = ref([
  {
    value: 'admin',
    label: 'Admin',
  },
  {
    value: 'edit',
    label: 'Edit',
  },
  {
    value: 'view',
    label: 'View',
  },
])

const users = ref();
const group = ref({
  name: '',
  describe: '',
});
const tag = ref();

const showCreate = ref(false);
const showDelete = ref(-1);
const newUser = ref({
  name: '',
  describe: '',
});

const createFormRef = ref();

onMounted(
  () => {
    if (!id) {
      return
    }
    request.get(`/api/v1/groups/${id}`).then((response) => {
      group.value = response.data.data;
      if (group.value.creatorId > 0) {
        tag.value = "Customer"
      } else {
        tag.value = "System"
      }
    })
    request.get(`/api/v1/groups/${id}/users`).then((response) => {
      users.value = response.data.data;
    })
  }
)

const addUser = () => {
  const form = unref(createFormRef)
  if (!form) {
    return
  }

  form.validate((valid) => {
    if (valid) {
      let user = newUser.value;
      request.post(`/api/v1/groups/${id}/users`, user).then(() => {
        ElMessage.success("Create success");
        users.value.push(user);
        showCreate.value = false;
      })
    } else {
      ElMessage.error("Input invalid");
    }
  });
};

const deleteUser = (row) => {
  request.delete(`/api/v1/groups/${id}/users?name=${row.name}&role=${row.role}`).then(() => {
    ElMessage.success("Delete success");
    const index = users.value.findIndex(v => v.name === row.name && v.role === row.role);
    users.value.splice(index, 1);
    showDelete.value = -1;
  })
};

</script>
