<template>
  <div class="w-full flex justify-center">
    <div class="flex flex-col w-full h-full px-[4rem] py-[2rem] space-y-[1rem]">
      <el-dialog v-model="showCreate" top="5vh" title="Create Group" width="50%">
        <el-form ref="createFormRef" :model="newGroup" label-position="top" label-width="auto">
          <el-form-item label="Name" prop="name" required>
            <el-input v-model="newGroup.name" />
            <span class="text-gray-400">The group name</span>
          </el-form-item>
          <el-form-item label="Describe" prop="describe" required>
            <el-input v-model="newGroup.describe" type="textarea" />
            <span class="text-gray-400">The group describe information</span>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button type="primary" @click="createGroup">Confirm</el-button>
            <el-button @click="showCreate = false">Cancel</el-button>
          </span>
        </template>
      </el-dialog>

      <el-dialog v-model="showUpdate" top="5vh" title="Update Group" width="50%">
        <el-form ref="updateFormRef" :model="updatedGroup" label-position="top" label-width="auto">
          <el-form-item label="Name" prop="name">
            <el-input v-model="updatedGroup.name" disabled />
          </el-form-item>
          <el-form-item label="Describe" prop="describe" required>
            <el-input v-model="updatedGroup.describe" placeholder="Group describe" />
            <span class="text-gray-400">The group describe information</span>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button type="primary" @click="updateGroup">Confirm</el-button>
            <el-button @click="showUpdate = false">Cancel</el-button>
          </span>
        </template>
      </el-dialog>

      <div class="flex overflow-hidden rounded-md shadow-md border">
        <div class="flex w-full h-[5rem] items-center">
          <Peoples class="ml-[1rem]" theme="filled" size="42" fill="#94A3B8" />
          <span class="m-[0.75rem] text-2xl font-600">Groups</span>
        </div>
      </div>

      <el-card class="w-full h-max">
        <template #header>
          <div class="flex justify-between space-x-[2rem]">
            <el-input v-model="search" placeholder="Type to search">
              <template #prefix>
                <el-icon>
                  <Search />
                </el-icon>
              </template>
            </el-input>
            <el-button type="primary" plain :icon="Peoples" @Click="showCreate = true">Create</el-button>
          </div>
        </template>
        <el-table :data="filter" class="w-full max-h-full">
          <el-table-column prop="name" label="Name">
            <template #default="scope">
              <router-link :to="getGroupUrl(scope.row.id)">
                <el-link type="primary">{{ scope.row.name }}</el-link>
              </router-link>
            </template>
          </el-table-column>
          <el-table-column prop="describe" label="Describe" />
          <el-table-column prop="creatorId" label="Creator" />
          <el-table-column prop="createdAt" label="CreateAt" min-width="120px" />
          <el-table-column label="Operation" min-width="120px">
            <template #default="scope">
              <el-button size="small" circle @click="editGroup(scope.row)" :icon="Edit" />

              <el-popover :visible="showDelete == scope.$index" placement="top" :width="180">
                <template #reference>
                  <el-button size="small" type="danger" @click="showDelete = scope.$index" :icon="Delete" circle
                    class="wl-[1rem]" />
                </template>
                <p>Are you sure to delete this group?</p>
                <div class="my-[0.5rem]">
                  <el-button size="small" text @click="showDelete = -1">cancel</el-button>
                  <el-button size="small" type="danger" @click="deleteGroup(scope.row)">confirm</el-button>
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
import { Edit, Delete, Peoples, Search } from '@icon-park/vue-next';
import { ref, unref, onMounted, computed } from 'vue';
import { ElMessage } from "element-plus";
import request from '@/axios';
import { getUser } from '@/utils';

const props = defineProps({
  listAll: Boolean,
})

const user = getUser();
const groups = ref([]);

const showCreate = ref(false);
const showUpdate = ref(false);
const showDelete = ref(-1);
const newGroup = ref({
  name: '',
  describe: '',
});

const updatedGroup = ref({});
const createFormRef = ref();
const updateFormRef = ref();

const search = ref('');
const filter = computed(() =>
  groups.value.filter(
    (data) =>
      !search.value ||
      data.name.toLowerCase().includes(search.value.toLowerCase())
  )
)

onMounted(
  () => {
    let id = user.id;
    if (!id) {
      return
    }

    let url = `/api/v1/users/${id}/groups`;
    if (props.listAll) {
      url = `/api/v1/groups`
    }

    request.get(url).then((response) => {
      groups.value = response.data.data;
    })
  }
)

const getGroupUrl = (id) => {
  return `/groups/${id}`
};

const createGroup = () => {
  const form = unref(createFormRef)
  if (!form) {
    return
  }

  form.validate((valid) => {
    if (valid) {
      request.post("/api/v1/groups", {
        name: newGroup.value.name,
        describe: newGroup.value.describe,
      }).then((response) => {
        ElMessage.success("Create success");
        groups.value.push(response.data.data);
        showCreate.value = false;
      })
    } else {
      ElMessage.error("Input invalid");
    }
  });
};

const editGroup = (row) => {
  updatedGroup.value = Object.assign({}, row);
  showUpdate.value = true;
}

const updateGroup = () => {
  const form = unref(updateFormRef);
  if (!form) {
    return
  }

  form.validate((valid) => {
    if (valid) {
      request.put("/api/v1/groups/" + updatedGroup.value.id, updatedGroup.value).then((response) => {
        ElMessage.success("Update success");
        const index = groups.value.findIndex(v => v.id === updatedGroup.value.id);
        groups.value[index] = updatedGroup.value;
        showUpdate.value = false;
      })
    } else {
      ElMessage.error("Input invalid");
    }
  });
};

const deleteGroup = (row) => {
  request.delete("/api/v1/groups/" + row.id).then(() => {
    ElMessage.success("Delete success");
    const index = groups.value.findIndex(v => v.id === row.id);
    groups.value.splice(index, 1);
    showDelete.value = -1;
  })
};

</script>
