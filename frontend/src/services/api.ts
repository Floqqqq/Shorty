import axios from "axios";

const API = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || "http://localhost:8080",
  timeout: 5000,
});

export const shorten = (url: string) =>
  API.post("/api/v1/shorten", { url });

export const getStats = (code: string) =>
  API.get(`/api/v1/stats/${code}`);