<template>
  <div class="container h-full flex justify-center align-center text-center">
    <div class="pt-2rem px-6rem">
      <h1 class="font-bold text-2xl mb-1rem">Hello Weave</h1>
      <h2 class="text-lg text-gray-500">Powerful server template with golang and vue</h2>
      <div class="my-2rem">
        <img src="@/assets/welcome.svg" class="w-18rem" >
      </div>
      <el-button :icon="Waves" type="primary" plain><router-link to="/about">About Weave</router-link></el-button>
      <el-button :icon="Github"><a :href="githuburl" target="_blank">Explore More</a></el-button>
    </div>
  </div>
</template>

<script setup>
import { Github, Waves } from '@icon-park/vue-next';

const githuburl = 'https://github.com/qingwave/weave';
</script>
