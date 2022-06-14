import { createRouter, createWebHistory } from "vue-router";
import Login from "../views/auth/Login.vue";
import Register from "../views/auth/Register.vue";
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
import User from "../views/user/Data.vue";
import UserCreate from "../views/user/Create.vue";
import UserEdit from "../views/user/Edit.vue";
import Profile from "../views/Profile.vue";
import NotFound from "../views/NotFound.vue";

const routes = [
  { path: "/", name: "login", component: Login },
  { path: "/register", name: "register", component: Register },
  { path: "/admin", name: "admin", component: Dashboard },
  { path: "/admin/transaction", name: "transaction", component: Transaction },
  { path: "/admin/transaction/create", name: "transaction.create", component: TransactionCreate },
  { path: "/admin/transaction/create/csv", name: "transaction.create.csv", component: TransactionCreateCSV },
  { path: "/admin/transaction/:no_transaction/edit", name: "transaction.edit", component: TransactionEdit },
  { path: "/admin/product", name: "product", component: Product },
  { path: "/admin/product/create", name: "product.create", component: ProductCreate },
  { path: "/admin/product/:code/edit", name: "product.edit", component: ProductEdit },
  { path: "/admin/apriori", name: "apriori", component: Apriori },
  { path: "/admin/apriori/create", name: "apriori.create", component: AprioriCreate },
  { path: "/admin/user", name: "user", component: User },
  { path: "/admin/user/create", name: "user.create", component: UserCreate },
  { path: "/admin/user/:id/edit", name: "user.edit", component: UserEdit },
  { path: "/admin/profile", name: "profile", component: Profile },
  { path: "/:pathMatch(.*)", name: "NotFound", component: NotFound },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
