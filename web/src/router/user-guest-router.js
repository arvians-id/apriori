import LandingPage from "@/views/user-guest/LandingPage";
import ProductList from "@/views/user-guest/ProductList";
import Cart from "@/views/user-guest/Cart";
import ProductDetail from "@/views/user-guest/ProductDetail";
import RecommendationDetail from "@/views/user-guest/RecommendationDetail";
import Order from "@/views/user-guest/Order";
import OrderDetail from "@/views/user-guest/OrderDetail";
import Profile from "@/views/user-guest/Profile";

const userGuestRouter = [
    { path: "/", name: "guest.index", component: LandingPage },
    { path: "/product-list", name: "guest.product", component: ProductList },
    { path: "/product-list/:code", name: "guest.product.detail", component: ProductDetail },
    { path: "/recommendation/:code/id/:id", name: "guest.product.recommendation", component: RecommendationDetail },
    { path: "/checkout", name: "guest.cart", component: Cart },
    { path: "/my-order", name: "member.order", component: Order },
    { path: "/my-order/:order_id", name: "member.order.detail", component: OrderDetail },
    { path: "/my-profile", name: "member.profile", component: Profile },
]

export default userGuestRouter;