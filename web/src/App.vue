<template>
  <v-app>
    <TopBar />

    <main>
      <v-main>
        <v-container fluid>
          <v-fade-transition mode="out-in">
            <router-view />
          </v-fade-transition>
          <v-footer class="pa-2" fixed app>
            <v-spacer></v-spacer>
            <div>
              Made with <v-icon color="red">mdi-heart</v-icon> by Lukas Malkmus,
              Philipp Alexander Händler & Robert Feuerhack &copy;
              {{ new Date().getFullYear() }}
            </div>
            <v-spacer></v-spacer>
          </v-footer>
        </v-container>
      </v-main>
    </main>
  </v-app>
</template>

<script>
import TopBar from "@/components/TopBar";

export default {
  name: "App",

  components: {
    TopBar
  },

  created() {
    this.$http.interceptors.response.use(undefined, error => {
      return new Promise(() => {
        if (
          error.status === 401 &&
          error.config &&
          !error.config.__isRetryRequest
        ) {
          this.$store.dispatch("auth/logout");
        }
        throw error;
      });
    });
  }
};
</script>
