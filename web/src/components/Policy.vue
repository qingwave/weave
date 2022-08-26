<template>
    <div class="flex flex-row w-full space-x-[2rem]">
        <el-input v-model="search" placeholder="Type to search">
            <template #prefix>
                <el-icon>
                    <Search />
                </el-icon>
            </template>
        </el-input>
        <el-button type="primary" plain @click="showCreate = true">Create Policy</el-button>
    </div>
    <el-dialog v-model="showCreate" center title="Create Policy" width="33%">
        <el-form ref="createFormRef" :model="newPolicy" label-position="left" label-width="auto">
            <el-form-item v-for="(label, i) in props.labels" :label="label">
                <el-input v-model="newPolicy[i]" />
            </el-form-item>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button type="primary" @click="createPolicy">Confirm</el-button>
                <el-button @click="showCreate = false">Cancel</el-button>
            </span>
        </template>
    </el-dialog>
    <el-dialog v-model="showUpdate" center title="Update Policy" width="33%">
        <el-form ref="updateFormRef" :model="updatedPolicy" label-position="left" label-width="auto">
            <el-form-item v-for="label, i in props.labels" :label="label">
                <el-input v-model="updatedPolicy[i]" />
            </el-form-item>
        </el-form>
        <template #footer>
            <span class="dialog-footer">
                <el-button type="primary" @click="updatePolicy(updateRow)">Confirm</el-button>
                <el-button @click="showUpdate = false">Cancel</el-button>
            </span>
        </template>
    </el-dialog>
    <el-table :data="filterPolicies" class="w-full mb-[2rem]">
        <el-table-column v-for="label, i in props.labels" :label="label">
            <template #default="scope">
                <span>{{ scope.row[i] }}</span>
            </template>
        </el-table-column>
        <el-table-column label="Operation" min-width="120px">
            <template #default="scope">
                <el-button size="small" circle @click="editPolicy(scope.row)" :icon="Edit"></el-button>
                <el-popover :visible="showDelete == scope.$index" placement="top" :width="180">
                    <template #reference>
                        <el-button size="small" type="danger" @click="showDelete = scope.$index" :icon="Delete" circle
                            class="wl-[1rem]" />
                    </template>
                    <p>Are you sure to delete this policy?</p>
                    <div class="my-[0.5rem]">
                        <el-button size="small" text @click="showDelete = -1">cancel</el-button>
                        <el-button size="small" type="danger" @click="deletePolicy(scope.row)">confirm</el-button>
                    </div>
                </el-popover>
            </template>
        </el-table-column>
    </el-table>
</template>

<style scoped>
</style>

<script setup>
import { Edit, Delete, Search } from '@icon-park/vue-next';
import { ref, unref, computed, onMounted } from 'vue';
import { ElMessage } from "element-plus";
import request from '@/axios'
import { updateItem, deleteItem } from '@/utils'

const props = defineProps({
    ptype: String,
    labels: Array,
})

const policies = ref([]);

const showCreate = ref(false);
const showUpdate = ref(false);
const showDelete = ref(-1);

const newPolicy = ref([]);
const updatedPolicy = ref([]);

const createFormRef = ref();
const updateFormRef = ref();
const updateRow = ref();

const search = ref('');
const filterPolicies = computed(() =>
    policies.value.filter(
        (data) =>
            !search.value ||
            data[0].toLowerCase().includes(search.value.toLowerCase()) ||
            data[1].toLowerCase().includes(search.value.toLowerCase())
    )
)

onMounted(() => {
    request.get(`/api/v1/policies?ptype=${props.ptype}`).then((response) => {
        policies.value = Array.from(response.data.data);
    })
});

const createPolicy = () => {
    const valid = newPolicy.value.length == props.labels.length
    if (valid) {
        request.post("/api/v1/policies", {
            type: props.ptype,
            action: "add",
            policy: [...newPolicy.value],
        }).then((response) => {
            ElMessage.success("Create success");
            policies.value.push(response.data.data);
            showCreate.value = false;
        })
    } else {
        ElMessage.error("Input invalid, all fields required");
    }
};

const editPolicy = (row) => {
    updatedPolicy.value = Array.from(row);
    updateRow.value = Array.from(row);
    showUpdate.value = true;
}

const updatePolicy = (row) => {
    const form = unref(updateFormRef);
    if (!form) {
        return
    }

    form.validate((valid) => {
        if (valid) {
            request.post("/api/v1/policies", {
                type: props.ptype,
                action: "update",
                policy: [...updatedPolicy.value],
                oldPolicy: [...row],
            }).then((response) => {
                ElMessage.success("Update success");
                updateItem(policies, row, updatedPolicy.value);
                showUpdate.value = false;
            })
        } else {
            ElMessage.error("Input invalid");
        }
    });
};

const deletePolicy = (row) => {
    request.post("/api/v1/policies", {
        type: props.ptype,
        action: "remove",
        policy: row,
    }).then(() => {
        ElMessage.success("Delete success");
        deleteItem(policies.value, row)
        showDelete.value = -1;
    })
};

</script>
