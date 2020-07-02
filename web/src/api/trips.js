import client from "./client";

export default {
  list() {
    return client.get(`/trips`);
  },
  get(id) {
    return client.get(`/trips/` + id);
  },
  create(trip) {
    return client.post(`/trips`, trip);
  },
  update(id, trip) {
    return client.put(`/trips/` + id, trip);
  },
  delete(id) {
    return client.delete(`/trips/` + id);
  }
};
