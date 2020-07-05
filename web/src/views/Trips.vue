<template>
  <v-container>
    <Alert />

    <v-card>
      <v-card-title>
        Trip Offers by other Cargonauts
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
          :items="tripOffers"
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
                <v-icon small class="mr-2" v-on="on" @click="bookTrip(item)"
                  >mdi-plus</v-icon
                >
              </template>
              <span>Book Trip</span>
            </v-tooltip>
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
  </v-container>
</template>

<script>
import { mapActions } from "vuex";
import { mapGetters } from "vuex";

import Alert from "@/components/Alert";

export default {
  name: "Trips",

  components: {
    Alert
  },

  computed: {
    ...mapGetters("auth", ["authId"]),
    ...mapGetters("trips", {
      trips: "trips",
      tripsLoading: "loading"
    }),
    ...mapGetters("vehicles", {
      vehicle: "vehicle",
      vehicles: "vehicles",
      vehiclesLoading: "loading"
    }),

    loading() {
      return this.tripsLoading && this.vehiclesLoading;
    },

    tripOffers() {
      return this.trips.filter(trip => {
        return trip.user_id != this.authId && !trip.rider_id;
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

      { text: "Actions", value: "action", sortable: false },
      { text: "", value: "data-table-expand" }
    ],
    search: ""
  }),

  methods: {
    ...mapActions("trips", { listTrips: "list" }),
    ...mapActions("users", { bookTripForRider: "bookTrip" }),
    ...mapActions("vehicles", { listVehicles: "list", getVehicle: "get" }),

    getTripVehicle({ item, value }) {
      if (!value) return;
      this.getVehicle(item.vehicle_id);
    },

    bookTrip(item) {
      this.bookTripForRider({
        userId: this.authId,
        tripId: item.id
      }).then(() => this.listTrips());
    }
  },

  created() {
    this.listTrips();
    this.listVehicles();
  }
};
</script>
