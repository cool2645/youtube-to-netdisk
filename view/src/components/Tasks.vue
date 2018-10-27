<template>
  <div>
    <h2>{{ title }}</h2>
    <div v-masonry transition-duration="0.3s" item-selector=".task-card-wrapper" class="row">
      <div v-masonry-tile class="task-card-wrapper col col-md-3 col-sm-6 col-xs-12" v-for="(task, index) in data">
        <div class="task-card" :class="classObject(task)">
          <div id="brief" v-show="isLoaded && !task.showMore">
            <p>{{ task.title }}</p>
            <hr v-if="task.file_name"/>
            <a target="_blank" v-if="task.file_name" :href="'/static/' + task.file_name">本地下载</a>
            <a target="_blank" v-if="task.subtitles"
               v-for="sub in getSubtitles(task.subtitles)"
               :href="'/static/' + getSubtitleFilename(sub, task.file_name)"
               :download="getSubtitleFilename(sub, task.file_name)"
            >{{ sub.lang + '字幕下载' }}</a>
            <span v-if="task.share_link"> | <a target="_blank" :href="getShareUrl(task.share_link)">网盘下载</a>
                            &nbsp;{{ getSharePwd(task.share_link) }}</span>
          </div>
          <div id="more" v-show="task.showMore">
            <div><i class="fa fa-id-card-o fa-fw"></i><span><strong>任务编号</strong>&nbsp;&nbsp;{{ task.id }}</span>
            </div>
            <div v-if="task.title.trim()"><i
                class="fa fa-bookmark-o fa-fw"></i><span><strong>稿件标题</strong>&nbsp;&nbsp;{{ task.title.trim() }}</span>
            </div>
            <div v-if="task.author.trim()"><i class="fa fa-paint-brush fa-fw"></i><span><strong>投稿频道</strong>&nbsp;&nbsp;{{ task.author.trim() }}</span>
            </div>
            <div v-if="enableKeyword"><i
                class="fa fa-question-circle-o fa-fw"></i><span><strong>任务理由</strong>&nbsp;&nbsp;{{ task.reason }}</span>
            </div>
            <div><i class="fa fa-link fa-fw"></i><span><strong>原始地址</strong>&nbsp;&nbsp;<a
                :href="task.url">{{ task.url.trim() }}</a></span>
            </div>
            <div><i class="fa fa-clock-o fa-fw"></i><span><strong>创建时间</strong>&nbsp;&nbsp;{{ formatDateTimeFromDatetimeString(task.created_at) }}</span>
            </div>
            <div><i class="fa fa-clock-o fa-fw"></i><span><strong>更新时间</strong>&nbsp;&nbsp;{{ formatDateTimeFromDatetimeString(task.updated_at) }}</span>
            </div>
            <div><i class="fa fa-file-text-o fa-fw"></i><span><a href="javascript:;"
                                                                 @click="showTaskLog(task.id)">显示日志</a></span>
            </div>
            <hr v-if="task.file_name"/>
            <a target="_blank" v-if="task.file_name" :href="'/static/' + task.file_name">本地下载</a>
            <a target="_blank" v-if="task.subtitles"
               v-for="sub in getSubtitles(task.subtitles)"
               :href="'/static/' + getSubtitleFilename(sub.ext, task.file_name)"
               :download="getSubtitleFilename(sub, task.file_name)"
            >{{ sub.lang + '字幕下载' }}</a>
            <span v-if="task.share_link"> | <a target="_blank" :href="getShareUrl(task.share_link)">网盘下载</a>
                            &nbsp;{{ getSharePwd(task.share_link) }}</span>
          </div>
          <i :class="{'fa-angle-double-right': !task.showMore,
                                'fa-angle-double-left': task.showMore}"
             class="task-card-more fa fa-1x"
             aria-hidden="true"
             @click.prevent="showMore(index)">
            <span v-if="!task.showMore">任务</span>
            <span v-else>隐藏</span>详情
          </i>
          <span :class="{'task-card-state': !task.showMore,
                                       'task-card-state-more': task.showMore}"
                class="task-card-state-text"
                aria-hidden="true">{{ task.state }}
                    </span>
        </div>
      </div>
    </div>
    <pagination :data="laravelData" :limit=2 v-on:pagination-change-page="onPageChange"></pagination>
    <h3 v-if="showLog">任务日志</h3>
    <pre v-if="showLog">
            {{ log }}
        </pre>
  </div>
</template>

<script>
  import config from '../config'
  import urlParam from '../utils/buildUrlParam'
  import formatDateTimeFromDatetimeString from "../utils/datetimeUtil"

  export default {
    name: "tasks",
    data() {
      return {
        jsonSource: {},
        isLoaded: false,
        data: [],
        lastDataLength: 0,
        page: 1,
        perPage: 10,
        showLog: false,
        showLogId: 0,
        log: ""
      }
    },
    watch: {
      $route() {
        this.page = 1
      }
    },
    computed: {
      enableKeyword() {
        return config.enableKeyword;
      },
      title() {
        if (this.$route.path !== "/reject-tasks")
          return "已启动任务";
        else
          return "已拒绝任务";
      },
      rejected() {
        return this.$route.path === "/reject-tasks"
      },
      total() {
        return this.jsonSource.result ? this.jsonSource.data.total : 0;
      },
      laravelData() {
        return {
          current_page: this.page,
          data: [],
          from: (this.page - 1) * this.perPage + 1,
          last_page: Math.ceil(this.total / this.perPage),
          next_page_url: null,
          per_page: this.perPage,
          prev_page_url: null,
          to: (this.page) * this.perPage,
          total: this.total,
        }
      },
    },
    methods: {
      getSubtitles(str) {
        if (!str) return [];
        try {
          const subtitles = JSON.parse(str);
          const res = [];
          for (const k in subtitles) {
            if (subtitles.hasOwnProperty(k)) {
              subtitles[k].lang = k;
              res.push(subtitles[k])
            }
          }
          return res;
        } catch (e) {
          return [];
        }
      },
      getSubtitleFilename(sub, filename) {
        return filename.split('.')[0] + '.' + sub.lang + '.' + sub.ext;
      },
      updatePerPage() {
        let w = document.body.clientWidth;
        if (w <= 768)
          this.perPage = 5;
        else if (w <= 990)
          this.perPage = 16;
        else
          this.perPage = 24;
        return this.perPage;
      },
      formatDateTimeFromDatetimeString(dt) {
        return formatDateTimeFromDatetimeString(dt)
      },
      highlight(text) {
        let reg = /(http:\/\/|https:\/\/)((\w|=|\?|\.|\/|&|-)+)/g;
        text = text.replace(reg, "<a href='$1$2'>$1$2</a><br/>");
        return text;
      },
      getShareUrl(text) {
        let reg = /(http:\/\/|https:\/\/)((\w|=|\?|\.|\/|&|-)+)/g;
        return text.match(reg)[0];
      },
      getSharePwd(text) {
        let reg = /密码：.{4}/g;
        return text.match(reg)[0];
      },
      updateData(recur = true) {
        fetch(config.urlPrefix + '/task?' + urlParam({
          page: this.page,
          perPage: this.updatePerPage(),
          order: 'desc',
          state: this.rejected ? "Rejected" : "%"
        }))
          .then(res => res.json())
          .then(
            res => {
              if (res.result) {
                let oldData = this.data;
                this.jsonSource = res;
                this.data = this.jsonSource.data.data;
                if (!this.isLoaded || (this.isLoaded && this.data.length !== this.lastDataLength)) {
                  this.data.forEach((d, i) => {
                    this.$set(this.data[i], 'showMore', false);
                  });
                } else {
                  this.data.forEach((d, i) => {
                    this.$set(this.data[i], 'showMore', oldData[i].showMore);
                  });
                }
                this.isLoaded = true;
                this.lastDataLength = this.jsonSource.data.data.length;
                setTimeout(() => {
                  this.$redrawVueMasonry();
                }, 1);
              }
              if (!this._isBeingDestroyed && recur) setTimeout(() => {
                this.updateData()
              }, 5000);
            }
          )
          .catch(() => {
            if (!this._isBeingDestroyed && recur) setTimeout(() => {
              this.updateData()
            }, 5000);
          });
        if (this.showLog) this.showTaskLog(this.showLogId);
      },
      onPageChange(page) {
        this.page = page;
        this.updateData(false);
      },
      showTaskLog(id) {
        fetch(config.urlPrefix + '/task/' + id)
          .then(res => res.json())
          .then(
            res => {
              if (res.result) {
                this.log = res.data.log;
                this.showLog = true;
                this.showLogId = id;
              }
            }
          );
      },
      shortenString(str) {
        if (str.length > config.maxStringLength) {
          return str.substring(0, config.maxStringLength) + "...";
        } else {
          return str;
        }
      },
      showMore(id) {
        this.data[id].showMore = !this.data[id].showMore;
        this.$nextTick(this.$redrawVueMasonry);
      },
      classObject(task) {
        return {
          'finished': task.state === 'Finished',
          'error': task.state === 'Exception' || /* legacy */ task.state === 'Error',
          'canceled': task.state === 'Canceled' || /*legacy */ task.state === 'Killed',
          'uploading': task.state === 'Uploading',
          'downloading': task.state === 'Downloading',
          'queuing': task.state === 'Queuing',
          'rejected': task.state === 'Rejected',
          'show-more': task.showMore
        };
      }
    },
    mounted() {
      this.updateData()
    }
  }
</script>

<style lang="stylus" scoped>
  pre {
    white-space:pre-wrap;
    white-space:-moz-pre-wrap;
    white-space:-pre-wrap;
    white-space:-o-pre-wrap;
    word-wrap:break-word;
  }

  a {
    word-break: break-all;
  }

  .row {
    margin-left: 0;
    margin-right: 0;
  }

  th, td {
    white-space: nowrap;
  }

  .task-card-wrapper {
    padding: 0.2em 0.2em 0.2em 0.2em;
  }

  .task-card {
    float: left;
    position: relative;
    width: 100%;
    margin: 0;
    padding: 0.2em 0.2em 1.4em 0.2em;
    border-style: none;
    overflow: auto;
    -webkit-border-radius: 0.2em;
    -moz-border-radius: 0.2em;
    border-radius: 0.4em;
    -ms-word-wrap: break-word;
    word-wrap: break-word;
    -webkit-transition: transform, -o-transform, -ms-transform, -moz-transform, -webkit-transform 0.5s, 0.5s, 0.5s, 0.5s, 0.5s;
    -moz-transition: transform, -o-transform, -ms-transform, -moz-transform, -webkit-transform 0.5s, 0.5s, 0.5s, 0.5s, 0.5s;
    -ms-transition: transform, -o-transform, -ms-transform, -moz-transform, -webkit-transform 0.5s, 0.5s, 0.5s, 0.5s, 0.5s;
    -o-transition: transform, -o-transform, -ms-transform, -moz-transform, -webkit-transform 0.5s, 0.5s, 0.5s, 0.5s, 0.5s;
    transition: transform, -o-transform, -ms-transform, -moz-transform, -webkit-transform 0.5s, 0.5s, 0.5s, 0.5s, 0.5s;
  }

  .show-more {
    /*width: 300px !important;*/
    /*min-height: 180px !important;*/
  }

  .task-card:hover {
    -webkit-transform: scale(1.02);
    -moz-transform: scale(1.02);
    -ms-transform: scale(1.02);
    -o-transform: scale(1.02);
    transform: scale(1.02);
  }

  .task-card hr {
    margin: 0 0 0.3em 0;
  }

  .task-card-more {
    position: absolute;
    float: right;
    bottom: 0;
    right: 0.2em;
    opacity: 1;
    cursor: pointer;
    -webkit-transition: all 0.5s;
    -moz-transition: all 0.5s;
    -ms-transition: all 0.5s;
    -o-transition: all 0.5s;
    transition: all 0.5s;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
  }

  .task-card-more:hover {
    right: 0;
    opacity: 0.5;
  }

  .task-card-state {
    position: absolute;
    float: left;
    bottom: 0;
    left: 0.2em;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
  }

  .task-card-state-more {
    position: absolute;
    float: right;
    top: 0;
    right: 0.2em;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
  }

  .finished {
    background: #98ebac;
  }

  .finished .task-card-state-text {
    color: white
  }

  .error {
    background: #ffa5a5;
  }

  .error .task-card-state-text {
    color: white
  }

  .queuing {
    background: #ffdaa5;
  }

  .queuing .task-card-state-text {
    color: white
  }

  .downloading {
    background: #a5d7ff;
  }

  .downloading .task-card-state-text {
    color: white
  }

  .uploading {
    background: #c8e8f9;
  }

  .uploading .task-card-state-text {
    color: white
  }

  .rejected {
    background: #f3cdff;
  }

  .rejected .task-card-state-text {
    color: white
  }

  .canceled {
    background: #f5f5f5;
  }

  .canceled .task-card-state-text {
    color: black
  }

  #more > div {
    margin-bottom: 0.5em;
  }

  #more > div > span {
    margin-left: 0.5em;
  }
</style>