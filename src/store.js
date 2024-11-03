import {
	createStore,
} from 'vuex';

import {
	ajaxGet,
	ajaxPost,
} from '@/utils/ajax.js';

const MOBILE_MAX = 900;

export const store = createStore({
	state() {
		return {
			user: null, // null means indeterminate
			langs: [],
			windowWidth: window.innerWidth,
		};
	},
	getters: {
		isMobile(state) {
			return state.windowWidth <= MOBILE_MAX;
		},
		passwordMinLength() {
			return window.appConstants.passwordMinLength;
		},
		treeMaxDepth() {
			return window.appConstants.treeMaxDepth;
		},
		maxLengths() {
			return window.appConstants.maxLengths;
		},
		loginLoaded(state) {
			return state.user !== null;
		},
		userIsAuthenticated(state) {
			return !!state.user;
		},
		currentUserId(state) {
			if (state.user) {
				return state.user.id;
			}
			return null;
		},
		currentUserHandle(state) {
			if (state.user && state.user.handle) {
				return state.user.handle;
			}
			return null;
		},
		currentUserDisplayName(state) {
			if (state.user) {
				return state.user.displayName;
			}
			return '';
		},
		currentUserRole(state) {
			if (state.user) {
				return state.user.role;
			}
			return null;
		},
	},
	mutations: {
		updateWindowWidth(state) {
			state.windowWidth = window.innerWidth;
		},
		setUser(state, user) {
			state.user = user;
		},
	},
	actions: {
		loadLogin(context) {
			context.commit('setUser', null);
			return ajaxGet('/ajax/load-login').then(user => {
				if (user.isAuthenticated) {
					context.commit('setUser', user);
				} else {
					context.commit('setUser', false);
				}
			});
		},
		logout(context) {
			return ajaxPost('/ajax/logout').then(() => {
				context.commit('setUser', false);
			});
		},
		loadLangs(context) {
			return ajaxGet('/ajax/load-langs').then(langs => {
				context.state.langs = langs;
			});
		},
	},
});

window.addEventListener('resize', () => {
	store.commit('updateWindowWidth');
});

export default store;
