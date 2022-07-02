import Apriori from "@/views/apriori/Data";
import AprioriCreate from "@/views/apriori/Create";
import AprioriDetail from "@/views/apriori/Detail";
import AprioriCodeDetail from "@/views/apriori/DetailApriori";
import AprioriEdit from "@/views/apriori/Edit";

const AprioriRouter = [
    { path: "/apriori", name: "apriori", component: Apriori },
    { path: "/apriori/create", name: "apriori.create", component: AprioriCreate },
    { path: "/apriori/:code", name: "apriori.detail", component: AprioriDetail },
    { path: "/apriori/:code/detail/:id", name: "apriori.code-detail", component: AprioriCodeDetail },
    { path: "/apriori/:code/edit/:id", name: "apriori.edit", component: AprioriEdit },
]

export default AprioriRouter;