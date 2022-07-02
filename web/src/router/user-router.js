import Dashboard from "@/views/Dashboard";
import User from "@/views/user/Data";
import UserCreate from "@/views/user/Create";
import UserEdit from "@/views/user/Edit";
import Profile from "@/views/user/Profile";

const UserRouter = [
    { path: "/home", name: "admin", component: Dashboard },
    { path: "/user", name: "user", component: User },
    { path: "/user/create", name: "user.create", component: UserCreate },
    { path: "/user/:id/edit", name: "user.edit", component: UserEdit },
    { path: "/profile", name: "profile", component: Profile },
]
export default UserRouter;