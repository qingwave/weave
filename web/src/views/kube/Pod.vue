<template>
    <div class="w-full justify-center">
        <el-dialog v-model="showUpdate" top="5vh" width="98%" title="View Yaml">
            <CodeEditor height="60vh" :value="obj2yaml(updatedPod)" readOnly></CodeEditor>
        </el-dialog>

        <el-dialog v-model="showLog" top="5vh" width="98%" title="View Log">
            <CodeEditor height="60vh" :value="logs" mode="log" readOnly light></CodeEditor>
        </el-dialog>

        <div class="flex flex-col h-full px-[4rem] py-[2rem] space-y-[1rem]">
            <div class="flex flex-col overflow-hidden rounded-md shadow-md border">
                <div class="flex w-full h-[5rem] items-center">
                    <CubeThree class="ml-[1rem]" theme="filled" size="42" fill="#94A3B8" />
                    <span class="m-[0.75rem] text-2xl font-600">Pods</span>
                </div>
            </div>

            <el-card class="h-max flex-row">
                <template #header>
                    <div class="flex w-full space-x-[2rem]">
                        <el-select class="w-1/3" v-model="currentNamespace" filterable
                            placeholder="please select namespace">
                            <el-option v-for="ns in namespaces" :label="ns.metadata.name" :value="ns.metadata.name" />
                        </el-select>

                        <el-input v-model="search" placeholder="Type to search">
                            <template #prefix>
                                <el-icon>
                                    <Search />
                                </el-icon>
                            </template>
                        </el-input>
                    </div>
                </template>
                <el-table :data="filter" height="360" class="w-full max-h-full">
                    <el-table-column prop="metadata.name" label="Name" sortable />
                    <el-table-column prop="metadata.namespace" label="Namespace" />
                    <el-table-column prop="status" label="Status">
                        <template #default="scope">
                            <el-tag :type="getStatusType(scope.row.status.phase)"> {{ scope.row.status.phase}} </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="spec.containers[0].image" label="Image" min-width="120px" />
                    <el-table-column prop="metadata.creationTimestamp" label="StartAt" sortable min-width="120px" />
                    <el-table-column label="Operation" min-width="120px">
                        <template #default="scope">
                            <div class="space-x-[0.75rem]">
                                <el-button size="small" type="primary" circle @click="execPod(scope.row)" :icon="Terminal" />
                                <el-button size="small" circle @click="logPod(scope.row)" :icon="Log" />
                                <el-button type="success" size="small" circle @click="viewPod(scope.row)" :icon="FileDisplayOne" />
                                <el-popover :visible="showDelete == scope.$index" placement="top" :width="180">
                                    <template #reference>
                                        <el-button size="small" type="danger" @click="showDelete = scope.$index"
                                            :icon="Delete" circle class="wl-[1rem]" />
                                    </template>
                                    <p>Are you sure to delete this app?</p>
                                    <div class="ml-[0.5rem]">
                                        <el-button size="small" text @click="showDelete = -1">cancel</el-button>
                                        <el-button size="small" type="danger" @click="deletePod(scope.row)">confirm
                                        </el-button>
                                    </div>
                                </el-popover>
                            </div>
                        </template>
                    </el-table-column>
                </el-table>
            </el-card>
        </div>
    </div>
</template>

<script setup>
import {
    FileDisplayOne, Delete, CubeThree, Search, Terminal, Log
} from '@icon-park/vue-next';
import { ref, onMounted, watchEffect, computed } from 'vue';
import { ElMessage } from "element-plus";
import request from '@/axios';
import { obj2yaml } from '@/utils/yaml.js';
import CodeEditor from '@/components/CodeEditor.vue';

const showUpdate = ref(false);
const showDelete = ref(-1);

const updatedPod = ref({});

const currentNamespace = ref();
const namespaces = ref([]);

const pods = ref([]);

const search = ref('');

const filter = computed(() =>
    pods.value.filter(
        (data) =>
            !search.value ||
            data.metadata.name.toLowerCase().includes(search.value.toLowerCase())
    )
);

onMounted(
    () => {
        request.get("/api/v1/namespaces").then((response) => {
            namespaces.value = response.data.data.items;
            namespaces.value.sort(function (a, b) {
                let x = a.metadata.name;
                let y = b.metadata.name;
                return x.localeCompare(y);
            })
        })
    }
);

const getPods = (namespace) => {
    request.get(`/api/v1/namespaces/${namespace}/pods`).then((response) => {
        pods.value = response.data.data.items;
    })
};

watchEffect(() => {
    if (!currentNamespace.value) {
        return []
    }
    getPods(currentNamespace.value)
});

const viewPod = (row) => {
    updatedPod.value = Object.assign({}, row);
    showUpdate.value = true;
};

const deletePod = (row) => {
    request.delete(`/api/v1/namespaces/${currentNamespace.value}/pods`).then(() => {
        ElMessage.success("Delete success");
        const index = pods.value.findIndex(v => v.metadata.name === row.metadata.name);
        pods.value.splice(index, 1);
        showDelete.value = -1;
    })
};

const execPod = (row) => {
    window.open(`/namespaces/${currentNamespace.value}/pods/${row.metadata.name}/exec`,'_blank');
};

const showLog = ref(false);
const logs = ref();
const logPod = (row) => {
  logs.value = "";
  request.get(`/api/v1/namespaces/${currentNamespace.value}/pods/${row.metadata.name}/log`).then((response) => {
    logs.value = response.data
  })
  showLog.value = true;
}

const getStatusType = (status) => {
    if (status == "Pending") {
        return ""
    } else if (status == "Running") {
        return "success"
    } else if (status == "Failed") {
        return "danger"
    } else {
        return "info"
    }
};
</script>
