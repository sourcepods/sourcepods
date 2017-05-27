<template>
	<div class="uk-container">

		<h3>Users</h3>

		<ul class="uk-list uk-list-divider">
			<li v-for="user in users">
				<div class="uk-grid">
					<div>
						<gravatar class="uk-border-circle" :email="user.attributes.email" :size="46"
								  default-img="mm"></gravatar>
					</div>
					<div class="uk-width-expand">
						<router-link class="uk-link-reset uk-text-bold" :to="`/${user.attributes.username}`">
							{{user.attributes.name}}
						</router-link>
						<br>
						<span>{{ user.attributes.email }}</span>
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
				let users = this.$store.getters.getUsers;
				// slice copies the array to not modify the one in vuex
				return users.slice().sort((a, b) => a.attributes.name > b.attributes.name)
			}
		},
		created(){
			this.$store.dispatch('fetchUsers');
		},
	}
</script>

<style>
</style>
