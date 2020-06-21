<template>
  <v-card :loading="loading" class="mx-auto my-12" max-width="374">
    <v-row justify="space-around">
      <v-avatar size="200" class="mt-4">
        <img
          :src="'/api/v1/users/' + id + '/avatar'"
          :alt="user.display_name"
        />
      </v-avatar>
    </v-row>

    <v-card-text>
      <div class="subtitle-1">
        {{ user.display_name }}
      </div>

      <div class="my-2 subtitle-2">
        {{ user.email }}
      </div>

      <v-row align="center" class="my-4 mx-0">
        <v-rating
          :value="4.5"
          color="amber"
          dense
          half-increments
          :readonly="id == authID"
          size="28"
        ></v-rating>

        <div class="grey--text ml-4">4.5 (413)</div>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script>
import { mapGetters } from "vuex";

export default {
  name: "User",

  props: ["id"],

  computed: {
    ...mapGetters("auth", ["authID"]),
    ...mapGetters("users", ["user"])
  },

  data: () => ({
    loading: false
  }),

  created() {
    this.loading = true;
    this.$store
      .dispatch("users/get", this.id)
      .finally(() => (this.loading = false));
  }
};
</script>
