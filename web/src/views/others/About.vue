<template>
  <div class="container w-full py-[2rem] px-[4rem]">
      <MarkDown :data=data></MarkDown>
  </div>
</template>

<script setup>
import { githubInfo } from '@/config.js';
import MarkDown from 'components/MarkDown.vue';

const data = `
# About

Weave is server template write by golang and vue.

The backend support multi-function, contains gin, postgres, redis, gorm, jwt and so on; The frontend support vue, element-plus, vite and so on.

- Github: [${githubInfo.project}](${githubInfo.project})
- Document: [${githubInfo.doc}](${githubInfo.doc})

<img src="/src/assets/snowman.svg" width=300 />

`
</script>
