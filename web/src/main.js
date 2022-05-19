import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import 'virtual:windi.css'
import 'element-plus/es/components/message/style/css'
import 'element-plus/es/components/notification/style/css'

const app = createApp(App).use(router);
app.mount('#app');
