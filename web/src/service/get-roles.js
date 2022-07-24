import axios from "axios";
import authHeader from "@/service/auth-header";

export default async function getRoles() {
    let getRole = null;
    if(authHeader()["Authorization"] !== undefined) {
        getRole = await axios.get(`${process.env.VUE_APP_SERVICE_URL}/profile`, { headers: authHeader() })
    }

    return getRole !== null ? getRole.data.data : null;
}