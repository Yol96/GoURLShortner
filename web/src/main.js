import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import axios from "axios";
import VueAxios from "vue-axios";

Vue.config.productionTip = false;

Vue.prototype.baseURL = "http://localhost:5000/api/";
Vue.prototype.baseStaticURL = "http://localhost:5000/static/"

axios.defaults.baseURL = Vue.prototype.baseURL;
axios.defaults.baseStaticURL = Vue.prototype.baseStaticURL

Vue.use(VueAxios, axios);

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
