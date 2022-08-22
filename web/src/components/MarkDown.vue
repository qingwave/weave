<template>
    <div class="prose max-w-none" v-html="parse()">
    </div>
</template>

<script setup>
import { marked } from 'marked';
import prism from 'prismjs';

const props = defineProps({
  data: { type: String },
})

const renderer = new marked.Renderer()

marked.setOptions({
  renderer: renderer,
  highlight: (code, lang) => {
    if (lang = 'golang') {
      lang = 'go'
    }
    if (prism.languages[lang]) {
      return prism.highlight(code, prism.languages[lang], lang);
    } else {
      return code;
    }
  },
  gfm: true,
  tables: true,
  breaks: true,
  pedantic: false,
  sanitize: false,
  smartLists: true,
  smartypants: false,
})

const parse = () => {
  return marked(props.data)
}

</script>
