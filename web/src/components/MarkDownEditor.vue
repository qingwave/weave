<template>
    <div class="flex flex-row w-full h-full">
        <div class="w-1/2 h-full border rounded">
            <CodeEditor ref="editorRef" class="text-base" height="100%" mode="text/x-markdown" light :value="data" @change="onChange" 
                 @scroll="escroll" @mousewheel="sid=1"></CodeEditor>
        </div>
        <div ref="mdRef" class="w-1/2 border rounded overflow-y-scroll bg-white" @scroll="scroll" @mousewheel="sid=2">
            <MarkDown class="m-4 h-full" :data=content></MarkDown>
        </div>
    </div>
</template>

<script setup>
import MarkDown from './MarkDown.vue'
import CodeEditor from './CodeEditor.vue';
import { ref } from 'vue';

const props = defineProps({
    data: { type: String },
})

const content = ref(props.data);

const onChange= (val, cm) => {
    content.value = val
}

const editor = ref()

const editorRef = ref()
const mdRef = ref()

const sid = ref(0)

const escroll = (cm) => {
    editor.value = cm
    if ( sid.value == 1) {
        let info = cm.getScrollInfo()
        mdRef.value.scrollTop = info.top * mdRef.value.offsetHeight / (info.height - info.clientHeight)
    }
}

const scroll = () => {
    if ( sid.value == 2 ) {
        if (editor.value) {
            editor.value.scrollTo(0, mdRef.value.scrollTop)
        }
    }
}

</script>

<style scoped>
::-webkit-scrollbar {
    display: none;
}

.codemirror-container.bordered {
    border-radius: 0;
    border: none;
}

.codemirror-scroll {
  overflow: visible;
}
</style>
