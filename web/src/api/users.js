import client from "./client";

export default {
  get(id) {
    return client.get(`/users/` + id);
  },
  listRatings(id) {
    return client.get(`/users/` + id + `/ratings`);
  },
  listVehicles(id) {
    return client.get(`/users/` + id + `/vehicles`);
  },
  bookTrip(userId, tripId) {
    return client.post(`/users/` + userId + `/trips/` + tripId);
  },
  cancelTrip(userId, tripId) {
    return client.put(`/users/` + userId + `/trips/` + tripId);
  }
};
