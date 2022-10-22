<template>
    <div class="p-4">
        <div class="flex flex-row w-full space-x-[2rem]">
            <el-input v-model="search" placeholder="Type to search">
                <template #prefix>
                    <el-icon>
                        <Search />
                    </el-icon>
                </template>
            </el-input>
            <el-button type="primary" plain @click="showCreate = true">Create Role</el-button>
        </div>
        <el-dialog v-model="showCreate" center title="Create Role" width="60%">
            <el-form ref="createFormRef" :model="newRole" label-position="left" label-width="auto">
                <el-form-item label="Name" prop="name" required>
                    <el-input v-model="newRole.name" />
                    <span class="text-gray-400">The role name</span>
                </el-form-item>
                <el-form-item label="Scope" prop="scope" class="w-1/2" required>
                    <el-select class="w-full" v-model="newRole.scope" placeholder="please select scope">
                        <el-option label="Cluster" value="cluster" />
                        <el-option label="Namespace" value="namespace" />
                    </el-select>
                </el-form-item>
                <el-form-item label="Namespace" prop="namespace" class="w-1/2" :required="newRole.scope=='namespace'">
                        <el-select class="w-full" v-model="newRole.namespace" filterable
                            placeholder="please select namespace" :disabled="newRole.scope == 'cluster'">
                            <el-option v-for="g in groups" :label="g.name" :value="g.name" v-bind:key="g.name"/>
                        </el-select>
                </el-form-item>

                <div class="py-4">Labels</div>

                <el-form-item v-for="(item, index) in newRole.rules">
                    <div class="flex flex-row w-full space-x-[1rem] justify-center">
                        <el-input v-model="item.resource" placeholder="resource" />
                        <el-input v-model="item.operation" placeholder="operation" />
                        <el-button link @click.prevent="removeRule(newRole, item)" :icon="Delete" />
                    </div>
                </el-form-item>
                <div class="my-[1rem] text-right">
                    <el-button @click="addRule(newRole)">Add</el-button>
                </div>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="createRole">Confirm</el-button>
                    <el-button @click="showCreate = false">Cancel</el-button>
                </span>
            </template>
        </el-dialog>
        <el-dialog v-model="showUpdate" center title="Update Role" width="33%">
            <el-form ref="updateFormRef" :model="updatedRole" label-position="left" label-width="auto">
                <el-form-item label="Name" prop="name" disabled>
                    <el-input v-model="updatedRole.name" disabled />
                </el-form-item>
                <el-form-item label="Scope" prop="scope" class="w-1/2" disabled>
                    <el-input v-model="updatedRole.scope" disabled />
                </el-form-item>
                <el-form-item label="Namespace" prop="namespace" class="w-1/2" disabled>
                    <el-input v-model="updatedRole.namespace" disabled />
                </el-form-item>

                <div class="py-4">Labels</div>
                <el-form-item v-for="(item, index) in updatedRole.rules">
                    <div class="flex flex-row w-full space-x-[1rem] justify-center">
                        <el-input v-model="item.resource" placeholder="resource" />
                        <el-input v-model="item.operation" placeholder="operation" />
                        <el-button link @click.prevent="removeRule(updatedRole, item)" :icon="Delete" type="danger" />
                    </div>
                </el-form-item>
                <div class="my-[1rem] text-right">
                    <el-button @click="addRule(updatedRole)">Add</el-button>
                </div>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="updateRole()">Confirm</el-button>
                    <el-button @click="showUpdate = false">Cancel</el-button>
                </span>
            </template>
        </el-dialog>
        <el-table :data="filterRoles" class="w-full mt-[0.75rem]">
            <el-table-column prop="name" label="Name" sortable />
            <el-table-column prop="scope" label="Scope" />
            <el-table-column prop="namespace" label="Namespace" />
            <el-table-column prop="rules" label="Rules">
                <el-table-column label="Resource">
                    <template #default="scope">
                        <div v-for="rule in scope.row.rules">
                            {{rule.resource}}
                        </div>
                    </template>
                </el-table-column>
                <el-table-column label="Operation">
                    <template #default="scope">
                        <div v-for="rule in scope.row.rules">
                            {{rule.operation}}
                        </div>
                    </template>
                </el-table-column>
            </el-table-column>
            <el-table-column label="Operation">
                <template #default="scope">
                    <el-button size="small" circle @click="editRole(scope.row)" :icon="Edit"></el-button>
                    <el-popover :visible="showDelete == scope.$index" placement="top" :width="180">
                        <template #reference>
                            <el-button size="small" type="danger" @click="showDelete = scope.$index" :icon="Delete"
                                circle class="wl-[1rem]" />
                        </template>
                        <p>Are you sure to delete this role?</p>
                        <div class="my-[0.5rem]">
                            <el-button size="small" text @click="showDelete = -1">cancel</el-button>
                            <el-button size="small" type="danger" @click="deleteRole(scope.row)">confirm</el-button>
                        </div>
                    </el-popover>
                </template>
            </el-table-column>
        </el-table>
    </div>
</template>

<style scoped>

</style>

<script setup>
import { Edit, Delete, Search } from '@icon-park/vue-next';
import { ref, unref, computed, onMounted } from 'vue';
import { ElMessage } from "element-plus";
import request from '@/axios'
import { updateItem, deleteItem } from '@/utils'

const roles = ref([]);
const groups = ref([]);

const showCreate = ref(false);
const showUpdate = ref(false);
const showDelete = ref(-1);

const newRole = ref({
    rules: [{}]
});
const updatedRole = ref({});
const updateRow = ref({});

const createFormRef = ref();
const updateFormRef = ref();

const search = ref('');
const filterRoles = computed(() =>
    roles.value.filter(
        (data) =>
            !search.value ||
            data.name.toLowerCase().includes(search.value.toLowerCase())
    )
)

onMounted(() => {
    request.get(`/api/v1/roles`).then((response) => {
        roles.value = Array.from(response.data.data);
    })

    request.get(`/api/v1/groups`).then((response) => {
        groups.value = Array.from(response.data.data);
    })
});

const createRole = () => {
    const form = unref(createFormRef);
    if (!form) {
        return
    }
    form.validate((valid, err) => {
        if (valid) {
            request.post("/api/v1/roles", newRole.value).then((response) => {
                ElMessage.success("Create success");
                roles.value.push(response.data.data);
                showCreate.value = false;
            })
        } else {
            console.log(err)
            ElMessage.error("Input invalid, all fields required");
        }
    });
};

const editRole = (row) => {
    updatedRole.value = row;
    updateRow.value = row;
    showUpdate.value = true;
}

const updateRole = () => {
    const form = unref(updateFormRef);
    if (!form) {
        return
    }

    form.validate((valid, err) => {
        if (valid) {
            request.put(`/api/v1/roles/${updatedRole.value.id}`, updatedRole.value).then((response) => {
                ElMessage.success("Update success");
                updateItem(roles, updateRow, updatedRole.value);
                showUpdate.value = false;
            })
        } else {
            console.log("Invalid update parameters", err)
            ElMessage.error(`Input invalid: ${err}`);
        }
    });
};

const deleteRole = (row) => {
    request.delete(`/api/v1/roles/${row.id}`).then(() => {
        ElMessage.success("Delete success");
        deleteItem(roles.value, row);
        showDelete.value = -1;
    })
};

const removeRule = (role, item) => {
    const index = role.rules.indexOf(item)
    if (index !== -1) {
        role.rules.splice(index, 1)
    }
}

const addRule = (role) => {
    role.rules.push({
        key: '',
        value: '',
    })
}

</script>
