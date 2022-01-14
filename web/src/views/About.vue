<template>
  <div class="container h-full">
    <div class="pt-2rem px-6rem">
      <h1 class="font-bold text-2xl mb-2rem">#About</h1>
      <p class="text-l mb-2rem" v-html="msg" />
      <div class="flex justify-center">
        <img src="@/assets/snowman.svg" class="w-18rem" >
      </div>
    </div>
  </div>
</template>

<script setup>

const msg = `Weave is server template write by golang and vue. <br > 
The backend support multi-function, contains gin, postgres, redis, gorm, jwt and so on; <br >
The frontend support vue, element-plus, vite and so on.`
</script>
