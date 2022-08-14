import axios from "axios";
import authHeader from "@/service/auth-header";

export default async function getRoles() {
    let getRole = null;
    if(authHeader()["Authorization"] !== undefined) {
        try {
            getRole = await axios.get(`${process.env.VUE_APP_SERVICE_URL}/profile`, { headers: authHeader() })
        } catch (error) {
            alert("You are not authorized to access this page, please login again");
            localStorage.removeItem("token")
            localStorage.removeItem("refresh-token")
            window.location.reload()
        }
    }

    return getRole !== null ? getRole.data.data : null;
}