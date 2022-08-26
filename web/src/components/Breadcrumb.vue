<template>
  <el-breadcrumb class="py-[1rem] px-[2rem]" separator="<" v-show="show()">
    <el-breadcrumb-item v-for="(item, index) in list" :to="{path: item.path}" :key="index" >{{item.name}}</el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script setup>
import { onMounted, watch, ref } from 'vue';
import { useRoute } from 'vue-router';

const router = useRoute();
const list = ref([]);

const show = () => {
  if (!list.value) {
    return false
  }
  return list.length > 2
};

onMounted(
  () => {
    let matched = router.matched;
    list.value = matched;
  }
)

watch(
  () => router.matched,
  (newVal, oldVal) =>  {
    list.value = newVal;
  }
)

</script>