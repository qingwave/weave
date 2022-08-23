<template>
    <div class="flex flex-col w-full h-full">
        <Header class="w-full h-full" />

        <div class="flex flex-row flex-grow w-full overflow-hidden">
            <div class="flex flex-1 h-full">
                <el-menu class="w-60 py-[1rem] flex-grow overflow-x-hidden overflow-y-auto" router>

                    <el-menu-item :index="menu.name" v-for="menu in menuList" :key="menu.title">
                        <span class="font-bold">{{ menu.title }}</span>
                    </el-menu-item>

                    <el-sub-menu :index="menu.name" v-for="menu in subMenuList" :key="menu.title">
                        <template #title>
                            <span class="font-bold">{{ menu.title }}</span>
                        </template>
                        <el-menu-item-group>
                            <el-menu-item :index="item.name" v-for="item in menu.children" :key="item.name">
                                {{ item.title }}
                            </el-menu-item>
                        </el-menu-item-group>
                    </el-sub-menu>
                </el-menu>
            </div>

            <div class="w-full overflow-y-scroll">
                <router-view :key=route.path />
            </div>
        </div>
    </div>
</template>

<script setup>
// @ is an alias to /src
import Header from '@/views/home/Header.vue';
import { useRoute } from 'vue-router';

const route = useRoute();

const menuList = [
  {
    name: '/docs/introduce',
    title: 'Introduce'
  }
]

const subMenuList = [
  {
    title: 'Usage',
    name: "/usage",
    children: [
      {
        name: '/docs/authentication',
        title: 'Authentication'
      },
      {
        name: '/docs/oauth',
        title: 'OAuth'
      }
    ]
  },
  {
    title: 'Blog',
    name: "/blog",
    children: [
      {
        name: '/docs/golang-distributed-system-x-cron',
        title: 'Golang Distributed with Cron'
      },
      {
        name: '/docs/golang-distributed-system-x-etcd',
        title: 'Golang Distributed with etcd'
      },
      {
        name: '/docs/golang-distributed-system-x-redis',
        title: 'Golang Distributed with Redis'
      },
      {
        name: '/docs/golang-distributed-system-x-zk',
        title: 'Golang Distributed with ZooKeeper'
      },
      {
        name: '/docs/how-to-write-k8s-cni',
        title: 'How to Write K8s CNI'
      },
      {
        name: '/docs/k8s-golang-design-pattern',
        title: 'K8s Golang Design Pattern'
      }
    ]
  }
]

</script>

<style scoped>
.el-menu-item {
  white-space: normal !important;
  line-height: 1rem;
}
</style>
