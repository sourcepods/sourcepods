import Vue from "vue";
import VueRouter from "vue-router";

import {store} from "./store";
import App from "./App.vue";
import UserList from "./UserList.vue";
import UserProfile from "./UserProfile.vue";
import UserEdit from "./UserEdit.vue";

Vue.use(VueRouter);

const router = new VueRouter({
    mode: 'history',
    routes: [
        {path: "/", component: UserList},
        {path: "/users/:username", component: UserProfile},
        {path: "/users/:username/edit", component: UserEdit},
    ],
});

new Vue({
    el: '#app',
    store: store,
    router: router,
    render: h => h(App)
});
