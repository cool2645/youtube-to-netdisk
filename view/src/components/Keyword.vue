<template>
  <div>
    <h2>关键字</h2>
    <ul>
      <li v-for="keyword in data">{{ keyword.keyword }}</li>
    </ul>
  </div>
</template>

<script>
  import config from '../config'
  import urlParam from '../utils/buildUrlParam'

  export default {
    name: "running",
    data() {
      return {
        data: []
      }
    },
    methods: {
      updateData() {
        fetch(config.urlPrefix + '/keyword')
          .then(res => res.json())
          .then(
            res => {
              if (res.result) {
                this.data = res.data;
              }
              if (!this._isBeingDestroyed) setTimeout(() => {
                this.updateData()
              }, 3000);
            }
          )
          .catch(error => {
            if (!this._isBeingDestroyed) setTimeout(() => {
              this.updateData()
            }, 3000);
          });
      },
    },
    mounted() {
      this.updateData()
    }
  }
</script>

<style lang="stylus" scoped>

</style>