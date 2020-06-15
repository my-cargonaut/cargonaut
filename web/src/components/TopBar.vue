<template>
  <div>
    <v-navigation-drawer
      v-model="drawer"
      :clipped="$vuetify.breakpoint.lgAndUp"
      app
    >
      <v-list dense>
        <v-list-item to="/" link>
          <v-list-item-action>
            <v-icon>mdi-view-dashboard</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Dashboard</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar
      :clipped-left="$vuetify.breakpoint.lgAndUp"
      app
      color="blue darken-3"
      dark
    >
      <v-app-bar-nav-icon v-show="isLoggedIn" @click.stop="drawer = !drawer" />
      <v-toolbar-title style="width: 300px" class="ml-0 pl-4">
        <span class="hidden-sm-and-down">My Cargonaut</span>
      </v-toolbar-title>

      <v-spacer></v-spacer>

      <v-menu offset-y transition="slide-y-transition">
        <template v-slot:activator="{ on }">
          <v-btn text v-show="isLoggedIn" v-on="on">
            <v-avatar size="35">
              <v-icon>mdi-account</v-icon>
            </v-avatar>
            <span class="ml-2">{{ authUsername }}</span>
            <v-icon class="ml-2">mdi-chevron-down</v-icon>
          </v-btn>
        </template>
        <v-list>
          <v-list-item @click="logout">
            <v-list-item-content>
              <v-list-item-title>Logout</v-list-item-title>
            </v-list-item-content>
            <v-list-item-action>
              <v-icon>mdi-logout</v-icon>
            </v-list-item-action>
          </v-list-item>
        </v-list>
      </v-menu>
    </v-app-bar>
  </div>
</template>

<script>
import { mapGetters } from "vuex";

export default {
  name: "TopBar",

  data: () => ({
    drawer: false
  }),

  computed: {
    ...mapGetters("auth", ["isLoggedIn", "authUsername"])
  },

  methods: {
    logout() {
      this.$store
        .dispatch("auth/logout")
        // .then(() => this.$router.push("/login"))
        // .catch(error => console.log(error.response.data.error))
        .finally(() => this.$router.push("/login"));
    }
  }
};
</script>
