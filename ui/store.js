import Vuex from "vuex";
import axios from "axios";

export const store = new Vuex.Store({
	strict: process.env.NODE_ENV !== 'production',
	state: {
		loading: false,
		user: null,
		users: [],
		repositories: {},
	},
	getters: {
		getUser: (state) => (username) => {
			const users = state.users.filter((user) => {
				return user.username === username;
			});
			if (users.length > 0) {
				return users[0];
			}
			return {};
		},
		getUserRepositories: (state) => (user_id) => {
			let userRepositories = [];

			const repositories = state.repositories;
			for (let id in repositories) {
				if (repositories.hasOwnProperty(id)) {
					const repository = repositories[id];
					if (repository.relationships.owner.data.id === user_id) {
						userRepositories.push(repository);
					}
				}
			}

			return userRepositories;
		},
		getRepository: (state) => (id) => {
			if (state.repositories[id] !== undefined) {
				return state.repositories[id];
			}
			return null;
		}
	},
	mutations: {
		loading(state, isLoading) {
			state.loading = isLoading;
		},
		setUser(state, user) {
			state.user = user;
		},
		setUsers(state, users) {
			state.users = users;
		},
		addUser(state, newUser) {
			for (let i = 0; i < state.users.length; i++) {
				if (state.users[i].username === newUser.username) {
					state.users[i] = newUser;
					return
				}
			}
			state.users.push(newUser);
		},
		updateUser(state, updatedUser) {
			for (let i = 0; i < state.users.length; i++) {
				if (state.users[i].id === updatedUser.id) {
					state.users[i] = updatedUser;

					// If the current user was updated, update it in the store too
					if (state.user.id === updatedUser.id) {
						state.user = updatedUser;
					}

					return
				}
			}
		},
		setRepository(state, repository) {
			state.repositories[repository.id] = repository;
		},
	},
	actions: {
		fetchAuthenticatedUser(ctx) {
			return new Promise((resolve, reject) => {
				axios.get(`${window.config.api}/user`)
					.then((res) => {
						ctx.commit('setUser', res.data);
						resolve(res.data);
					})
					.catch((err) => {
						reject(err);
					})
			})
		},
		fetchUserRepositories(ctx, username) {
			return new Promise((resolve, reject) => {
				axios.get(`${window.config.api}/users/${username}/repositories`)
					.then((res) => {
						res.data.data.forEach((repository) => {
							ctx.commit('setRepository', repository);
						});
						resolve(res.data);
					})
					.catch((err) => {
						reject(err);
					})
			})
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
				axios.get(`${window.config.api}/users`)
					.then((res) => {
						ctx.commit('setUsers', res.data);
						resolve(res.data);
					})
					.catch((err) => {
						reject(err);
					})
			})
		},
		fetchUser(ctx, username) {
			return new Promise((resolve, reject) => {
				axios.get(`${window.config.api}/users/${username}`)
					.then((res) => {
						ctx.commit('addUser', res.data);
						resolve(res.data);
					})
					.catch((err) => {
						reject(err);
					})
			})
		},
		updateUser(ctx, user) {
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
			return new Promise((resolve, reject) => {
				axios.get(`${window.config.api}/repositories/${data.owner}/${data.repository}`)
					.then((res) => {
						ctx.commit('setRepository', res.data.data);
						resolve(res.data.data);
					})
					.catch((err) => {
						reject(err);
					})
			})
		}
	},
});
