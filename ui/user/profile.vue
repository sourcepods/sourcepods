<template>
	<div v-if="loading"></div>
	<div class="uk-container" v-else>

		<div class="uk-grid-small" uk-grid v-if="user!==null">
			<div class="uk-flex-top uk-padding-small uk-width-2-5@m uk-width-1-4@l">
				<gravatar :email="user.attributes.email" default-img="mm" :size="512"
						  class="uk-border-rounded"></gravatar>
				<h3 class="user-name">{{user.attributes.name}}</h3>
				<h4 class="uk-text-muted user-username">@{{user.attributes.username}}</h4>

				<hr class="uk-divider-icon">

				<ul class="uk-list user-details">
					<li>
						<span class="uk-icon-link" uk-icon="icon: mail"></span>
						<a :href="`mailto:${user.attributes.email}`">{{user.attributes.email}}</a>
					</li>
					<li>
						<span uk-icon="icon: clock"></span>
						<span>Joined on {{user.attributes.created_at}}</span>
					</li>
					<li></li>
				</ul>
			</div>

			<div class="uk-width-3-5@m uk-width-3-4@l">

				<div uk-sticky>
					<ul class="uk-child-width-expand profile-tab" uk-tab>
						<li class="uk-active"><a href="#">Repositories</a></li>
						<li><a href="#">Activity</a></li>
						<li><a href="#">Stars</a></li>
					</ul>
				</div>

				<div>
					<ul class="uk-list uk-list-large uk-list-divider repository-list">
						<li v-for="repository in repositories">
							<div class="uk-flex">
								<div class="uk-flex-auto">
									<h4 class="repository-name">
										<router-link :to="`/${user.attributes.username}/${repository.attributes.name}`">
											{{repository.attributes.name}}
										</router-link>
									</h4>
									<p class="repository-description">{{repository.attributes.description}}</p>
								</div>
								<div class="uk-text-right">
									<span uk-icon="icon: star"></span>
									<span>{{repository.attributes.stars}}</span>
									<span uk-icon="icon: git-fork"></span>
									<span>{{repository.attributes.forks}}</span>
								</div>
							</div>
						</li>
					</ul>
				</div>

			</div>
		</div>

	</div>
</template>

<script>
	import Gravatar from 'vue-gravatar';

	export default {
		components: {
			gravatar: Gravatar,
		},
		data() {
			return {
				loading: true,
			}
		},
		created() {
			this.setLoading(true);
			Promise.all([
				this.$store.dispatch('fetchUser', this.$route.params.username),
				this.$store.dispatch('fetchUserRepositories', this.$route.params.username),
			])
				.then((responses) => {
					this.setLoading(false);
				})
				.catch(() => {
					this.setLoading(false);
				})
		},
		computed: {
			user() {
				return this.$store.getters.getUserByUsername(this.$route.params.username);
			},
			repositories() {
				return this.$store.getters.getUserRepositories(this.user.id);
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
	h3.user-name {
		margin-top: 16px;
		margin-bottom: 0;
	}

	h4.user-username {
		margin: 0;
	}

	ul.user-details {
		margin-top: 0;
	}

	ul.profile-tab {
		background-color: white;
	}

	ul.profile-tab li a {
		padding-top: 20px;
		padding-bottom: 20px;
	}

	ul.repository-list .repository-name {
		margin-bottom: 0;
	}

	ul.repository-list .repository-description {
		margin: 0;
	}
</style>
