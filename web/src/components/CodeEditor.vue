<template>
  <Codemirror :value="props.value" :options="cmOptions" border placeholder="yaml code here" :height="props.height"
    @change="onChange" />
</template>

<script setup>
import Codemirror from "codemirror-editor-vue3";
import "codemirror/mode/yaml/yaml.js";
import "codemirror/theme/monokai.css";
import "codemirror/addon/fold/brace-fold.js";
import "codemirror/addon/fold/foldcode.js";
import "codemirror/addon/fold/comment-fold.js";
import "codemirror/addon/fold/indent-fold.js";
import "codemirror/addon/fold/foldgutter.js";
import "codemirror/addon/fold/foldgutter.css";
import "codemirror/addon/lint/lint.js";
import "codemirror/addon/lint/lint.css";
import "codemirror/addon/lint/yaml-lint.js";
import yaml from 'js-yaml';
import { ref } from 'vue';

window.jsyaml = yaml;

const props = defineProps({
  value: { type: String },
  height: { type: String, default: "60vh" },
  mode: { type: String, default: "text/yaml" },
  light: { type: Boolean, default: false },
  readOnly: { type: Boolean, default: false },
})

const emit = defineEmits(['change']);

const onChange = (val, cm) => {
  emit("change", val, cm)
}

const cmOptions = ref({
  mode: props.mode,
  theme: props.light? "": "monokai",
  lineNumbers: true,
  smartIndent: true,
  indentUnit: 2,
  foldGutter: true,
  lint: true,
  styleActiveLine: true,
  readOnly: props.readOnly,
  gutters: ['CodeMirror-linenumbers', 'CodeMirror-foldgutter', 'CodeMirror-lint-markers']
})

</script>
