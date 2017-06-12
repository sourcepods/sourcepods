<template>
	<div v-if="loading || repository === null"></div>
	<div v-else>
		<div class="repository-nav">
			<div class="uk-container">
				<div uk-grid>
					<div>
						<h3><span uk-icon="icon: lock"></span></h3>
					</div>

					<div style="padding-left: 16px;">
						<h3>
							<router-link :to="`/${owner_name}`">{{owner_name}}</router-link>
							<span>/</span>
							<router-link :to="`/${owner_name}/${repository.name}`">
								{{repository.name}}
							</router-link>
						</h3>
						<span>{{repository.description}}</span>
					</div>
				</div>

				<ul class="uk-child-width-expand" uk-tab>
					<li>
						<a href="#">Project</a>
					</li>
					<li>
						<a href="#">Repository</a>
					</li>
					<li>
						<a href="#">
							Issues
							<span class="uk-badge">{{repository.issue_stats.open}}</span>
						</a>
					</li>
					<li>
						<a href="#">
							Pull Requests
							<span class="uk-badge">{{repository.pull_request_stats.open}}</span>
						</a>
					</li>
					<li>
						<a href="#">Pipelines</a>
					</li>
					<li>
						<a href="#">Settings</a>
					</li>
				</ul>
			</div>
		</div>
	</div>
</template>

<script>
	export default {
		data() {
			return {
				loading: true,
				repository_id: '',
			}
		},
		created() {
			this.setLoading(true);

			const payload = {
				owner: this.$route.params.owner,
				repository: this.$route.params.repository,
			};

			this.$store.dispatch('fetchRepository', payload)
				.then((repository) => {
					this.repository_id = repository.id;
					this.setLoading(false);
				});
		},
		computed: {
			owner_name() {
				return this.$route.params.owner;
			},
			repository() {
				return this.$store.getters.getRepository(this.repository_id);
			},
		},
		methods: {
			setLoading(isLoading) {
				this.loading = isLoading;
				this.$store.commit('loading', isLoading)
			}
		},
	}
</script>

<style scoped>
	.repository-nav {
		margin-top: -20px;
		margin-bottom: 20px;
		background-color: #f8f8f8;
	}

	.repository-nav h3 {
		margin: 0;
	}

	ul.uk-child-width-expand a {
		line-height: 28px;
		text-transform: none;
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
