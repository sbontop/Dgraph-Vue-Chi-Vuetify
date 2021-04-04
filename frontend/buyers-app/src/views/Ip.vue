<template>
  <div>
    <v-card class="mx-auto">
      <v-list two-line>
        <v-list-item-group active-class="green--text" multiple>
          <template v-for="(b, index) in buyers">
            <v-list-item :key="`${b.buyer_id}`">
              <template>
                <v-list-item-content>
                  <v-list-item-title
                    v-text="b.buyer_name"
                  ></v-list-item-title>
                  <v-list-item-subtitle
                    class="text--primary"
                    v-text="b.buyer_age"
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
  </div>
</template>

<script>
import axios from "axios";
export default {
  name: "Ip",
  data: () => ({
    buyers: null,
  }),
  mounted() {
    axios
      .get(`http://localhost:3333/realbuyers/ip/${this.$route.params.id}`)
      .then((result) => {
        console.log(result.data.getBuyerByIpAddress);
        this.buyers = result.data.getBuyerByIpAddress;
      });
  },
};
</script>
