<template>
	<div>
		<navbar></navbar>
		<router-view></router-view>
	</div>
</template>

<script>
	import Navbar from './navbar.vue';

	export default {
		components: {
			Navbar,
		},
		created() {
			this.$store.dispatch('fetchAuthenticatedUser')
				.then((res) => {
				})
				.catch((err) => {
					this.$router.push('/login');
				})
		}
	}
</script>

<style scoped>
</style>
