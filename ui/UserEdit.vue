<template>
    <table>
        <tr>
            <td>id</td>
            <td>
                <input type="text" :value="user.id">
            </td>
        </tr>
        <tr>
            <td>username</td>
            <td>
                <input type="text" :value="user.username">
            </td>
        </tr>
        <tr>
            <td>name</td>
            <td>
                <input type="text" :value="user.name">
            </td>
        </tr>
        <tr>
            <td>email</td>
            <td>
                <input type="email" :value="user.email">
            </td>
        </tr>
        <tr>
            <td colspan="2">
                <button @click="saveUser">Save</button>
            </td>
        </tr>
    </table>
</template>

<script>
    export default {
        created() {
            this.$store.dispatch('fetchUser', this.$route.params.username);
        },
        computed: {
            user() {
                return this.$store.getters.getUser(this.$route.params.username);
            }
        }
    }
</script>

<style>
</style>
