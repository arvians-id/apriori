import axios from "axios";

export default function authHeader() {
    let token = localStorage.getItem('token');

    if (token) {
        axios.get(`${process.env.VUE_APP_SERVICE_URL}/auth/token`, {
            headers: {
                'Authorization': `Bearer ${token}`,
                'X-API-KEY': '31dd69599f4928d1500020a94b6ac391'
            }
        }).catch(() => {
            alert("You are not authorized to access this page, please login again");
            localStorage.removeItem("token")
            localStorage.removeItem("refresh-token")
            window.location.reload()
        })
    }

    if (token) {
        return {
            'Authorization': 'Bearer ' + token,
            'X-API-KEY': '31dd69599f4928d1500020a94b6ac391'
        };
    } else {
        return { 'X-API-KEY': '31dd69599f4928d1500020a94b6ac391' };
    }
}

