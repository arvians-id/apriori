import Product from "@/views/product/Data";
import ProductCreate from "@/views/product/Create";
import ProductDetail from "@/views/product/Detail";
import ProductEdit from "@/views/product/Edit";

const productRouter = [
    { path: "/product", name: "product", component: Product },
    { path: "/product/create", name: "product.create", component: ProductCreate },
    { path: "/product/:code", name: "product.detail", component: ProductDetail },
    { path: "/product/:code/edit", name: "product.edit", component: ProductEdit },
]

export default productRouter;