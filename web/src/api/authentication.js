import client from "./client";

export default {
  login(username, password) {
    return client.post(`/auth/login`, {
      username: username,
      password: password
    });
  },
  refresh(token) {
    return client.patch(`/auth/refresh`, {
      token: token
    });
  },
  logout(token) {
    return client.post(`/auth/logout`, {
      token: token
    });
  }
};
