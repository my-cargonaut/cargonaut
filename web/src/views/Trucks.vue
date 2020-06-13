<template>
  <div>
    <v-container>
      <v-card>
        <v-card-title>
          Trucks
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
            :items="trucks"
            :search="search"
          >
            <template v-slot:item.action="{ item }">
              <v-icon small class="mr-2" @click="editTruck(item)"
                >mdi-pencil</v-icon
              >
              <v-icon small @click="deleteTruck(item)">mdi-delete</v-icon>
            </template>
          </v-data-table>

          <v-dialog v-model="dialog" max-width="750px">
            <template v-slot:activator="{ on }">
              <v-btn color="primary" class="mb-2" v-on="on">New Truck</v-btn>
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
                          v-model.trim="editedTruck.manufacturer"
                          label="Manufacturer"
                          placeholder="Mercedes-Benz"
                          :rules="rules"
                          required
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12" sm="6" md="6">
                        <v-text-field
                          v-model.trim="editedTruck.model"
                          label="Model"
                          placeholder="Actros (BM 963)"
                          :rules="rules"
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
                  @click="saveTruck"
                  :disabled="!valid"
                  >Save</v-btn
                >
              </v-card-actions>
            </v-card>
          </v-dialog>
        </v-card-text>
      </v-card>
    </v-container>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from "vue-property-decorator";
import { Getter } from "vuex-class";
import { Alert, Truck } from "../types";
import TruckAPI from "../services/TruckAPI";

@Component
export default class Trucks extends Vue {
  $refs!: {
    form: HTMLFormElement;
  };

  @Getter("trucks", { namespace: "truck" }) trucks!: Truck[];
  @Getter("loading", { namespace: "fuelTank" }) loading!: boolean;

  private headers = [
    {
      text: "Manufacturer",
      value: "manufacturer",
      sortable: true,
      align: "start"
    },
    { text: "Model", value: "model", sortable: true },
    { text: "Actions", value: "action", sortable: false }
  ];
  private search = "";
  private dialog = false;
  private valid = false;
  private rules = [
    (v: any) => !!v || "Field is required",
    (v: any) => (v && v.length) >= 3 || "Field must be at least 3 characters"
  ];
  private editedTruck: Truck = {} as Truck;
  private editedIndex = -1;

  saveTruck() {
    if (this.editedIndex > -1) {
      TruckAPI.update(this.editedTruck)
        .then(() => {
          this.$store.dispatch("truck/listTrucks");
        })
        .catch(e => this.error(e));
    } else {
      TruckAPI.create(this.editedTruck)
        .then(() => {
          this.$store.dispatch("truck/listTrucks");
        })
        .catch(e => this.error(e));
    }
    this.close();
  }

  editTruck(truck: Truck) {
    this.editedIndex = this.trucks.indexOf(truck);
    this.editedTruck = Object.assign({}, truck);
    this.dialog = true;
  }

  deleteTruck(truck: Truck) {
    confirm("Are you sure you want to delete this truck?") &&
      TruckAPI.delete(truck.id)
        .then(() => {
          this.$store.dispatch("truck/listTrucks");
          this.$store.dispatch("fuelTank/listFuelTanks");
        })
        .catch(e => this.error(e));
  }

  close() {
    this.dialog = false;
    setTimeout(() => {
      this.editedTruck = Object.assign({}, {} as Truck);
      this.editedIndex = -1;
      this.$refs.form.reset();
    }, 300);
  }

  error(e: any) {
    const alert = {
      type: "error",
      message: "Something went wrong!"
    } as Alert;
    if (e.response && e.response.data && e.response.data.error) {
      alert.message = e.response.data.error;
    }
    this.$store.dispatch("alert/setAlert", alert);
  }

  get formTitle() {
    return this.editedIndex === -1 ? "New Truck" : "Edit Truck";
  }

  @Watch("dialog")
  onDialogChanged(value: boolean) {
    value || this.close();
  }
}
</script>
