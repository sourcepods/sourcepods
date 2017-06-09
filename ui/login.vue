<template>
	<div class="uk-flex uk-flex-center">
		<div class="uk-card uk-card-default uk-card-hover uk-card-body">

			<img src="/img/logo.svg" class="uk-align-center logo"/>

			<h3 class="uk-heading-line"><span>Sign in to GitPod</span></h3>

			<!-- TODO: set failed = false on @hide -->
			<div class="uk-alert-danger" uk-alert v-if="failed">
				<a class="uk-alert-close" uk-close></a>
				<p>Incorrect username or password.</p>
			</div>

			<form class="uk-form-stacked" v-on:submit.prevent="login">
				<div class="uk-margin">
					<label class="uk-form-label" for="email">Email</label>
					<div class="uk-form-controls">
						<input class="uk-input" id="email" type="email"
							   autofocus="autofocus" autocapitalize="off" autocorrect="off" required v-model="email">
					</div>
				</div>

				<div class="uk-margin">
					<label class="uk-form-label" for="password">Password</label>
					<div class="uk-form-controls">
						<input class="uk-input" type="password" id="password" required v-model="password"/>
					</div>
				</div>

				<div class="uk-margin">
					<button class="uk-button uk-button-primary uk-width-1-1 uk-margin-small-bottom">Sign in</button>
				</div>

				<div class="uk-text-right">
					<a href="">Forgot your password?</a>
				</div>
			</form>

		</div>
	</div>

</template>

<script>
	import axios from 'axios';

	export default {
		data() {
			return {
				email: null,
				password: null,
				failed: false,
			}
		},
		methods: {
			login() {
				this.failed = false; // reset the alert for a new try
				this.$store.dispatch('authenticateUser', {email: this.email, password: this.password})
					.then((user) => {
						// We actually need to reload the page, so we can use the cookie.
						window.location.replace('/');
					})
					.catch((err) => {
						this.failed = true;
					})
			},
		}
	}
</script>

<style scoped>
	.uk-card {
		width: 380px;
	}

	.logo {
		width: 64px;
	}
</style>
