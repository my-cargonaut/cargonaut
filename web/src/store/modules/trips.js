import tripsAPI from "@/api/trips";

const trips = {
  namespaced: true,

  state: {
    loading: false,
    trips: [],
    rating: {}
  },

  mutations: {
    SET_LOADING(state, loading) {
      state.loading = loading;
    },
    SET_TRIPS(state, trips) {
      state.trips = trips;
    },
    SET_RATING(state, rating) {
      state.rating = rating;
    }
  },

  actions: {
    list({ commit }) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        tripsAPI
          .list()
          .then(response => {
            const trips = response.data;
            commit("SET_TRIPS", trips);
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

    create({ commit }, trip) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        tripsAPI
          .create(trip)
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

    update({ commit }, { id, trip }) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        tripsAPI
          .update(id, trip)
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

    delete({ commit }, id) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        tripsAPI
          .delete(id)
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

    getRating({ commit }, id) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        tripsAPI
          .getRating(id)
          .then(response => {
            const rating = response.data;
            commit("SET_RATING", rating);
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

    createRating({ commit }, { id, rating }) {
      return new Promise((resolve, reject) => {
        commit("SET_LOADING", true);
        tripsAPI
          .createRating(id, rating)
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
    trips: state => (state.trips ? state.trips : []),
    rating: state => (state.rating ? state.rating : {})
  }
};

export default trips;

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
