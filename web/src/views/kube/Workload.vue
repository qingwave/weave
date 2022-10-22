<template>
    <div class="w-full justify-center">
        <el-dialog v-model="showCreate" top="5vh" title="Create Application" width="80%">
            <el-steps :active="createActive" finish-status="success">
                <el-step title="Basic Information" />
                <el-step title="Detail Settings" />
                <el-step title="Advanced Settings" />
            </el-steps>

            <el-form class="my-[2rem] " v-show="createActive == 0" ref="formRef0" :model="newApp" label-position="top"
                label-width="auto">
                <div class="flex flex-row w-full space-x-[2rem]">
                    <el-form-item class="w-1/2" label="Name" prop="name" required>
                        <el-input v-model="newApp.name" />
                        <span class="text-gray-400">The application name</span>
                    </el-form-item>
                    <el-form-item class="w-1/2" label="Namespace" prop="namespace" required>
                        <el-select class="w-full" v-model="newApp.namespace" filterable
                            placeholder="please select namespace">
                            <el-option v-for="ns in namespaces" :label="ns.metadata.name" :value="ns.metadata.name" v-bind:key="ns.metadata.name"/>
                        </el-select>
                        <span class="text-gray-400">The application namespace</span>
                    </el-form-item>
                </div>
                <el-form-item label="Describe" prop="describe" required>
                    <el-input v-model="newApp.describe" type="textarea" />
                    <span class="text-gray-400">The application describe information</span>
                </el-form-item>
            </el-form>

            <el-form class="my-[2rem]" v-show="createActive == 1" ref="formRef1" :model="newApp" label-position="top"
                label-width="auto">
                <div class="flex flex-col w-full space-y-[1rem]">
                    <div class="w-full border rounded px-[1rem]">
                        <el-form-item label="Replicas" prop="replicas" required>
                            <el-input-number v-model="newApp.replicas" :min="1" placeholder="Replicas" />
                        </el-form-item>
                    </div>
                    <div class="w-full border rounded px-[1rem]">
                        <div class="py-2 font-bold">Container Spec</div>
                        <div class="flex flex-row w-full space-x-[2rem]">
                            <el-form-item class="w-1/2" label="Container" prop="container" required>
                                <el-input v-model="newApp.container" />
                                <span class="text-gray-400">The container name</span>
                            </el-form-item>
                            <el-form-item class="w-1/2" label="Image" prop="image" required>
                                <el-input v-model="newApp.image" />
                                <span class="text-gray-400">The container image</span>
                            </el-form-item>
                        </div>
                        <el-form-item label="Command" prop="command">
                            <el-input v-model="newApp.command" type="textarea" />
                            <span class="text-gray-400">The container startup command</span>
                        </el-form-item>
                        <el-form-item label="Args" prop="args">
                            <el-input v-model="newApp.args" type="textarea" />
                            <span class="text-gray-400">The container startup args</span>
                        </el-form-item>
                    </div>
                </div>
            </el-form>

            <el-form class="my-[2rem]" v-show="createActive == 2" ref="formRef2" :model="newApp" label-position="top"
                label-width="auto">
                <div class="flex flex-col w-full space-y-[1rem]">
                    <div class="w-full border rounded px-[1rem]">
                        <div class="py-2 font-bold">Metadata</div>
                        <div class="py-1">Labels</div>
                        <el-form-item v-for="(item, index) in newApp.labels">
                            <div class="flex flex-row w-full space-x-[1rem] justify-center">
                                <el-input v-model="item.key" placeholder="key" />
                                <el-input v-model="item.value" placeholder="value" />
                                <el-button link @click.prevent="removeLabel(labelType, item)" :icon="Delete" />
                            </div>
                        </el-form-item>
                        <div class="my-[1rem] text-right">
                            <el-button @click="addLabel(labelType)">Add</el-button>
                        </div>

                        <div class="py-1">Annotations</div>
                        <el-form-item v-for="(item, index) in newApp.annotations">
                            <div class="flex flex-row w-full space-x-[1rem] justify-center">
                                <el-input v-model="item.key" placeholder="key" />
                                <el-input v-model="item.value" placeholder="value" />
                                <el-button link @click.prevent="removeLabel(annotationType, item)" :icon="Delete" />
                            </div>
                        </el-form-item>
                        <div class="my-[1rem] text-right">
                            <el-button @click="addLabel(annotationType)">Add</el-button>
                        </div>
                    </div>
                </div>
            </el-form>

            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="showCreate = false">Cancel</el-button>
                    <el-button v-show="createActive > 0" @click="createActive--">Previous</el-button>
                    <el-button type="primary" v-show="createActive < 2" @click="nextStep">Next</el-button>
                    <el-button type="primary" v-show="createActive == 2" @click="createApp">Confirm</el-button>
                </span>
            </template>
        </el-dialog>

        <el-dialog v-model="showUpdate" top="5vh" width="98%" title="Edit Yaml">
            <CodeEditor height="60vh" :value="obj2yaml(updatedApp)" @change="onChange"></CodeEditor>
            <template #footer>
                <span>
                    <el-button @click="showUpdate = false">Cancel</el-button>
                    <el-button type="primary" @click="updateApp">Confirm</el-button>
                </span>
            </template>
        </el-dialog>

        <div class="flex flex-col h-full px-[4rem] py-[2rem] space-y-[1rem]">
            <div class="flex flex-col overflow-hidden rounded-md shadow-md border">
                <div class="flex w-full h-[5rem] items-center">
                    <application-one class="ml-[1rem]" theme="filled" size="42" fill="#94A3B8" />
                    <span class="m-[0.75rem] text-2xl font-600">Workloads</span>
                </div>
                <div class="flex h-[3rem] items-center bg-slate-100">
                    <el-radio-group class="ml-4" v-model="appTypeIndex">
                        <el-radio-button v-for="(a, index) in appTypes" plain :label="index">{{ a.name }}
                        </el-radio-button>
                    </el-radio-group>
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

                        <el-button type="primary" plain :icon="ApplicationOne"
                            @Click="newApp.namespace = currentNamespace; showCreate = true">Create
                        </el-button>
                    </div>
                </template>
                <el-table :data="filter" height="360" class="w-full max-h-full">
                    <el-table-column prop="metadata.name" label="Name" sortable />
                    <el-table-column prop="metadata.namespace" label="Namespace" />
                    <el-table-column prop="status" label="Status">
                        <template #default="scope">
                            <el-tag :type="getStatusType(scope.row.status)">
                                {{ getAppStatusType(scope.row.status) }}({{ scope.row.status.readyReplicas || 0 }} / {{
                                        scope.row.status.replicas || 0
                                }})
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="spec.template.spec.containers[0].image" label="Image" min-width="120px" />
                    <el-table-column prop="metadata.creationTimestamp" label="StartAt" sortable min-width="120px" />
                    <el-table-column label="Operation" min-width="120px">
                        <template #default="scope">
                            <div class="space-x-[0.75rem]">
                                <el-button class="" size="small" circle @click="editApp(scope.row)" :icon="Edit" />

                                <el-popover :visible="showDelete == scope.$index" placement="top" :width="180">
                                    <template #reference>
                                        <el-button size="small" type="danger" @click="showDelete = scope.$index"
                                            :icon="Delete" circle class="wl-[1rem]" />
                                    </template>
                                    <p>Are you sure to delete this app?</p>
                                    <div class="ml-[0.5rem]">
                                        <el-button size="small" text @click="showDelete = -1">cancel</el-button>
                                        <el-button size="small" type="danger" @click="deleteApp(scope.row)">confirm
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
    Edit, Delete, ApplicationOne, Search
} from '@icon-park/vue-next';
import { ref, unref, onMounted, watchEffect, computed, reactive } from 'vue';
import { ElMessage } from "element-plus";
import request from '@/axios';
import { obj2yaml, yaml2obj } from '@/utils/yaml.js';
import CodeEditor from '@/components/CodeEditor.vue';

const appTypes = [
    { name: "Deployment", resource: "deployments" },
    { name: "StatefulSet", resource: "statefulsets" },
    { name: "DaemonSet", resource: "daemonsets" }
]

const appTypeIndex = ref(0);

const createActive = ref(0);

const showCreate = ref(false);
const showUpdate = ref(false);
const showDelete = ref(-1);
const newApp = reactive({
    name: '',
    image: '',
    cmd: [],
    labels: [{}],
    annotations: [{}],
});

const updatedApp = ref({
    name: '',
    image: '',
    cmd: []
});

const formRef0 = ref();
const formRef1 = ref();
const formRef2 = ref();
const createFormRef = reactive([formRef0, formRef1, formRef2])

const currentNamespace = ref();
const namespaces = ref([]);

const apps = ref([]);

const search = ref('');

const filter = computed(() =>
    apps.value.filter(
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

const getApps = (namespace, resource) => {
    request.get(`/api/v1/namespaces/${namespace}/${resource}`).then((response) => {
        apps.value = response.data.data.items;
    })
};

watchEffect(() => {
    if (!currentNamespace.value) {
        return []
    }
    getApps(currentNamespace.value, appTypes[appTypeIndex.value].resource)
});

const getCommand = (cmd) => {
    if (!cmd) {
        return []
    }

    if (Array.isArray(cmd)) {
        return cmd
    }

    if (typeof cmd === "string") {
        return cmd.trim().split(/\s+/)
    }

    return [cmd]
};

const nextStep = () => {
    let form = unref(createFormRef[createActive.value]);
    form.validate((valid, err) => {
        if (valid) {
            createActive.value++;
        } else {
            ElMessage.error("Input invalid");
        }
    })
};

const annotationType = "annotation";
const labelType = "label";
const removeLabel = (type, item) => {
    if (type == annotationType) {
        const index = newApp.annotations.indexOf(item)
        if (index !== -1) {
            newApp.annotations.splice(index, 1)
        }
    } else if (type == labelType) {
        const index = newApp.labels.indexOf(item)
        if (index !== -1) {
            newApp.labels.splice(index, 1)
        }
    }
}

const addLabel = (type) => {
    if (type == annotationType) {
        newApp.annotations.push({
            key: '',
            value: '',
        })
    } else if (type == labelType) {
        newApp.labels.push({
            key: '',
            value: '',
        })
    }
}


const createApp = () => {
    const app = {
        apiVersion: "v1",
        kind: appTypes[appTypeIndex.value].name,
        metadata: {
            name: newApp.name,
            namespace: newApp.namespace,
            annotations: {
                "weave.io/platform": "true",
                "weave.io/describe": newApp.describe,
            },
            labels: {
                "weave.io/platform": "true",
                "weave.io/app": newApp.name,
            }
        },
        spec: {
            replicas: newApp.replicas,
            selector: {
                matchLabels: { "weave.io/app": newApp.name }
            },
            template: {
                metadata: {
                    labels: { "weave.io/app": newApp.name }
                },
                spec: {
                    containers: [{
                        name: newApp.container,
                        image: newApp.image,
                        command: getCommand(newApp.command),
                        args: getCommand(newApp.args),
                    }]
                }
            }
        }
    }

    for (const item of newApp.labels) {
        if (item.key && item.value) {
            app.metadata.labels[item.key] = item.value
        }
    }

    for (const item of newApp.annotations) {
        if (item.key && item.value) {
            app.metadata.annotations[item.key] = item.value
        }
    }

    request.post(`/api/v1/namespaces/${newApp.namespace}/${appTypes[appTypeIndex.value].resource}`, app).then((response) => {
        ElMessage.success("Create success");
        apps.value.push(response.data.data);
        showCreate.value = false;
    })
};

const editApp = (row) => {
    updatedApp.value = Object.assign({}, row);
    newUpdatedApp = updatedApp.value;
    showUpdate.value = true;
}

let newUpdatedApp = {};
const onChange = (val, cm) => {
    newUpdatedApp = val
};

const updateApp = () => {
    if (newUpdatedApp instanceof Object) {
        ElMessage.info("Nothing to update")
        return
    } else {
        newUpdatedApp = yaml2obj(newUpdatedApp)
        if (!newUpdatedApp) {
            ElMessage.warning("Input invaild")
            return
        }
    }

    request.put(`/api/v1/namespaces/${currentNamespace.value}/${appTypes[appTypeIndex.value].resource}/${updatedApp.value.metadata.name}`,
        newUpdatedApp).then((response) => {
            ElMessage.success("Update success");
            const index = apps.value.findIndex(v => v.metadata.name === updatedApp.value.metadata.name);
            apps.value[index] = newUpdatedApp;
            showUpdate.value = false;
        })
};

const deleteApp = (row) => {
    request.delete(`/api/v1/namespaces/${currentNamespace.value}/${appTypes[appTypeIndex.value].resource}/${row.metadata.name}`).then(() => {
        ElMessage.success("Delete success");
        const index = apps.value.findIndex(v => v.metadata.name === row.metadata.name);
        apps.value.splice(index, 1);
        showDelete.value = -1;
    })
};

const getAppStatusType = (status) => {
    if (status.observedGeneration <= 1 || !status.readyReplicas || status.replicas == 0) {
        if (status.updatedReplicas < status.replicas) {
            return "creating"
        }
    }

    if (status.updatedReplicas < status.replicas) {
        return "updating"
    } else if (status.unavailableReplicas > 0 || status.readyReplicas < status.replicas) {
        return "notAllReady"
    } else {
        return "running"
    }
};

const getStatusType = (status) => {
    const code = getAppStatusType(status)
    if (code == "creating") {
        return ""
    } else if (code == "running") {
        return "success"
    } else if (code == "notAllReady") {
        return "danger"
    } else {
        return "info"
    }
};
</script>
