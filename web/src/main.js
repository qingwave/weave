import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import 'virtual:windi.css'

const app = createApp(App).use(router);
app.mount('#app');
