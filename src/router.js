import {
	createRouter,
	createWebHistory,
} from 'vue-router';

import DashboardPage from '@/pages/dashboard.vue';
import LoginPage from '@/pages/user-account/login.vue';
import SignupPage from '@/pages/user-account/signup.vue';
import SignupVerifyPage from '@/pages/user-account/signup-verify.vue';

import LangsView from '@/pages/languages.vue';
import TagsView from '@/pages/tags.vue';
import NodeView from '@/pages/nodes/view.vue';
import CreateTag from '@/pages/nodes/create-tag.vue';
import CreateCategory from '@/pages/nodes/create-category.vue';
import CreatePost from '@/pages/nodes/create-post.vue';
import CreateComment from '@/pages/nodes/create-comment.vue';

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
			path: '/langs',
			name: 'langs-view',
			component: LangsView,
		},
		{
			path: '/tags',
			name: 'tags-view',
			component: TagsView,
		},
		{
			path: '/node/:id',
			name: 'node-view',
			component: NodeView,
		},
		{
			path: '/create/tag',
			name: 'create-tag',
			component: CreateTag,
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
		{
			path: '/create/comment',
			name: 'create-comment',
			component: CreateComment,
		},
	],
});
