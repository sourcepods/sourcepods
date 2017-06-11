<template>
	<div class="uk-navbar-container uk-margin">
		<nav class="uk-container" uk-navbar>

			<div class="uk-position-top-center uk-overlay" v-if="loading">
				<div uk-spinner></div>
			</div>

			<div class="uk-navbar-left">
				<router-link class="uk-navbar-item uk-logo" to="/">
					<img src="/img/logo.svg" alt="" style="height: 46px;"/>
				</router-link>

				<main-nav class="uk-visible@s"></main-nav>

			</div>

			<div class="uk-navbar-right uk-flex-last@s">
				<!--<div class="uk-navbar-item">-->
				<!--<form class="uk-search uk-search-navbar">-->
				<!--<span class="uk-search-icon-flip" uk-search-icon></span>-->
				<!--<input class="uk-search-input" type="search" placeholder="Search...">-->
				<!--</form>-->
				<!--</div>-->
				<ul class="uk-navbar-nav" v-if="user !== null">
					<li>
						<div>
							<router-link :to="`/${user.username}`">
								<gravatar class="uk-border-circle" :email="user.email" :size="46"
										  default-img="mm"></gravatar>
							</router-link>
							<span uk-icon="icon: more-vertical"></span>
						</div>
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

		<nav class="uk-container uk-hidden@s" uk-navbar>
			<div class="uk-navbar-left">
				<main-nav class="uk-hidden@s"></main-nav>
			</div>
		</nav>

	</div>
</template>

<script>
	import Gravatar from 'vue-gravatar';
	import MainNav from './navbar_mainnav.vue';

	export default {
		components: {
			Gravatar,
			MainNav,
		},
		created() {
			this.$store.dispatch('fetchAuthenticatedUser')
				.then((res) => {
				})
				.catch((err) => {
					// TODO: Make this an apollo middleware
					this.$router.push('/login');
				});
		},
		computed: {
			user() {
				return this.$store.getters.getAuthUser;
			},
			loading() {
				return this.$store.state.loading;
			}
		},
	}
</script>
