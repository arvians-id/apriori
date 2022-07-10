export default function authHeader() {
    let token = localStorage.getItem('token');
    if (token) {
        return {
            'Authorization': 'Bearer ' + token,
            'X-API-KEY': '31dd69599f4928d1500020a94b6ac391'
        };
    } else {
        return { 'X-API-KEY': '31dd69599f4928d1500020a94b6ac391' };
    }
}

