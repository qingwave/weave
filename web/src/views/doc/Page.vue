<template>
  <div class="flex flex-col w-full">
    <article class="w-full px-16">
      <div class="prose">
        <h1 class="my-6"> {{ title }} </h1>
        <div v-if="date" class="mb-4"> {{ date }} </div>
      </div>

      <MarkDown class="w-full mt-4 mb-8" :data="data" toc></MarkDown>
    </article>
  </div>
</template>

<script setup>
import MarkDown from 'components/MarkDown.vue';
import axios from "axios";
import { ref, onMounted, watch } from 'vue';
import fm from 'front-matter';

const props = defineProps({
  config: Object,
})

const title = ref(props.config.title)
const date = ref(props.config.date)
const data = ref("")

const parse = (content) => {
  if (!content) {
    return
  }

  let page = fm(content)
  data.value = page.body.replace(/#.+\n/, '')

  if (!title.value) {
    title.value = page.attributes.title
  }

  if (!date.value && page.attributes.date) {
    date.value = page.attributes.date
  }

  if (date.value) {
    date.value = (new Date(date.value)).toLocaleDateString('en', { year: "numeric", month: "short", day: "numeric" })
  }
}

watch(
  () => props.config,
  (val) => {
    getData(val)
  }
)

const getData = (config) => {
  title.value = config.title
  date.value = config.date

  if (config.content) {
    parse(config.content)
  } else if (config.url) {
    axios.get(config.url).then((response) => {
      parse(response.data)
    })
  }
}

onMounted(
  () => getData(props.config)
)

</script>
