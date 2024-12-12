import {
	createRouter,
	createWebHistory,
} from 'vue-router';

import store from './store.js';

import DashboardPage from '@/pages/dashboard.vue';
import LoginPage from '@/pages/user-account/login.vue';
import SignupPage from '@/pages/user-account/signup.vue';
import SignupVerifyPage from '@/pages/user-account/signup-verify.vue';

import SpacePage from './pages/space.vue';

const router = createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: '/',
			name: 'dashboard',
			component: DashboardPage,
		},
		{
			path: '/login',
			name: 'login',
			component: LoginPage,
		},
		{
			path: '/signup',
			name: 'signup',
			component: SignupPage,
		},
		{
			path: '/verify-signup',
			name: 'signup-verify',
			component: SignupVerifyPage,
		},
		{
			path: '/space/:spaceId',
			name: 'space',
			component: SpacePage,
		},
	],
});

router.beforeEach((to, from, next) => {
	store.commit('setLoading', true);
	next();
});

router.afterEach(() => {
	store.commit('setLoading', false);
});

export default router;
