<template>
  <div class="flex w-full h-full">
    <Menu/>
    <div class="flex-col w-full min-h-screen overflow-hidden">
      <Header/>
      <div class="h-full w-full">
          <router-view class="overflow-y-scroll" />
      </div>
    </div>
  </div>
</template>

<script step>
// @ is an alias to /src
import Menu from 'components/Menu.vue'
import Header from 'components/Header.vue'
</script>
