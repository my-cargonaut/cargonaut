import usersAPI from "@/api/users";

const users = {
  namespaced: true,

  state: {
    loading: false,
    user: {
      email: "",
      display_name: "",
      birthday: "",
      created_at: "",
      updated_at: ""
    },
    vehicles: [],
    ratings: {
      ratings: [],
      count: 0,
      average: 0
    }
  },

  mutations: {
    SET_LOADING(state, loading) {
      state.loading = loading;
    },
    SET_USER(state, user) {
      state.user = user;
    },
    SET_VEHICLES(state, vehicles) {
      state.vehicles = vehicles;
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
        commit("SET_LOADING", true);
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
          })
          .finally(() => {
            commit("SET_LOADING", false);
          });
      });
    },

    listRatings({ commit }, id) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
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
          })
          .finally(() => {
            commit("SET_LOADING", false);
          });
      });
    },

    rate({ commit }, { id, value }) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        usersAPI
          .createRating(id, "", value)
          .then(response => {
            commit("INC_RATING", value);
            resolve(response);
          })
          .catch(e => {
            commit("alert/SET", getAlert(e), { root: true });
            reject(e);
          })
          .finally(() => {
            commit("SET_LOADING", false);
          });
      });
    },

    listVehicles({ commit }, id) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        usersAPI
          .listVehicles(id)
          .then(response => {
            const vehicles = response.data;
            commit("SET_VEHICLES", vehicles);
            resolve(response);
          })
          .catch(e => {
            commit("alert/SET", getAlert(e), { root: true });
            reject(e);
          })
          .finally(() => {
            commit("SET_LOADING", false);
          });
      });
    },

    bookTrip({ commit }, { userId, tripId }) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        usersAPI
          .bookTrip(userId, tripId)
          .then(response => {
            resolve(response);
          })
          .catch(e => {
            commit("alert/SET", getAlert(e), { root: true });
            reject(e);
          })
          .finally(() => {
            commit("SET_LOADING", false);
          });
      });
    },

    cancelTrip({ commit }, { userId, tripId }) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        usersAPI
          .cancelTrip(userId, tripId)
          .then(response => {
            resolve(response);
          })
          .catch(e => {
            commit("alert/SET", getAlert(e), { root: true });
            reject(e);
          })
          .finally(() => {
            commit("SET_LOADING", false);
          });
      });
    }
  },

  getters: {
    loading: state => (state.loading ? state.loading : false),
    user: state => (state.user ? state.user : {}),
    vehicles: state => (state.vehicles ? state.vehicles : []),
    ratings: state => (state.ratings ? state.ratings : [])
  }
};

export default users;

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
