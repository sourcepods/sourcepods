import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";

Vue.use(Vuex);

export const store = new Vuex.Store({
    state: {
        users: [],
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
        }
    },
    mutations: {
        addUsers(state, users) {
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
        }
    },
    actions: {
        fetchUsers(ctx){
            axios.get('/api/users')
                .then((res) => {
                    ctx.commit('addUsers', res.data);
                })
                .catch((err) => {
                    console.log(err);
                })
        },
        fetchUser(ctx, username) {
            axios.get(`/api/users/${username}`)
                .then((res) => {
                    ctx.commit('addUser', res.data);
                })
                .catch((err) => {
                    alert(err);
                })
        },
        deleteUser(ctx, username){
            axios.delete(`/api/users/${username}`)
                .then((res) => {
                    ctx.dispatch('fetchUsers');
                })
                .catch((err) => {
                    console.log(err);
                })
        }
    },
});
