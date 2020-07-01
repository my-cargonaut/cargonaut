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
  listVehicles(userID) {
    return client.get(`/users/` + userID + `/vehicles`);
  },
  createVehicle(userID, vehicle) {
    return client.post(`/users/` + userID + `/vehicles`, vehicle);
  },
  updateVehicle(userID, vehicleId, vehicle) {
    return client.put(`/users/` + userID + `/vehicles/` + vehicleId, vehicle);
  },
  deleteVehicle(userID, vehicleId) {
    return client.delete(`/users/` + userID + `/vehicles/` + vehicleId);
  }
};
