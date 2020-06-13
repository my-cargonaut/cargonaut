<template>
  <div>
    <v-alert
      class="text-center"
      :value="!!alert.message"
      :type="alert.type"
      outlined
      dismissible
      @input="closeAlert"
      >{{ alert.message }}</v-alert
    >
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import { Getter } from "vuex-class";
import { Alert } from "../types";

@Component
export default class AlertBox extends Vue {
  @Getter("alert", { namespace: "alert" }) alert!: Alert;

  closeAlert() {
    this.$store.dispatch("alert/setAlert", {} as Alert);
  }
}
</script>
