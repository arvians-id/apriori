import { createRouter, createWebHistory } from "vue-router";
import NotFound from "../views/NotFound.vue";
import authHeader from "@/service/auth-header";
import axios from "axios";
import AuthRouter from "@/router/auth-router";
import TransactionRouter from "@/router/transaction-router";
import ProductRouter from "@/router/product-router";
import AprioriRouter from "@/router/apriori-router";
import UserRouter from "@/router/user-router";
import UserGuestRouter from "@/router/user-guest-router";

const routes = [
    ...UserGuestRouter,
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
  if (authHeader()["Authorization"] === undefined && to.name.split(".")[0] !== "auth" && to.name.split(".")[0] !== "guest") {
      return { name: 'auth.login'}
  } else if (authHeader()["Authorization"] && to.name.split(".")[0] === "auth") {
      return { name: 'admin' }
  }

  if (localStorage.getItem("role") === "2" &&
      (to.name.split(".")[0] === "transaction" ||
      to.name.split(".")[0] === "product" ||
      to.name.split(".")[0] === "apriori" ||
      to.name.split(".")[0] === "user" ||
      to.name.split(".")[0] === "admin" ||
      to.name.split(".")[0] === "profile")) {
      return { name: 'guest.index' }
  }

  setTimeout(() => {
      if (to.name.split(".")[0] !== "auth" && authHeader()["Authorization"]) {
          axios.get(`${process.env.VUE_APP_SERVICE_URL}/auth/token`, { headers: authHeader() })
              .catch(() => {
                  let refreshToken = {
                      refresh_token: localStorage.getItem("refresh-token")
                  }
                  axios.post(`${process.env.VUE_APP_SERVICE_URL}/auth/refresh`,refreshToken,{ headers: authHeader() })
                      .then(response => {
                          let token = response.data.data.access_token
                          let refreshToken = response.data.data.refresh_token
                          localStorage.setItem("token", token)
                          localStorage.setItem("refresh-token", refreshToken)

                          axios.get(`${process.env.VUE_APP_SERVICE_URL}/profile`, { headers: authHeader() })
                              .then(response => {
                                  localStorage.setItem("user", response.data.data.id_user)
                                  localStorage.setItem("name", response.data.data.name)
                              }).catch(error => {
                              console.log(error.response.data.status)
                          })
                      }).catch(() => {
                      localStorage.removeItem("token")
                      localStorage.removeItem("refresh-token")
                      localStorage.removeItem("user")
                      localStorage.removeItem("name")
                      this.$router.push({ name: 'auth.login' })
                  })
              })
      }
  }, 5000);
})

export default router;
