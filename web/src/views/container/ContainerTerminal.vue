<template>
  <Terminal type="docker" :uri="uri" />
</template>

<script setup>
import Terminal from '@/components/Terminal.vue';
import { useRoute } from 'vue-router';
import { ref } from 'vue';

const route = useRoute();
const containerId = route.params.id;

const uri = ref(`/api/v1/containers/${containerId}/exec`);

</script>