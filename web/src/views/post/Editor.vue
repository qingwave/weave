<template>
    <div class="flex flex-col w-full h-full">
        <el-dialog v-model="showCreate" title="Create Post" width="50%">
            <el-form :model="newPost" label-position="top" label-width="auto">
                <el-form-item label="Name" prop="name" required>
                    <el-input v-model="newPost.name" />
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button type="primary" @click="save">Confirm</el-button>
                    <el-button @click="showCreate = false">Cancel</el-button>
                </span>
            </template>
        </el-dialog>

        <div class="flex flex-col w-full h-full p-4">
            <div class="flex w-full h-10 space-x-6 border justify-between items-center px-2">
                <input class="w-4/5 border-none outline-none font-bold text-lg" v-model="newPost.name" type="text" placeholder="Please input article title..." />
                <el-button type="primary" @click="showCreate = true">Save</el-button>
            </div>
            <MarkDownEditor ref="me" :data="newPost.content" class="w-full h-full flex-1"></MarkDownEditor>
        </div>
    </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import MarkDownEditor from 'components/MarkDownEditor.vue';
import request from '@/axios';
import { ElMessage } from "element-plus";
import { usePostStore } from '@/store/post';
import { useRoute } from 'vue-router';

const route = useRoute()
const id = route.params.id

const showCreate = ref(false)

const newPost = ref({})
const me = ref()

onMounted(()=>{
    if (id) {
        const postStore = usePostStore()
        newPost.value = postStore.data
    }
})

const create = () => {
    request.post("/api/v1/posts", {
        name: newPost.value.name,
        content: me.value.content,
    }).then((response) => {
        ElMessage.success("Create success");
        showCreate.value = false;
    })
}

const update = () => {
    newPost.value.content = me.value.content
    request.put(`/api/v1/posts/${newPost.value.id}`, newPost.value).then((response) => {
        ElMessage.success("Update success");
        showCreate.value = false;
    })
}

const save = () => {
    if (newPost.value.id) {
        update()
    } else {
        create()
    }
}

</script>