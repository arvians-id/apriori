import axios from "axios";
import authHeader from "@/service/auth-header";

export default async function getRoles() {
    let getRole = null;
    if(authHeader()["Authorization"] !== undefined) {
        getRole = await axios.get("/api/profile", {headers: authHeader()}).catch((err) => {
            console.log(err);
        });
    }

    return getRole !== null ? getRole.data.data : null;
}