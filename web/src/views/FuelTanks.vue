<template>
  <div>
    <v-container>
      <v-card>
        <v-card-title>
          Fuel Tanks
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
            :loading="fuelTanksLoading"
            :headers="headers"
            :items="fuelTanks"
            :search="search"
          >
            <template v-slot:item.action="{ item }">
              <v-icon small class="mr-2" @click="editFuelTank(item)"
                >mdi-pencil</v-icon
              >
              <v-icon small @click="deleteFuelTank(item)">mdi-delete</v-icon>
            </template>
          </v-data-table>

          <v-dialog v-model="dialog" max-width="750px">
            <template v-slot:activator="{ on }">
              <v-btn color="primary" class="mb-2" v-on="on"
                >New Fuel Tank</v-btn
              >
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
                        <v-select
                          :loading="trucksLoading"
                          :items="filteredTrucks"
                          item-value="id"
                          item-text="name"
                          label="Truck"
                          prepend-icon="mdi-truck"
                          v-model="editedFuelTank.truckId"
                          :rules="truckRules"
                          required
                        ></v-select>
                      </v-col>
                      <v-col cols="12" sm="6" md="6">
                        <v-text-field
                          label="Name"
                          placeholder="250L"
                          v-model.trim="editedFuelTank.name"
                          :rules="nameRules"
                          required
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12" sm="6" md="6">
                        <v-text-field
                          type="number"
                          label="Width"
                          placeholder="50"
                          suffix="cm"
                          v-model.number="editedFuelTank.width"
                          :rules="numberRules"
                          required
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12" sm="6" md="6">
                        <v-text-field
                          type="number"
                          label="Length"
                          placeholder="100"
                          suffix="cm"
                          v-model.number="editedFuelTank.length"
                          :rules="numberRules"
                          required
                        ></v-text-field>
                      </v-col>
                    </v-row>
                  </v-container>
                </v-form>
              </v-card-text>

              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="blue darken-1" text @click="close">Cancel</v-btn>
                <v-btn
                  color="blue darken-1"
                  text
                  :disabled="!valid"
                  @click="saveFuelTank"
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
import { Alert, FuelTank, Truck } from "../types";
import FuelTankAPI from "../services/FuelTankAPI";

@Component
export default class FuelTanks extends Vue {
  $refs!: {
    form: HTMLFormElement;
  };

  @Getter("fuelTanks", { namespace: "fuelTank" }) fuelTanks!: FuelTank[];
  @Getter("loading", { namespace: "fuelTank" }) fuelTanksLoading!: boolean;
  @Getter("trucks", { namespace: "truck" }) trucks!: Truck[];
  @Getter("loading", { namespace: "truck" }) trucksLoading!: boolean;

  private headers = [
    { text: "Name", value: "name", sortable: true, align: "start" },
    { text: "Length (cm)", value: "length", sortable: true },
    { text: "Width (cm)", value: "width", sortable: true },
    { text: "Actions", value: "action", sortable: false }
  ];
  private search = "";
  private dialog = false;
  private valid = false;
  private truckRules = [(v: any) => !!v || "Field is required"];
  private nameRules = [
    (v: any) => !!v || "Field is required",
    (v: any) => (v && v.length) >= 3 || "Field must be at least 3 characters"
  ];
  private numberRules = [
    (v: any) => !!v || "Field is required",
    (v: any) => v > 0 || "Field value must be greater than 0"
  ];
  private editedFuelTank: FuelTank = {} as FuelTank;
  private editedIndex = -1;

  saveFuelTank() {
    if (this.editedIndex > -1) {
      FuelTankAPI.update(this.editedFuelTank)
        .then(() => {
          this.$store.dispatch("fuelTank/listFuelTanks");
        })
        .catch(e => this.error(e));
    } else {
      FuelTankAPI.create(this.editedFuelTank)
        .then(() => {
          this.$store.dispatch("fuelTank/listFuelTanks");
        })
        .catch(e => this.error(e));
    }
    this.close();
  }

  editFuelTank(fuelTank: FuelTank) {
    this.editedIndex = this.fuelTanks.indexOf(fuelTank);
    this.editedFuelTank = Object.assign({}, fuelTank);
    this.dialog = true;
  }

  deleteFuelTank(fuelTank: FuelTank) {
    confirm("Are you sure you want to delete this fuel tank?") &&
      FuelTankAPI.delete(fuelTank.id)
        .then(() => {
          this.$store.dispatch("fuelTank/listFuelTanks");
        })
        .catch(e => this.error(e));
  }

  close() {
    this.dialog = false;
    setTimeout(() => {
      this.editedFuelTank = Object.assign({}, {} as FuelTank);
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
    return this.editedIndex === -1 ? "New Fuel Tank" : "Edit Fuel Tank";
  }

  get filteredTrucks() {
    return this.trucks.map((truck: Truck) => {
      return {
        id: truck.id,
        name: truck.manufacturer + " " + truck.model
      };
    });
  }

  @Watch("dialog")
  onDialogChanged(value: boolean) {
    value || this.close();
  }
}
</script>
