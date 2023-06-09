import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import HomeView from '../views/HomeView.vue';

const routes: Array<RouteRecordRaw> = [
    {
        path: '/',
        name: 'home',
        component: HomeView,
    },
    {
        path: '/new-game',
        name: 'new-game',
        component: () => import('../views/NewGameView.vue'),
    },
    {
        path: '/game',
        name: 'game',
        component: () => import('../views/GameView.vue'),
        props: true,
    },
];

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes,
});

export default router;
