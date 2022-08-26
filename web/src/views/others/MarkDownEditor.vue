<template>
    <div class="flex w-full h-full">
        <div class="w-full m-2">
            <MarkDownEditor class="w-full h-full" :data=md ></MarkDownEditor>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
import MarkDownEditor from 'components/MarkDownEditor.vue';

const md = ref(`
# Header 1
## Header 2
### Header 3

This is a **markdown** test *demo*.

* xxx
* aaa 

- first item
- second item
- third item

> quote message here <br >
> quote line two

Code demo
\`\`\`go
func main() {
    a, b := 2, 3
    fmt.Println(a+b)
}
\`\`\`

Link example
[qingwave.github.io/](https://qingwave.github.io/)

`)

</script>