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
import CreateEmptySpacePage from './pages/create-space/empty.vue';
import CreateTitleSpacePage from './pages/create-space/title.vue';
import CreateTagSpacePage from './pages/create-space/tag.vue';
import CreateTextSpacePage from './pages/create-space/text.vue';

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
		{
			path: '/space/create/empty',
			name: 'create-empty-space',
			component: CreateEmptySpacePage,
		},
		{
			path: '/space/create/title',
			name: 'create-title',
			component: CreateTitleSpacePage,
		},
		{
			path: '/space/create/tag',
			name: 'create-tag',
			component: CreateTagSpacePage,
		},
		{
			path: '/space/create/text',
			name: 'create-text',
			component: CreateTextSpacePage,
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
