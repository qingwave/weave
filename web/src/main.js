import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import './styles/index.css'
import 'element-plus/es/components/message/style/css'
import 'element-plus/es/components/notification/style/css'

const app = createApp(App).use(router).use(createPinia());
app.mount('#app');
