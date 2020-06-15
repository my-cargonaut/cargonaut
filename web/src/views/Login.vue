<template>
  <v-container fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <Alert />

        <Credentials @login="login" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Alert from "@/components/Alert";
import Credentials from "@/components/login/Credentials";

export default {
  name: "Login",

  components: {
    Alert,
    Credentials
  },

  methods: {
    login(username, password) {
      this.loading = true;
      this.$store
        .dispatch("auth/login", { username, password })
        .then(() => this.$router.push("/"))
        // .catch(error => console.log(error.response.data.error))
        .finally(() => (this.loading = false));
    }
  }
};
</script>
