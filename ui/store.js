import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

export const store = new Vuex.Store({
    state: {
        users: [],
    },
    mutations: {
        addUsers(state, users) {
            state.users = state.users.concat(users);
        },
        addUser(state, username) {
            state.users.push(username);
        }
    }
});
