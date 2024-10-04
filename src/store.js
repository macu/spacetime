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
		userIsAuthenticated(state) {
			return !!state.user;
		},
		currentUserDisplayName(state) {
			if (state.user) {
				return state.user.displayName;
			}
			return '';
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
	},
});

window.addEventListener('resize', () => {
	store.commit('updateWindowWidth');
});

export default store;
