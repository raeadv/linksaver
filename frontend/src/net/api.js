import { useAuthStore } from "@/store/authStore";
import { netState } from "@/store/net-state";
import { router } from "@/router/index";
import axios from "axios";
import { differenceInSeconds } from "date-fns";
import { createDiscreteApi } from "naive-ui";

const getNetState = () => netState();



const baseUrl = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';
const api = axios.create({ baseURL: baseUrl })
api.showErrMessage = true

let isRefreshing = false;
let queue = [];


const { message, loadingBar } = createDiscreteApi(["message", "loadingBar"]);

api.interceptors.request.use(async (c) => {
    const net = getNetState();
    net.startLoading()
    loadingBar.start();

    const auth = useAuthStore();
    if (auth.access_token) {
        c.headers.Authorization = `Bearer ${auth.access_token}`;
        c.headers['Access-Control-Allow-Origin'] = "http://localhost:3000"

        // renew if token exp is less than 2 minutes (120 secs), skip on auth routes
        const isAuthRoute = ['/login', '/register'].some(p => c.url?.includes(p))
        const tokenExpiration = localStorage.getItem('expires_at')
        const secondDiffs = differenceInSeconds(tokenExpiration, new Date())
        if (secondDiffs < 120 && !isAuthRoute) {
            await renewToken();
            c.headers.Authorization = `Bearer ${auth.access_token}`;
        }
    }

    return c
}, err => {
    const net = getNetState();
    setTimeout(() => {
        net.stopLoading()
        loadingBar.finish();
    }, 300);
    return err
})


api.interceptors.response.use((ok) => {
    const net = getNetState();
    setTimeout(() => {
        net.stopLoading()
        loadingBar.finish();
    }, 300);

    return ok
}, (err) => {
    const net = getNetState();
    setTimeout(() => {
        loadingBar.finish();
    }, 300);

    const isAuthRoute = ['/login', '/register'].some(p => err.config?.url?.includes(p))
    if (err.response?.status === 401 && !isAuthRoute) {
        renewToken(err)
    }

    const { response } = err
    if (response?.data) {
        const responseData = response.data
        api?.showErrMessage && message.error(responseData.message + ': ' + response.status)
    } else {
        api?.showErrMessage && message.error(err.message + ': ' + err?.status)
    }
    net.stopLoading()

    return err
})



export async function renewToken(err) {
    const auth = useAuthStore();

    if (err) {
        const original = err.config;

        if (err.response?.status !== 401 || original?._retry) {
            return Promise.reject(err);
        }

        if (isRefreshing) {
            return new Promise((resolve) => {
                queue.push((token) => {
                    original.headers.Authorization = `Bearer ${token}`;
                    resolve(api(original));
                });
            });
        }

        original._retry = true;
    } else {
        // Proactive refresh from request interceptor — queue if already in progress
        if (isRefreshing) {
            return new Promise((resolve) => {
                queue.push(() => resolve());
            });
        }
    }

    isRefreshing = true;

    try {
        const r = await axios.post(
            `${baseUrl}/api/auth/refresh-token`,
            {},
            { withCredentials: true }
        );

        const newAccessToken = (r.data?.token?.access_token ?? r.data?.access_token);
        const newExpiresAt = (r.data?.token?.expires_at ?? r.data?.expires_at);
        auth.access_token = newAccessToken;
        localStorage.setItem("access_token", newAccessToken);
        localStorage.setItem("expires_at", newExpiresAt);

        queue.forEach((cb) => cb(newAccessToken));
        queue = [];

        if (err) {
            const original = err.config
            original.headers.Authorization = `Bearer ${newAccessToken}`;
            return api(original);
        }
    } catch (refreshErr) {
        auth.clearAuth();
        message.info("Sesi Login habis.")
        router.replace('/login')
        return Promise.reject(refreshErr);
    } finally {
        isRefreshing = false;
    }
}


document.addEventListener('visibilitychange', async () => {
    if (document.visibilityState !== 'visible') return
    const auth = useAuthStore()
    if (!auth.access_token) return
    const expiresAt = localStorage.getItem('expires_at')
    if (!expiresAt) return
    const secondDiffs = differenceInSeconds(expiresAt, new Date())
    if (secondDiffs < 120) {
        await renewToken()
    }
})


export default api