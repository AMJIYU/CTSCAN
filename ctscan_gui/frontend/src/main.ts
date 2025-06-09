import {createApp} from 'vue'
import App from './App.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css' // 引入样式
import * as ElementPlusIconsVue from '@element-plus/icons-vue' // 引入图标

const app = createApp(App)

// 全局注册 Element Plus 组件和图标
app.use(ElementPlus)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.mount('#app')
