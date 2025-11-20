import 'element-plus/theme-chalk/dark/css-vars.css'
import './assets/main.scss'
import 'animate.css'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import slideIn from '@/hook/slideIn.ts'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.directive('slide-in', slideIn)

app.mount('#app')
