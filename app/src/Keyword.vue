<template>
    <div>
        <h2>关键字</h2>
        <ul>
            <li v-for="keyword in data">{{ keyword.Keyword }}</li>
        </ul>
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
                fetch(config.urlPrefix + '/keyword')
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

</style>