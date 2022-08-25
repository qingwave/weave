<template>
    <div class="flex w-full justify-center">
        <div class="max-w-fit">
            <header class="py-16 sm:text-center">
                <h1 class="mb-4 text-3xl sm:text-4xl tracking-tight text-slate-900 font-extrabold dark:text-slate-200">
                    Latest Updates</h1>
                <p class="text-lg text-slate-700 dark:text-slate-400">All the latest blogs from <a target="__blank" href="https://qingwave.github.io">QingWave</a>.
                  &nbsp;&nbsp;<el-button type="primary"><router-link to="/posts/editor" target="_blank">New Post</router-link></el-button>
                  </p>
            </header>
            <div
                class="relative sm:pb-12 sm:ml-[calc(2rem+1px)] md:ml-[calc(3.5rem+1px)] lg:ml-[max(calc(14.5rem+1px),calc(100%-48rem))]">
                <div
                    class="hidden absolute top-3 bottom-0 right-full mr-7 md:mr-[3.25rem] w-px bg-slate-200 dark:bg-slate-800 sm:block">
                </div>
                <div class="space-y-16">
                    <article class="relative group" v-for="p in posts">
                        <div
                            class="absolute -inset-y-2.5 -inset-x-4 md:-inset-y-4 md:-inset-x-6 sm:rounded-2xl group-hover:bg-slate-50/70 dark:group-hover:bg-slate-800/50">
                        </div><svg viewBox="0 0 9 9"
                            class="hidden absolute right-full mr-6 top-2 text-slate-200 dark:text-slate-600 md:mr-12 w-[calc(0.5rem+1px)] h-[calc(0.5rem+1px)] overflow-visible sm:block">
                            <circle cx="4.5" cy="4.5" r="4.5" stroke="currentColor"
                                class="fill-white dark:fill-slate-900" stroke-width="2"></circle>
                        </svg>
                        <div class="relative">
                            <h3
                                class="text-base font-semibold tracking-tight text-slate-900 dark:text-slate-200 pt-8 lg:pt-0">
                                {{ p.name }}</h3>
                            <div
                                class="mt-2 mb-4 prose prose-slate prose-a:relative prose-a:z-10 dark:prose-dark line-clamp-2">
                                <p> {{ p.summary }}</p>
                            </div>
                            <dl class="absolute left-0 top-0 lg:left-auto lg:right-full lg:mr-[calc(6.5rem+1px)]">
                                <dt class="sr-only">Date</dt>
                                <dd class="whitespace-nowrap text-sm leading-6 dark:text-slate-400"><time
                                        :datetime="p.createdAt">{{ timeFormat(p.createdAt) }}</time></dd>
                            </dl>
                        </div>
                        <router-link class="flex items-center text-sm text-emerald-400 font-medium"
                            :to="`/posts/${p.id}`"><span
                                class="absolute -inset-y-2.5 -inset-x-4 md:-inset-y-4 md:-inset-x-6 sm:rounded-2xl"></span><span
                                class="relative">Read more</span><svg
                                class="relative mt-px overflow-visible ml-2.5 text-emerald-200 dark:text-sky-700"
                                width="3" height="6" viewBox="0 0 3 6" fill="none" stroke="currentColor"
                                stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                <path d="M0 0L3 3L0 6"></path>
                            </svg></router-link>
                    </article>

                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import request from '@/axios';

const posts = ref([]);

onMounted(() => {
    request.get('/api/v1/posts').then(
        (response) => {
            posts.value = response.data.data
        }
    )
})

const timeFormat = (date) => {
    return (new Date(date)).toLocaleDateString('en', { year: "numeric", month: "short", day: "numeric" })
}
</script>
