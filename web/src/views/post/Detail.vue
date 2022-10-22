<template>
  <div class="flex w-full h-full justify-center">
    <div class="max-w-fit h-full">
      <Page :config="item" class="flex-1"> </Page>
      <div class="pr-48 mx-16">
        <div class="flex justify-center w-full space-x-4 h-auto">
          <div class="flex space-x-1">
            <button>
              <Like theme="two-tone" size="24" :fill="likeFill()" @click="like" />
            </button>
            <p>{{ post.likes }}</p>
          </div>
          <div class="flex space-x-1">
            <Planet theme="two-tone" size="24" :fill="['#333', '#38bdf8']" />
            <p>{{ post.views }}</p>
          </div>
          <div class="flex space-x-1">
            <MessageOne theme="two-tone" size="24" :fill="['#333', '#fbbf24']" />
            <p>{{ getCommentCount() }}</p>
          </div>
        </div>

        <div class="flex flex-col antialiased mx-8 my-4">
          <form class="flex my-2">
            <div class="mb-4 w-full bg-gray-50 rounded-lg border border-gray-200">
              <div class="py-2 px-4 bg-white rounded-t-lg">
                <textarea rows="4" v-model="commentText"
                  class="px-0 w-full text-sm text-gray-900 bg-white border-0 outline-none"
                  placeholder="Write a comment..." required=""></textarea>
              </div>
              <div class="flex justify-end items-center py-2 px-3 border-t">
                <el-button type="primary" @click="addComment">Post comment</el-button>
              </div>
            </div>
          </form>

          <div class="space-y-4">
            <div class="flex" v-for="c in comments">
              <div class="flex-shrink-0 mr-3">
                <Panda size="32" :fill="['#333' ,'#2F88FF']" />
              </div>
              <div class="flex-1 border rounded-lg px-4 py-2 sm:px-6 sm:py-4 leading-relaxed">
                <strong> {{ c.user.name }}</strong> <span class="text-xs text-gray-400">{{ getTime(c.createdAt)
                }}</span>
                <p class="text-sm"> {{ c.content }} </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <footer v-show="item.content" class="flex-2 h-16 max-w-fit pl-16 text-sm leading-6 mt-16">
        <div class="pt-10 pb-28 border-t border-slate-200 sm:flex justify-between text-slate-500">
          <div class="mb-6 sm:mb-0 sm:flex">
            <p>Copyright © 2022 Qinng.</p>
          </div>
          <div class="mb-6 sm:mb-0 sm:flex">
            <p class="sm:ml-4 sm:pl-4">
              <el-popover :visible="showDelete" placement="top" :width="180">
                <template #reference>
                  <button class="hover:text-slate-700" @click="showDelete = true">Remove</button>
                </template>
                <p>Are you sure to delete this article?</p>
                <span class="ml-[0.5rem]">
                  <el-button size="small" text @click="showDelete = false">cancel</el-button>
                  <el-button size="small" type="danger" @click="deletePost">confirm</el-button>
                </span>
              </el-popover>
            </p>
            <p class="sm:ml-4 sm:pl-4 sm:border-l sm:border-slate-200">
              <router-link class="hover:text-slate-700" :to="`/posts/editor/${id}`" target="__blank">Edit Page
              </router-link>
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
import { onMounted, ref, computed } from 'vue';
import request from '@/axios';
import { ElMessage } from "element-plus";
import { usePostStore } from '@/store/post';
import { Like, Planet, MessageOne, Panda, Dog } from '@icon-park/vue-next';
import { getUser } from '@/utils'

const user = getUser();
const showPage = ref(false)
const route = useRoute();
const id = route.params.post;
const item = ref({});
const post = ref({});
const postStore = usePostStore();

const showDelete = ref(false)

const defaultTime = { timeout: "5000" }
const comments = ref([])

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
    comments.value = post.value.comments
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

const unLiked = ['#333', '#ffffff']
const Liked = ['#333', '#f87171']

const likeFill = () => {
  if (post.value.userLiked) {
    return Liked
  }
  return unLiked
}

const like = () => {
  if (post.value.userLiked) {
    request.delete(`/api/v1/posts/${id}/like`).then(() => {
      post.value.userLiked = false
      post.value.likes = post.value.likes - 1
      likeFill.value = unLiked
    })
  } else {
    request.post(`/api/v1/posts/${id}/like`).then(() => {
      post.value.userLiked = true
      post.value.likes = post.value.likes + 1
      likeFill.value = Liked
    })
  }
}

const getCommentCount = () => {
  if (comments.value) {
    return comments.value.length
  }
  return []
}


const commentText = ref('');

const addComment = () => {
  request.post(`/api/v1/posts/${id}/comment`, {
    postId: Number(id),
    userId: user.id,
    content: commentText.value
  }).then((response) => {
    comments.value.push(response.data.data)
  })
}

const getTime = (date) => {
  return (new Date(date)).toLocaleDateString('en', { year: "numeric", month: "short", day: "numeric" })
}

</script>