<template>
    <iframe :src="proxyURL" width="100%" height="100%" allowfullscreen></iframe>
</template>

<style scoped>
</style>

<script setup>
import {ref, onMounted} from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const proxyURL = ref('');

onMounted(
    () => {
        let id = route.params.id;
        if (id.length == 0) {
            return
        }
        let path = "/api/v1/containers/"+id+"/proxy" +route.params.proxyPath;
        proxyURL.value = path
    }
)

</script>
