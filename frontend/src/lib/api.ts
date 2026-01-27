import axios from "axios";
import { goto } from "$app/navigation";

export const api = axios.create({
    baseURL: "http://localhost:8080",
    withCredentials: true
})

api.interceptors.response.use(
    response => response,
    error => {
        if (error.response && error.response.status === 401) {
            goto('/login');
        }
        return Promise.reject(error);
    }
)
