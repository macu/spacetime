<template>
<form-layout class="signup-verify-page page-width-sm">

	<template #title>Verify sign up</template>

	<loading-message v-if="loading" message="Loading signup request..."/>

	<template v-else-if="signupRequest">

		<p>Finish setting up your account below.</p>

		<form-field title="Email address" required>
			<el-input v-model="signupRequest.email" type="email"
				readonly autocomplete="username"/>
		</form-field>

		<form-field title="Handle">
			<p>This handle will make your user profile accessible via URL.</p>
			<p>Format: letters A-Z, numbers, and underscores, with no spaces.</p>
			<el-input v-model="handle" type="text"
				:maxlength="25" show-word-limit
				autocapitalize="none" autocomplete="off"/>
		</form-field>

		<form-field title="Display name" required>
			<label></label>
			<p>This name will be displayed on your user profile and posts you create.</p>
			<el-input v-model="displayName" type="text"
				:maxlength="50" show-word-limit
				autocapitalize="none" autocomplete="off"/>
		</form-field>

		<form-field title="Password">
			<p>Minimum length is {{$store.getters.passwordMinLength}} characters.</p>
			<el-input v-model="password" type="password"
				:maxlength="100" show-password
				autocapitalize="none" autocomplete="new-password"/>
		</form-field>

		<form-field title="Verify password">
			<el-input v-model="verifyPassword" type="password"
				:maxlength="100" show-password
				autocapitalize="none" autocomplete="new-password"/>
		</form-field>

		<form-field title="Message to admin (optional; say Hi)">
			<el-input type="textarea" :maxlength="200"
				:autosize="{minRows: 2}" show-word-limit
				autocapitalize="none" autocomplete="off"/>
		</form-field>

		<form-actions>
			<el-button @click="create()" :disabled="createDisabled">Create account</el-button>
		</form-actions>

		<loading-message v-if="submitting" message="Creating account..."/>

	</template>

	<el-alert v-else title="Invalid request" type="error" show-icon :closable="false">
		<p>The signup request is invalid or expired.</p>
		<p>Please <router-link :to="{name: 'signup'}">register again</router-link>.</p>
	</el-alert>

</form-layout>
</template>

<script>
import store from '@/store.js';

import {
	ajaxGet,
	ajaxPost,
} from '@/utils/ajax.js';

const handlePattern = /^[a-zA-Z0-9_]+$/;

export default {
	data() {
		return {
			loading: true,
			signupRequest: null,

			password: '',
			verifyPassword: '',
			handle: '',
			displayName: '',
			submitting: false,
		};
	},
	computed: {
		createDisabled() {
			return this.loading || !this.signupRequest ||
				this.submitting || !this.displayName.trim() ||
				(!!this.handle.trim() && !handlePattern.test(this.handle.trim())) ||
				!this.password.trim() ||
				this.password.length < this.$store.getters.passwordMinLength ||
				this.password !== this.verifyPassword;
		},
	},
	beforeRouteEnter(to, from, next) {
		if (store.getters.userIsAuthenticated) {
			next({name: 'dashboard'});
		} else {
			next(vm => {
				if (vm.$route.query.token) {
					vm.loadSignupRequest(vm.$route.query.token);
				} else {
					vm.loading = false;
				}
			});
		}
	},
	methods: {
		loadSignupRequest(token) {
			ajaxGet('/ajax/load-signup', {
				token,
			}, {
				'invalid-token': 'Invalid token.',
				'token-expired': 'This signup request has expired.',
			}).then(response => {
				this.signupRequest = response;
				this.loading = false;
			}).catch(() => {
				this.loading = false;
			});
		},
		create() {
			if (this.createDisabled) {
				return;
			}
			this.submitting = true;
			ajaxPost('/ajax/signup-verify', {
				token: this.signupRequest.token,
				password: this.password,
				handle: this.handle.trim(),
				displayName: this.displayName.trim(),
				message: this.message,
			}, {
				'invalid-token': 'Invalid token.',
				'token-expired': 'This signup request has expired.',
				'invalid-handle': 'The given handle is invalid.',
				'handle-exists': 'The given handle already exists.',
				'email-exists': 'A user with the given email address already exists.',
			}).then(() => {
				this.$store.dispatch('loadLogin');
				this.$router.replace({
					name: 'dashboard',
				});
			}).catch(() => {
				this.submitting = false;
			});
		},
	},
};
</script>
