<template>
  <div>
    <v-card class="mx-auto">
      <v-list two-line>
        <v-list-item-group active-class="green--text" multiple>
          <template v-for="(r, index) in recommendations">
            <v-list-item :key="`${r.product_id}`">
              <template>
                <v-list-item-content>
                  <v-list-item-title v-text="r.product_name"></v-list-item-title>
                  <v-list-item-subtitle
                    class="text--primary"
                    v-text="r.product_price"
                  ></v-list-item-subtitle>
                </v-list-item-content>
              </template>
            </v-list-item>
            <v-divider
              v-if="index < recommendations.length - 1"
              :key="index"
            ></v-divider>
          </template>
        </v-list-item-group>
      </v-list>
    </v-card>
  </div>
</template>

<script>
import axios from "axios";
export default {
  name: "Recommendations",
  data: () => ({
    recommendations: null,
  }),
  mounted() {
    axios
      .get(
        `http://localhost:3333/realbuyers/recommendations/${this.$route.params.id}`
      )
      .then((result) => {
        console.log(result.data.productRecom[0].product);
        this.recommendations = result.data.productRecom[0].product;
      });
  },
};
</script>
