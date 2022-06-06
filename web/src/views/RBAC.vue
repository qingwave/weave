<template>
  <div class="w-full justify-center">
    <div class="flex flex-col h-full px-4rem py-2rem space-y-1rem">
      <div class="flex flex-col overflow-hidden rounded-lg shadow-lg">
        <div class="flex w-full h-5rem bg-white items-center">
          <ListView class="ml-1rem" theme="filled" size="42" fill="#94A3B8" />
          <span class="m-0.75rem text-2xl font-600">RBAC Policies</span>
        </div>
        <div class="flex h-3rem items-center">
          <el-button-group v-for="(tab, index) in policyTabs" class="ml-1rem">
            <el-button plain text @click="currentTab = index">{{ tab.name }}</el-button>
          </el-button-group>
        </div>
      </div>

      <el-card class="h-max w-full">
        <div class="flex flex-col w-full" v-for="(tab, index) in policyTabs">
          <Policy class="w-full mx-1rem" v-if="currentTab == index" :ptype="tab.ptype" :labels="tab.labels" />
        </div>
      </el-card>
    </div>
  </div>
</template>

<style scoped>
</style>

<script setup>
import { ref } from 'vue';
import Policy from '@/components/Policy.vue';
import { ListView } from '@icon-park/vue-next';

const policyTabs = [
  {
    name: "Policy",
    ptype: "p",
    labels: ["Role", "Group", "Resource", "ResourceName", "Action"]
  },
  {
    name: "UserRole",
    ptype: "g",
    labels: ["User", "Role", "Group"]
  },
  {
    name: "ResourceGroup",
    ptype: "g2",
    labels: ["Resource", "ResourceGroup"]
  },
]

const currentTab = ref(0);

const handlePolicy = (tab) => {
  currentTab.value = tab;
}
</script>
