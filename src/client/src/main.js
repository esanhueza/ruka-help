import Vue from 'vue';
import App from './App.vue';
import Home from './components/Home.vue';
import Workspace from './components/Workspace.vue';
import axios from 'axios';
import VueRouter from 'vue-router';
import VueAxios from 'vue-axios';
import Unicon from 'vue-unicons'
import { uniCommentsAlt, uniUsersAlt, uniMessage, uniArrow, uniDatabase, uniPlay, uniExclamationTriangle, uniPadlock, uniHeartAlt, uniStar, uniTable, uniListUl, uniMicrophone, uniMicrophoneSlash, uniWebcam, uniSlidersVAlt, uniStopwatch, uniStopCircle } from 'vue-unicons/src/icons'

Unicon.add([uniCommentsAlt, uniUsersAlt, uniMessage, uniArrow, uniDatabase, uniPlay, uniExclamationTriangle, uniPadlock, uniHeartAlt, uniStar, uniTable, uniListUl, uniMicrophone, uniMicrophoneSlash, uniWebcam, uniSlidersVAlt, uniStopwatch, uniStopCircle])


Vue.config.productionTip = false;
Vue.use(VueAxios, axios);
Vue.use(VueRouter)
Vue.use(Unicon)

const router = new VueRouter({
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/workspace/:id',
      name: 'workspace',
      component: Workspace
    }
  ]
});


new Vue({
  router: router,
  el: '#app',
  render: h => h(App),
}).$mount('#app')
