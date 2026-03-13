import './app.css'
import 'vfonts/Lato.css'
import 'vfonts/FiraCode.css'


import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import { router } from './router'

const isProduction = import.meta.env.MODE === 'production'
if (isProduction) {
    window.console = {
        ...window.console,
        log: (l) => null,
        info: (l) => null,
    }
}



const app = createApp(App)

app.use(createPinia())
app.use(router)
router.isReady().then(() => {
    app.mount('#app')
})
