import { createRouter, createWebHistory } from "vue-router";
import NotFound from "../views/NotFound.vue";
import authHeader from "@/service/auth-header";
import axios from "axios";
import AuthRouter from "@/router/auth-router";
import TransactionRouter from "@/router/transaction-router";
import ProductRouter from "@/router/product-router";
import CategoryRouter from "@/router/category-router";
import AprioriRouter from "@/router/apriori-router";
import UserRouter from "@/router/user-router";
import UserGuestRouter from "@/router/user-guest-router";
import UserOrder from "@/router/user-order";
import getRoles from "@/service/get-roles";

const routes = [
    ...UserGuestRouter,
    ...AuthRouter,
    ...TransactionRouter,
    ...ProductRouter,
    ...CategoryRouter,
    ...AprioriRouter,
    ...UserRouter,
    ...UserOrder,
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

  let getRole = await getRoles();
  if (getRole != null && (getRole.role == "2" &&
      (to.name.split(".")[0] === "transaction" ||
          to.name.split(".")[0] === "product" ||
          to.name.split(".")[0] === "apriori" ||
          to.name.split(".")[0] === "user" ||
          to.name.split(".")[0] === "admin" ||
          to.name.split(".")[0] === "profile" ||
          to.name.split(".")[0] === "user-order"))) {
      return { name: 'guest.index' }
  }

  if (to.name.split(".")[0] !== "auth" && to.name.split(".")[0] !== "guest" && authHeader()["Authorization"]) {
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/auth/token`, { headers: authHeader() })
          .catch(() => {
              alert("You are not authorized to access this page, please login again");
              localStorage.removeItem("token")
              localStorage.removeItem("refresh-token")
              window.location.reload()
          })
  }
})

export default router;
