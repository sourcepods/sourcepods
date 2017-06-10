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
	}
</script>

<style scoped>
</style>
