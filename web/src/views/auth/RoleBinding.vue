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
            <el-button type="primary" plain @click="showCreate = true">Create RoleBinding</el-button>
        </div>
        <el-dialog v-model="showCreate" center title="Create RoleBinding" width="30rem">
            <el-form ref="createFormRef" :model="newRoleBinding" label-position="left" label-width="auto">
                <el-form-item :label="props.subject" prop="subject" required>
                        <el-select class="w-full" v-model="newRoleBinding.subject" filterable
                            placeholder="please select subject">
                            <el-option v-for="sub in subjects" :label="sub.name" :value="sub.id" v-bind:key="sub.id"/>
                        </el-select>
                </el-form-item>
                <el-form-item label="Role" prop="role" required>
                    <el-select class="w-full" v-model="newRoleBinding.role" filterable
                        placeholder="please select role">
                        <el-option v-for="role in roles" :label="role.name" :value="role.id"  v-bind:key="role.id"/>
                    </el-select>
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="createRoleBinding">Confirm</el-button>
                    <el-button @click="showCreate = false">Cancel</el-button>
                </span>
            </template>
        </el-dialog>

        <el-table :data="filterRoleBindings" class="w-full mt-[0.75rem]">
            <el-table-column prop="subject.name" :label="props.subject" sortable />
            <el-table-column prop="role.name" label="Role" />
            <el-table-column prop="role.namespace" label="RoleNamespace" />
            <el-table-column prop="role.rules" label="Rules">
                <el-table-column label="Resource">
                    <template #default="scope">
                        <div v-for="rule in scope.row.role.rules" v-bind:key="rule">
                            {{rule.resource}}
                        </div>
                    </template>
                </el-table-column>
                <el-table-column label="Operation">
                    <template #default="scope">
                        <div v-for="rule in scope.row.role.rules" v-bind:key="rule">
                            {{rule.operation}}
                        </div>
                    </template>
                </el-table-column>
            </el-table-column>
    
            <el-table-column label="Operation">
                <template #default="scope">
                    <el-popover :visible="showDelete == scope.$index" placement="top" :width="180">
                        <template #reference>
                            <el-button size="small" type="danger" @click="showDelete = scope.$index" :icon="Delete"
                                circle class="wl-[1rem]" />
                        </template>
                        <p>Are you sure to delete this role binding?</p>
                        <div class="my-[0.5rem]">
                            <el-button size="small" text @click="showDelete = -1">cancel</el-button>
                            <el-button size="small" type="danger" @click="deleteRoleBindings(scope.row)">confirm</el-button>
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
import { Delete, Search } from '@icon-park/vue-next';
import { ref, unref, computed, onMounted } from 'vue';
import { ElMessage } from "element-plus";
import request from '@/axios';
import { updateItem, deleteItem } from '@/utils';

const props = defineProps({
    resource: {type: String},
    subject: {type: String}
})

const roleBindings = ref([]);
const roles = ref([]);
const subjects = ref([]);

const showCreate = ref(false);
const showDelete = ref(-1);

const newRoleBinding = ref({});

const createFormRef = ref();

const search = ref('');
const filterRoleBindings = computed(() =>
    roleBindings.value.filter(
        (data) =>
            !search.value ||
            data.name.toLowerCase().includes(search.value.toLowerCase())
    )
)

onMounted(() => {
    request.get(`/api/v1/${props.resource}`).then((response) => {
        subjects.value = Array.from(response.data.data);
        for (let sub of Array.from(response.data.data)) {
            for (let role of Array.from(sub.roles)) {
                roleBindings.value.push({
                    "subject": sub,
                    "role": role,
                })
            }
        }
    })
    request.get(`/api/v1/roles`).then((response) => {
        roles.value = Array.from(response.data.data)
    })
});

const createRoleBinding = () => {
    const form = unref(createFormRef);
    if (!form) {
        return
    }
    form.validate((valid, err) => {
        if (valid) {
            request.post(`/api/v1/${props.resource}/${newRoleBinding.value.subject}/roles/${newRoleBinding.value.role}`).then((response) => {
                ElMessage.success("Create success");
                let subject = roleBindings.value.find(item => item.subject.id == newRoleBinding.value.subject).subject;
                let role = roles.value.find(item => item.id == newRoleBinding.value.role);
                roleBindings.value.push({
                    "subject": subject,
                    "role": role,
                })
                showCreate.value = false;
            })
        } else {
            console.log(err)
            ElMessage.error("Input invalid, all fields required");
        }
    });
};

const deleteRoleBindings = (row) => {
    request.delete(`/api/v1/${props.resource}/${row.subject.id}/roles/${row.role.id}`).then(() => {
        ElMessage.success("Delete success");
        deleteItem(roleBindings.value, row);
        showDelete.value = -1;
    })
};

</script>
