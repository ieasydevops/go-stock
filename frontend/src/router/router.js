import {createMemoryHistory, createRouter, createWebHashHistory, createWebHistory} from 'vue-router'

import stockView from '../components/stock.vue'
import settingsView from '../components/settings.vue'
import about from "../components/about.vue";
import fundView from "../components/fund.vue";
import market from "../components/market.vue";

const routes = [
    { path: '/', component: stockView,name: 'stock'},
    { path: '/fund', component: fundView,name: 'fund' },
    { path: '/settings', component: settingsView,name: 'settings' },
    { path: '/about', component: about,name: 'about' },
    { path: '/market', component: market,name: 'market' },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

export default router