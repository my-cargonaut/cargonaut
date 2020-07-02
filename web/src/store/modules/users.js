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
      ratings: [],
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
          .catch(e => {
            commit("alert/SET", getAlert(e), { root: true });
            reject(e);
          });

        usersAPI
          .listRatings(id)
          .then(response => {
            const ratings = response.data;
            commit("SET_RATINGS", ratings);
            resolve(response);
          })
          .catch(e => {
            commit("alert/SET", getAlert(e), { root: true });
            reject(e);
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
          .catch(e => {
            commit("alert/SET", getAlert(e), { root: true });
            reject(e);
          });
      });
    }
  },

  getters: {
    user: state => (state.user ? state.user : {}),
    ratings: state => (state.ratings ? state.ratings : [])
  }
};

function getAlert(e) {
  const alert = {
    type: "error",
    message: "Something went wrong!",
    title: true
  };
  if (e.response && e.response.data && e.response.data.error) {
    alert.message = e.response.data.error;
  }
  return alert;
}

export default users;
