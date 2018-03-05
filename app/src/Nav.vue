<template>
    <div>
        <ul>
            <li v-for="route in routes" v-if="route.title"><a href="javascript:;" @click="$router.push(route.path)" >{{route.title}}</a></li>
        </ul>
    </div>
</template>

<script>

    export default {
        props: ['routes'],
    }
</script>

<style scoped>
    ul {
        margin-top: 20px;
        margin-bottom: 20px;
    }
    li {
        display: inline;
        list-style-type: none;
        padding-right: 20px;
    }
</style>
