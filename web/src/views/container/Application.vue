<template>
  <div class="w-full flex justify-center">
    <el-dialog v-model="showCreate" top="5vh" title="Create Application" width="60%">
      <el-form ref="createFormRef" :model="newApp" label-position="top" label-width="auto">
        <div class="flex flex-row w-full space-x-[2rem]">
          <el-form-item label="Name" prop="name" class="w-1/2" required>
            <el-input v-model="newApp.name" />
            <span class="text-gray-400">The application name contains only lowercase letters, numbers, and hyphens
              (-)</span>
          </el-form-item>
          <el-form-item label="Describe" prop="describe" class="w-1/2">
            <el-input v-model="newApp.describe" type="textarea" />
            <span class="text-gray-400">The application describe information</span>
          </el-form-item>
        </div>
        <el-form-item label="Image" prop="image" required>
          <el-input v-model="newApp.image" />
          <span class="text-gray-400">The application image</span>
        </el-form-item>
        <el-form-item label="Command" prop="cmd">
          <el-input v-model="newApp.cmd" type="textarea" />
          <span class="text-gray-400">The application runs command on startup</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="createApp">Confirm</el-button>
          <el-button @click="showCreate = false">Cancel</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="showUpdate" top="5vh" title="Update Application" width="60%">
      <el-form ref="updateFormRef" :model="updatedApp" label-position="top" label-width="auto">
        <div class="flex flex-row w-full space-x-[2rem]">
          <el-form-item class="w-1/2" label="Name" prop="name">
            <el-input v-model="updatedApp.name" disabled />
          </el-form-item>
          <el-form-item class="w-1/2" label="Describe" prop="describe">
            <el-input v-model="updatedApp.describe" type="textarea" />
            <span class="text-gray-400">The application describe information</span>
          </el-form-item>
        </div>
        <el-form-item label="Image" prop="image" required>
          <el-input v-model="updatedApp.image" />
          <span class="text-gray-400">The application image</span>
        </el-form-item>
        <el-form-item label="Command" prop="cmd">
          <el-input v-model="updatedApp.cmd" />
          <span class="text-gray-400">The application runs command on startup</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="updateApp()">Confirm</el-button>
          <el-button @click="showUpdate = false">Cancel</el-button>
        </span>
      </template>
    </el-dialog>

    <div class="flex flex-col w-full h-full px-[4rem] py-[2rem] space-y-[1rem]">

      <div class="flex overflow-hidden rounded-md shadow-md border">
        <div class="flex w-full h-[5rem] items-center">
          <ApplicationOne class="ml-[1rem]" theme="filled" size="42" fill="#94A3B8" />
          <span class="m-[0.75rem] text-2xl font-600">Applications</span>
        </div>
      </div>

      <el-card class="h-max">
        <template #header>
          <div class="flex justify-between space-x-[2rem]">
            <el-input v-model="search" placeholder="Type to search">
              <template #prefix>
                <el-icon>
                  <Search />
                </el-icon>
              </template>
            </el-input>
            <el-button type="primary" plain :icon="ApplicationOne" @Click="showCreate = true">Create</el-button>
          </div>
        </template>
        <el-table :data="filter" height="360" class="w-full max-h-full">
          <el-table-column prop="name" label="Name" sortable />
          <el-table-column prop="image" label="Image" />
          <el-table-column prop="status" label="Status">
            <template #default="scope">
              <el-tag :type="getAppStatusType(scope.row.status)"> {{ scope.row.status }} </el-tag>
            </template>

          </el-table-column>
          <el-table-column prop="cmd" label="Command" min-width="120px" />
          <el-table-column prop="startAt" label="StartAt" sortable min-width="120px" />
          <el-table-column prop="Operation" label="Operation" min-width="120px">
            <template #default="scope">
              <el-button size="small" type="primary" circle @click="execApp(scope.row)" :icon="Terminal" />
              <el-button class="ml-[0.5rem]" size="small" circle @click="editApp(scope.row)" :icon="Edit" />

              <el-popover :visible="showDelete == scope.$index" placement="top" :width="180">
                <template #reference>
                  <el-button size="small" type="danger" @click="showDelete = scope.$index" :icon="Delete" circle
                    class="wl-[1rem]" />
                </template>
                <p>Are you sure to delete this app?</p>
                <span class="ml-[0.5rem]">
                  <el-button size="small" text @click="showDelete = -1">cancel</el-button>
                  <el-button size="small" type="danger" @click="deleteApp(scope.row)">confirm</el-button>
                </span>
              </el-popover>
              <el-dropdown class="ml-[0.5rem]" trigger="click">
                <el-button size="small" circle :icon="More" />
                <el-dialog v-model="showLog" top="5vh" width="98%" title="View Log">
                  <CodeEditor height="60vh" :value="logs" mode="log" readOnly light></CodeEditor>
                </el-dialog>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item :icon="Log" @click="logApp(scope.row)">Log</el-dropdown-item>
                    <el-dropdown-item :icon="ApiApp" @click="proxyApp(scope.row)">Proxy</el-dropdown-item>
                    <el-dropdown-item :icon="Browser" @click="openApp(scope.row)">Open</el-dropdown-item>
                    <el-dropdown-item :icon="PauseOne" @click="operateApp(scope.row, 'stop')">Stop</el-dropdown-item>
                    <el-dropdown-item :icon="Power" @click="operateApp(scope.row, 'start')">Start</el-dropdown-item>
                    <el-dropdown-item :icon="RefreshOne" @click="operateApp(scope.row, 'restart')">Restart
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import {
  Edit, Delete, Terminal, Log, More, ApplicationOne,
  PauseOne, Power, RefreshOne, ApiApp, Browser, Search
} from '@icon-park/vue-next';
import { ref, unref, onMounted, computed } from 'vue';
import { ElMessage } from "element-plus";
import request from '@/axios'
import { useRouter } from 'vue-router';
import CodeEditor from '@/components/CodeEditor.vue';

const router = useRouter();

const apps = ref([]);
const showCreate = ref(false);
const showUpdate = ref(false);
const showLog = ref(false);
const showDelete = ref(-1);
const logs = ref("");
const newApp = ref({
  name: '',
  image: '',
  cmd: []
});

const updatedApp = ref({
  name: '',
  image: '',
  cmd: []
});
const createFormRef = ref();
const updateFormRef = ref();

const search = ref('');
const filter = computed(() =>
  apps.value.filter(
    (data) =>
      !search.value ||
      data.name.toLowerCase().includes(search.value.toLowerCase())
  )
)

const defaultTime = { timeout: "10000" }
onMounted(
  () => {
    request.get("/api/v1/containers").then((response) => {
      apps.value = response.data.data;
    })
  }
)

const getCommand = (cmd) => {
  if (Array.isArray(cmd)) {
    return cmd
  }

  return cmd.trim().split(/\s+/)
};

const createApp = () => {
  const form = unref(createFormRef)
  if (!form) {
    return
  }

  form.validate((valid) => {
    if (valid) {
      request.post("/api/v1/containers", {
        name: newApp.value.name,
        image: newApp.value.image,
        cmd: getCommand(newApp.value.cmd),
        describe: newApp.describe,
      }, defaultTime).then((response) => {
        ElMessage.success("Create success");
        apps.value.push(response.data.data);
        showCreate.value = false;
      })
    } else {
      ElMessage.error("Input invalid");
    }
  });
};

const editApp = (row) => {
  updatedApp.value = Object.assign({}, row);
  showUpdate.value = true;
}

const updateApp = () => {
  const form = unref(updateFormRef);
  if (!form) {
    return
  }

  form.validate((valid, err) => {
    updatedApp.value.cmd = getCommand(updatedApp.value.cmd);

    if (valid) {
      request.put("/api/v1/containers/" + updatedApp.value.id, updatedApp.value, defaultTime).then((response) => {
        ElMessage.success("Update success");
        const index = apps.value.findIndex(v => v.id === updatedApp.value.id);
        apps.value[index] = updatedApp.value;
        showUpdate.value = false;
      })
    } else {
      ElMessage.error("Input invalid", err);
    }
  });
};

const deleteApp = (row) => {
  request.delete("/api/v1/containers/" + row.id).then(() => {
    ElMessage.success("Delete success");
    const index = apps.value.findIndex(v => v.id === row.id);
    apps.value.splice(index, 1);
    showDelete.value = -1;
  })
};

const operateApp = (row, verb) => {
  request.post(`/api/v1/containers/${row.id}?verb=${verb}`).then(() => {
    ElMessage.success(verb + " success");
    const index = apps.value.findIndex(v => v.id === row.id);
    apps.value[index].status = verb;
  })
}

const getShortID = (row) => {
  return row.id.substring(0, 12)
}

const execApp = (row) => {
  window.open("/apps/" + getShortID(row) + "/exec", '_blank');
}

const proxyApp = (row) => {
  router.push("/apps/" + getShortID(row) + "/proxy")
}

const openApp = (row) => {
  const uri = "/api/v1/containers/" + getShortID(row) + "/proxy"
  window.open(uri, '_blank')
}

const logApp = (row) => {
  showLog.value = true
  const id = getShortID(row)
  request.get(`/api/v1/containers/${id}/log`).then((response) => {
    logs.value = response.data
  })
}

const getAppStatusType = (status) => {
  if (status == "running" || status == "start") {
    return "success"
  } else if (status == "stop" || status == "exited") {
    return "warning"
  } else if (status == "dead") {
    return "danger"
  } else {
    return "info"
  }
}

</script>
