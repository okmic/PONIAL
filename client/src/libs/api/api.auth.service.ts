import axios, { type AxiosInstance } from 'axios'
import appconfig from '../../appconfig'
import store from '../../store/store'
import { authStatus } from '../../store/slices/auth.slice'
import { getAuthHeader, handlerError } from './api.util'

class ApiAuthService {
    private axiosInstance: AxiosInstance
    constructor() {
        this.axiosInstance = axios.create({
            baseURL: appconfig.backendUrl,
            headers: {
                'Content-Type': 'application/json',
            },
        })
    }
    async signin(payload: { email: string, password: string }) {
        return await this.axiosInstance.post<any>(
            `/api/v1/auth/signin`,
            payload
        ).then(r => {
            store.dispatch(
                authStatus({ status: "loading", user: r.data.user })
            )
            localStorage.setItem('token', r.data.token)
            return r.data
        })
        .catch(e => handlerError(e))
    }
    async signup(payload: { name: string, email: string, password: string, role: "admin" | "user", vin: string }) {
        return await this.axiosInstance.post<any>(
            `/api/v1/auth/signup`,
            {...payload, adminSecret: appconfig.adminSecret }
        ).then(r => {
            store.dispatch(
                authStatus({ status: "loading", user: r.data.user })
            )
            localStorage.setItem('token', r.data.token)
        })
            .catch(e => handlerError(e))
    }
    public async geMyInfo() {
        return await this.axiosInstance.get<any>(
            `/api/v1/users/me`,
            { headers: getAuthHeader() },
        )
            .then(r => r.data)
            .catch(e => handlerError(e))
    }
}

export default new ApiAuthService()
