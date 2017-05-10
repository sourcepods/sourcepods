<template>
	<div class="uk-container">

		<div class="uk-grid-match uk-grid-small" uk-grid>
			<div class="uk-width-2-5@m uk-width-1-4@l">
				<gravatar :email="user.email" default-img="mm" size="256"></gravatar>
				<h2>{{user.name}}</h2>
				<h4>{{user.username}}</h4>

				<hr>

				<a :href="`mailto:${user.email}`">{{user.email}}</a>

			</div>
			<div class="uk-width-3-5@m uk-width-3-4@l">
				<div class="uk-card uk-card-default uk-card-body">
					<table>
						<tr>
							<td>id</td>
							<td>{{user.id}}</td>
						</tr>
						<tr>
							<td>name</td>
							<td>{{user.name}}</td>
						</tr>
						<tr>
							<td>email</td>
							<td>{{user.email}}</td>
						</tr>
					</table>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
	import Gravatar from 'vue-gravatar';

	export default {
		components: {
			gravatar: Gravatar,
		},
		created() {
			this.$store.dispatch('fetchUser', this.$route.params.username);
		},
		computed: {
			user() {
				return this.$store.getters.getUser(this.$route.params.username);
			}
		},
	}
</script>

<style>
	h2, h4 {
		margin-top: 10px;
		margin-bottom: 10px;
	}
</style>
