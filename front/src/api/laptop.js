import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8081",
});

export const getLaptopList = async () => {
  const res = await api.get("/laptop/list");
  console.log(res);
  return res.data.data;
};