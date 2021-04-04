<template>
    <v-card class="mx-auto">
      <v-list two-line>
        <v-list-item-group active-class="green--text" multiple>
          <template v-for="(buyer, index) in buyers">
            <v-list-item :key="`${buyer.buyer_id}`">
              <template>
                <v-list-item-content @click="verDetalleComprador(buyer.buyer_id)">
                  <v-list-item-title
                    v-text="buyer.buyer_name"
                  ></v-list-item-title>
                  <v-list-item-subtitle
                    class="text--primary"
                    v-text="buyer.buyer_age"
                  ></v-list-item-subtitle>

                  <v-list-item-subtitle
                    v-text="buyer.Ip"
                  ></v-list-item-subtitle>
                </v-list-item-content>
              </template>
            </v-list-item>
            <v-divider
              v-if="index < buyers.length - 1"
              :key="index"
            ></v-divider>
          </template>
        </v-list-item-group>
      </v-list>
    </v-card>
</template>

<script>
import axios from "axios";
export default {
  name: "Buyers",
  data: () => ({
    buyers: null,
  }),
  methods: {
    verDetalleComprador: function (buyer_id) {
      console.log(buyer_id);
      this.$router.push(`/buyers/${buyer_id}`)
    },
  },
  mounted() {
    axios.get("http://localhost:3333/realbuyers").then((result) => {
      this.buyers = result.data;
    });
  },
};
</script>
