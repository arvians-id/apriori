import Data from "@/views/user-order/Data";
import Detail from "@/views/user-order/Detail";

const UserRouter = [
    { path: "/order", name: "user-order", component: Data },
    { path: "/order/:order_id", name: "user-order.detail", component: Detail },
]
export default UserRouter;