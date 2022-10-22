<template>
  <div class="w-full justify-center">
    <div class="flex flex-col h-full px-[4rem] py-[2rem] space-y-[1rem]">
      <div class="flex flex-col overflow-hidden rounded-md shadow-md border">
        <div class="flex w-full h-[5rem] items-center">
          <ListView class="ml-[1rem]" theme="filled" size="42" fill="#94A3B8" />
          <span class="m-[0.75rem] text-2xl font-600">RBAC Policies</span>
        </div>
        <div class="flex h-[3rem] items-center pl-[1rem] bg-slate-100">
          <el-radio-group v-model="currentTab">
            <el-radio-button v-for="(tab, index) in rbacTabs" plain :label="index" v-bind:key="tab">{{ tab.name }}</el-radio-button>
          </el-radio-group>
        </div>
      </div>

      <el-card class="h-max w-full">
        <div class="flex flex-col w-full">
          <template v-for="(tab, index) in rbacTabs">
            <component v-if="currentTab == index" :is="tab.component" v-bind:key="index" :resource=tab.resource :subject=tab.subject />
          </template>
        </div>
      </el-card>
    </div>
  </div>
</template>

<style scoped>

</style>

<script setup>
import { ref } from 'vue';
import Role from './Role.vue';
import RoleBinding from './RoleBinding.vue';
import { ListView } from '@icon-park/vue-next';

const rbacTabs = [
  {
    name: "Roles",
    component: Role
  },
  {
    name: "GroupRoleBinding",
    component: RoleBinding,
    resource: "groups",
    subject: "Group"
  },
  {
    name: "UserRoleBinding",
    component: RoleBinding,
    resource: "users",
    subject: "User"
  }
]

const currentTab = ref(0);

const handlePolicy = (tab) => {
  currentTab.value = tab;
}
</script>
