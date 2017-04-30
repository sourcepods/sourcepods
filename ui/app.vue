<template>
    <div>
        <h1>{{ name }}</h1>

        <router-view></router-view>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                name: 'gitpod',
            }
        }
    }
</script>

<style scoped>
    div {
        font-family: sans-serif;
    }
</style>
