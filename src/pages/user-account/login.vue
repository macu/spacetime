<template>
<form-layout class="login-page page-width-sm">

	<template #title>Log in</template>

	<form-field title="Email address">
		<el-input v-model="email" type="email" :maxlength="50"
			autocapitalize="none" autocomplete="username"
			@keyup.enter.native="login()"
			/>
	</form-field>

	<form-field title="Password">
		<el-input v-model="password" type="password" :maxlength="100"
			autocapitalize="none" autocomplete="current-password"
			@keyup.enter.native="login()"
			/>
	</form-field>

	<form-actions>
		<el-button @click="login()" :disabled="loginDisabled">Log in</el-button>
		<router-link :to="{name: 'signup'}">Register</router-link>
	</form-actions>

	<loading-message v-if="loading" message="Logging in..."/>

</form-layout>
</template>

<script>
import store from '@/store.js';

import {
	ajaxPost,
} from '@/utils/ajax.js';

export default {
	data() {
		return {
			email: '',
			password: '',
			loading: false,
		};
	},
	computed: {
		loginDisabled() {
			return this.loading || !this.email.trim() || !this.password.trim();
		},
	},
	beforeRouteEnter(to, from, next) {
		if (store.getters.userIsAuthenticated) {
			next({name: 'dashboard'});
		} else {
			next();
		}
	},
	methods: {
		login() {
			if (this.loginDisabled) {
				return;
			}
			this.loading = true;
			ajaxPost('/ajax/login', {
				email: this.email.trim(),
				password: this.password,
			}, {
				403: 'Invalid email or password.',
			}).then(() => {
				this.$store.dispatch('loadLogin');
				this.$router.replace({
					name: 'dashboard',
				});
			}).catch(() => {
				this.loading = false;
			});
		},
	},
};
</script>
