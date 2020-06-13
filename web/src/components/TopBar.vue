<template>
  <div>
    <v-navigation-drawer app fixed v-model="showMenu">
      <v-list dense>
        <v-list-item to="/" link>
          <v-list-item-action>
            <v-icon>mdi-calculator</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Calculate</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item to="/tanks" link>
          <v-list-item-action>
            <v-icon>mdi-gas-station</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Fuel Tanks</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item to="/trucks" link>
          <v-list-item-action>
            <v-icon>mdi-truck</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Trucks</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app fixed color="primary">
      <v-app-bar-nav-icon @click.stop="toggleMenu"></v-app-bar-nav-icon>

      <div class="d-flex align-center">
        <v-toolbar-title>
          <strong>CARGONAUT</strong>
        </v-toolbar-title>
      </div>

      <v-spacer></v-spacer>

      <v-btn href="https://www.si-elektronik.eu/" target="_blank" text>
        <span class="mr-2">SI-Elektronik GmbH</span>
        <v-icon>mdi-open-in-new</v-icon>
      </v-btn>
    </v-app-bar>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class TopToolbar extends Vue {
  private showMenu = false;

  toggleMenu(): void {
    this.showMenu = !this.showMenu;
  }
}
</script>
