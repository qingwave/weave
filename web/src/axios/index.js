import axios from "axios";

// vite proxy, need not cors
// axios.defaults.baseURL = "http://127.0.0.1:8080";
// axios.defaults.withCredentials = true;

const request = axios.create({
    timeout: 3000,
    headers: {
        'Content-Type': 'application/json;charset=UTF-8'
    }
});

export default request;