<template>
    <div class="flex flex-row w-full h-auto overflow-x-hidden">
        <div class="w-1/2 max-h-full border rounded">
            <CodeEditor v-model="content" class="text-base max-w-full" height="100%" mode="text/x-markdown" light :value="content"
                @change="onChange"></CodeEditor>
        </div>
        <div class="w-1/2 border rounded max-h-full">
            <MarkDown ref="mdRef" class="p-4" :data=content></MarkDown>
        </div>
    </div>
</template>

<script setup>
import MarkDown from './MarkDown.vue'
import CodeEditor from './CodeEditor.vue';
import { ref, computed, watch } from 'vue';

const props = defineProps({
    data: { type: String },
})

const content = ref(props.data);

watch(() => props.data, 
    (val) => {
        content.value = val
    }
)

const onChange = (v, cm) => {
    content.value = v
}

defineExpose({
    content,
})

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
