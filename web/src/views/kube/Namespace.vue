<template>
  <div class="w-full flex justify-center">
    <div class="flex flex-col w-full h-full mx-4rem my-2rem space-y-1rem">

      <div class="flex overflow-hidden rounded-lg shadow-lg">
        <div class="flex w-full h-5rem bg-white items-center">
          <BookmarkThree class="ml-1rem" theme="filled" size="42" fill="#94A3B8" />
          <span class="m-0.75rem text-2xl font-600">Namespaces</span>
        </div>
      </div>

      <el-card class="w-full h-max">
        <template #header>
          <div class="flex justify-between space-x-10rem">
            <el-input v-model="search" placeholder="Type to search">
              <template #prefix>
                <el-icon>
                  <Search />
                </el-icon>
              </template>
            </el-input>
            <el-button type="primary" round plain :icon="BookmarkThree" @Click="showCreate = true">Create</el-button>
          </div>
          <el-dialog v-model="showCreate" center title="Create Namespace" width="33%">
            <el-select class="w-full" v-model="newNamespace.group" filterable placeholder="please select related group">
              <el-option v-for="(g, index) in groups" :label="getNsLabel(g)" :value="index" />
            </el-select>
            <template #footer>
              <span class="dialog-footer">
                <el-button type="primary" @click="createNamespace">Confirm</el-button>
                <el-button @click="showCreate = false">Cancel</el-button>
              </span>
            </template>
          </el-dialog>
        </template>
        <el-table :data="filter" class="w-full max-h-full">
          <el-table-column prop="metadata.name" label="Name" />
          <el-table-column prop="metadata.labels['weave.io/group']" label="Group" />
          <el-table-column prop="metadata.annotations['weave.io/describe']" label="Describe" />
          <el-table-column prop="metadata.creationTimestamp" label="CreateAt" min-width="120px" />
          <el-table-column label="Operation">
            <template #default="scope">
              <el-popover :visible="showDelete == scope.$index" placement="top" :width="160">
                <template #reference>
                  <el-button size="small" type="danger" @click="showDelete = scope.$index" :icon="Delete" circle
                    class="wl-1rem" />
                </template>
                <p>Are you sure to delete this namespace?</p>
                <div class="my-0.5rem">
                  <el-button size="small" text @click="showDelete = -1">cancel</el-button>
                  <el-button size="small" type="danger" @click="deleteNamespace(scope.row)">confirm</el-button>
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
import { Delete, Search, BookmarkThree } from '@icon-park/vue-next';
import { ref, onMounted, computed } from 'vue';
import { ElMessage } from "element-plus";
import request from '@/axios';
import { getUser } from '@/utils';

const user = getUser();
const groups = ref([]);
const namespaces = ref([])

const showCreate = ref(false);
const showDelete = ref(-1);
const newNamespace = ref({});

const search = ref('');
const filter = computed(() =>
  namespaces.value.filter(
    (data) =>
      !search.value ||
      data.metadata.name.toLowerCase().includes(search.value.toLowerCase())
  )
)

onMounted(
  () => {
    let id = user.id;
    if (!id) {
      return
    }

    let url = `/api/v1/users/${id}/groups`;
    if (user.inRoot) {
      url = `/api/v1/groups`
    }

    request.get(url).then((response) => {
      groups.value = response.data.data;
      if (groups.value.length == 0) {
        ElMessage.failed("No group find, please create it at first");
      }
    })

    request.get(`/api/v1/namespaces`).then((response) => {
      namespaces.value = response.data.data.items;
    })
  }
)

const getNsLabel = (group) => {
  return `${group.name}(${group.describe})`
};

const createNamespace = () => {
  let group = groups.value[newNamespace.value.group]
  if (!group) {
    ElMessage.failed("Group must selected");
    return
  }
  console.log(group.name)
  request.post("/api/v1/namespaces", {
    apiVersion: "v1",
    kind: "Namespace",
    metadata: {
      name: group.name,
      annotations: {
        "weave.io/platform": "true",
        "weave.io/describe": group.describe,
        "weave.io/group": group.name,
        "weave.io/group_id": `${group.id}`
      },
      labels: {
        "weave.io/platform": "true",
        "weave.io/group": group.name,
      }
    }
  }).then((response) => {
    ElMessage.success("Create success");
    namespaces.value.push(response.data.data);
    showCreate.value = false;
  })
}

const deleteNamespace = (row) => {
  request.delete("/api/v1/namespaces/" + row.metadata.name).then(() => {
    ElMessage.success("Delete success");
    const index = namespaces.value.findIndex(v => v.metadata.name === row.metadata.name);
    namespaces.value.splice(index, 1);
    showDelete.value = -1;
  })
};

</script>
