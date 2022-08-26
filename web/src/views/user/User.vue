<template>
  <div class="w-full flex justify-center">
    <div class="flex flex-col w-full h-full px-[4rem] py-[2rem] space-y-[1rem]">
      <el-dialog v-model="showUpdate" top="5vh" title="Update User" width="50%">
        <el-form ref="updateFormRef" :model="updatedUser" label-position="top" label-width="auto">
          <el-form-item label="Name" prop="name">
            <el-input v-model="updatedUser.name" disabled />
          </el-form-item>
          <el-form-item label="Email" prop="email" required>
            <el-input v-model="updatedUser.email" />
            <span class="text-gray-400">The user email address</span>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button type="primary" @click="updateUser">Confirm</el-button>
            <el-button @click="showUpdate = false">Cancel</el-button>
          </span>
        </template>
      </el-dialog>

      <div class="flex overflow-hidden rounded-md shadow-md border">
        <div class="flex w-full h-[5rem] items-center">
          <User class="ml-[1rem]" theme="filled" size="42" fill="#94A3B8" />
          <span class="m-[0.75rem] text-2xl font-600">Users</span>
        </div>
      </div>

      <el-card class="h-max">
        <template #header>
          <div class="flex">
            <el-input v-model="search" placeholder="Type to search">
              <template #prefix>
                <el-icon>
                  <Search />
                </el-icon>
              </template>
            </el-input>
          </div>
        </template>
        <el-table :data="filter" class="w-full max-h-full">
          <el-table-column prop="name" label="Name">
            <template #default="scope">
              <router-link :to="getUserUrl(scope.row.id)">
                <el-link type="primary">{{ scope.row.name }}</el-link>
              </router-link>
            </template>
          </el-table-column>
          <el-table-column prop="email" label="Email" />
          <el-table-column prop="createdAt" label="CreateAt" min-width="120px" />
          <el-table-column label="Operation" min-width="120px">
            <template #default="scope">
              <el-button size="small" circle @click="editUser(scope.row)" :icon="Edit"></el-button>
              <el-popover :visible="showDelete == scope.$index" placement="top" :width="180">
                <template #reference>
                  <el-button size="small" type="danger" @click="showDelete = scope.$index" :icon="Delete" circle
                    class="wl-[1rem]" />
                </template>
                <p>Are you sure to delete this user?</p>
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
import { Edit, Delete, Search, User } from '@icon-park/vue-next';
import { ref, unref, onMounted, computed } from 'vue';
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

const search = ref('');
const filter = computed(() =>
  users.value.filter(
    (data) =>
      !search.value ||
      data.name.toLowerCase().includes(search.value.toLowerCase())
  )
)


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

const updateUser = () => {
  const form = unref(updateFormRef);
  if (!form) {
    return
  }

  form.validate((valid) => {
    if (valid) {
      request.put("/api/v1/users/" + updatedUser.value.id, updatedUser.value).then((response) => {
        ElMessage.success("Update success");
        const index = users.value.findIndex(v => v.id === updatedUser.value.id);
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
