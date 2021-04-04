<template>
  <div id="buyer-detail">
    <div id="nav">
      <router-link :to="{ name: 'PurchaseHistory', params: { id: this.$route.params.id }}">Historial de Compras</router-link> |      
      <router-link :to="{ name: 'Ip', params: { id: this.$route.params.id }}">Ip</router-link> |
      <router-link :to="{ name: 'Recommendations', params: { id: this.$route.params.id }}">Recomendaciones</router-link>
    </div>
  </div>
</template>

<style>
#buyer-detail {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}
</style>

<script>
import axios from "axios";
export default {
  name: "BuyerDetail",
  data: () => ({
    buyers: null,
  }),
  mounted() {
    axios.get("http://localhost:3333/realbuyers").then((result) => {
      this.buyers = result.data;
    });
  },
};
</script>