<template>
  <div>
    <v-navigation-drawer
      v-show="isLoggedIn"
      :clipped="$vuetify.breakpoint.lgAndUp"
      app
      permanent
      expand-on-hover
    >
      <v-list>
        <v-list-item class="px-2" :to="'/users/' + this.authId" color="primary">
          <v-list-item-avatar>
            <v-img :src="'/api/v1/users/' + authId + '/avatar'"></v-img>
          </v-list-item-avatar>

          <v-list-item-content class="py-1">
            <v-list-item-title class="title">{{ authName }}</v-list-item-title>
            <v-list-item-subtitle>{{ authEmail }}</v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
      </v-list>

      <v-divider></v-divider>

      <v-list nav dense>
        <v-list-item-group color="primary">
          <v-list-item to="/">
            <v-list-item-icon>
              <v-icon>mdi-view-dashboard</v-icon>
            </v-list-item-icon>
            <v-list-item-title>Dashboard</v-list-item-title>
          </v-list-item>
          <v-list-item to="/trips">
            <v-list-item-icon>
              <v-icon>mdi-road-variant</v-icon>
            </v-list-item-icon>
            <v-list-item-title>Trip Offers</v-list-item-title>
          </v-list-item>
          <v-list-item to="/vehicles">
            <v-list-item-icon>
              <v-icon>mdi-car</v-icon>
            </v-list-item-icon>
            <v-list-item-title>Vehicles</v-list-item-title>
          </v-list-item>
        </v-list-item-group>
      </v-list>

      <v-divider></v-divider>

      <v-list nav dense>
        <v-list-item @click="logout" color="primary">
          <v-list-item-icon>
            <v-icon>mdi-logout</v-icon>
          </v-list-item-icon>
          <v-list-item-title>Logout</v-list-item-title>
        </v-list-item>
      </v-list>

      <template v-slot:append>
        <v-list nav dense>
          <v-list-item-group color="primary">
            <v-list-item to="/disclaimer">
              <v-list-item-icon>
                <v-icon>mdi-information</v-icon>
              </v-list-item-icon>
              <v-list-item-title>Disclaimer</v-list-item-title>
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </template>
    </v-navigation-drawer>
  </div>
</template>

<script>
import { mapActions } from "vuex";
import { mapGetters } from "vuex";

export default {
  name: "Sidebar",

  computed: {
    ...mapGetters("auth", ["isLoggedIn", "authId", "authName", "authEmail"])
  },

  methods: {
    ...mapActions("auth", {
      logoutUser: "logout"
    }),

    logout() {
      this.logoutUser().finally(() => this.$router.push("/login"));
    }
  }
};
</script>
