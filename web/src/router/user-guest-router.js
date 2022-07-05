import LandingPage from "@/views/user-guest/LandingPage";
import ProductList from "@/views/user-guest/ProductList";
import Cart from "@/views/user-guest/Cart";

const userGuestRouter = [
    { path: "/", name: "guest.index", component: LandingPage },
    { path: "/product-list", name: "guest.product", component: ProductList },
    { path: "/my-cart", name: "guest.cart", component: Cart },
]

export default userGuestRouter;