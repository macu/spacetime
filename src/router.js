import {
	createRouter,
	createWebHistory,
} from 'vue-router';

import DashboardPage from './pages/dashboard.vue';
import LoginPage from './pages/user-account/login.vue';
import SignupPage from './pages/user-account/signup.vue';
import SignupVerifyPage from './pages/user-account/signup-verify.vue';

import NodeView from './pages/nodes/view.vue';
import CreateCategory from './pages/nodes/create-category.vue';
import CreatePost from './pages/nodes/create-post.vue';

export default createRouter({
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
			path: '/node/:id',
			name: 'node-view',
			component: NodeView,
		},
		{
			path: '/create/category',
			name: 'create-category',
			component: CreateCategory,
		},
		{
			path: '/create/post',
			name: 'create-post',
			component: CreatePost,
		},
	],
});
