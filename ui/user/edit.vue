<template>
    <form v-on:submit.prevent="saveUser">
        <table>
            <tr>
                <td>id</td>
                <td>
                    <input type="text" v-model="user.id"/>
                </td>
            </tr>
            <tr>
                <td>username</td>
                <td>
                    <input type="text" v-model="user.username"/>
                </td>
            </tr>
            <tr>
                <td>name</td>
                <td>
                    <input type="text" v-model="user.name"/>
                </td>
            </tr>
            <tr>
                <td>email</td>
                <td>
                    <input type="email" v-model="user.email"/>
                </td>
            </tr>
            <tr>
                <td colspan="2">
                    <button type="submit">Save</button>
                </td>
            </tr>
        </table>
    </form>
</template>

<script>
    export default {
        data(){
            return {}
        },
        created() {
            this.$store.dispatch('fetchUser', this.$route.params.username);
        },
        computed: {
            user() {
                const user = this.$store.getters.getUser(this.$route.params.username);
                return Object.assign({}, user);
            }
        },
        methods: {
            saveUser() {
                this.$store.dispatch('updateUser', this.user);
            }
        }
    }
</script>

<style>
</style>
