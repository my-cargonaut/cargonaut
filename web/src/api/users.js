import client from "./client";

export default {
  get(id) {
    return client.get(`/users/` + id);
  },
  listRatings(id) {
    return client.get(`/users/` + id + `/ratings`);
  },
  createRating(id, comment, value) {
    return client.post(`/users/` + id + `/ratings`, {
      comment: comment,
      value: value
    });
  },
  listVehicles(id) {
    return client.get(`/users/` + id + `/vehicles`);
  }
};
