import {
	createStore,
} from 'vuex';

const MOBILE_MAX = 900;

export const store = createStore({
	state() {
		return {
			user: null,
			windowWidth: window.innerWidth,
		};
	},
	getters: {
		isAuthenticated(state) {
			return !!state.user;
		},
		isMobile(state) {
			return state.windowWidth <= MOBILE_MAX;
		},
	},
	mutations: {
		updateWindowWidth(state) {
			state.windowWidth = window.innerWidth;
		},
	},
});

window.addEventListener('resize', () => {
	store.commit('updateWindowWidth');
});

export default store;
