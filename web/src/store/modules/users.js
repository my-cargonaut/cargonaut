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
    },

    ratings: {
      ratings: null,
      count: 0,
      average: 0
    }
  },

  mutations: {
    SET_USER(state, user) {
      state.user = user;
    },
    SET_RATINGS(state, ratings) {
      state.ratings.ratings = ratings;
      state.ratings.count = ratings.length;

      if (ratings.length > 0)
        state.ratings.average =
          ratings.reduce((total, next) => total + next.value, 0) /
          ratings.length;
    },
    INC_RATING(state, value) {
      state.ratings.count++;

      if (state.ratings.ratings.length > 0) {
        state.ratings.ratings.push({
          value: value
        });
        state.ratings.average =
          state.ratings.ratings.reduce((total, next) => total + next.value, 0) /
          state.ratings.ratings.length;
      }
    }
  },

  actions: {
    get({ commit }, id) {
      return new Promise((resolve, reject) => {
        usersAPI
          .get(id)
          .then(response => {
            const user = response.data;
            commit("SET_USER", user);
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

        usersAPI
          .listRatings(id)
          .then(response => {
            const ratings = response.data;
            commit("SET_RATINGS", ratings);
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
    },

    rate({ commit }, { id, value }) {
      return new Promise((resolve, reject) => {
        usersAPI
          .createRating(id, "", value)
          .then(response => {
            commit("INC_RATING", value);
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
    user: state => (state.user ? state.user : null),
    ratings: state => (state.ratings ? state.ratings : null)
  }
};

export default users;
