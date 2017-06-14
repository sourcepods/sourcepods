import Vuex from "vuex";
import axios from "axios";
import Lokka from "lokka";
import {Transport} from "lokka-transport-http";
import {normalize, schema} from "normalizr";

const client = new Lokka({
	transport: new Transport(`${window.config.api}/query`)
});

export const store = new Vuex.Store({
	strict: process.env.NODE_ENV !== 'production',
	state: {
		loading: false,
		user_id: null,
		users: {},
		repositories: {},
	},
	getters: {
		getAuthUser(state) {
			if (state.user_id === null) {
				return null;
			}
			return state.users[state.user_id];
		},
		getUsers: (state) => {
			return state.users;
		},
		getUserByUsername: (state) => (username) => {
			const users = Object.keys(state.users).map((id) => state.users[id]);
			return users.filter((user) => user.username === username)[0];
		},
		getUserRepositories: (state) => (user_id) => {
			const user = state.users[user_id];

			if (user.repositories === undefined) {
				return [];
			}

			let repositories = [];
			user.repositories.forEach((repo_id) => {
				repositories.push(state.repositories[repo_id]);
			});

			return repositories;
		},
		getRepository: (state) => (id) => {
			return state.repositories[id];
		}
	},
	mutations: {
		loading(state, isLoading) {
			state.loading = isLoading;
		},
		setAuthUser(state, user_id) {
			state.user_id = user_id;
		},
		setUsers(state, users) {
			state.users = Object.assign({}, state.users, users);
		},
		setRepositories(state, repositories) {
			state.repositories = Object.assign({}, state.repositories, repositories);
		},
	},
	actions: {
		fetchAuthenticatedUser(ctx) {
			const userSchema = new schema.Entity('me');

			return new Promise((resolve, reject) => {
				const query = `
					query me {
						me {
							id
							username
							name
							email
							created_at
							updated_at
						}
					}`;
				client.query(query).then((res) => {
					const data = normalize(res.me, userSchema);

					ctx.commit('setUsers', data.entities.me);
					ctx.commit('setAuthUser', data.result);

					resolve(res);
				}).catch((err) => {
					reject(err);
				})
			});
		},
		authenticateUser(ctx, payload) {
			return new Promise((resolve, reject) => {
				axios.post(`${window.config.api}/authorize`, payload)
					.then((res) => {
						ctx.commit('setUsers', res.data);
						resolve(res.data);
					})
					.catch((err) => {
						reject(err);
					})
			})
		},
		fetchUsers(ctx) {
			const userSchema = new schema.Entity('users');

			return new Promise((resolve, reject) => {
				const query = `
					query allUsers {
						users {
							id
							username
							name
							email
						}
					}`;
				client.query(query).then((res) => {
					const data = normalize(res.users, [userSchema]);

					ctx.commit('setUsers', data.entities.users);

					resolve(res);
				}).catch((err) => {
					reject(err);
				});
			});
		},
		fetchUserProfile(ctx, username) {
			const repositorySchema = new schema.Entity('repositories');
			const userSchema = new schema.Entity('user', {
				repositories: [repositorySchema],
			});

			return new Promise((resolve, reject) => {
				const query = `
					query userProfile($username: String) {
						user(username: $username) {
							id
							username
							name
							email
							created_at
							updated_at
							repositories {
								id
								name
								description
								forks
								stars
							}
						}
					}`;
				client.query(query,
					{
						username: username,
					}
				).then((res) => {
					const data = normalize(res.user, userSchema);

					ctx.commit('setUsers', data.entities.user);
					ctx.commit('setRepositories', data.entities.repositories);

					resolve(res);
				}).catch((err) => {
					reject(err);
				});
			})
		},
		updateUser(ctx, user) {
			const userSchema = new schema.Entity('user');

			return new Promise((resolve, reject) => {
				const query = `
				($id: ID!, $user: updatedUser!) {
					updateUser(id: $id, user: $user) {
						id
						email
						username
						name
						created_at
						updated_at
					}
				}`;
				client.mutate(query,
					{
						id: user.id,
						user: {
							name: user.name,
						}
					}
				).then((res) => {
					const data = normalize(res.updateUser, userSchema);

					ctx.commit('setUsers', data.entities.user);

					resolve(res.updateUser);
				}).catch((err) => {
					reject(err);
				});
			})
		},
		deleteUser(ctx, username) {
			return new Promise((resolve, reject) => {
				axios.delete(`${window.config.api}/users/${username}`)
					.then((res) => {
						ctx.dispatch('fetchUsers');
						resolve(res.data.data);
					})
					.catch((err) => {
						reject(err);
					})
			})
		},
		fetchRepository(ctx, data) {
			const repositorySchema = new schema.Entity('repository');

			return new Promise((resolve, reject) => {
				const query = `
					query Repository($owner: String, $name: String) {
						repository(owner: $owner, name: $name) {
							id
							name
							description
							website
							private
							stars
							forks
							issue_stats {
								open
							}
							pull_request_stats {
								open
							}
						}
					}`;
				client.query(query,
					{
						owner: data.owner,
						name: data.repository,
					}
				).then((res) => {
					const data = normalize(res.repository, repositorySchema);

					ctx.commit('setRepositories', data.entities.repository);

					resolve(res.repository);
				}).catch((err) => {
					reject(err);
				});
			})
		}
	},
});
