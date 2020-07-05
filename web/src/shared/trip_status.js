export default {
  status: Object.freeze({
    WAITING_FOR_RIDER: "Waiting for rider",
    WAITING_FOR_START: "Waiting for driver to start the trip",
    WAITING_FOR_STOP: "In transit",
    COMPLETED: "Completed",
    UNKNOWN: "Unkown status"
  }),

  get(trip) {
    if (this.isWaitingForRider(trip)) {
      return this.status.WAITING_FOR_RIDER;
    } else if (this.isWaitingForStart(trip)) {
      return this.status.WAITING_FOR_START;
    } else if (this.isWaitingForStop(trip)) {
      return this.status.WAITING_FOR_STOP;
    } else if (this.isCompleted(trip)) {
      return this.status.COMPLETED;
    }
    return this.status.UNKNOWN;
  },
  isWaitingForRider(trip) {
    return (
      !trip.rider_id &&
      +new Date(trip.depature) < 0 &&
      +new Date(trip.arrival) < 0
    );
  },
  isWaitingForStart(trip) {
    return (
      trip.rider_id &&
      +new Date(trip.depature) < 0 &&
      +new Date(trip.arrival) < 0
    );
  },
  isWaitingForStop(trip) {
    return (
      trip.rider_id &&
      +new Date(trip.depature) > 0 &&
      +new Date(trip.arrival) < 0
    );
  },
  isCompleted(trip) {
    return (
      trip.rider_id &&
      +new Date(trip.depature) > 0 &&
      +new Date(trip.arrival) > 0
    );
  }
};
