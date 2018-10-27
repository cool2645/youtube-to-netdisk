<template>
    <center>
        <ul>
            <li v-for="route in routes" v-if="route.title"><a href="javascript:;" @click="$router.push(route.path)" >{{route.title}}</a></li>
        </ul>
    </center>
</template>

<script>

    export default {
        props: ['routes'],
    }
</script>

<style lang="stylus" scoped>
    ul {
        padding: 0;
        margin: 20px auto;
    }
    li {
        display: inline;
        list-style-type: none;
        padding-right: 20px;
    }
</style>
