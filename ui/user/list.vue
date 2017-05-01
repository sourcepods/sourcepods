<template>
    <div class="uk-container">
        <p>users:</p>
        <ul>
            <li v-for="user in users">
                {{ user.name }} - {{ user.email }} -


                <router-link :to="`/users/${user.username}`">profile</router-link>
                <router-link :to="`/users/${user.username}/edit`">edit</router-link>
                <span @click="deleteUser(user)">delete</span>
            </li>
        </ul>
    </div>
</template>

<script>
    import axios from 'axios';

    export default {
        name: 'user-list',
        data() {
            return {}
        },
        computed: {
            users() {
                return this.$store.state.users;
            }
        },
        created(){
            axios.get('/api/users')
                .then((res) => {
                    this.$store.commit('addUsers', res.data);
                })
                .catch((err) => {
                    console.log(err);
                })
        },
        methods: {
            deleteUser(user) {
                if (confirm(`Do you really want to delete ${user.username}?`)) {
                    this.$store.dispatch('deleteUser', user.username);
                }
            }
        }
    }
</script>

<style>
</style>
