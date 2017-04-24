<template>
    <div>
        <h3>{{ name }}</h3>

        <p>users:</p>
        <ul>
            <li v-for="user in users">{{ user.name }} - {{ user.email }}</li>
        </ul>
    </div>
</template>

<script>
    import axios from 'axios';

    export default {
        name: 'app',
        data () {
            return {
                name: 'gitloud',
            }
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
        }
    }
</script>

<style scoped>
    h3, p, li {
        font-family: sans-serif;
    }
</style>
