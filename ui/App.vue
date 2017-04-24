<template>
    <div>
        <h3>{{ name }}</h3>

        <router-view></router-view>
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
        }
    }
</script>

<style scoped>
    h3, p, li {
        font-family: sans-serif;
    }
</style>
