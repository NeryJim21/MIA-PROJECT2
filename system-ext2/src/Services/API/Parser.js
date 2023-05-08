import axios from 'axios';

const instance = axios.create(
    {
        baseURL: 'http://192.168.0.24:5000',
        timeout: 15000,
        headers: {
            'content-Type': 'application/json'
        }
    }
)

export const parse = async (value) => {
    console.log(value)
    const { data } = await instance.post("/ejecutar", { cmd: value })
    console.log(data)
    return data
}

export const ping = async () => {
    const { data } = await instance.get("/ping")
    return data
}

export const ejecutar = async (value) => {
    const { data } = await instance.post("/ejecutar", { peticion: value })
    return data
}