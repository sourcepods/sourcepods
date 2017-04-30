<template>
    <div>
        <h4>{{user.username}}</h4>

        <gravatar :email="user.email" default-img="mm"></gravatar>

        <table>
            <tr>
                <td>id</td>
                <td>{{user.id}}</td>
            </tr>
            <tr>
                <td>name</td>
                <td>{{user.name}}</td>
            </tr>
            <tr>
                <td>email</td>
                <td>{{user.email}}</td>
            </tr>
        </table>
    </div>
</template>

<script>
    import Gravatar from 'vue-gravatar';

    export default {
        components:{
            gravatar: Gravatar,
        },
        created() {
            this.$store.dispatch('fetchUser', this.$route.params.username);
        },
        computed: {
            user() {
                return this.$store.getters.getUser(this.$route.params.username);
            }
        },
    }
</script>

<style>
</style>
