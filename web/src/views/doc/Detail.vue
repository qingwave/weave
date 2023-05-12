<template>
    <Page :config="item"> </Page>
</template>

<script setup>
import Page from './Page.vue'
import { useRoute } from 'vue-router';
import { ref } from 'vue';

const route = useRoute();
const page = route.params.page;

const pages = {
    "introduce": {
        title: "Introduce",
        url: 'https://raw.githubusercontent.com/qingwave/writings/main/src/content/post/golang-vue-starter.md'
    },
    "authentication": {
        title: "Authentication",
        url: "https://raw.githubusercontent.com/qingwave/weave/master/document/authentication.md"
    },
    "oauth": {
        title: "OAuth",
        url: "https://raw.githubusercontent.com/qingwave/weave/master/document/oauth.md"
    },
    "golang-distributed-system-x-cron": {
        title: "Golang Distributed System: Cron",
        url: 'https://raw.githubusercontent.com/qingwave/writings/main/src/content/post/golang-distributed-system-x-cron.md'
    },
    "golang-distributed-system-x-etcd": {
        title: "Golang Distributed System: Etcd",
        url: 'https://raw.githubusercontent.com/qingwave/writings/main/src/content/post/golang-distributed-system-x-etcd.md'
    },
    "golang-distributed-system-x-redis": {
        title: "Golang Distributed System: Redis",
        url: 'https://raw.githubusercontent.com/qingwave/writings/main/src/content/post/golang-distributed-system-x-redis.md'
    },
    "golang-distributed-system-x-zk": {
        title: "Golang Distributed System: ZooKeeper",
        url: 'https://raw.githubusercontent.com/qingwave/writings/main/src/content/post/golang-distributed-system-x-zk.md'
    },
    "how-to-write-k8s-cni": {
        title: "How to Write K8s CNI",
        url: 'https://raw.githubusercontent.com/qingwave/writings/main/src/content/post/how-to-write-k8s-cni.md'
    },
    "k8s-golang-design-pattern": {
        title: "K8s Golang Design Pattern",
        url: 'https://raw.githubusercontent.com/qingwave/writings/main/src/content/post/k8s-golang-design-pattern.md'
    }
}

const item = ref(pages[page])

</script>