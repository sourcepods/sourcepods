<template>
	<div class="uk-container">
		<h3>
			<router-link :to="`/${owner}`">{{owner}}</router-link>
			<span>/</span>
			<router-link :to="`/${owner}/${repository}`">{{repository}}</router-link>
		</h3>
	</div>
</template>

<script>
	export default {
		computed: {
			owner() {
				return this.$route.params.owner;
			},
			repository() {
				return this.$route.params.repository;
			},
		},
	}
</script>

<style scoped>
</style>
