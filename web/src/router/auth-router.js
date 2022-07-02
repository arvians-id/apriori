import Login from "@/views/auth/Login";
import Register from "@/views/auth/Register";
import ForgotPassword from "@/views/auth/ForgotPassword";
import ResetPassword from "@/views/auth/ResetPassword";

const authRouter = [
    { path: "/", name: "auth.login", component: Login },
    { path: "/register", name: "auth.register", component: Register },
    { path: "/forgot-password", name: "auth.forgot-password", component: ForgotPassword },
    { path: "/reset-password", name: "auth.reset-password", component: ResetPassword },
]

export default authRouter;