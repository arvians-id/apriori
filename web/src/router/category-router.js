import Category from "@/views/category/Data";
import CategoryCreate from "@/views/category/Create";
import CategoryEdit from "@/views/category/Edit";

const categoryRouter = [
    { path: "/category", name: "category", component: Category },
    { path: "/category/create", name: "category.create", component: CategoryCreate },
    { path: "/category/:id/edit", name: "category.edit", component: CategoryEdit },
]

export default categoryRouter;