// const TOKEN_KEY = 'Authorization';

const TokenService = {
    removeToken() {
        localStorage.removeItem(token)
    }
}

export default TokenService