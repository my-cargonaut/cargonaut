import client from "./client";

export default {
  list() {
    return client.get(`/vehicles`);
  },
  get(id) {
    return client.get(`/vehicles/` + id);
  },
  create(vehicle) {
    return client.post(`/vehicles`, vehicle);
  },
  update(id, vehicle) {
    return client.put(`/vehicles/` + id, vehicle);
  },
  delete(id) {
    return client.delete(`/vehicles/` + id);
  }
};
