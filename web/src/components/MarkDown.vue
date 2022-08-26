<template>
  <div class="flex flex-row w-full h-full space-x-12 overflow-x-hidden">
    <div class="prose break-words w-max-none h-full" v-html="parse()">
    </div>

    <div v-if="props.toc" class="hidden w-48 lg:flex">
      <div  class="fixed top-24 bottom-0 overflow-y-auto no-scrollbar">
      <nav class="no-scrollbar font-medium leading-loose border-solid border-l-4 pl-4 mr-8">
        <ul v-for="(item, i) in tocList.c">
          <li>
            <a class="hover:text-emerald-500" :class="{ 'text-emerald-500': tocIndex == i }" @click="tocIndex = i"
              :href=getHref(item.n)> {{ item.n }} </a>
            <ul class="ml-6" v-for="(t, j) in item.c">
              <li><a class="hover:text-emerald-500" :class="{ 'text-emerald-500': tocIndex == (i + 1) * 1000 + j }"
                  @click="tocIndex = (i + 1) * 1000 + j" :href=getHref(t.n)> {{ t.n }}</a></li>
            </ul>
          </li>
        </ul>
      </nav>
      </div>
    </div>
  </div>
</template>

<script setup>
import MarkDownIt from 'markdown-it';
import prism from 'prismjs';
import MarkDownItAnchor from 'markdown-it-anchor';
import MarkFownDoneRight from 'markdown-it-toc-done-right';
import uslug from 'uslug';
import { ref } from 'vue';

const props = defineProps({
  data: { type: String },
  toc: { type: Boolean, default: false },
})

const tocList = ref([]);
const tocIndex = ref(-1);

const uslugify = (s) => {
  return uslug(s);
}

const md = MarkDownIt({
  breaks: true,
  html: true,
  linkify: true,
  typographer: true,
  highlight: (code, lang) => {
    if (lang = 'golang') {
      lang = 'go'
    }

    if (prism.languages[lang]) {
      return prism.highlight(code, prism.languages[lang], lang);
    } else {
      return code;
    }
  }
}).use(MarkDownItAnchor, {
  permalink: true,
  permalinkBefore: true,
  permalinkSymbol: '#',
  slugify: uslugify,
}).use(MarkFownDoneRight, {
  slugify: uslugify,
  listType: 'ul',
  callback: function (html, ast) {
    if (tocList.value.length == 0) {
      tocList.value = ast
    }
  }
})

const parse = () => {
  if (!props.data) {
    return ''
  }
  return md.render(props.data)
}

const getHref = (target) => {
  return '#' + uslug(target)
}

</script>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
    display: none;
}
</style>
