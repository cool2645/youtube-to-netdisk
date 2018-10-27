<template>
  <div class="center">
    <ul>
      <li v-for="route in routes" v-if="route.title && !route.hide"><a href="javascript:;" @click="$router.push(route.path)">{{route.title}}</a>
      </li>
    </ul>
  </div>
</template>

<script>

  export default {
    props: ['routes'],
  }
</script>

<style lang="stylus" scoped>
  .center {
    display flex
    justify-content center
  }

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
