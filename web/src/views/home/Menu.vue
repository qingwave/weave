<template>
  <div class="flex flex-col h-full">
    <el-menu class="w-[14rem] py-[1rem] flex-grow overflow-x-hidden overflow-y-auto"
      text-color="#4b5563" active-text-color="#10b981" unique-opened :collapse="isCollapse" router>

      <el-menu-item :index="menu.name" v-for="menu in menuList" :key="menu.title">
        <el-icon size="14">
          <component :is="menu.icon" />
        </el-icon>
        <span class="font-bold">{{ menu.title }}</span>
      </el-menu-item>
      
      <el-sub-menu :index="menu.name" v-for="menu in subMenuList" v-show="menu.show" :key="menu.title">
        <template #title>
          <el-icon size="14">
            <component :is="menu.icon" />
          </el-icon>
          <span class="font-bold">{{ menu.title }}</span>
        </template>
        <el-menu-item-group>
          <el-menu-item :index="item.name" v-for="item in menu.children" :key="item.name" class="font-bold">
          {{item.title}}
          </el-menu-item>
        </el-menu-item-group>
      </el-sub-menu>
    </el-menu>
    
    <div class="flex flex-grow-1 w-full my-[1rem] pl-[1.5rem]">
      <div v-if="isCollapse">
        <menu-fold-one v-model="isCollapse" @click="collapseMenu" />
      </div>
      <div v-else>
        <menu-unfold-one v-model="isCollapse" @click="collapseMenu" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import {
  DashboardCar, ApplicationTwo, Peoples,
  MoreFour, MenuFoldOne, MenuUnfoldOne, CategoryManagement,
  SailboatOne
} from '@icon-park/vue-next';
import { isAdmin } from '@/utils';

const isRoot = isAdmin();

const isCollapse = ref(false);

const menuList = [
  {
    name: '/dashboard',
    icon: DashboardCar,
    title: 'Dashboard',
  },
  {
    name: '/apps',
    icon: ApplicationTwo,
    title: 'Applications',
  },
  {
    name: '/user_groups',
    icon: Peoples,
    title: 'Groups',
  }
]

const subMenuList = [
  {
    title: 'Kubernetes',
    icon: SailboatOne,
    name: "/kubernetes",
    show: isRoot,
    children: [
      {
        name: '/namespaces',
        title: 'Namespaces'
      },
      {
        name: '/workloads',
        title: 'Workloads'
      },
      {
        name: '/pods',
        title: 'Pods'
      },
      {
        name: '/services',
        title: 'Services'
      },
      {
        name: '/ingresses',
        title: 'Ingresses'
      }
    ]
  },
  {
    title: 'Admin',
    icon: CategoryManagement,
    name: "/admin",
    show: isRoot,
    children: [
      {
        name: '/users',
        title: 'Users'
      },
      {
        name: '/groups',
        title: 'Groups'
      },
      {
        name: '/rbac',
        title: 'RBAC'
      },
    ]
  },
  {
    title: 'Others',
    show: true,
    name: '/others',
    icon: MoreFour,
    children: [
      {
        name: '/about',
        title: 'About'
      },
      {
        name: '/markdown',
        title: 'MarkDown'
      },
      {
        name: '/webcode',
        title: 'WebCode'
      },
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
.el-menu {
  border: none;
}
</style>
