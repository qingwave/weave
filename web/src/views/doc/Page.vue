<template>
  <div class="flex flex-col w-full">
    <article class="px-16 prose">
      <h1 class="my-6"> {{ title }} </h1>

      <div v-if="date" class="mb-4"> {{ date }} </div>

      <MarkDown class="w-full mt-4 mb-8" :data=data></MarkDown>
    </article>
  </div>
</template>

<script setup>
import MarkDown from 'components/MarkDown.vue';
import axios from "axios";
import { ref, onMounted } from 'vue';
import fm from 'front-matter';

const props = defineProps({
  config: Object,
})

const title = ref(props.config.title)
const date = ref(props.config.date)
const data = ref("")

const parse = (content) => {
  let page = fm(content)
  data.value = page.body.replace(/#.+\n/, '')

  if (!title.value) {
    title.value = page.attributes.title
  }

  if (!date.value && page.attributes.date) {
    date.value = (new Date(page.attributes.date)).toLocaleDateString('en', { year: "numeric", month: "short", day: "numeric" })
  }
}

onMounted(
  () => {
    if (props.config.url) {
      axios.get(props.config.url).then((response) => {
        parse(response.data)
      })
    } else {
      parse(props.config.content)
    }
  }
)

</script>
