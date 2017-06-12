<template>
	<div v-if="loading || user === null"></div>
	<div class="uk-container" v-else>

		<div class="uk-grid-small" uk-grid v-if="user!==null">
			<div class="uk-flex-top uk-padding-small uk-width-2-5@m uk-width-1-4@l">
				<gravatar :email="user.email" default-img="mm" :size="512"
						  class="uk-border-rounded"></gravatar>
				<h3 class="user-name">{{user.name}}</h3>
				<h4 class="uk-text-muted user-username">@{{user.username}}</h4>

				<hr class="uk-divider-icon">

				<ul class="uk-list user-details">
					<li>
						<span class="uk-icon-link" uk-icon="icon: mail"></span>
						<a :href="`mailto:${user.email}`">{{user.email}}</a>
					</li>
					<li>
						<span uk-icon="icon: clock"></span>
						<span
							:title="new Date(user.created_at*1000) | moment('YYYY MMMM DD, HH:mm:ss')">
							Joined on {{new Date(user.created_at * 1000) | moment("from")}}
						</span>
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
						<li>
							<div uk-grid>
								<div class="uk-width-expand uk-margin">
									<form class="uk-search uk-search-default" style="width: 100%;">
										<a href="" class="uk-search-icon-flip" uk-search-icon></a>
										<input class="uk-search-input" type="search" placeholder="Search..."
											   v-model="search">
									</form>
								</div>
								<div class="uk-width-auto">
									<router-link class="uk-button uk-button-primary" :to="`/new`">New</router-link>
								</div>
							</div>
						</li>
						<li v-for="repository in repositories">
							<div class="uk-flex">
								<div class="uk-flex-auto">
									<h4 class="repository-name">
										<router-link :to="`/${user.username}/${repository.name}`">
											{{repository.name}}
										</router-link>
									</h4>
									<p class="repository-description">{{repository.description}}</p>
								</div>
								<div class="uk-text-right">
									<span uk-icon="icon: star"></span>
									<span>{{repository.stars}}</span>
									<span uk-icon="icon: git-fork"></span>
									<span>{{repository.forks}}</span>
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
				search: null,
			}
		},
		created() {
			this.setLoading(true);
			this.$store.dispatch('fetchUserProfile', this.$route.params.username)
				.then(this.setLoading(false));
		},
		computed: {
			user() {
				return this.$store.getters.getUserByUsername(this.$route.params.username);
			},
			repositories() {
				let repositories = this.$store.getters.getUserRepositories(this.user.id);
				if (this.search === null) {
					return repositories;
				}
				return repositories.filter((repository) => repository.name.includes(this.search));
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
