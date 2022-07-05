<template>
  <nav class="navbar navbar-top navbar-expand navbar-dark bg-primary border-bottom">
    <div class="container-fluid">
      <div class="collapse navbar-collapse" id="navbarSupportedContent">
        <!-- Navbar links -->
        <ul class="navbar-nav align-items-center ml-md-auto">
          <li class="nav-item d-xl-none">
            <!-- Sidenav toggler -->
            <div class="pr-3 sidenav-toggler sidenav-toggler-dark" data-action="sidenav-pin" data-target="#sidenav-main">
              <div class="sidenav-toggler-inner">
                <i class="sidenav-toggler-line"></i>
                <i class="sidenav-toggler-line"></i>
                <i class="sidenav-toggler-line"></i>
              </div>
            </div>
          </li>
        </ul>
        <ul class="navbar-nav align-items-center ml-auto ml-md-0" v-if="!isLoggedIn">
          <li class="nav-item dropdown">
            <router-link class="nav-link pr-0" :to="{ name: 'auth.login' }" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
              <div class="media align-items-center">
                <div class="media-body ml-2">
                  <span class="mb-0 text-sm text-white font-weight-bold">Login</span>
                </div>
              </div>
            </router-link>
          </li>
        </ul>
        <ul class="navbar-nav align-items-center ml-auto ml-md-0" v-if="isLoggedIn">
          <li class="nav-item dropdown">
            <router-link class="nav-link pr-0" :to="{ name: 'admin' }" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
              <div class="media align-items-center">
                <div class="media-body ml-2">
                  <span class="mb-0 text-sm text-white font-weight-bold">{{ name }}</span>
                </div>
              </div>
            </router-link>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script>

import authHeader from "@/service/auth-header";
import axios from "axios";

export default {
  mounted() {
    this.checkLogin()
    if(this.isLoggedIn) {
      this.fetchData()
    }
  },
  data() {
    return {
      name: "",
      isLoggedIn: false
    }
  },
  methods: {
    checkLogin() {
      if(Object.keys(authHeader()).length > 0) {
        this.isLoggedIn = true
      }
    },
    fetchData() {
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/profile`, { headers: authHeader() })
          .then(response => {
            this.name = response.data.data.name
          }).catch(error => {
            console.log(error.response.data.status)
          })
    },
  }
}
</script>