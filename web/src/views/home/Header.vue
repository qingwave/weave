<template>
  <el-header>
    <el-row class="flex h-full w-full" justify="center" align="middle">
      <el-col :span="4" class="flex text-center items-center content-center">
            <img class="w-[2.5rem] mx-[0.5rem]" src="@/assets/weave.png" />
            <span class="font-bold font-mono text-2xl pl-[0.5rem]">Weave</span>
      </el-col>
      <el-col :span="14" class="flex-col text-right content-center">
        <el-menu mode="horizontal" class="font-bold" active-text-color="#000000">
          <el-menu-item index="1">
            <router-link to="/">Home</router-link>
          </el-menu-item>
          <el-menu-item index="2">
            <router-link to="/docs">Document</router-link>
          </el-menu-item>
          <el-menu-item index="3">
            <router-link to="/posts">Post</router-link>
          </el-menu-item>
        </el-menu>
      </el-col>
      <el-col :span="6" class="text-right pr-[1rem] space-x-[1rem]">
        <el-button link @click="notImplement('Search')">
          <search theme="outline" size="18" fill="#333" />
        </el-button>
        <el-button link>
          <a :href="githubInfo.project" target="_blank">
            <github-one theme="outline" size="18" :fill="['#333']" />
          </a>
        </el-button>
        <el-dropdown placement="bottom">
          <el-button link>
            <el-avatar v-if="user.avatar" :size="50" :src="user.avatar" />
            <me v-else theme="two-tone" size="18" :fill="['#333', '#50e3c2']" />
          </el-button>
          <template #dropdown>
            <span class="flex items-center content-center text-center mt-[1rem] mb-[0.5rem] mx-4 font-bold text-l">
              <SunOne theme="two-tone" size="24" :fill="['#333', '#f8e71c']" />
              <span class="ml-2">Hi {{ user.name }}</span>
            </span>
            <el-dropdown-menu>
              <el-dropdown-item :icon="Info" @click="toUserInfo">
                UserInfo
              </el-dropdown-item>
            </el-dropdown-menu>
            <el-dropdown-menu>
              <el-dropdown-item :icon="SettingOne" @click="notImplement('Setting')">
                Setting
              </el-dropdown-item>
            </el-dropdown-menu>

            <el-dropdown-menu>
              <el-dropdown-item :icon="Logout" @click="logout">
                Logout
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-col>
    </el-row>
  </el-header>
</template>

<script setup>
import { Info, SettingOne, Logout, SunOne, Search, GithubOne, Me } from '@icon-park/vue-next';
import { getUser, delUser, isAnonymous } from '@/utils';
import request from '@/axios';
import { ElMessage, ElNotification } from "element-plus";
import { useRouter } from 'vue-router';
import { githubInfo } from '@/config.js';

const user = getUser();
const router = useRouter();

function logout() {
  let lg = function () {
    console.log("logout success")
    ElNotification.success({
      title: 'Logout Success',
      message: 'Bye~ ' + user.name,
      showClose: true,
      duration: 1500,
    })
    delUser();
    router.push('/login');
  }

  if (isAnonymous(user.name)) {
    lg();
    return
  }
  request.delete("/api/v1/auth/token").then(() => {
    lg();
  }).catch((error) => {
    console.log(error)
  })
};

function toUserInfo() {
  router.push(`/users/${user.id}`)
};

function notImplement(name) {
  ElMessage({
    message: name + ' Coming Soon',
    type: 'warning',
    duration: 1000,
  })
}

</script>

<style scoped>
.el-header {
  border-bottom: 1px solid #d1d5db;
}

.el-menu {
  border: none;
}

.el-menu-item {
  border: none;
}

.el-menu-item.is-active {
  border: none;
}
</style>