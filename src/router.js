import {
	createRouter,
	createWebHistory,
} from 'vue-router';

import store from './store.js';

import DashboardPage from '@/pages/dashboard.vue';
import LoginPage from '@/pages/user-account/login.vue';
import SignupPage from '@/pages/user-account/signup.vue';
import SignupVerifyPage from '@/pages/user-account/signup-verify.vue';

import UserPage from './pages/user.vue';
import BookmarksPage from '@/pages/bookmarks.vue';

import SpacePage from './pages/space/index.vue';
import CreateBranchSpacePage from './pages/create-space/branch.vue';
import CreateTagSpacePage from './pages/create-space/tag.vue';
import CreateTextSpacePage from './pages/create-space/text.vue';
import CreateLinkSpacePage from './pages/create-space/link.vue';

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
			path: '/user/:id',
			name: 'user',
			component: UserPage,
		},
		{
			path: '/bookmarks',
			name: 'bookmarks',
			component: BookmarksPage,
		},
		{
			path: '/space/:spaceId',
			name: 'space',
			component: SpacePage,
		},
		{
			path: '/space/create/branch',
			name: 'create-branch',
			component: CreateBranchSpacePage,
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
		{
			path: '/space/create/link',
			name: 'create-space-link',
			component: CreateLinkSpacePage,
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
