<template>
  <v-card class="mx-auto">
    <v-list two-line>
      <v-list-item-group active-class="green--text" multiple>
        <template v-for="(p, index) in purchaseHistory">
          <v-list-item :key="`${p.product_id}`">
            <template>
              <v-list-item-content>
                <v-list-item-title v-text="p.product_name"></v-list-item-title>
                <v-list-item-subtitle
                  class="text--primary"
                  v-text="p.product_price"
                ></v-list-item-subtitle>
              </v-list-item-content>
            </template>
          </v-list-item>
          <v-divider v-if="index < purchaseHistory.length - 1" :key="index"></v-divider>
        </template>
      </v-list-item-group>
    </v-list>
  </v-card>
</template>

<script>
import axios from "axios";
export default {
  name: "Purchase History",
  data: () => ({
    purchaseHistory: null,
  }),
  mounted() {
    axios
      .get(
        `http://localhost:3333/realbuyers/purchaseHistory/${this.$route.params.id}`
      )
      .then((result) => {
        console.log(result.data.purchasesHistory[0].product);
        this.purchaseHistory = result.data.purchasesHistory[0].product;
      });
  },
};
</script>
