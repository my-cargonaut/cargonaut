<template>
  <v-card>
    <v-card-title>
      My Trips
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
        :loading="tripsLoading"
        :headers="headers"
        :items="userTrips"
        :search="search"
      >
        <template v-slot:item.vehicle="{ item }">
          {{ tripVehicle(item.vehicle_id).brand }}
          {{ tripVehicle(item.vehicle_id).model }}
        </template>
        <template v-slot:item.depature="{ item }">
          {{
            +new Date(item.depature) > 0
              ? new Date(item.depature).toLocaleString()
              : "-"
          }}
        </template>
        <template v-slot:item.arrival="{ item }">
          {{
            +new Date(item.arrival) > 0
              ? new Date(item.arrival).toLocaleString()
              : "-"
          }}
        </template>
        <template v-slot:item.action="{ item }">
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-icon
                small
                class="mr-2"
                v-on="on"
                :disabled="!isWaitingForStart(item)"
                @click="startTrip(item)"
                >mdi-map-marker-right</v-icon
              >
            </template>
            <span>Start Trip</span>
          </v-tooltip>
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-icon
                small
                class="mr-2"
                v-on="on"
                :disabled="!isWaitingForStop(item)"
                @click="endTrip(item)"
                >mdi-map-marker-left</v-icon
              >
            </template>
            <span>End Trip</span>
          </v-tooltip>
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-icon
                small
                class="mr-2"
                v-on="on"
                :disabled="!isWaitingForRider(item)"
                @click="editTrip(item)"
                >mdi-pencil</v-icon
              >
            </template>
            <span>Edit Trip</span>
          </v-tooltip>
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-icon
                small
                class="mr-2"
                v-on="on"
                :disabled="!isWaitingForRider(item)"
                @click="deleteTrip(item)"
                >mdi-delete</v-icon
              >
            </template>
            <span>Delete Trip</span>
          </v-tooltip>
        </template>
        <template v-slot:item.status="{ item }">
          {{ getTripStatus(item) }}
        </template>
        <template v-slot:item.rating="{ item }">
          {{ item.rating ? item.rating : "-" }}
          {{ item.rating ? "/5" : "" }}
        </template>
      </v-data-table>

      <v-dialog v-model="dialog" max-width="750px">
        <template v-slot:activator="{ on }">
          <v-btn color="primary" class="mb-2" v-on="on">New Trip</v-btn>
        </template>

        <v-card>
          <v-card-title>
            <span class="headline">{{ formTitle }}</span>
          </v-card-title>

          <v-card-text>
            <v-form ref="form" v-model="valid">
              <v-container>
                <v-row>
                  <v-col cols="12" sm="6" md="6">
                    <v-text-field
                      v-model.trim="editedTrip.start"
                      label="From"
                      prepend-icon="mdi-map-marker-right"
                      :rules="rules"
                      required
                    ></v-text-field>
                  </v-col>
                  <v-col cols="12" sm="6" md="6">
                    <v-text-field
                      v-model.trim="editedTrip.destination"
                      label="To"
                      prepend-icon="mdi-map-marker-left"
                      :rules="rules"
                      required
                    ></v-text-field>
                  </v-col>
                </v-row>
                <v-row>
                  <v-col cols="12" sm="6" md="6">
                    <v-text-field
                      type="number"
                      v-model.number="editedTrip.price"
                      label="Price"
                      prepend-icon="mdi-currency-eur"
                      :rules="numRules"
                      required
                    ></v-text-field>
                  </v-col>
                  <v-col cols="12" sm="6" md="6">
                    <v-select
                      :loading="usersLoading"
                      :items="vehicles"
                      item-value="id"
                      item-text="name"
                      v-model="editedTrip.vehicle_id"
                      label="Vehicle"
                      prepend-icon="mdi-car"
                      :rules="rules"
                      required
                    >
                      <template v-slot:item="{ item }">
                        {{ item.brand }} {{ item.model }}
                      </template>
                      <template v-slot:selection="{ item }">
                        {{ item.brand }} {{ item.model }}
                      </template>
                    </v-select>
                  </v-col>
                </v-row>
              </v-container>
            </v-form>
          </v-card-text>

          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="primary" text @click="close">Cancel</v-btn>
            <v-btn color="primary" text @click="saveTrip" :disabled="!valid"
              >Save</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-card-text>
  </v-card>
</template>

<script>
import { mapActions } from "vuex";
import { mapGetters } from "vuex";

import TripStatus from "@/shared/trip_status";

export default {
  name: "UserTrips",

  computed: {
    ...mapGetters("auth", ["authId"]),
    ...mapGetters("trips", {
      trips: "trips",
      tripsLoading: "loading"
    }),
    ...mapGetters("users", {
      vehicles: "vehicles",
      usersLoading: "loading"
    }),

    userTrips() {
      return this.trips.filter(trip => {
        return trip.user_id == this.authId;
      });
    },

    formTitle() {
      return this.editedIndex === -1 ? "New Trip" : "Edit Trip";
    }
  },

  data: () => ({
    headers: [
      {
        text: "From",
        value: "start",
        sortable: true,
        align: "start"
      },
      { text: "To", value: "destination", sortable: true },
      { text: "Price (€)", value: "price", sortable: true },
      { text: "Vehicle", value: "vehicle", sortable: true },
      { text: "Depature", value: "depature", sortable: true },
      { text: "Arrival", value: "arrival", sortable: true },
      { text: "Status", value: "status", sortable: false },
      { text: "Actions", value: "action", sortable: false },
      { text: "Rating", value: "rating", sortable: false },
      { text: "", value: "data-table-expand" }
    ],
    search: "",
    dialog: false,
    valid: false,
    rules: [
      v => !!v || "Field is required",
      v => (v && v.length) >= 3 || "Field must be at least 3 characters"
    ],
    numRules: [
      v => !!v || "Field is required",
      v => v >= 0 || "Field value must be greater or equal than 0"
    ],
    editedTrip: {},
    editedIndex: -1
  }),

  methods: {
    ...mapActions("trips", ["list", "create", "update", "delete"]),
    ...mapActions("users", { listUserVehicles: "listVehicles" }),

    tripVehicle(id) {
      const vehicle = this.vehicles.find(vehicle => {
        return vehicle.id == id;
      });
      return vehicle ? vehicle : {};
    },

    saveTrip() {
      if (this.editedIndex > -1) {
        this.update({
          id: this.editedTrip.id,
          trip: this.editedTrip
        }).then(() => this.list());
      } else {
        this.create(this.editedTrip).then(() => this.list());
      }
      this.close();
    },

    editTrip(trip) {
      this.editedIndex = this.trips.indexOf(trip);
      this.editedTrip = Object.assign({}, trip);
      this.dialog = true;
    },

    deleteTrip(trip) {
      confirm("Are you sure you want to delete this trip?") &&
        this.delete(trip.id).then(() => this.list());
    },

    startTrip(trip) {
      const tripToStart = Object.assign({}, trip);
      tripToStart.depature = new Date().toISOString();
      this.update({
        id: tripToStart.id,
        trip: tripToStart
      }).then(() => this.list());
    },

    endTrip(trip) {
      const tripToEnd = Object.assign({}, trip);
      tripToEnd.arrival = new Date().toISOString();
      this.update({
        id: tripToEnd.id,
        trip: tripToEnd
      }).then(() => this.list());
    },

    close() {
      this.dialog = false;
      setTimeout(() => {
        this.editedTrip = Object.assign({});
        this.editedIndex = -1;
        this.$refs.form.reset();
      }, 300);
    },

    getTripStatus(trip) {
      return TripStatus.get(trip);
    },

    isWaitingForRider(trip) {
      return TripStatus.isWaitingForRider(trip);
    },

    isWaitingForStart(trip) {
      return TripStatus.isWaitingForStart(trip);
    },

    isWaitingForStop(trip) {
      return TripStatus.isWaitingForStop(trip);
    }
  },

  created() {
    this.listUserVehicles(this.authId);
  }
};
</script>
