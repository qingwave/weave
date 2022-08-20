<template>
    <div class="prose h-full" v-html="parse()">
    </div>
</template>

<script setup>
import { marked } from 'marked';
import prism from 'prismjs';

const props = defineProps({
    data: { type: String },
})

marked.setOptions({
    highlight: (code, lang) => {
    if (prism.languages[lang]) {
      return prism.highlight(code, prism.languages[lang], lang);
    } else {
      return code;
    }
  },
})

const parse = () => {
    return marked(props.data)
}

</script>
