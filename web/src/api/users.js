import client from "./client";

export default {
  get(id) {
    return client.get(`/users/` + id);
  }
};
