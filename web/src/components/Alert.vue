<template>
  <div>
    <v-alert
      border="right"
      colored-border
      elevation="2"
      class="text-center"
      transition="scale-transition"
      dismissible
      :value="!!alert['message']"
      :type="alert['type']"
      @input="closeAlert"
      >{{
        alert["title"] ? titleCase(alert["message"]) : alert["message"]
      }}</v-alert
    >
  </div>
</template>

<script>
import { mapState, mapActions } from "vuex";

export default {
  name: "Alert",

  computed: {
    ...mapState("alert", ["alert"])
  },

  methods: {
    ...mapActions("alert", {
      setAlert: "set"
    }),
    closeAlert() {
      this.setAlert({});
    },
    titleCase(sentence) {
      if (!sentence) {
        return sentence;
      }
      sentence = sentence.toLowerCase().split(" ");
      for (let i = 0; i < sentence.length; i++) {
        sentence[i] = sentence[i][0].toUpperCase() + sentence[i].slice(1);
      }
      return sentence.join(" ");
    }
  }
};
</script>
