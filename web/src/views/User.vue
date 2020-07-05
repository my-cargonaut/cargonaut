<template>
  <v-container>
    <Alert />

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
            v-model="ratings.average"
            color="amber"
            dense
            hover
            half-increments
            readonly
            size="28"
          ></v-rating>

          <div class="grey--text ml-4">
            {{ ratings.average.toFixed(1) }} ({{ ratings.count }})
          </div>
        </v-row>
      </v-card-text>
    </v-card>

    <v-card
      class="mx-auto my-12"
      max-width="374"
      v-if="ratings.ratings && ratings.ratings.length > 0"
    >
      <v-card-title>
        Ratings by other Cargonauts
      </v-card-title>

      <v-list three-line v-for="rating in ratings.ratings" :key="rating.id">
        <v-divider></v-divider>
        <v-list-item>
          <v-list-item-avatar class="mt-7">
            <v-img
              :src="'/api/v1/users/' + rating.author_id + '/avatar'"
              alt="Rating author"
            ></v-img>
          </v-list-item-avatar>

          <v-list-item-content>
            <v-list-item-title>
              <v-rating
                :value="rating.value"
                color="amber"
                dense
                hover
                half-increments
                readonly
                size="28"
              ></v-rating>
            </v-list-item-title>
            <v-list-item-subtitle>{{ rating.comment }}</v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-card>
  </v-container>
</template>

<script>
import { mapActions } from "vuex";
import { mapGetters } from "vuex";

import Alert from "@/components/Alert";

export default {
  name: "User",

  components: {
    Alert
  },

  props: ["id"],

  computed: {
    ...mapGetters("auth", ["authId"]),
    ...mapGetters("users", ["user", "ratings", "loading"])
  },

  methods: {
    ...mapActions("users", ["get", "listRatings"])
  },

  created() {
    this.get(this.id);
    this.listRatings(this.id);
  }
};
</script>
