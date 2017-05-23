<template>
	<div class="uk-container">

		<h3>Users</h3>

		<ul class="uk-list uk-list-divider">
			<li v-for="user in users">
				<div class="uk-grid">
					<div>
						<gravatar class="uk-border-circle" :email="user.email" :size="46" default-img="mm"></gravatar>
					</div>
					<div class="uk-width-expand">
						<router-link class="uk-link-reset uk-text-bold" :to="`/${user.username}`">
							{{user.name}}
						</router-link>
						<br>
						<span>{{ user.email }}</span>
					</div>
				</div>
			</li>
		</ul>
	</div>
</template>

<script>
	import axios from 'axios';
	import Gravatar from 'vue-gravatar';

	export default {
		components: {
			Gravatar,
		},
		data() {
			return {}
		},
		computed: {
			users() {
				return this.$store.state.users;
			}
		},
		created(){
			axios.get('/api/users')
				.then((res) => {
					this.$store.commit('addUsers', res.data);
				})
				.catch((err) => {
					console.log(err);
				})
		},
	}
</script>

<style>
</style>
