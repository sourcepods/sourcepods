<template>
    <div class="container-fluid">
        <h3>{{ name }}</h3>
    </div>
</template>

<script>
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
    .container-fluid h3 {
        font-family: sans-serif;
    }
</style>
