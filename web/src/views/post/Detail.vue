<template>
  <div class="flex w-full justify-center">
    <div class="max-w-fit">
      <Page :config="item" v-if="showPage"> </Page>
      <footer v-show="item.content" class="max-w-fit pl-16 text-sm leading-6 mt-16">
        <div
          class="pt-10 pb-28 border-t border-slate-200 sm:flex justify-between text-slate-500">
          <div class="mb-6 sm:mb-0 sm:flex">
            <p>Copyright © 2022 Qinng.</p>
          </div>
          <div class="mb-6 sm:mb-0 sm:flex">
            <p class="sm:ml-4 sm:pl-4">
              <router-link class="hover:text-slate-700" :to="`/posts/editor/${id}`" target="__blank" >Edit Page</router-link></p>
              <p class="sm:ml-4 sm:pl-4 sm:border-l sm:border-slate-200">
                <el-popover :visible="showDelete" placement="top" :width="180">
                <template #reference>
                  <button class="hover:text-slate-700" @click="showDelete = true">Remove Page</button>
                </template>
                <p>Are you sure to delete this article?</p>
                <span class="ml-[0.5rem]">
                  <el-button size="small" text @click="showDelete = false">cancel</el-button>
                  <el-button size="small" type="danger" @click="deletePost">confirm</el-button>
                </span>
              </el-popover>
            </p>
          </div>
        </div>
      </footer>
    </div>

  </div>
</template>

<script setup>
import Page from 'views/doc/Page.vue'
import { useRoute } from 'vue-router';
import { onMounted, ref } from 'vue';
import request from '@/axios';
import { ElMessage } from "element-plus";
import { usePostStore } from '@/store/post';

const showPage = ref(false)
const route = useRoute();
const id = route.params.post;
const item = ref({});
const post = ref({});
const postStore = usePostStore();

const showDelete = ref(false)

const defaultTime = { timeout: "5000" }

const getPost = () => {
  request.get(`/api/v1/posts/${id}`, defaultTime).then((response) => {
    post.value = response.data.data;
    item.value = {
      title: post.value.name,
      content: post.value.content,
      author: post.value.creator.name,
      date: post.value.createdAt,
      likes: post.value.likes,
      views: post.value.views,
    }
    postStore.set(post.value)
    showPage.value = true
  })
}

onMounted(
  getPost
)

const deletePost = () => {
  request.delete(`/api/v1/posts/${id}`).then(() => {
    ElMessage.success("Delete success");
    route.push("/posts/list")
  })
}
</script>