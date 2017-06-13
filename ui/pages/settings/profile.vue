<template>
	<div class="uk-container">
		<form class="uk-form-stacked" v-on:submit.prevent="saveUser">

			<div class="uk-margin">
				<label class="uk-form-label" for="form-stacked-text">Name</label>
				<div class="uk-form-controls">
					<input class="uk-input" id="form-stacked-text" type="text" required v-model="user.name">
				</div>
			</div>

			<div class="uk-margin">
				<button class="uk-button uk-button-primary" type="submit">Save</button>
			</div>
		</form>
	</div>
</template>

<script>
	import UIkit from 'uikit';

	export default {
		created() {
			this.$store.dispatch('fetchAuthenticatedUser');
		},
		computed: {
			user() {
				const user = this.$store.getters.getAuthUser;
				return Object.assign({}, user);
			}
		},
		methods: {
			saveUser() {
				this.$store.dispatch('updateUser', this.user)
					.then((res) => {
						this.$router.push(`/${this.user.username}`);
						UIkit.notification('User updated', {status: 'success', pos: 'bottom-center'});
					})
					.catch((err) => {
						UIkit.notification('Updating user failed', {status: 'danger', pos: 'bottom-center'})
					})
			}
		}
	}
</script>

<style>
</style>
