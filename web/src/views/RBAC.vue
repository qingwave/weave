<template>
  <div class="w-full h-full flex justify-center">
    <el-card class="mx-4rem my-2rem w-4/5 h-max" shadow="never">
      <el-tabs @tab-change="handlePolicy">
        <el-tab-pane v-for="(tab, index) in policyTabs" :label="tab.name">
          <Policy v-if="currentTab == index" :ptype="tab.ptype" :labels="tab.labels" />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<style scoped>
</style>

<script setup>
import {ref} from 'vue';
import Policy from '@/components/Policy.vue';

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
