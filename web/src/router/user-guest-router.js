import LandingPage from "@/views/user-guest/LandingPage";
import ProductList from "@/views/user-guest/ProductList";
import Cart from "@/views/user-guest/Cart";
import ProductDetail from "@/views/user-guest/ProductDetail";
import RecommendationDetail from "@/views/user-guest/RecommendationDetail";

const userGuestRouter = [
    { path: "/", name: "guest.index", component: LandingPage },
    { path: "/product-list", name: "guest.product", component: ProductList },
    { path: "/product-list/:code", name: "guest.product.detail", component: ProductDetail },
    { path: "/recommendation/:code/id/:id", name: "guest.recommendation.detail", component: RecommendationDetail },
    { path: "/my-cart", name: "guest.cart", component: Cart },
]

export default userGuestRouter;