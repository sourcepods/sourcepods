import Vuex from "vuex";
import axios from "axios";
import Lokka from "lokka";
import {Transport} from "lokka-transport-http";

const client = new Lokka({
	transport: new Transport(`${window.config.api}/query`)
});

export const store = new Vuex.Store({
		strict: process.env.NODE_ENV !== 'production',
		state: {
			loading: false,
			user_id: null,
			users: [],
			repositories: [],
		},
		getters: {
			getAuthUser(state) {
				if (state.user_id === null) {
					return null;
				}

				let index = state.users.findIndex((user) => user.id === state.user_id);
				if (index >= 0) {
					return state.users[index];
				}
				return null;
			},
			getUsers: (state) => {
				return state.users;
			},
			getUserByUsername: (state) => (username) => {
				let index = state.users.findIndex((user) => user.username === username);
				if (index >= 0) {
					return state.users[index];
				}
				return null;
			},
			getUserRepositories: (state) => (user_id) => {
				let userIndex = state.users.findIndex((user) => user.id === user_id);

				if (state.users[userIndex].repositories === undefined) {
					return null;
				}

				return state.repositories.filter(function (e) {
					return this.indexOf(e.id) >= 0;
				}, state.users[userIndex].repositories);
			},
			getRepository: (state) => (id) => {
				let index = state.repositories.findIndex((repository) => repository.id === id);
				if (index >= 0) {
					return state.repositories[index]
				}
				return null;
			}
		},
		mutations: {
			loading(state, isLoading) {
				state.loading = isLoading;
			},
			setAuthUser(state, user_id) {
				state.user_id = user_id;
			},
			setUser(state, user) {
				let index = state.users.findIndex((u) => u.id === user.id);
				if (index >= 0) {
					state.users[index] = Object.assign({}, state.users[index], user);
				} else {
					state.users.push(user);
				}
			},
			setRepository(state, repository) {
				let index = state.repositories.findIndex((r) => r.id === repository.id);
				if (index >= 0) {
					state.repositories[index] = Object.assign({}, state.repositories[index], repository);
				} else {
					state.repositories.push(repository);
				}
			},
		},
		actions: {
			fetchAuthenticatedUser(ctx) {
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
						ctx.commit('setUser', res.me);
						ctx.commit('setAuthUser', res.me.id);
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
							ctx.commit('setUser', res.data);
							resolve(res.data);
						})
						.catch((err) => {
							reject(err);
						})
				})
			},
			fetchUsers(ctx) {
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
						res.users.forEach((user) => {
							ctx.commit('setUser', user);
						});
						resolve(res);
					}).catch((err) => {
						reject(err);
					});
				});
			},
			fetchUserProfile(ctx, username) {
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
						let user = Object.assign({}, res.user);
						user.repositories.forEach((repository) => {
							ctx.commit('setRepository', repository);
						});
						user.repositories = user.repositories.map(repository => repository.id);
						ctx.commit('setUser', user);
						setTimeout(resolve(res), 100);
					}).catch((err) => {
						reject(err);
					});
				})
			},
			updateUser(ctx, user)
			{
				return new Promise((resolve, reject) => {
					axios.put(`${window.config.api}/users/${user.username}`, user)
						.then((res) => {
							ctx.commit('updateUser', res.data);
							resolve(res.data);
						})
						.catch((err) => {
							reject(err);
						})
				})
			}
			,
			deleteUser(ctx, username)
			{
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
			}
			,
			fetchRepository(ctx, data)
			{
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
						ctx.commit('setRepository', res.repository);
						resolve(res.repository);
					}).catch((err) => {
						reject(err);
					});
				})
			}
		},
	})
;
