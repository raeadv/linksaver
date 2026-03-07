import { netState } from "@/store/net-state"
import { createDiscreteApi } from "naive-ui"
import { useRouter } from "vue-router"


export function appUtils() {
    const router = useRouter()
    const net = netState()

    const loading = net.isLoading
    const { message, notification, modal, dialog, loadingBar } = createDiscreteApi(
        ['message', 'dialog', 'notification', 'loadingBar', 'modal'],
    )

    return { message, notification, modal, dialog, loading, loadingBar, router }
}