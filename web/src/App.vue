<template>
  <v-app>
    <TopBar />

    <v-content>
      <v-container fluid>
        <AlertBox />
        <v-fade-transition mode="out-in">
          <router-view></router-view>
        </v-fade-transition>
      </v-container>
    </v-content>
  </v-app>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import AlertBox from "./components/AlertBox.vue";
import TopBar from "./components/TopBar.vue";

@Component({
  components: {
    AlertBox,
    TopBar
  }
})
export default class extends Vue {
  mounted() {
    this.$store.dispatch("fuelTank/listFuelTanks");
    this.$store.dispatch("truck/listTrucks");
  }
}
</script>
