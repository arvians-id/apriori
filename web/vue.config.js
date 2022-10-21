module.exports = {
    devServer: {
        proxy: process.env.VUE_APP_PROXY_URL,
        port: 8081, // CHANGE YOUR PORT HERE!
    }
}