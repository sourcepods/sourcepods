<template>
	<div>
		<nav class="uk-navbar-container uk-margin" uk-navbar>
			<div class="uk-navbar-left">

				<router-link class="uk-navbar-item uk-logo" to="/">
					<img src="/img/logo.svg" alt="" style="height: 46px;">
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

				<ul class="uk-navbar-nav">
					<li>
						<router-link to="/users/metalmatze">
							<img src="https://www.gravatar.com/avatar/915d80f0d9b6678fad4d1ab36dfc8960?s=46&d=mm" alt=""
								 class="uk-border-circle">
						</router-link>
					</li>
				</ul>

			</div>
		</nav>

		<router-view></router-view>
	</div>
</template>

<script>
	export default {}
</script>

<style scoped>
</style>
