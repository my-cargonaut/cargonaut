<template>
  <div>
    <v-card class="elevation-12">
      <v-toolbar color="primary" dark flat>
        <v-toolbar-title>Login</v-toolbar-title>
      </v-toolbar>

      <v-fade-transition mode="out-in">
        <v-container v-if="loading">
          <div class="text-xs-center">
            <v-progress-linear indeterminate></v-progress-linear>
          </div>
        </v-container>
      </v-fade-transition>

      <v-card-text>
        <v-form ref="form" v-model="valid" @keyup.native.enter="login">
          <v-text-field
            v-model.trim="username"
            label="Username"
            prepend-icon="mdi-account"
            :rules="[v => !!v || 'Username is required']"
            required
          />
          <v-text-field
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            label="Password"
            prepend-icon="mdi-lock"
            :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
            @click:append="showPassword = !showPassword"
            :rules="[v => !!v || 'Password is required']"
            required
          />
        </v-form>
      </v-card-text>

      <v-divider></v-divider>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="info" @click="login" :disabled="!valid || loading"
          >Login</v-btn
        >
      </v-card-actions>
    </v-card>
  </div>
</template>

<script>
export default {
  name: "Credentials",

  data: () => ({
    loading: false,
    valid: false,
    username: "",
    password: "",
    showPassword: false
  }),

  methods: {
    login() {
      if (!this.valid) return;
      this.$emit("login", this.username, this.password);
    }
  }
};
</script>
