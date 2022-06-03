<template>
  <Terminal type="k8s" :uri="uri" />
</template>

<script setup>
import Terminal from '@/components/Terminal.vue';
import { useRoute } from 'vue-router';
import { ref } from 'vue';

const route = useRoute();
const namespace = route.params.namespace;
const pod = route.params.pod;
let container = route.query.container;
let command = route.query.command;

if (!container) {
  container = "";
}

if (!command) {
  command = "sh"
}

const uri = ref(`/api/v1/namespaces/${namespace}/pods/${pod}/exec?command=${command}&container=${container}&tty=true`)

</script>