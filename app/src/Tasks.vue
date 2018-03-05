<template>
    <div>
        <h2>{{ title }}</h2>
        <div class="table-responsive">
            <table class="table">
                <tbody>
                <tr>
                    <th>ID</th>
                    <th>标题</th>
                    <th>投稿</th>
                    <th>原始 URL</th>
                    <th>状态</th>
                    <th>任务理由</th>
                    <th>本地下载地址</th>
                    <th>百度网盘地址</th>
                    <th>任务日志</th>
                </tr>
                <tr v-for="task in data">
                    <td>{{ task.ID }}</td>
                    <td>{{ task.Title }}</td>
                    <td>{{ task.Author }}</td>
                    <td><a :href="task.URL">{{ task.URL }}</a> </td>
                    <td>{{ task.State }}</td>
                    <td>{{ task.Reason }}</td>
                    <td><a :href="'/static/' + task.FileName">{{ task.FileName }}</a></td>
                    <td>{{ task.ShareLink }}</td>
                    <td><a href="javascript:;" @click="showTaskLog(task.ID)">显示日志</a></td>
                </tr>
                </tbody>
            </table>
        </div>
        <pagination :data="laravelData" :limit=2 v-on:pagination-change-page="onPageChange"></pagination>
        <h3 v-if="showLog">任务日志</h3>
        <pre v-if="showLog">
            {{ log }}
        </pre>
    </div>
</template>

<script>
    import config from './config'
    import urlParam from './buildUrlParam'
    export default {
        name: "tasks",
        data() {
            return {
                jsonSource: {},
                page: 1,
                perPage: 10,
                showLog: false,
                log: ""
            }
        },
        computed: {
            title() {
                if (this.$route.path === "/tasks")
                    return "已启动任务";
                else
                    return "已拒绝任务";
            },
            rejected() {
                return this.$route.path !== "/tasks"
            },
            data() {
                return this.jsonSource.result ? this.jsonSource.data.data : []
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
            updateData() {
                fetch(config.urlPrefix + '/task?' + urlParam({
                    page: this.page,
                    order: 'desc',
                    state: this.rejected ? "Rejected" : "%"
                }))
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    this.jsonSource = res;
                                }
                            }
                        );
                        setTimeout(() => {this.updateData()}, 3000);
                    })
                    .catch(error => {
                        setTimeout(() => {this.updateData()}, 3000);
                    });
            },
            onPageChange(page) {
                this.page = page;
                this.updateData();
            },
            showTaskLog(id) {
                fetch(config.urlPrefix + '/task/' + id)
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    this.log = res.data.Log;
                                    this.showLog = true;
                                }
                            }
                        );
                    })
            }
        },
        mounted() {
            this.updateData()
        }
    }
</script>

<style scoped>
    td {
        white-space: nowrap;
    }
</style>