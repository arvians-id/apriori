import LandingPage from "@/views/user-guest/LandingPage";
import ProductList from "@/views/user-guest/ProductList";
import Cart from "@/views/user-guest/Cart";
import ProductDetail from "@/views/user-guest/ProductDetail";
import RecommendationDetail from "@/views/user-guest/RecommendationDetail";
import Order from "@/views/user-guest/Order";
import OrderDetail from "@/views/user-guest/OrderDetail";
import Profile from "@/views/user-guest/Profile";
import HistoryOrder from "@/views/user-guest/HistoryOrder";
import HistoryOrderRate from "@/views/user-guest/HistoryOrderRate";
import Notification from "@/views/user-guest/Notification";

const userGuestRouter = [
    { path: "/", name: "guest.index", component: LandingPage },
    { path: "/product-list", name: "guest.product", component: ProductList },
    { path: "/product-list/:code", name: "guest.product.detail", component: ProductDetail },
    { path: "/recommendation/:code/id/:id", name: "guest.product.recommendation", component: RecommendationDetail },
    { path: "/checkout", name: "guest.cart", component: Cart },
    { path: "/my-order", name: "member.order", component: Order },
    { path: "/my-order/:order_id", name: "member.order.detail", component: OrderDetail },
    { path: "/my-profile", name: "member.profile", component: Profile },
    { path: "/my-history-order", name: "member.history", component: HistoryOrder },
    { path: "/my-history-order/:id_order/rate", name: "member.history.rate", component: HistoryOrderRate },
    { path: "/notification", name: "member.notification", component: Notification },
]

export default userGuestRouter;