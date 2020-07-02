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
        >
          <template v-slot:item.action="{ item }">
            <v-icon small class="mr-2" @click="editTrip(item)"
              >mdi-pencil</v-icon
            >
            <v-icon small @click="deleteTrip(item)">mdi-delete</v-icon>
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
                        :rules="rules"
                        required
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" sm="6" md="6">
                      <v-select
                        :loading="vehiclesLoading"
                        :items="combinedVehicleNamesList"
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
    ...mapGetters("trips", {
      trips: "trips",
      tripsLoading: "loading"
    }),
    ...mapGetters("vehicles", {
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
          name: vehicle.brand + " " + vehicle.model
        };
      });
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

      { text: "Actions", value: "action", sortable: false }
    ],
    search: "",
    dialog: false,
    valid: false,
    rules: [
      v => !!v || "Field is required",
      v => (v && v.length) >= 3 || "Field must be at least 3 characters"
    ],
    editedTrip: {},
    editedIndex: -1
  }),

  methods: {
    ...mapActions("trips", ["list", "create", "update", "delete"]),
    ...mapActions("vehicles", { listVehicles: "list" }),

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
