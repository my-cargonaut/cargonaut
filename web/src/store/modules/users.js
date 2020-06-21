import usersAPI from "@/api/users";

const users = {
  namespaced: true,

  state: {
    user: {
      email: "",
      display_name: "",
      birthday: "",
      created_at: "",
      updated_at: ""
    }
  },

  mutations: {
    SET(state, user) {
      state.user = user;
    }
  },

  actions: {
    get({ commit }, id) {
      return new Promise((resolve, reject) => {
        usersAPI
          .get(id)
          .then(response => {
            const user = response.data;
            commit("SET", user);
            resolve(response);
          })
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
          });
      });
    }
  },

  getters: {
    user: state => (state.user ? state.user : null)
  }
};

export default users;
