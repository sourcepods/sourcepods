<template>
	<nav class="uk-navbar-container uk-margin" uk-navbar>
		<div class="uk-navbar-left">

			<router-link class="uk-navbar-item uk-logo" to="/">
				<img src="/img/logo.svg" alt="" style="height: 46px;"/>
			</router-link>

			<ul class="uk-navbar-nav">
				<li class="uk-active">
					<router-link to="/pulls">Pull Requests</router-link>
				</li>
				<li>
					<router-link to="/issues">Issues</router-link>
				</li>
			</ul>

		</div>

		<div class="uk-navbar-right">

			<div class="uk-navbar-item">
				<form class="uk-search uk-search-navbar">
					<span class="uk-search-icon-flip" uk-search-icon></span>
					<input class="uk-search-input" type="search" placeholder="Search...">
				</form>
			</div>

			<ul class="uk-navbar-nav" v-if="user !== null">
				<li>
					<router-link :to="`/${user.username}`">
						<gravatar class="uk-border-circle" :email="user.email" :size="46" default-img="mm"></gravatar>
						<span uk-icon="icon: more-vertical"></span>
					</router-link>
					<div class="uk-navbar-dropdown">
						<ul class="uk-nav uk-navbar-dropdown-nav">
							<li>
								<router-link :to="`/${user.username}`">Profile</router-link>
							</li>
							<li>
								<router-link to="/settings/profile">Settings</router-link>
							</li>
							<li class="uk-nav-divider"></li>
							<li>
								<router-link to="/login">Sign out</router-link>
							</li>
						</ul>
					</div>
				</li>
			</ul>
		</div>
	</nav>
</template>

<script>
	import Gravatar from 'vue-gravatar';

	export default {
		components: {
			Gravatar,
		},
		computed: {
			user() {
				return this.$store.state.user;
			}
		},
	}
</script>

<style scoped>
</style>
