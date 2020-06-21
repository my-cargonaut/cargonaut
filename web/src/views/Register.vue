<template>
  <v-container fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <Alert />

        <v-card class="elevation-12">
          <v-toolbar color="primary" dark flat>
            <v-toolbar-title>Register</v-toolbar-title>
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
                v-model.trim="email"
                label="E-Mail"
                placeholder="Your E-Mail address will be your username"
                prepend-icon="mdi-email"
                :rules="emailRules"
                required
              />
              <v-text-field
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                label="Password"
                prepend-icon="mdi-lock"
                :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                @click:append="showPassword = !showPassword"
                :rules="passwordRules"
                required
              />
              <v-text-field
                v-model.trim="display_name"
                label="Name"
                prepend-icon="mdi-account"
                :rules="[v => !!v || 'Name is required']"
                required
              />
              <v-menu
                ref="birthdayMenu"
                v-model="birthdayMenu"
                :close-on-content-click="false"
                transition="scale-transition"
                offset-y
                min-width="290px"
              >
                <template v-slot:activator="{ on, attrs }">
                  <v-text-field
                    v-model.trim="birthday"
                    label="Birthday"
                    placeholder="Pick a date"
                    prepend-icon="mdi-calendar"
                    :rules="[v => !!v || 'Birthday is required']"
                    v-bind="attrs"
                    v-on="on"
                    readonly
                    required
                  ></v-text-field>
                </template>
                <v-date-picker
                  ref="picker"
                  v-model="birthday"
                  :max="new Date().toISOString().substr(0, 10)"
                  min="1950-01-01"
                  @change="save"
                ></v-date-picker>
              </v-menu>
              <v-file-input
                v-model="avatar"
                label="Avatar"
                placeholder="Pick an image"
                prepend-icon="mdi-camera"
                accept="image/png, image/jpeg"
                :rules="avatarRules"
                show-size
                required
              ></v-file-input>
            </v-form>
          </v-card-text>

          <v-divider></v-divider>

          <v-card-actions>
            <v-btn color="warning" to="/login">Back to login</v-btn>
            <v-spacer></v-spacer>
            <v-btn color="info" @click="register" :disabled="!valid || loading"
              >Register</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import Alert from "@/components/Alert";

import authAPI from "@/api/auth";

export default {
  name: "Register",

  components: {
    Alert
  },

  data: () => ({
    loading: false,
    valid: false,
    email: "",
    password: "",
    display_name: "",
    birthday: null,
    avatar: null,
    showPassword: false,
    birthdayMenu: false,
    emailRules: [
      v => !!v || "E-Mail is required",
      v => /.+@.+/.test(v) || "E-mail must be valid"
    ],
    passwordRules: [
      v => !!v || "Password is required",
      v => (v && v.length >= 8) || "Password must have 8 or more characters"
    ],
    avatarRules: [
      v => !!v || "Avatar is required",
      v => !v || v.size < 2e6 || "Avatar size should be less than 2 MB!"
    ]
  }),

  watch: {
    birthdayMenu(val) {
      val && setTimeout(() => (this.$refs.picker.activePicker = "YEAR"));
    }
  },

  methods: {
    save(birthday) {
      this.$refs.birthdayMenu.save(birthday);
    },

    register() {
      if (!this.valid) return;

      const birthday = new Date(this.birthday).toISOString();

      this.loading = true;
      getBase64(this.avatar).then(avatar => {
        authAPI
          .register(
            this.email,
            this.password,
            this.display_name,
            birthday,
            btoa(avatar)
          )
          .then(() => {
            this.$store
              .dispatch("alert/set", {
                message: "Registered successfully!",
                type: "success",
                title: true
              })
              .then(() => this.$router.push("/login"));
          })
          .catch(error => {
            this.$store.dispatch("alert/set", {
              message: error.response.data.error,
              type: "error",
              title: true
            });
          })
          .finally(() => (this.loading = false));
      });
    }
  }
};

function getBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsBinaryString(file);
    reader.onload = () => resolve(reader.result);
    reader.onerror = error => reject(error);
  });
}
</script>
