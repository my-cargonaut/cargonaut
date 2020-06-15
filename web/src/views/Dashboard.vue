<template>
  <div>
    <v-container>
      <Alert />

      <v-card>
        <v-card-title>
          Dashboard
          <v-spacer></v-spacer>
          <v-text-field
            v-model="search"
            append-icon="mdi-magnify"
            label="Search"
            single-line
            hide-details
          ></v-text-field>
        </v-card-title>

        <v-card-text>
          <v-data-table
            :loading="loading"
            :headers="headers"
            :items="stations"
            :search="search"
          >
            <template v-slot:item.action="{ item }">
              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                  <v-icon
                    small
                    class="mr-2"
                    v-on="on"
                    :disabled="activeStation(item) && loading"
                    @click="
                      activeStation(item) ? stopStream() : playStream(item)
                    "
                    >{{
                      activeStation(item) ? "mdi-pause" : "mdi-play"
                    }}</v-icon
                  >
                </template>
                <span>{{
                  activeStation(item) ? "Pause Station" : "Play Station"
                }}</span>
              </v-tooltip>

              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                  <v-icon
                    small
                    class="mr-2"
                    v-on="on"
                    @click="editStation(item)"
                    >mdi-pencil</v-icon
                  >
                </template>
                <span>Edit Station</span>
              </v-tooltip>

              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                  <v-icon small v-on="on" @click="deleteStation(item)"
                    >mdi-delete</v-icon
                  >
                </template>
                <span>Delete Station</span>
              </v-tooltip>
            </template>
          </v-data-table>
        </v-card-text>
      </v-card>
    </v-container>
  </div>
</template>

<script>
import { Howl } from "howler";

import { mapActions } from "vuex";

import Alert from "@/components/Alert";

export default {
  name: "Dashboard",

  components: {
    Alert
  },

  data: () => ({
    headers: [
      {
        text: "Name",
        value: "name",
        sortable: true,
        align: "start"
      },
      { text: "Actions", value: "action", sortable: false, align: "right" }
    ],
    loading: false,
    search: "",
    player: null,
    playerId: "",
    stations: [
      {
        id: "12",
        name: "HR3",
        url:
          "https://hr-edge-10ab-fra-dtag-cdn.cast.addradio.de/hr/hr3/live/mp3/128/stream.mp3"
      },
      {
        id: "34",
        name: "FFH",
        url: "http://mp3.ffh.de/radioffh/hqlivestream.mp3"
      },
      {
        id: "56",
        name: "Planet",
        url: "http://mp3.ffh.de/radio/hqlivestream.mp3"
      },
      {
        id: "78",
        name: "YouFM",
        url: "http://mp3.ffh.de/radioffh/livestream.mp3"
      }
    ]
  }),

  methods: {
    ...mapActions("alert", {
      setAlert: "set"
    }),
    activeStation(station) {
      return this.player ? this.playerId == station.id : false;
    },
    playStream(station) {
      this.stopStream();
      this.player = new Howl({
        src: station.url,
        html5: true,
        onplay: () => {
          if (this.activeStation(station)) this.loading = false;
        },
        onplayerror: () => {
          this.setAlert({
            message: "Failed to play " + station.name,
            type: "error"
          });
          this.stopStream();
        },
        onloaderror: () => {
          this.setAlert({
            message: "Failed to load " + station.name,
            type: "error"
          });
          this.stopStream();
        }
      });
      this.player.play();
      this.playerId = station.id;
      this.loading = true;
    },
    stopStream() {
      if (this.player) this.player.stop();
      this.playerId = "";
      this.loading = false;
    },
    editStation(item) {
      this.stopStream();
      console.log("edit: " + item);
    },
    deleteStation(item) {
      this.stopStream();
      console.log("delete: " + item);
    }
  },

  beforeDestroy() {
    this.stopStream();
  }
};
</script>
