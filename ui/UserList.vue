<template>
    <div>
        <p>users:</p>
        <ul>
            <li v-for="user in users">{{ user.name }} - {{ user.email }}</li>
        </ul>
    </div>
</template>

<script>
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
        }
    }
</script>

<style>
</style>
