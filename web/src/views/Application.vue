<template>
  <div class="container h-full flex justify-center">
    <el-card class="mx-4rem my-2rem w-4/5" shadow="never">
      <template #header>
        <div class="flex justify-between">
          <span>Applications</span>
          <el-button type="primary" plain :icon="ApplicationOne" @Click="showCreate = true">Create</el-button>
        </div>
        <el-dialog v-model="showCreate" center title="Create Application" width="33%">
          <el-form ref="createFormRef" :model="newApp" label-position="left" label-width="auto">
            <el-form-item label="Name" prop="name" required>
              <el-input v-model="newApp.name" placeholder="App name" />
            </el-form-item>
            <el-form-item label="Image" prop="image" required>
              <el-input v-model="newApp.image" placeholder="App image" />
            </el-form-item>
            <el-form-item label="Command" prop="cmd">
              <el-input v-model="newApp.cmd" placeholder="App command" />
            </el-form-item>
          </el-form>
          <template #footer>
            <span class="dialog-footer">
              <el-button type="primary" @click="createApp">Confirm</el-button>
              <el-button @click="showCreate = false">Cancel</el-button>
            </span>
          </template>
        </el-dialog>
      </template>
      <el-table :data="apps" height="320" class="w-full max-h-full">
        <el-table-column prop="name" label="Name" />
        <el-table-column prop="image" label="Image" />
        <el-table-column prop="status" label="Status" />
        <el-table-column prop="cmd" label="Command" min-width="120px" />
        <el-table-column prop="startAt" label="StartAt" min-width="120px" />
        <el-table-column prop="Operation" label="Operation" min-width="120px">
          <template #default="scope">
            <el-button size="small" type="primary" circle @click="execApp(scope.$index)" :icon="Terminal" />
            <el-button size="small" circle @click="editApp(scope.$index)" :icon="Edit">
              <el-dialog v-model="showUpdate" center title="Update Application" width="33%">
                <el-form
                  ref="updateFormRef"
                  :model="updatedApp"
                  label-position="left"
                  label-width="auto"
                >
                  <el-form-item label="Name" prop="name">
                    <el-input v-model="updatedApp.name" disabled />
                  </el-form-item>
                  <el-form-item label="Image" prop="image" required>
                    <el-input v-model="updatedApp.image" placeholder="App image" />
                  </el-form-item>
                  <el-form-item label="Command" prop="cmd">
                    <el-input v-model="updatedApp.cmd" placeholder="App command" />
                  </el-form-item>
                </el-form>
                <template #footer>
                  <span class="dialog-footer">
                    <el-button type="primary" @click="updateApp(scope.$index)">Confirm</el-button>
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
              <p>Are you sure to delete this app?</p>
              <div class="my-0.5rem">
                <el-button size="small" type="text" @click="showDelete = -1">cancel</el-button>
                <el-button size="small" type="danger" @click="deleteApp(scope.$index)">confirm</el-button>
              </div>
            </el-popover>
            <el-dropdown class="pl-10px">
              <el-button size="small" circle :icon="More" />
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :icon="Log" @click="logApp(scope.$index)">Log</el-dropdown-item>
                  <el-dropdown-item :icon="ApiApp" @click="proxyApp(scope.$index)">Proxy</el-dropdown-item>
                  <el-dropdown-item :icon="Browser" @click="openApp(scope.$index)">Open</el-dropdown-item>
                  <el-dropdown-item :icon="PauseOne" @click="operateApp(scope.$index, 'stop')">Stop</el-dropdown-item>
                  <el-dropdown-item :icon="Power" @click="operateApp(scope.$index, 'start')">Start</el-dropdown-item>
                  <el-dropdown-item :icon="RefreshOne" @click="operateApp(scope.$index, 'restart')">Restart</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { Edit, Delete, Terminal, Log, More, ApplicationOne,
 PauseOne, Power, RefreshOne, ApiApp, Browser } from '@icon-park/vue-next';
import { ref, unref, reactive, onMounted } from 'vue';
import { ElMessage } from "element-plus";
import request from '@/axios'
import { useRouter } from 'vue-router';

const router = useRouter();

const apps = ref([
  {
    name: "app1",
    image: "image1",
    cmd: "aaa",
  }
]);
const showCreate = ref(false);
const showUpdate = ref(false);
const showDelete = ref(-1);
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

onMounted(
  () => {
    request.get("/api/v1/containers").then((response) => {
      apps.value = response.data.data;
    }).catch((error) => {
      if (error.response) {
        console.log("list failed: " + error.response.data.msg);
      }
      console.log(error)
    })
  }
)

const getCommand = (cmd) => {
  if (Array.isArray(cmd)) {
    return cmd
  }
  let commands = []
  if (cmd.indexOf(" ") == -1) {
    commands.push(cmd)
  } else {
    commands.push("sh", "-c", cmd)
  }

  return commands
}

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
      }).then((response) => {
        ElMessage.success("Create success");
        apps.value.push(response.data.data);
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

const editApp = (index) => {
  updatedApp.value = Object.assign({}, apps.value[index]);
  showUpdate.value = true;
}

const updateApp = (index) => {
  const form = unref(updateFormRef);
  if (!form) {
    return
  }

  let app = Object.assign({}, apps.value[index]);
  app.cmd = getCommand(updatedApp.value.cmd)

  form.validate((valid) => {
    if (valid) {
      request.put("/api/v1/containers/" + updatedApp.value.id, app).then((response) => {
        ElMessage.success("Update success");
        apps.value[index] = app;
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

const deleteApp = (index) => {
  request.delete("/api/v1/containers/" + apps.value[index].id).then(() => {
    ElMessage.success("Delete success");
    apps.value.splice(index, 1);
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

const operateApp = (index, verb) => {
  request.post("/api/v1/containers/" + apps.value[index].id + "?verb="+verb).then(() => {
    ElMessage.success(verb + " success");
    apps.value[index].status = verb;
  }).catch((error) => {
    let msg = verb + 'failed'
    if (error.response) {
      msg += ": " + error.response.data.msg
    } else {
      console.log(error)
    }
    ElMessage.error(msg);
  })
}

const getShortID = (index) => {
  return apps.value[index].id.substring(0,12)
}

const execApp = (index) => {
  router.push("/apps/"+ getShortID(index) + "/exec")
}

const proxyApp = (index) => {
  router.push("/apps/"+ getShortID(index) + "/proxy")
}

const openApp = (index) => {
  const uri = "/api/v1/containers/"+getShortID(index)+"/proxy"
  window.open(uri,'_blank')
}

</script>
