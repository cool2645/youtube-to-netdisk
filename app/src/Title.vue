<template>
    <div>
        <a class="ribbon" href="https://github.com/cool2645/youtube-to-netdisk"><img style="position: absolute; top: 0; right: 0; border: 0;" src="https://s3.amazonaws.com/github/ribbons/forkme_right_gray_6d6d6d.png" alt="Fork me on GitHub"></a>
        <p v-if="error">{{ error }}</p>
    </div>
</template>

<script>
    export default {
        data() {
            return {
                error: ""
            }
        },
        mounted() {
            if (!fetch) {
                this.error = "内核太旧尚未支持，请升级或更换浏览器或帮助完善代码来正常浏览本页面(>ω<)"
            }
        }
    }
</script>

<style scoped>
    @media (max-width: 600px) {
        .ribbon {
            display: none;
        }
    }
    p {
        margin-top: 20px;
        margin-bottom: 20px;
    }
</style>