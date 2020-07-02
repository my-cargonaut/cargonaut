<template>
  <v-container>
    <Alert />

    <v-card>
      <v-card-title>
        Vehicles
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
          :items="vehicles"
          :search="search"
        >
          <template v-slot:item.action="{ item }">
            <v-icon small class="mr-2" @click="editVehicle(item)"
              >mdi-pencil</v-icon
            >
            <v-icon small @click="deleteVehicle(item)">mdi-delete</v-icon>
          </template>
        </v-data-table>

        <v-dialog v-model="dialog" max-width="750px">
          <template v-slot:activator="{ on }">
            <v-btn color="primary" class="mb-2" v-on="on">New Vehicle</v-btn>
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
                        v-model.trim="editedVehicle.brand"
                        label="Manufacturer"
                        :rules="rules"
                        required
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" sm="6" md="6">
                      <v-text-field
                        v-model.trim="editedVehicle.model"
                        label="Model"
                        :rules="rules"
                        required
                      ></v-text-field>
                    </v-col>
                  </v-row>
                  <v-row>
                    <v-col cols="12" sm="4" md="4">
                      <v-text-field
                        type="number"
                        v-model.number="editedVehicle.passengers"
                        label="Max. Passengers"
                        required
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" sm="4" md="4">
                      <v-text-field
                        type="number"
                        v-model.number="editedVehicle.loading_area_length"
                        suffix="cm"
                        label="Loading Area Length"
                        required
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" sm="4" md="4">
                      <v-text-field
                        type="number"
                        v-model.number="editedVehicle.loading_area_width"
                        suffix="cm"
                        label="Loading Area Width"
                        required
                      ></v-text-field>
                    </v-col>
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
                @click="saveVehicle"
                :disabled="!valid"
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
  name: "Vehicles",

  components: {
    Alert
  },

  computed: {
    ...mapGetters("vehicles", ["vehicles", "loading"]),

    formTitle() {
      return this.editedIndex === -1 ? "New Vehicle" : "Edit Vehicle";
    }
  },

  data: () => ({
    headers: [
      {
        text: "Brand",
        value: "brand",
        sortable: true,
        align: "start"
      },
      { text: "Model", value: "model", sortable: true },
      { text: "Passengers", value: "passengers", sortable: true },
      {
        text: "Loading Area Length (cm)",
        value: "loading_area_length",
        sortable: true
      },
      {
        text: "Loading Area Width (cm)",
        value: "loading_area_width",
        sortable: true
      },
      { text: "Actions", value: "action", sortable: false }
    ],
    search: "",
    dialog: false,
    valid: false,
    rules: [
      v => !!v || "Field is required",
      v => (v && v.length) >= 3 || "Field must be at least 3 characters"
    ],
    editedVehicle: {},
    editedIndex: -1
  }),

  methods: {
    ...mapActions("vehicles", ["list", "create", "update", "delete"]),

    saveVehicle() {
      if (this.editedIndex > -1) {
        this.update({
          id: this.editedVehicle.id,
          vehicle: this.editedVehicle
        }).then(() => this.list());
      } else {
        this.create(this.editedVehicle).then(() => this.list());
      }
      this.close();
    },

    editVehicle(vehicle) {
      this.editedIndex = this.vehicles.indexOf(vehicle);
      this.editedVehicle = Object.assign({}, vehicle);
      this.dialog = true;
    },

    deleteVehicle(vehicle) {
      confirm("Are you sure you want to delete this vehicle?") &&
        this.delete(vehicle.id).then(() => this.list());
    },

    close() {
      this.dialog = false;
      setTimeout(() => {
        this.editedVehicle = Object.assign({});
        this.editedIndex = -1;
        this.$refs.form.reset();
      }, 300);
    }
  },

  created() {
    this.list();
  }
};
</script>
