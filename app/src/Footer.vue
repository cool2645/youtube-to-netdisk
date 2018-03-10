<template>
    <footer>
        <p v-if="telegram"><a :href="telegram">添加此 Telegram Bot </a>来订阅通知（实验功能）</p>
    </footer>
</template>

<script>
    import config from './config'
    export default {
        computed: {
            telegram() {
                return config.telegram ? "https://t.me/" + config.telegram : ""
            }
        }
    }
</script>

<style scoped>
    footer {
        margin: 20px auto;
    }
</style>