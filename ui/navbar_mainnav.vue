<template>
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
</template>

<style scoped>
	span.uk-badge {
		height: 15px;
		width: 15px;
		min-width: 15px;
		line-height: 15px;
		font-size: 0.6rem;
		margin-left: 5px;
		margin-bottom: 10px;
	}
</style>
