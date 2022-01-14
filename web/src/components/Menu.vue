<template>
  <div class="flex overflow-hidden h-full flex-col">
    <el-menu
      class="w-16rem h-full overflow-x-hidden overflow-y-hidden"
      background-color="#0f172a"
      text-color="#a1a1aa"
      active-text-color="#ffffff"
      unique-opened
      :collapse="isCollapse"
      router
    >
      <el-menu-item index="/" class="h-16 center">
        <el-icon size="22">
          <img src="@/assets/weave.png" />
        </el-icon>
        <span class="font-bold text-xl">Weave</span>
      </el-menu-item>

      <el-menu-item :index="menu.name" v-for="menu in menuList" :key="menu.title">
        <el-icon size="14">
          <icon-park :type="menu.icon" />
        </el-icon>
        <span class>{{ menu.title }}</span>
      </el-menu-item>

      <el-sub-menu :index="menu.name" v-for="menu in subMenuList" :key="menu.title">
        <template #title>
          <el-icon size="14">
            <icon-park :type="menu.icon" />
          </el-icon>
          <span class>{{ menu.title }}</span>
        </template>
        <el-menu-item-group>
          <el-menu-item
            :index="item.name"
            v-for="item in menu.children"
            :key="item.name"
          >{{ item.title }}</el-menu-item>
        </el-menu-item-group>
      </el-sub-menu>

      <div class="absolute h-8 w-full bottom-0 text-gray-400 pl-20px my-2">
        <div v-if="isCollapse">
          <icon-park type="menu-fold-one" v-model="isCollapse" @click="collapseMenu" />
        </div>
        <div v-else>
          <icon-park type="menu-unfold-one" v-model="isCollapse" @click="collapseMenu" />
        </div>
      </div>
    </el-menu>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { IconPark } from '@icon-park/vue-next/es/all';

const isCollapse = ref(false);

const menuList = [
  {
    name: '/dashboard',
    icon: 'dashboard-car',
    title: 'Dashboard',
  },
  {
    name: '/apps',
    icon: 'application-two',
    title: 'Applications',
  },
  {
    name: '/users',
    icon: 'user',
    title: 'Users',
  },
  {
    name: '/about',
    icon: 'like',
    title: 'About',
  }
]

const subMenuList = [
  {
    name: '/others',
    icon: 'more-four',
    title: 'Others',
    children: [
      {
        name: '/404',
        title: '404'
      }
    ]
  }
]

function collapseMenu() {
  isCollapse.value = !isCollapse.value;
}

</script>

<style scoped>
.el-menu-item.is-active {
  background-color: #334155 !important;
}

.lg-fixed {
    position: fixed;
    top: 0;
    left: 0;
    z-index: 100;
    height: 100%;
    overflow: auto;
    overflow-x: hidden;
}
</style>
