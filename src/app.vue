<template>
<div class="treetime-app">

	<header class="flex-row-md">
		<h1 class="flex-1">TreeTime</h1>
		<template v-if="showUser">
			<span class="flex-row nowrap">
				<material-icon icon="account_circle"/>
				<span v-text="currentUserDisplayName"/>
			</span>
			<el-button @click="logout()">Log out</el-button>
		</template>
		<template v-else-if="showLogin">
			<router-link :to="{name: 'login'}">Log in</router-link>
			<router-link :to="{name: 'signup'}">Register</router-link>
		</template>
	</header>

	<div class="body">
		<router-view/>
	</div>

</div>
</template>

<script>
export default {
	computed: {
		loginLoaded() {
			return this.$store.state.user !== null;
		},
		showUser() {
			return this.loginLoaded && this.$store.getters.userIsAuthenticated;
		},
		showLogin() {
			return this.loginLoaded && !this.$store.getters.userIsAuthenticated;
		},
		currentUserDisplayName() {
			return this.$store.getters.currentUserDisplayName;
		},
	},
	mounted() {
		// Load current user
		this.$store.dispatch('loadLogin');
	},
	methods: {
		logout() {
			this.$confirm('Are you sure you want to log out?', 'Log out', {
				confirmButtonText: 'Log out',
				cancelButtonText: 'Cancel',
				type: 'warning',
			}).then(() => {
				this.$store.dispatch('logout');
			}).catch(() => {
				// Do nothing
			});
		},
	},
};
</script>

<style lang="scss">
.treetime-app {
	>header {
		margin: 0;
		padding: 20px;
		background-color: #f0f0f0;
		border-bottom: 1px solid #ccc;
		>h1 {
			font-size: 2em;
		}
	}
	>.body {
		padding: 40px;
	}
}
</style>
