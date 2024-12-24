<template>
<form-layout class="signup-page page-width-sm">

	<template #title>Sign up</template>

	<template v-if="submitted">

		<p>Your request has been submitted. Please check your email for a verification link, to finish creating your account.</p>

	</template>

	<template v-else>

		<p>Please provide a valid email address to create your account.</p>

		<form-field title="Email address">
			<el-input v-model="email" type="email"
				:maxlength="50" show-word-limit
				autocapitalize="none" autocomplete="email"/>
		</form-field>

		<form-actions>
			<el-button @click="register()" :disabled="registerDisabled">Register</el-button>
			<router-link :to="{name: 'login'}">Log in</router-link>
		</form-actions>

		<loading-message v-if="loading" message="Submitting..."/>

	</template>

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
			loading: false,
			submitted: false,
		};
	},
	computed: {
		registerDisabled() {
			return this.loading || this.submitted || !this.email.trim();
		},
	},
	beforeRouteEnter(to, from, next) {
		if (store.getters.authenticated) {
			next('/');
		} else {
			next();
		}
	},
	methods: {
		register() {
			if (this.registerDisabled) {
				return;
			}
			this.loading = true;
			ajaxPost('/ajax/signup', {
				email: this.email.trim(),
			}, {
				'invalid-email': 'The given email address is invalid.',
				'email-exists': 'An account with this email address already exists.',
			}).then(response => {
				if (typeof response === 'object' && response.token) {
					this.$router.replace({
						name: 'signup-verify',
						query: {
							token: response.token,
						},
					});
				} else {
					this.submitted = true;
					this.loading = false;
				}
			}).catch(() => {
				this.loading = false;
			});
		},
	},
};
</script>
