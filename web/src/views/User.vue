<template>
  <div class="w-full h-full flex justify-center">
    <el-card class="mx-4rem my-2rem w-4/5 h-max" shadow="never">
      <template #header>
        <div class="flex justify-between">
          <span>Users</span>
        </div>
      </template>
      <el-table :data="users" height="360" class="w-full max-h-full">
        <el-table-column prop="name" label="Name">
          <template #default="scope">
            <router-link :to="getUserUrl(scope.row.id)">
              <el-link type="primary">{{ scope.row.name }}</el-link>
            </router-link>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="Email" />
        <el-table-column prop="createAt" label="CreateAt" min-width="120px" />
        <el-table-column label="Operation" min-width="120px">
          <template #default="scope">
            <el-dialog v-model="showUpdate" center title="Update User" width="33%">
              <el-form ref="updateFormRef" :model="updatedUser" label-position="left" label-width="auto">
                <el-form-item label="Name" prop="name">
                  <el-input v-model="updatedUser.name" disabled />
                </el-form-item>
                <el-form-item label="Email" prop="email" required>
                  <el-input v-model="updatedUser.email" placeholder="User email" />
                </el-form-item>
              </el-form>
              <template #footer>
                <span class="dialog-footer">
                  <el-button type="primary" @click="updateUser(scope.row)">Confirm</el-button>
                  <el-button @click="showUpdate = false">Cancel</el-button>
                </span>
              </template>
            </el-dialog>
            <el-button size="small" circle @click="editUser(scope.row)" :icon="Edit"></el-button>

            <el-popover :visible="showDelete == scope.$index" placement="top" :width="160">
              <template #reference>
                <el-button size="small" type="danger" @click="showDelete = scope.$index" :icon="Delete" circle
                  class="wl-1rem" />
              </template>
              <p>Are you sure to delete this user?</p>
              <div class="my-0.5rem">
                <el-button size="small" type="text" @click="showDelete = -1">cancel</el-button>
                <el-button size="small" type="danger" @click="deleteUser(scope.row)">confirm</el-button>
              </div>
            </el-popover>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<style scoped>
</style>

<script setup>
import { Edit, Delete } from '@icon-park/vue-next';
import { ref, unref, onMounted } from 'vue';
import { ElMessage } from "element-plus";
import request from '@/axios'

const users = ref([]);
const showCreate = ref(false);
const showUpdate = ref(false);
const showDelete = ref(-1);
const newUser = ref({
  name: '',
  image: '',
  cmd: []
});

const updatedUser = ref({});
const createFormRef = ref();
const updateFormRef = ref();

onMounted(
  () => {
    request.get("/api/v1/users").then((response) => {
      users.value = response.data.data;
    })
  }
)

const getUserUrl = (id) => {
  return `/users/${id}`
};

const createUser = () => {
  const form = unref(createFormRef)
  if (!form) {
    return
  }

  form.validate((valid) => {
    if (valid) {
      request.post("/api/v1/users", {
        name: newUser.value.name,
        email: newUser.value.email,
      }).then((response) => {
        ElMessage.success("Create success");
        users.value.push(response.data.data);
        showCreate.value = false;
      })
    } else {
      ElMessage.error("Input invalid");
    }
  });
};

const editUser = (row) => {
  updatedUser.value = Object.assign({}, row);
  showUpdate.value = true;
}

const updateUser = (row) => {
  const form = unref(updateFormRef);
  if (!form) {
    return
  }

  form.validate((valid) => {
    if (valid) {
      request.put("/api/v1/users/" + row.id, updatedUser.value).then((response) => {
        ElMessage.success("Update success");
        const index = users.value.findIndex(v => v.id === row.id);
        users.value[index] = updatedUser.value;
        showUpdate.value = false;
      })
    } else {
      ElMessage.error("Input invalid");
    }
  });
};

const deleteUser = (row) => {
  request.delete("/api/v1/users/" + row.id).then(() => {
    ElMessage.success("Delete success");
    const index = users.value.findIndex(v => v.id === row.id);
    users.value.splice(index, 1);
    showDelete.value = -1;
  })
};

</script>
