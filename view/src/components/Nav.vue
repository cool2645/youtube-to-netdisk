<template>
  <div class="center">
    <ul>
      <li v-for="route in routes" v-if="route.title && !route.hide"><a href="javascript:;"
          @click="pushHistory(route.path)">{{route.title}}</a>
      </li>
    </ul>
  </div>
</template>

<script>
  export default {
    props: ['routes'],
    methods: {
      pushHistory(path) {
        if (this.$route.query['token']) {
          this.$router.push({
            path, query: { token: this.$route.query['token'] }
          })
        } else {
          this.$router.push({
            path
          })
        }
      }
    }
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
