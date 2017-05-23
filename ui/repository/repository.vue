<template>
	<div>
		<div class="repository-nav">
			<div class="uk-container">
				<h3>
					<router-link :to="`/${owner}`">{{owner}}</router-link>
					<span>/</span>
					<router-link :to="`/${owner}/${repository}`">{{repository}}</router-link>
				</h3>
				<span>git flying in the cloud</span>

				<ul class="uk-child-width-expand" uk-tab>
					<li><a href="#">Project</a></li>
					<li><a href="#">Repository</a></li>
					<li><a href="#">
						Issues
						<span class="uk-badge">3</span>
					</a></li>
					<li><a href="#">
						Pull Requests
						<span class="uk-badge">0</span>
					</a></li>
					<li><a href="#">Pipelines</a></li>
					<li><a href="#">Settings</a></li>
				</ul>
			</div>
		</div>
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
	.repository-nav {
		margin-top: -20px;
		margin-bottom: 20px;
		padding-top: 20px;
		background-color: #f8f8f8;
	}

	ul.uk-child-width-expand a {
		line-height: 28px;
	}

	ul.uk-child-width-expand span.uk-badge {
		height: 15px;
		width: 15px;
		min-width: 15px;
		line-height: 15px;
		font-size: 0.6rem;
		margin-left: 5px;
		margin-bottom: 10px;
	}
</style>
