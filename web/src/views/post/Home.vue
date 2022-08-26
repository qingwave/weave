<template>
    <div class="flex flex-col w-full h-full">
        <Header class="w-full" />
        <div class="w-full h-full overflow-y-scroll">
            <router-view :key=route.path />
        </div>
    </div>
</template>

<script setup>
// @ is an alias to /src
import Header from '@/views/home/Header.vue';
import { useRoute } from 'vue-router';

const route = useRoute();
</script>
