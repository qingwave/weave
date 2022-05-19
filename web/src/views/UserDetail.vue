<template>
  <div class="w-full h-full flex justify-center">
    <el-card class="mx-4rem my-2rem w-4/5 h-max" shadow="never">
      <el-descriptions title="User Info" :model="user" size="large" :column=1>
        <el-descriptions-item label="Name">{{ user.name }}</el-descriptions-item>
        <el-descriptions-item label="Email">{{ user.email }}</el-descriptions-item>
        <el-descriptions-item label="CreateAt">{{ user.createAt }}</el-descriptions-item>
        <div v-for="auth in user.authInfos">
          <el-descriptions-item :label="auth.authType">
            {{ auth.authId }}
          </el-descriptions-item>
        </div>

      </el-descriptions>

      <div class="mt-1rem">
        <span class="my-0.5rem">
          <el-button @Click="showUpdate = true">Update</el-button>
        </span>

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
              <el-button type="primary" @click="updateUser">Confirm</el-button>
              <el-button @click="showUpdate = false">Cancel</el-button>
            </span>
          </template>
        </el-dialog>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
</style>

<script setup>
import { ref, unref, onMounted } from 'vue';
import { ElMessage } from "element-plus";
import request from '@/axios';
import { useRoute } from 'vue-router';

const route = useRoute();

const id = route.params.id;

const user = ref({
  name: '',
  describe: '',
});

const showUpdate = ref(false);
const updatedUser = ref({
  name: '',
  describe: '',
});

const updateFormRef = ref();

onMounted(
  () => {
    if (!id) {
      return
    }
    request.get(`/api/v1/users/${id}`).then((response) => {
      user.value = response.data.data;
      updatedUser.value = response.data.data;
    })
  }
)

const updateUser = () => {
  const form = unref(updateFormRef);
  if (!form) {
    return
  }

  form.validate((valid) => {
    if (valid) {
      request.put("/api/v1/users/" + updatedUser.value.id, updatedUser.value).then((response) => {
        ElMessage.success("Update success");
        users.value[index] = updatedUser.value;
        showUpdate.value = false;
      })
    } else {
      ElMessage.error("Input invalid");
    }
  });
};

</script>
