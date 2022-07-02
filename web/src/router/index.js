import { createRouter, createWebHistory } from "vue-router";
import NotFound from "../views/NotFound.vue";
import authHeader from "@/service/auth-header";
import axios from "axios";
import AuthRouter from "@/router/auth-router";
import TransactionRouter from "@/router/transaction-router";
import ProductRouter from "@/router/product-router";
import AprioriRouter from "@/router/apriori-router";
import UserRouter from "@/router/user-router";

const routes = [
    ...AuthRouter,
    ...TransactionRouter,
    ...ProductRouter,
    ...AprioriRouter,
    ...UserRouter,
  { path: "/:pathMatch(.*)", name: "NotFound", component: NotFound },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

router.beforeEach( async (to) => {
  if (Object.keys(authHeader()).length === 0 && to.name.split(".")[0] !== "auth") {
    return { name: 'auth.login' }
  } else if (Object.keys(authHeader()).length > 0 && to.name.split(".")[0] === "auth") {
      return { name: 'admin' }
  }

  if (to.name.split(".")[0] !== "auth") {
        axios.get(`${process.env.VUE_APP_SERVICE_URL}/auth/token`, { headers: authHeader() })
            .catch(() => {
                let refreshToken = {
                    refresh_token: localStorage.getItem("refresh-token")
                }
                axios.post(`${process.env.VUE_APP_SERVICE_URL}/auth/refresh`,refreshToken)
                    .then(response => {
                        let token = response.data.data.access_token
                        let refreshToken = response.data.data.refresh_token
                        localStorage.setItem("token", token)
                        localStorage.setItem("refresh-token", refreshToken)
                    }).catch(() => {
                    localStorage.removeItem("token")
                    localStorage.removeItem("refresh-token")
                })
            })
    }
})

export default router;
