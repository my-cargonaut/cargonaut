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
  },
  register(email, password, display_name, birthday, avatar) {
    return client.post(`/auth/register`, {
      email: email,
      password: password,
      display_name: display_name,
      birthday: birthday,
      avatar: avatar
    });
  }
};
