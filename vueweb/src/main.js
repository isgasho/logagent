// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import logs from '../static/js/logs'
import VueError from '../static/js/errorcatch.js'

Vue.config.productionTip = false
Vue.use(VueError)

Vue.prototype.logs = logs

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})
