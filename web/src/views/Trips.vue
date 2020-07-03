<template>
  <v-container>
    <Alert />

    <v-card>
      <v-card-title>
        Trips
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
          :items="trips"
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
            <v-icon small class="mr-2" @click="editTrip(item)"
              >mdi-pencil</v-icon
            >
            <v-icon small @click="deleteTrip(item)">mdi-delete</v-icon>
          </template>
          <template v-slot:expanded-item="{ headers }">
            <td :colspan="headers.length">
              <v-chip class="ma-2" color="indigo" text-color="white">
                <v-avatar left>
                  <v-icon>mdi-car</v-icon>
                </v-avatar>
                {{ tripVehicle.brand }} {{ tripVehicle.model }}
              </v-chip>
              <v-chip
                class="ma-2"
                color="green"
                text-color="white"
                :active="tripVehicle.passengers > 0"
              >
                <v-avatar left>
                  <v-icon>mdi-account-circle</v-icon>
                </v-avatar>
                Passengers
                <v-avatar right class="green darken-4">
                  {{ tripVehicle.passengers }}
                </v-avatar>
              </v-chip>
              <v-chip
                class="ma-2"
                color="orange"
                text-color="white"
                :active="tripVehicle.loading_area_length > 0"
              >
                {{ tripVehicle.loading_area_length }}cm x
                {{ tripVehicle.loading_area_width }}cm
                <v-icon right>mdi-crop</v-icon>
              </v-chip>
            </td>
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
                        :loading="vehiclesLoading"
                        :items="userVehicles"
                        item-value="id"
                        item-text="name"
                        label="Vehicle"
                        prepend-icon="mdi-car"
                        v-model="editedTrip.vehicle_id"
                      ></v-select>
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
      tripVehicle: "vehicle",
      vehicles: "vehicles",
      vehiclesLoading: "loading"
    }),

    formTitle() {
      return this.editedIndex === -1 ? "New Trip" : "Edit Trip";
    },

    combinedVehicleNamesList() {
      return this.vehicles.map(vehicle => {
        return {
          id: vehicle.id,
          user_id: vehicle.user_id,
          name: vehicle.brand + " " + vehicle.model,
          passengers: vehicle.passengers,
          loading_area_length: vehicle.loading_area_length,
          loading_area_width: vehicle.loading_area_width,
          created_at: vehicle.created_at,
          updated_at: vehicle.updated_at
        };
      });
    },

    userVehicles() {
      return this.combinedVehicleNamesList.filter(vehicle => {
        return vehicle.user_id == this.authId;
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
    ...mapActions("vehicles", { listVehicles: "list", getVehicle: "get" }),

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

    getTripVehicle({ item, value }) {
      if (!value) return;
      this.getVehicle(item.vehicle_id);
    },

    close() {
      this.dialog = false;
      setTimeout(() => {
        this.editedTrip = Object.assign({});
        this.editedIndex = -1;
        this.$refs.form.reset();
      }, 300);
    }
  },

  created() {
    this.list();
    this.listVehicles();
  }
};
</script>
