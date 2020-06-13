<template>
  <div>
    <v-snackbar v-model="copiedSnack" color="success" top right timeout="3000">
      Copied "{{ level }}" to clipboard!
      <v-btn dark text @click="copiedSnack = false">Close</v-btn>
    </v-snackbar>

    <v-container>
      <v-card>
        <v-card-title>Calculate</v-card-title>

        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12" sm="4" md="4">
                <v-select
                  :loading="trucksLoading"
                  :items="filteredTrucks"
                  item-value="id"
                  item-text="name"
                  label="Truck"
                  prepend-icon="mdi-truck"
                  v-model="truckId"
                ></v-select>
              </v-col>
              <v-col cols="12" sm="4" md="4">
                <v-select
                  :disabled="truckId == ''"
                  :loading="fuelTanksLoading"
                  :items="fuelTanksForTruck"
                  item-value="id"
                  item-text="name"
                  label="Fuel Tank"
                  prepend-icon="mdi-gas-station"
                  v-model="fuelTankId"
                ></v-select>
              </v-col>
              <v-col cols="12" sm="4" md="4">
                <v-text-field
                  :disabled="fuelTankId == ''"
                  type="number"
                  label="Fill Height"
                  prepend-icon="mdi-car-coolant-level"
                  v-model="height"
                  suffix="cm"
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="12" md="12">
                <v-text-field
                  ref="levelLiter"
                  label="Level"
                  readonly="true"
                  append-icon="mdi-content-copy"
                  v-model="level"
                  suffix="liter"
                  @click:append="copyToClipboard"
                ></v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
      </v-card>
    </v-container>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { Getter } from "vuex-class";
import { FuelTank, Truck } from "../types";

@Component
export default class Calculate extends Vue {
  $refs!: {
    levelLiter: HTMLFormElement;
  };

  @Getter("fuelTankByID", { namespace: "fuelTank" }) fuelTank!: (
    id: string
  ) => FuelTank;
  @Getter("fuelTanksForTruck", { namespace: "fuelTank" }) fuelTanks!: (
    truckId: string
  ) => FuelTank[];
  @Getter("loading", { namespace: "fuelTank" }) fuelTanksLoading!: boolean;
  @Getter("trucks", { namespace: "truck" }) trucks!: Truck[];
  @Getter("loading", { namespace: "truck" }) trucksLoading!: boolean;

  private copiedSnack = false;
  private truckId = "";
  private fuelTankId = "";
  private height = 0;

  copyToClipboard() {
    const textToCopy = this.$refs.levelLiter.$el.querySelector("input");
    textToCopy.select();
    textToCopy.setSelectionRange(0, 99999);
    document.execCommand("copy");
    this.copiedSnack = true;
  }

  get fuelTanksForTruck() {
    return this.fuelTanks(this.truckId);
  }

  get level() {
    const fuelTank = this.fuelTank(this.fuelTankId);
    if (!fuelTank) {
      return 0;
    } else if (this.height <= 0) {
      return 0;
    }
    return (fuelTank.length * fuelTank.width * this.height) / 1000.0;
  }

  get filteredTrucks() {
    return this.trucks.map((truck: Truck) => {
      return {
        id: truck.id,
        name: truck.manufacturer + " " + truck.model
      };
    });
  }
}
</script>
