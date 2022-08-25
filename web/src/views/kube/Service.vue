<template>
    <div class="w-full justify-center">
        <el-dialog v-model="showCreate" title="Create Service" width="98%">
            <CodeEditor height="60vh" :value="obj2yaml(newService)" @change="onCreateChange"></CodeEditor>
            <template #footer>
                <span>
                    <el-button @click="showCreate = false">Cancel</el-button>
                    <el-button type="primary" @click="createService">Confirm</el-button>
                </span>
            </template>
        </el-dialog>
        <el-dialog v-model="showUpdate" top="5vh" width="98%" title="Edit Yaml">
            <CodeEditor height="60vh" :value="obj2yaml(updatedService)" @change="onChange"></CodeEditor>
            <template #footer>
                <span>
                    <el-button @click="showUpdate = false">Cancel</el-button>
                    <el-button type="primary" @click="updateService">Confirm</el-button>
                </span>
            </template>
        </el-dialog>

        <div class="flex flex-col h-full px-[4rem] py-[2rem] space-y-[1rem]">
            <div class="flex flex-col overflow-hidden rounded-md shadow-md border">
                <div class="flex w-full h-[5rem] items-center">
                    <CircularConnection class="ml-[1rem]" theme="filled" size="42" fill="#94A3B8" />
                    <span class="m-[0.75rem] text-2xl font-600">Services</span>
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

                        <el-button type="primary" plain :icon="CircularConnection" @Click="clickCreate">
                            Create
                        </el-button>
                    </div>
                </template>
                <el-table :data="filter" height="360" class="w-full max-h-full">
                    <el-table-column prop="metadata.name" label="Name" sortable />
                    <el-table-column prop="metadata.namespace" label="Namespace" />
                    <el-table-column prop="spec.clusterIP" label="ClusterIP" />
                    <el-table-column prop="spec.ports" label="Ports" min-width="120px">
                        <template #default="scope">
                            {{ getPorts(scope.row.spec.ports) }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="metadata.creationTimestamp" label="StartAt" sortable min-width="120px" />
                    <el-table-column label="Operation" min-width="120px">
                        <template #default="scope">
                            <div class="space-x-[0.75rem]">
                                <el-button class="" size="small" circle @click="editService(scope.row)" :icon="Edit" />

                                <el-popover :visible="showDelete == scope.$index" placement="top" :width="180">
                                    <template #reference>
                                        <el-button size="small" type="danger" @click="showDelete = scope.$index"
                                            :icon="Delete" circle class="wl-[1rem]" />
                                    </template>
                                    <p>Are you sure to delete this app?</p>
                                    <div class="ml-[0.5rem]">
                                        <el-button size="small" text @click="showDelete = -1">cancel</el-button>
                                        <el-button size="small" type="danger" @click="deleteService(scope.row)">confirm
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
    Edit, Delete, CircularConnection, Search
} from '@icon-park/vue-next';
import { ref, onMounted, watchEffect, computed, toRaw } from 'vue';
import { ElMessage } from "element-plus";
import request from '@/axios';
import { obj2yaml, yaml2obj } from '@/utils/yaml.js';
import CodeEditor from '@/components/CodeEditor.vue';

const showCreate = ref(false);
const showUpdate = ref(false);
const showDelete = ref(-1);

const updatedService = ref({});

const currentNamespace = ref();
const namespaces = ref([]);

const services = ref([]);

const search = ref('');

const filter = computed(() =>
    services.value.filter(
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

const getServices = (namespace) => {
    request.get(`/api/v1/namespaces/${namespace}/services`).then((response) => {
        services.value = response.data.data.items;
    })
};

watchEffect(() => {
    if (!currentNamespace.value) {
        return []
    }
    getServices(currentNamespace.value)
});

const getPorts = (data) => {
    const ports = toRaw(data)
    let arr = new Array();
    for (let p of Array.from(ports)) {
        arr.push(`${p.port}/${p.protocol}`)
    }
    return arr.join(",")
};

const newService = {
        apiVersion: "v1",
        kind: "Service",
        metadata: {
            name: {},
            annotations: {
                "weave.io/platform": "true",
            },
            labels: {
                "weave.io/platform": "true",
                "weave.io/service": {},
            }
        },
        spec: {
            selector: {
                "weave.io/app": {}
            },
            ports: [
                {
                    port: {},
                    protocol: "TCP",
                }
            ]
        }
    }

const clickCreate = () => {
    newCreateService = {};
    showCreate.value = true;
}

const createService = () => {
    if (newCreateService instanceof Object) {
        ElMessage.info("Nothing to create")
        return
    } else {
        newCreateService = yaml2obj(newCreateService)
        if (!newCreateService) {
            ElMessage.warning("Input invaild")
            return
        }
    }

    const namespace = newCreateService.metadata.namespace ? newCreateService.metadata.namespace : currentNamespace.value
    request.post(`/api/v1/namespaces/${namespace}/services`, newCreateService).then((response) => {
        ElMessage.success("Create success");
        services.value.push(response.data.data);
        newCreateService = {}
        showCreate.value = false;
    })
};

const editService = (row) => {
    updatedService.value = Object.assign({}, row);
    newUpdatedService = updatedService.value;
    showUpdate.value = true;
}

let newCreateService = {};
const onCreateChange = (val, cm) => {
    newCreateService = val
};

let newUpdatedService = {};
const onChange = (val, cm) => {
    newUpdatedService = val
};

const updateService = () => {
    if (newUpdatedService instanceof Object) {
        ElMessage.info("Nothing to update")
        return
    } else {
        newUpdatedService = yaml2obj(newUpdatedService)
        if (!newUpdatedService) {
            ElMessage.warning("Input invaild")
            return
        }
    }
    const namespace = updatedService.metadata.namespace ? updatedService.metadata.namespace : currentNamespace.value
    request.put(`/api/v1/namespaces/${namespace}/services/${updatedService.value.metadata.name}`,
        newUpdatedService).then((response) => {
            ElMessage.success("Update success");
            const index = services.value.findIndex(v => v.metadata.name === updatedService.value.metadata.name);
            services.value[index] = newUpdatedService;
            showUpdate.value = false;
        })
};

const deleteService = (row) => {
    const namespace = row.metadata.namespace ? row.metadata.namespace : currentNamespace.value
    request.delete(`/api/v1/namespaces/${namespace}/services/${row.metadata.name}`).then(() => {
        ElMessage.success("Delete success");
        const index = services.value.findIndex(v => v.metadata.name === row.metadata.name);
        services.value.splice(index, 1);
        showDelete.value = -1;
    })
};
</script>
