<template>
  <v-container>
    <Alert />

    <v-row>
      <v-col cols="12" md="6">
        <UserTrips />
      </v-col>

      <v-col cols="12" md="6">
        <UserRides />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mapActions } from "vuex";

import Alert from "@/components/Alert";
import UserTrips from "@/components/UserTrips";
import UserRides from "@/components/UserRides";

export default {
  name: "Trips",

  components: {
    Alert,
    UserTrips,
    UserRides
  },

  methods: {
    ...mapActions("trips", { listTrips: "list" })
  },

  created() {
    this.listTrips();
  }
};
</script>
