<template>
	<div class="uk-navbar-container uk-margin">
		<nav class="uk-container" uk-navbar>
			<div class="uk-navbar-left">

				<router-link class="uk-navbar-item uk-logo" to="/">
					<img src="/img/logo.svg" alt="" style="height: 46px;"/>
				</router-link>

				<ul class="uk-navbar-nav">
					<li>
						<router-link to="/issues">
							Issues
							<!--<span class="uk-badge">12</span>-->
						</router-link>
					</li>
					<li>
						<router-link to="/pulls">
							Pull Requests
							<!--<span class="uk-badge">1</span>-->
						</router-link>
					</li>
				</ul>

				<div uk-spinner v-if="loading"></div>

			</div>

			<div class="uk-navbar-right">
				<!--<div class="uk-navbar-item">-->
				<!--<form class="uk-search uk-search-navbar">-->
				<!--<span class="uk-search-icon-flip" uk-search-icon></span>-->
				<!--<input class="uk-search-input" type="search" placeholder="Search...">-->
				<!--</form>-->
				<!--</div>-->
				<ul class="uk-navbar-nav" v-if="user !== null">
					<li>
						<router-link :to="`/${user.username}`">
							<gravatar class="uk-border-circle" :email="user.email" :size="46"
									  default-img="mm"></gravatar>
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
	</div>
</template>

<script>
	import Gravatar from 'vue-gravatar';

	export default {
		components: {
			Gravatar,
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

<style scoped>
	ul.uk-navbar-nav span.uk-badge {
		height: 15px;
		width: 15px;
		min-width: 15px;
		line-height: 15px;
		font-size: 0.6rem;
		margin-left: 5px;
		margin-bottom: 10px;
	}
</style>
