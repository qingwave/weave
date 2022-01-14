<template>
  <div class="container h-full flex justify-center">
    <el-card class="mx-4rem my-2rem w-4/5" shadow="never">
      <template #header>
        <div class="flex justify-between">
          <span>Users</span>
        </div>
      </template>
      <el-table :data="users" height="320" class="w-full max-h-full">
        <el-table-column prop="name" label="Name" />
        <el-table-column prop="email" label="Email" />
        <el-table-column prop="create_time" label="CreateAt" min-width="120px" />
        <el-table-column prop="Operation" label="Operation" min-width="120px">
          <template #default="scope">
            <el-button size="small" circle @click="editUser(scope.$index)" :icon="Edit">
              <el-dialog v-model="showUpdate" center title="Update User" width="33%">
                <el-form
                  ref="updateFormRef"
                  :model="updatedUser"
                  label-position="left"
                  label-width="auto"
                >
                  <el-form-item label="Name" prop="name">
                    <el-input v-model="updatedUser.name" disabled />
                  </el-form-item>
                  <el-form-item label="Email" prop="email" required>
                    <el-input v-model="updatedUser.image" placeholder="User email" />
                  </el-form-item>
                </el-form>
                <template #footer>
                  <span class="dialog-footer">
                    <el-button type="primary" @click="updateUser(scope.$index)">Confirm</el-button>
                    <el-button @click="showUpdate = false">Cancel</el-button>
                  </span>
                </template>
              </el-dialog>
            </el-button>

            <el-popover :visible="showDelete == scope.$index" placement="top" :width="160">
              <template #reference>
                <el-button
                  size="small"
                  type="danger"
                  @click="showDelete = scope.$index"
                  :icon="Delete"
                  circle
                  class="wl-1rem"
                />
              </template>
              <p>Are you sure to delete this user?</p>
              <div class="my-0.5rem">
                <el-button size="small" type="text" @click="showDelete = -1">cancel</el-button>
                <el-button size="small" type="danger" @click="deleteUser(scope.$index)">confirm</el-button>
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

const users = ref([
  {
    name: "fakeuser",
    email: "fakeuser@email.com"
  }
]);
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
    }).catch((error) => {
      if (error.response) {
        console.log("list failed: " + error.response.data.msg);
      }
      console.log(error)
    })
  }
)

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
      }).catch((error) => {
        if (error.response) {
          ElMessage.error("Create failed: " + error.response.data.msg);
        } else if (error.request) {
          console.log(error.request)
        } else {
          ElMessage.error("Create failed")
        }
      })
    } else {
      ElMessage.error("Input invalid");
    }
  });
};

const editUser = (index) => {
  updatedUser.value = Object.assign({}, users.value[index]);
  showUpdate.value = true;
}

const updateUser = (index) => {
  const form = unref(updateFormRef);
  if (!form) {
    return
  }

  let user = Object.assign({}, users.value[index]);
  user.cmd = getCommand(updatedUser.value.cmd)

  form.validate((valid) => {
    if (valid) {
      request.put("/api/v1/users/" + updatedUser.value.id, user).then((response) => {
        ElMessage.success("Update success");
        users.value[index] = user;
        showUpdate.value = false;
      }).catch((error) => {
        let msg = 'Update failed'
        if (error.response) {
          msg += ": " + error.response.data.msg
        } else {
          console.log(error)
        }
        ElMessage.error(msg);
      })
    } else {
      ElMessage.error("Input invalid");
    }
  });
};

const deleteUser = (index) => {
  request.delete("/api/v1/users/" + users.value[index].id).then(() => {
    ElMessage.success("Delete success");
    users.value.splice(index, 1);
    showDelete.value = -1;
  }).catch((error) => {
    let msg = 'Delete failed'
    if (error.response) {
      msg += ": " + error.response.data.msg
    } else {
      console.log(error)
    }
    ElMessage.error(msg);
  })
};

</script>
