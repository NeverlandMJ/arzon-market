import TokenService from "./TokenServis"

const http = axios.create({
    baseURL: process.env.ARZON_MARKET_API_URL,
    headers: {
        'Content-type': "application/json"
    }
});

http.inteceptors.request.use((config) => {
    const token = localStorage.getItem('token')

    if (token) {
        config.headers["Authorization"] = `Bearer ${token}`
        config.headers['accept'] = 'application/json'
    }

    return config
}, (error) => Promise.reject(error))


http.interceptors.response.use(
    (res) => res,
    (error) => {
        if (
            error.response &&
            (error.response.status === 401 || error.response.status === 403)
        ) {
            TokenService.removeToken();
            router.push({
                name: "login"
            });
        }
        return Promise.reject(error);
    }
);

http.interceptors.response.use((response) => response, (error) => error);

export default http;