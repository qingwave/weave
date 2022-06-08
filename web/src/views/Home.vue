<template>
  <div class="flex flex-col w-full h-full bg-slate">
      <Header class="w-full h-full bg-slate"/>
      <div class="flex flex-row flex-grow w-full overflow-hidden">
        <Menu/>
        <div class="w-full overflow-y-scroll">
          <router-view />
        </div>
      </div>
  </div>
</template>

<style scoped>
.bg-slate {
  --tw-bg-opacity: 1;
  background-color: #f1f5f9;
}
</style>

<script setup>
// @ is an alias to /src
import Menu from '@/views/home/Menu.vue'
import Header from '@/views/home/Header.vue'
</script>
