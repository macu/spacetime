import {
	createRouter,
	createWebHistory,
} from 'vue-router';

import DashboardPage from './pages/dashboard.vue';

export default createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: '/',
			name: 'dashboard',
			component: DashboardPage,
		},
	],
});
