import Login from "@/views/auth/Login";
import Register from "@/views/auth/Register";
import ForgotPassword from "@/views/auth/ForgotPassword";
import ResetPassword from "@/views/auth/ResetPassword";

const authRouter = [
    { path: "/auth/login", name: "auth.login", component: Login },
    { path: "/auth/register", name: "auth.register", component: Register },
    { path: "/auth/forgot-password", name: "auth.forgot-password", component: ForgotPassword },
    { path: "/auth/reset-password", name: "auth.reset-password", component: ResetPassword },
]

export default authRouter;