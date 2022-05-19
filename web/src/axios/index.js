import axios from "axios";
import { ElMessage } from "element-plus";

// vite proxy, need not cors
// axios.defaults.baseURL = "http://127.0.0.1:8080";
// axios.defaults.withCredentials = true;

const request = axios.create({
    timeout: 3000,
    headers: {
        'Content-Type': 'application/json;charset=UTF-8'
    }
});

request.interceptors.response.use(response => {
    return response
}, error => {
    let msg;
    if (error.response) {
        if (error.response.data.msg) {
            msg = `${error.response.data.msg}`;
        } else if (error.response.status) {
            switch (error.response.status) {
                case 400:
                    msg = `bad request`;
                    break;
                case 401:
                    msg = `request unauthorized`;
                    break;
                case 403:
                    msg = `request forbidden`;
                    break;
                case 404:
                    msg = `request not found`;
                    break;
                case 405:
                    msg = `unsupported request`;
                    break;
                case 500:
                    msg = `internal server error`;
                    break;
                default:
                    msg = `request failed with ${error.response.statusText}`;
            }
        }
    } else if (error.code) {
        switch (error.code) {
            case "ECONNABORTED":
                msg = `request timeout`;
                break;
            default:
                msg = `request failed with code ${error.code}`;
        }
    } else {
        msg = `network failed with error ${error}`
    }

    ElMessage.error(msg);

    return Promise.reject(error);
});

export default request;
