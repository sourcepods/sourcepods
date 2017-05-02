<template>
	<div class="uk-container">
		<form v-on:submit.prevent="saveUser">
			<table>
				<tr>
					<td>id</td>
					<td>
						<input class="uk-input" type="text" v-model="user.id"/>
					</td>
				</tr>
				<tr>
					<td>username</td>
					<td>
						<input class="uk-input" type="text" v-model="user.username"/>
					</td>
				</tr>
				<tr>
					<td>name</td>
					<td>
						<input class="uk-input" type="text" v-model="user.name"/>
					</td>
				</tr>
				<tr>
					<td>email</td>
					<td>
						<input class="uk-input" type="email" v-model="user.email"/>
					</td>
				</tr>
				<tr>
					<td colspan="2">
						<button class="uk-button uk-button-primary" type="submit">Save</button>
					</td>
				</tr>
			</table>
		</form>
	</div>
</template>

<script>
	import UIkit from 'uikit';

	export default {
		data(){
			return {}
		},
		created() {
			this.$store.dispatch('fetchUser', this.$route.params.username);
		},
		computed: {
			user() {
				const user = this.$store.getters.getUser(this.$route.params.username);
				return Object.assign({}, user);
			}
		},
		methods: {
			saveUser() {
				this.$store.dispatch('updateUser', this.user)
					.then((user) => {
						this.$router.push(`/users/${user.username}`);
						UIkit.notification('User updated', {status: 'success', pos: 'bottom-center'});
					})
					.catch((err) => {
						UIkit.notification('Updating user failed', {status: 'danger', pos:'bottom-center'})
					})
			}
		}
	}
</script>

<style>
</style>
