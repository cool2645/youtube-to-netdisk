<template>
  <div class="landing">
    <div id="msg-error" class="alert alert-danger alert-dismissable" v-if="showAlert">
      <button type="button" class="close" @click="close" aria-hidden="true">&times;</button>
      <h4><i class="icon fa fa-warning"></i> 噢！太糟糕了！</h4>

      <p id="msg-error-p">{{ error }}</p>
    </div>
    <form class="form-horizontal w" @submit="submit">
      <div class="form-group">
        <label for="url" class="col-sm-2 control-label">视频地址</label>
        <div class="col-sm-10">
          <input type="text" class="form-control" id="url" v-model="form.url"
                 placeholder="例如：https://youtu.be/GmkeBgTqtjk"
                 required
          />
        </div>
      </div>
      <div class="form-group">
        <div class="col-sm-offset-2 col-sm-10">
          <button v-if="!waiting" type="submit" class="btn btn-primary">创建搬运任务</button>
          <button v-else disabled type="submit" class="btn btn-primary">提交中，请稍候...</button>
        </div>
      </div>
    </form>
  </div>
</template>

<script>
  import config from '../config';
  import urlEncode from '../utils/buildUrlParam';

  export default {
    data() {
      return {
        form: {
          url: '',
          token: '',
        },
        error: '',
        showAlert: false,
        waiting: false
      }
    },
    mounted() {
      if (this.$route.query['token']) this.form.token = this.$route.query['token'];
      if (this.$route.query['url']) this.form.url = this.$route.query['url'];
    },
    methods: {
      submit(e) {
        e.preventDefault();
        this.waiting = true;
        fetch(config.urlPrefix + '/trigger', {
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
          },
          method: 'POST',
          body: urlEncode(this.form),
        })
          .then(res => res.json())
          .then(res => {
            this.waiting = false;
            if (res && res.result) {
              const next = { path: '/tasks' };
              if (this.$route.query['token']) next.query = { token: this.$route.query['token'] };
              this.$router.push(next);
            } else {
              this.alert(res.msg)
            }
          })
      },
      alert(msg) {
        this.error = msg;
        this.showAlert = true;
      },
      close() {
        this.showAlert = false;
      }
    }
  }
</script>

<style lang="stylus" scoped>
  .landing
    display flex
    flex-direction column
    justify-content center
    align-items center
    position absolute
    left 0
    top 50%
    transform translateY(-50%)
    width 100%
  .w
    width 600px
    max-width 90%
</style>