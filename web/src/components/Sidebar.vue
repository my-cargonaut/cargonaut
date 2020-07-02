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
        <v-list-item class="px-2" :to="'/users/' + this.authId">
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
        <v-list-item to="/trips">
          <v-list-item-icon>
            <v-icon>mdi-road-variant</v-icon>
          </v-list-item-icon>
          <v-list-item-title>Trips</v-list-item-title>
        </v-list-item>
        <v-list-item to="/vehicles">
          <v-list-item-icon>
            <v-icon>mdi-car</v-icon>
          </v-list-item-icon>
          <v-list-item-title>Vehicles</v-list-item-title>
        </v-list-item>
      </v-list>

      <v-divider></v-divider>

      <v-list nav dense>
        <v-list-item @click="logout">
          <v-list-item-icon>
            <v-icon>mdi-logout</v-icon>
          </v-list-item-icon>
          <v-list-item-title>Logout</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
  </div>
</template>

<script>
import { mapGetters } from "vuex";

export default {
  name: "Sidebar",

  computed: {
    ...mapGetters("auth", ["isLoggedIn", "authId", "authName", "authEmail"])
  },

  methods: {
    logout() {
      this.$store
        .dispatch("auth/logout")
        .finally(() => this.$router.push("/login"));
    }
  }
};
</script>
