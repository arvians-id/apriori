import { createApp, h } from "vue";
import App from "./App.vue";
import router from "./router";
import { ApolloClient, createHttpLink, InMemoryCache } from '@apollo/client/core'
import { createApolloProvider } from '@vue/apollo-option'
import authHeader from "@/service/auth-header";

const httpLink = createHttpLink({
    // You should use an absolute URL here
    uri: 'http://localhost:3000/query',
    headers: authHeader()
})


// Create the apollo client
const apolloClient = new ApolloClient({
    link: httpLink,
    cache: new InMemoryCache(),
})

const apolloProvider = createApolloProvider({
    defaultClient: apolloClient,
})

const app = createApp({
    render: () => h(App),
})

app.use(router)
app.use(apolloProvider)
app.mount("#app");