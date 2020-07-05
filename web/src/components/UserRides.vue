<template>
  <v-card>
    <v-card-title>
      My Rides
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        append-icon="mdi-magnify"
        label="Search"
        single-line
        hide-details
      ></v-text-field>
    </v-card-title>

    <v-card-text>
      <v-data-table
        :loading="loading"
        :headers="headers"
        :items="rides"
        :search="search"
        show-expand
        single-expand
        @item-expanded="getTripVehicle"
      >
        <template v-slot:item.user_id="{ item }">
          <v-btn icon :to="'/users/' + item.user_id">
            <v-avatar size="36px">
              <v-img
                :src="'/api/v1/users/' + item.user_id + '/avatar'"
                alt="Trip user avatar"
              ></v-img>
            </v-avatar>
          </v-btn>
        </template>
        <template v-slot:item.action="{ item }">
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-icon
                small
                class="mr-2"
                v-on="on"
                :disabled="!isWaitingForStart(item)"
                @click="cancelTrip(item)"
                >mdi-cancel</v-icon
              >
            </template>
            <span>Cancel Trip</span>
          </v-tooltip>
          <v-dialog v-model="dialog" max-width="750px">
            <template v-slot:activator="{ on }">
              <v-icon
                small
                class="mr-2"
                v-on="on"
                :disabled="!isCompleted(item)"
                @click="rateTrip(item)"
                >mdi-star</v-icon
              >
            </template>

            <v-card>
              <v-card-title>
                <span class="headline">Rate Trip</span>
              </v-card-title>

              <v-card-text>
                <v-form ref="form" v-model="valid">
                  <v-container>
                    <v-row>
                      <v-rating
                        v-model="value"
                        color="amber"
                        dense
                        hover
                        half-increments
                        size="28"
                      ></v-rating>
                    </v-row>
                    <v-row class="mt-4">
                      <v-textarea
                        v-model="comment"
                        label="Comment"
                        :rules="rules"
                        counter
                        maxlength="200"
                        required
                      ></v-textarea>
                    </v-row>
                  </v-container>
                </v-form>
              </v-card-text>

              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="primary" text @click="close">Cancel</v-btn>
                <v-btn
                  color="primary"
                  text
                  @click="saveRating"
                  :disabled="!valid"
                  >Save</v-btn
                >
              </v-card-actions>
            </v-card>
          </v-dialog>
        </template>
        <template v-slot:item.status="{ item }">
          {{ getTripStatus(item) }}
        </template>
        <template v-slot:expanded-item="{ headers }">
          <td :colspan="headers.length">
            <v-chip class="ma-2" color="indigo" text-color="white">
              <v-avatar left>
                <v-icon>mdi-car</v-icon>
              </v-avatar>
              {{ vehicle.brand }} {{ vehicle.model }}
            </v-chip>
            <v-chip
              class="ma-2"
              color="green"
              text-color="white"
              :active="vehicle.passengers > 0"
            >
              <v-avatar left>
                <v-icon>mdi-account-circle</v-icon>
              </v-avatar>
              Passengers
              <v-avatar right class="green darken-4">
                {{ vehicle.passengers }}
              </v-avatar>
            </v-chip>
            <v-chip
              class="ma-2"
              color="orange"
              text-color="white"
              :active="vehicle.loading_area_length > 0"
            >
              {{ vehicle.loading_area_length }}cm x
              {{ vehicle.loading_area_width }}cm
              <v-icon right>mdi-crop</v-icon>
            </v-chip>
          </td>
        </template>
      </v-data-table>
    </v-card-text>
  </v-card>
</template>

<script>
import { mapActions } from "vuex";
import { mapGetters } from "vuex";

import TripStatus from "@/shared/trip_status";

export default {
  name: "Trips",

  computed: {
    ...mapGetters("auth", ["authId"]),
    ...mapGetters("trips", {
      trips: "trips",
      tripsLoading: "loading"
    }),
    ...mapGetters("users", {
      usersLoading: "loading"
    }),
    ...mapGetters("vehicles", {
      vehicle: "vehicle",
      vehicles: "vehicles",
      vehiclesLoading: "loading"
    }),

    loading() {
      return this.tripsLoading && this.vehiclesLoading && this.usersLoading;
    },

    rides() {
      return this.trips.filter(trip => {
        return trip.rider_id === this.authId;
      });
    }
  },

  data: () => ({
    headers: [
      {
        text: "User",
        value: "user_id",
        sortable: true,
        align: "start"
      },
      {
        text: "From",
        value: "start",
        sortable: true,
        align: "start"
      },
      { text: "To", value: "destination", sortable: true },
      { text: "Price (€)", value: "price", sortable: true },
      { text: "Status", value: "status", sortable: false },
      { text: "Actions", value: "action", sortable: false },
      { text: "", value: "data-table-expand" }
    ],
    search: "",
    dialog: false,
    valid: false,
    rules: [
      v => !!v || "Field is required",
      v => (v && v.length) >= 25 || "Field must be at least 25 characters"
    ],
    ratedIndex: -1,
    value: 0,
    comment: ""
  }),

  methods: {
    ...mapActions("trips", { listTrips: "list", createRating: "createRating" }),
    ...mapActions("users", { cancelTripForRider: "cancelTrip" }),
    ...mapActions("vehicles", { getVehicle: "get" }),

    getTripVehicle({ item, value }) {
      if (!value) return;
      this.getVehicle(item.vehicle_id);
    },

    cancelTrip(item) {
      this.cancelTripForRider({
        userId: this.authId,
        tripId: item.id
      }).then(() => this.listTrips());
    },

    rateTrip(trip) {
      this.ratedIndex = this.trips.indexOf(trip);
      this.dialog = true;
    },

    saveRating() {
      this.createRating({
        id: this.trips[this.ratedIndex].id,
        rating: {
          comment: this.comment,
          value: this.value
        }
      });
      this.close();
    },

    close() {
      this.dialog = false;
      setTimeout(() => {
        this.ratedIndex = -1;
        this.$refs.form.reset();
        this.rating = 0;
      }, 300);
    },

    getTripStatus(trip) {
      return TripStatus.get(trip);
    },

    isWaitingForStart(trip) {
      return TripStatus.isWaitingForStart(trip);
    },

    isCompleted(trip) {
      return TripStatus.isCompleted(trip);
    }
  }
};
</script>
