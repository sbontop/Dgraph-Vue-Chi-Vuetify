import Vue from "vue";
import VueRouter from "vue-router";
import Home from "../views/Home.vue";
import Buyers from "../views/Buyers.vue";
import BuyerDetail from "../views/BuyerDetail.vue";
import Recommendations from "../views/Recommendations.vue";
import Ip from "../views/Ip.vue";
import PurchaseHistory from "../views/PurchaseHistory.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/about",
    name: "About",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/About.vue"),
  },
  {
    path: "/buyers",
    name: "Buyers",
    component: Buyers,
  },
  {
    path: "/buyers/:id",
    name: "BuyerDetail",
    component: BuyerDetail,
  },
  {
    path: "/purchaseHistory/:id",
    name: "PurchaseHistory",
    component: PurchaseHistory,
  },
  {
    path: "/ip/:id",
    name: "Ip",
    component: Ip,
  },
  {
    path: "/recommendations/:id",
    name: "Recommendations",
    component: Recommendations,
  },    
];

const router = new VueRouter({
  routes,
});

export default router;
