<template>
    <v-container fluid>
        <v-row align="center" justify="center">
            <v-col cols="12" sm="12" md="10" lg="8">
                <Alert/>

                <v-card class="elevation-12">
                    <v-toolbar color="primary" dark flat>
                        <v-toolbar-title>Fahrt erstellen</v-toolbar-title>
                    </v-toolbar>

                    <v-fade-transition mode="out-in">
                        <v-container v-if="loading">
                            <div class="text-xs-center">
                                <v-progress-linear indeterminate></v-progress-linear>
                            </div>
                        </v-container>
                    </v-fade-transition>

                    <v-card-text>
                        <v-form ref="form" v-model="valid" @keyup.native.enter="create">
                            <v-container>
                                <v-row>
                                    <v-text-field
                                            v-model.trim="email"
                                            label="Von:"
                                            placeholder="Z.B. Berlin"
                                            prepend-icon="mdi-home-circle"
                                    />
                                    <v-text-field
                                            v-model="password"
                                            label="Nach:"
                                            placeholder="Z.B. Hamburg"
                                            prepend-icon="mdi-home-circle-outline"
                                    />
                                </v-row>
                                <v-row>
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
                                                    label="Datum"
                                                    placeholder="Pick a date"
                                                    prepend-icon="mdi-calendar"
                                                    :rules="[v => !!v || 'Date is required']"
                                                    v-bind="attrs"
                                                    v-on="on"
                                                    readonly
                                                    required
                                            ></v-text-field>
                                        </template>
                                        <v-date-picker
                                                ref="picker"
                                                v-model="birthday"
                                                max="2030-01-01"
                                                :min="new Date().toISOString().substr(0, 10)"
                                                @change="save"
                                        ></v-date-picker>
                                    </v-menu>
                                    <v-text-field
                                            v-model.trim="display_name"
                                            label="Uhrzeit"
                                            prepend-icon="mdi-clock-outline"
                                    />
                                </v-row>
                                <v-row>
                                    <v-text-field
                                            v-model.trim="email"
                                            label="Sitzplätze:"
                                            prepend-icon="mdi-seat"
                                    />
                                    <v-text-field
                                            v-model="password"
                                            label="Ladefläche (m³/l):"
                                            prepend-icon="mdi-bag-checked"
                                    />
                                </v-row>
                                <v-row>
                                    <v-text-field
                                            v-model.trim="email"
                                            label="Preis pro Sitzplatz:"
                                            prepend-icon="mdi-currency-usd"
                                    />
                                    <v-text-field
                                            v-model="password"
                                            label="Preis pro m³/l:"
                                            prepend-icon="mdi-currency-eur"
                                    />
                                </v-row>
                            </v-container>
                        </v-form>
                    </v-card-text>

                    <v-divider></v-divider>

                    <v-card-actions>
                        <v-spacer></v-spacer>
                        <v-btn color="info" @click="register" :disabled="!valid || loading"
                        >Erstellen
                        </v-btn
                        >
                    </v-card-actions>
                </v-card>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
    export default {
        name: "CreateTour"
    }
</script>

<style scoped>

</style>