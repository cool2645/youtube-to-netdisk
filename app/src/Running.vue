<template>
    <div>
        <h2>正在运行</h2>
        <div v-for="task in data" class="table-responsive">
            <table class="table">
                <tbody>
                <tr>
                    <th>任务 ID</th>
                    <td>{{ task.ID }}</td>
                </tr>
                <tr>
                    <th>标题</th>
                    <td>{{ task.Title }}</td>
                </tr>
                <tr>
                    <th>投稿</th>
                    <td>{{ task.Author }}</td>
                </tr>
                <tr>
                    <th>原始 URL</th>
                    <td><a :href="task.URL">{{ task.URL }}</a> </td>
                </tr>
                <tr>
                    <th>状态</th>
                    <td>{{ task.State }}</td>
                </tr>
                <tr>
                    <th>任务理由</th>
                    <td>{{ task.Reason }}</td>
                </tr>
                <tr>
                    <th>任务日志</th>
                    <td><pre>{{ task.Log }}</pre></td>
                </tr>
                </tbody>
            </table>
        </div>
    </div>
</template>

<script>
    import config from './config'
    import urlParam from './buildUrlParam'
    export default {
        name: "running",
        data() {
            return {
                data: []
            }
        },
        methods: {
            updateData() {
                fetch(config.urlPrefix + '/running')
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    this.data = res.data;
                                }
                            }
                        );
                        setTimeout(() => {this.updateData()}, 3000);
                    })
                    .catch(error => {
                        setTimeout(() => {this.updateData()}, 3000);
                    });
            },
        },
        mounted() {
            this.updateData()
        }
    }
</script>

<style scoped>
    pre {
        white-space: pre-wrap;
        word-wrap: break-word;
        white-space: -moz-pre-wrap;
    }
</style>