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
			loading: false,
			user: null, // null means indeterminate
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
		setLoading(state, loading) {
			state.loading = loading;
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
	},
});

window.addEventListener('resize', () => {
	store.commit('updateWindowWidth');
});

export default store;
