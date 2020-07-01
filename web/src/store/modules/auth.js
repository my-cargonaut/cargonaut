import authAPI from "@/api/auth";
import client from "@/api/client"; // Move this to Vue.$http

const auth = {
  namespaced: true,

  state: {
    token: localStorage.getItem("token") || "",
    expiry: +new Date(localStorage.getItem("token_expiry")) || 0,
    user: parseUserFromJWT(localStorage.getItem("token"))
  },

  mutations: {
    LOGIN(state, token) {
      state.token = token.token;
      state.expiry = +new Date(token.token_expiry);
      state.user = parseUserFromJWT(token.token);
    },
    LOGOUT(state) {
      state.token = "";
      state.expiry = 0;
      state.user = null;
    }
  },

  actions: {
    login({ commit }, auth) {
      return new Promise((resolve, reject) => {
        authAPI
          .login(auth.username, auth.password)
          .then(response => {
            const token = response.data;
            localStorage.setItem("token", token.token);
            localStorage.setItem("token_expiry", token.token_expiry);
            client.defaults.headers.common[
              "Authorization"
            ] = `Bearer ${token.token}`;
            commit("LOGIN", token);
            resolve(response);
          })
          .catch(error => {
            localStorage.removeItem("token");
            localStorage.removeItem("token_expiry");
            delete client.defaults.headers.common["Authorization"];
            commit(
              "alert/SET",
              {
                message: error.response.data.error,
                type: "error",
                title: true
              },
              { root: true }
            );
            reject(error);
          });
      });
    },
    refresh({ commit, state }) {
      return new Promise((resolve, reject) => {
        authAPI
          .refresh(state.token)
          .then(response => {
            const token = response.data;
            localStorage.setItem("token", token.token);
            localStorage.setItem("token_expiry", token.token_expiry);
            client.defaults.headers.common[
              "Authorization"
            ] = `Bearer ${token.token}`;
            commit("LOGIN", token);
            resolve(response);
          })
          .catch(error => {
            localStorage.removeItem("token");
            localStorage.removeItem("token_expiry");
            delete client.defaults.headers.common["Authorization"];
            commit(
              "alert/SET",
              {
                message: error.response.data.error,
                type: "error",
                title: true
              },
              { root: true }
            );
            reject(error);
          });
      });
    },
    logout({ commit, state }) {
      return new Promise((resolve, reject) => {
        authAPI
          .logout(state.token)
          .then(resolve())
          .catch(error => {
            commit(
              "alert/SET",
              {
                message: error.response.data.error,
                type: "error",
                title: true
              },
              { root: true }
            );
            reject(error);
          })
          .finally(() => {
            localStorage.removeItem("token");
            localStorage.removeItem("token_expiry");
            delete client.defaults.headers.common["Authorization"];
            commit("LOGOUT");
          });
      });
    }
  },

  getters: {
    isLoggedIn: state => !!state.token && !!state.expiry,
    authId: state => (state.user ? state.user.id : ""),
    authEmail: state => (state.user ? state.user.email : ""),
    authName: state => (state.user ? state.user.name : ""),
    tokenExpiry: state => state.expiry
  }
};

export default auth;

function parseUserFromJWT(token) {
  if (!token) {
    return null;
  }
  let base64Url = token.split(".")[1];
  let base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
  let jsonPayload = decodeURIComponent(
    atob(base64)
      .split("")
      .map(c => {
        return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
      })
      .join("")
  );
  let payload = JSON.parse(jsonPayload);
  return payload.user;
}
