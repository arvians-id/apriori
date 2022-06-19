import { createRouter, createWebHistory } from "vue-router";
import Login from "../views/auth/Login.vue";
import Register from "../views/auth/Register.vue";
import ForgotPassword from "../views/auth/ForgotPassword.vue";
import ResetPassword from "../views/auth/ResetPassword.vue";
import Dashboard from "../views/Dashboard.vue";
import Transaction from "../views/transaction/Data.vue";
import TransactionCreate from "../views/transaction/Create.vue";
import TransactionCreateCSV from "../views/transaction/CreateFile.vue";
import TransactionEdit from "../views/transaction/Edit.vue";
import Product from "../views/product/Data.vue";
import ProductCreate from "../views/product/Create.vue";
import ProductEdit from "../views/product/Edit.vue";
import Apriori from "../views/apriori/Data.vue";
import AprioriCreate from "../views/apriori/Create.vue";
import AprioriDetail from "../views/apriori/Detail.vue";
import User from "../views/user/Data.vue";
import UserCreate from "../views/user/Create.vue";
import UserEdit from "../views/user/Edit.vue";
import Profile from "../views/user/Profile.vue";
import NotFound from "../views/NotFound.vue";
import authHeader from "@/service/auth-header";
import axios from "axios";

const routes = [
  { path: "/", name: "login", component: Login },
  { path: "/register", name: "register", component: Register },
  { path: "/forgot-password", name: "forgot-password", component: ForgotPassword },
  { path: "/reset-password", name: "reset-password", component: ResetPassword },
  { path: "/home", name: "admin", component: Dashboard },
  { path: "/transaction", name: "transaction", component: Transaction },
  { path: "/transaction/create", name: "transaction.create", component: TransactionCreate },
  { path: "/transaction/create/csv", name: "transaction.create.csv", component: TransactionCreateCSV },
  { path: "/transaction/:no_transaction/edit", name: "transaction.edit", component: TransactionEdit },
  { path: "/product", name: "product", component: Product },
  { path: "/product/create", name: "product.create", component: ProductCreate },
  { path: "/product/:code/edit", name: "product.edit", component: ProductEdit },
  { path: "/apriori", name: "apriori", component: Apriori },
  { path: "/apriori/create", name: "apriori.create", component: AprioriCreate },
  { path: "/apriori/:code", name: "apriori.detail", component: AprioriDetail },
  { path: "/user", name: "user", component: User },
  { path: "/user/create", name: "user.create", component: UserCreate },
  { path: "/user/:id/edit", name: "user.edit", component: UserEdit },
  { path: "/profile", name: "profile", component: Profile },
  { path: "/:pathMatch(.*)", name: "NotFound", component: NotFound },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

router.beforeEach( async (to) => {
  if (Object.keys(authHeader()).length === 0 && to.name !== 'login' && to.name !== 'register' && to.name !== 'forgot-password' && to.name !== 'reset-password') {
    return { name: 'login' }
  } else if (Object.keys(authHeader()).length > 0 && to.name === 'login') {
      return { name: 'admin' }
  } else if (Object.keys(authHeader()).length > 0 && to.name === 'register') {
      return { name: 'admin' }
  } else if (Object.keys(authHeader()).length > 0 && to.name === 'forgot-password') {
      return { name: 'admin' }
  } else if (Object.keys(authHeader()).length > 0 && to.name === 'reset-password') {
      return { name: 'admin' }
  }

  if (to.name !== 'login' && to.name !== 'register' && to.name !== 'forgot-password' && to.name !== 'reset-password') {
        axios.get("http://localhost:3000/api/auth/token", { headers: authHeader() })
            .catch(() => {
                let refreshToken = {
                    refresh_token: localStorage.getItem("refresh-token")
                }
                axios.post("http://localhost:3000/api/auth/refresh",refreshToken)
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
