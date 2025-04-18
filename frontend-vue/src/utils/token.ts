// JWT Token管理工具类

const TOKEN_KEY = 'jwt_token'

export const getToken = (): string | null => {
    return localStorage.getItem(TOKEN_KEY)
}

export const setToken = (token: string): void => {
    localStorage.setItem(TOKEN_KEY, token)
}

export const removeToken = (): void => {
    localStorage.removeItem(TOKEN_KEY)
}

export const hasToken = (): boolean => {
    return !!getToken()
}