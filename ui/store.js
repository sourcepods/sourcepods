import Vuex from "vuex";
import axios from "axios";

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
			let index = state.users.findIndex((user) => user.attributes.username === username);
			if (index >= 0) {
				return state.users[index];
			}
			return null;
		},
		getUserRepositories: (state) => (user_id) => {
			return state.repositories.filter((repository) => repository.relationships.owner.data.id === user_id);
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
				state.users[index] = user;
			} else {
				state.users.push(user);
			}
		},
		setRepository(state, repository) {
			let index = state.repositories.findIndex((r) => r.id === repository.id);
			if (index >= 0) {
				state.repositories[index] = repository;
			} else {
				state.repositories.push(repository);
			}
		},
	},
	actions: {
		fetchAuthenticatedUser(ctx) {
			return new Promise((resolve, reject) => {
				axios.get(`${window.config.api}/user`)
					.then((res) => {
						ctx.commit('setUser', res.data.data);
						ctx.commit('setAuthUser', res.data.data.id);
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
						res.data.data.forEach((user) => {
							ctx.commit('setUser', user);
						});
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
						ctx.commit('setUser', res.data.data);
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
