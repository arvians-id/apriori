import Data from "@/views/user-order/Data";
import Detail from "@/views/user-order/Detail";
import AddReceiptNumber from "@/views/user-order/AddReceiptNumber";

const UserRouter = [
    { path: "/order", name: "user-order", component: Data },
    { path: "/order/:order_id", name: "user-order.detail", component: Detail },
    { path: "/order/:order_id/add-receipt-number", name: "user-order.add-receipt-number", component: AddReceiptNumber },
]
export default UserRouter;