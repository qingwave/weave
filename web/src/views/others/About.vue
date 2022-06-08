<template>
  <div class="container w-full">
    <div class="py-2rem px-4rem h-full space-y-2rem">
      <h1 class="font-bold text-2xl">#About</h1>

      <div class="space-y-2">
        <p v-html="msg" />
        <ul>
          <li>Github: <el-link type="success" :href="githubInfo.project" target="_blank">{{ githubInfo.project }}
            </el-link>
          </li>
          <li>Document: <el-link type="success" :href="githubInfo.doc" target="_blank">{{ githubInfo.doc }}</el-link>
          </li>
        </ul>
      </div>
      <div class="flex justify-center">
        <img src="@/assets/snowman.svg" class="w-18rem">
      </div>
    </div>
  </div>
</template>

<script setup>

import { githubInfo } from '@/config.js';

const msg = `Weave is server template write by golang and vue. <br > 
The backend support multi-function, contains gin, postgres, redis, gorm, jwt and so on; <br >
The frontend support vue, element-plus, vite and so on.`
</script>
