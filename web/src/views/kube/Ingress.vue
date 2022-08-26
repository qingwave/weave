<template>
    <div class="w-full justify-center">
        <el-dialog v-model="showCreate" title="Create Ingress" width="98%">
            <CodeEditor height="60vh" :value="obj2yaml(newIngress)" @change="onCreateChange"></CodeEditor>
            <template #footer>
                <span>
                    <el-button @click="showCreate = false">Cancel</el-button>
                    <el-button type="primary" @click="createIngress">Confirm</el-button>
                </span>
            </template>
        </el-dialog>
        <el-dialog v-model="showUpdate" top="5vh" width="98%" title="Edit Yaml">
            <CodeEditor height="60vh" :value="obj2yaml(updatedIngress)" @change="onChange"></CodeEditor>
            <template #footer>
                <span>
                    <el-button @click="showUpdate = false">Cancel</el-button>
                    <el-button type="primary" @click="updateIngress">Confirm</el-button>
                </span>
            </template>
        </el-dialog>

        <div class="flex flex-col h-full px-[4rem] py-[2rem] space-y-[1rem]">
            <div class="flex flex-col overflow-hidden rounded-md shadow-md border">
                <div class="flex w-full h-[5rem] items-center">
                    <NetworkDrive class="ml-[1rem]" theme="filled" size="42" fill="#94A3B8" />
                    <span class="m-[0.75rem] text-2xl font-600">Ingresses</span>
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

                        <el-button type="primary" plain :icon="NetworkDrive" @Click="clickCreate">
                            Create
                        </el-button>
                    </div>
                </template>
                <el-table :data="filter" height="360" class="w-full max-h-full">
                    <el-table-column prop="metadata.name" label="Name" sortable />
                    <el-table-column prop="metadata.namespace" label="Namespace" />
                    <el-table-column label="Hosts" min-width="120px">
                        <template #default="scope">
                            {{ getHosts(scope.row.spec.rules) }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="metadata.creationTimestamp" label="StartAt" sortable min-width="120px" />
                    <el-table-column label="Operation" min-width="120px">
                        <template #default="scope">
                            <div class="space-x-[0.75rem]">
                                <el-button class="" size="small" circle @click="editIngress(scope.row)" :icon="Edit" />

                                <el-popover :visible="showDelete == scope.$index" placement="top" :width="180">
                                    <template #reference>
                                        <el-button size="small" type="danger" @click="showDelete = scope.$index"
                                            :icon="Delete" circle class="wl-[1rem]" />
                                    </template>
                                    <p>Are you sure to delete this app?</p>
                                    <div class="ml-[0.5rem]">
                                        <el-button size="small" text @click="showDelete = -1">cancel</el-button>
                                        <el-button size="small" type="danger" @click="deleteIngress(scope.row)">confirm
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
    Edit, Delete, NetworkDrive, Search
} from '@icon-park/vue-next';
import { ref, onMounted, watchEffect, computed, toRaw } from 'vue';
import { ElMessage } from "element-plus";
import request from '@/axios';
import { obj2yaml, yaml2obj } from '@/utils/yaml.js';
import CodeEditor from '@/components/CodeEditor.vue';

const showCreate = ref(false);
const showUpdate = ref(false);
const showDelete = ref(-1);

const updatedIngress = ref({});

const currentNamespace = ref();
const namespaces = ref([]);

const ingresses = ref([]);

const search = ref('');

const filter = computed(() =>
    ingresses.value.filter(
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

const getIngresses = (namespace) => {
    request.get(`/api/v1/namespaces/${namespace}/ingresses`).then((response) => {
        ingresses.value = response.data.data.items;
    })
};

watchEffect(() => {
    if (!currentNamespace.value) {
        return []
    }
    getIngresses(currentNamespace.value)
});

const getHosts = (data) => {
    const rules = toRaw(data)
    let arr = new Array();
    for (let x of Array.from(rules)) {
        arr.push(x.host)
    }
    return arr.join(",")
};

const newIngress = {
        apiVersion: "v1",
        kind: "Ingress",
        metadata: {
            name: {},
            namespace: {},
            labels: {
                "weave.io/platform": "true",
            }
        },
        spec: {
            rules: [
                {
                    host: {},
                    http: {
                        paths: [
                            {
                                backend: {
                                    service: {
                                        name: {},
                                        port: {
                                            number: {}
                                        }
                                    }
                                },
                                path: "/"
                            }
                        ]
                    }
                }
            ]
        }
    }

const clickCreate = () => {
    newCreateIngress = {};
    showCreate.value = true;
}

const createIngress = () => {
    if (newCreateIngress instanceof Object) {
        ElMessage.info("Nothing to create")
        return
    } else {
        newCreateIngress = yaml2obj(newCreateIngress)
        if (!newCreateIngress) {
            ElMessage.warning("Input invaild")
            return
        }
    }

    const namespace = newCreateIngress.metadata.namespace ? newCreateIngress.metadata.namespace : currentNamespace.value
    request.post(`/api/v1/namespaces/${namespace}/ingresses`, newCreateIngress).then((response) => {
        ElMessage.success("Create success");
        ingresses.value.push(response.data.data);
        newCreateIngress = {}
        showCreate.value = false;
    })
};

const editIngress = (row) => {
    updatedIngress.value = Object.assign({}, row);
    newUpdatedIngress = updatedIngress.value;
    showUpdate.value = true;
}

let newCreateIngress = {};
const onCreateChange = (val, cm) => {
    newCreateIngress = val
};

let newUpdatedIngress = {};
const onChange = (val, cm) => {
    newUpdatedIngress = val
};

const updateIngress = () => {
    if (newUpdatedIngress instanceof Object) {
        ElMessage.info("Nothing to update")
        return
    } else {
        newUpdatedIngress = yaml2obj(newUpdatedIngress)
        if (!newUpdatedIngress) {
            ElMessage.warning("Input invaild")
            return
        }
    }
    const namespace = updatedIngress.metadata.namespace ? updatedIngress.metadata.namespace : currentNamespace.value
    request.put(`/api/v1/namespaces/${namespace}/ingresses/${updatedIngress.value.metadata.name}`,
        newUpdatedIngress).then((response) => {
            ElMessage.success("Update success");
            const index = ingresses.value.findIndex(v => v.metadata.name === updatedIngress.value.metadata.name);
            ingresses.value[index] = newUpdatedIngress;
            showUpdate.value = false;
        })
};

const deleteIngress = (row) => {
    const namespace = row.metadata.namespace ? row.metadata.namespace : currentNamespace.value
    request.delete(`/api/v1/namespaces/${namespace}/ingresses/${row.metadata.name}`).then(() => {
        ElMessage.success("Delete success");
        const index = ingresses.value.findIndex(v => v.metadata.name === row.metadata.name);
        ingresses.value.splice(index, 1);
        showDelete.value = -1;
    })
};
</script>
