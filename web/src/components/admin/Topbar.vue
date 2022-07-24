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
        <ul class="navbar-nav align-items-center ml-auto ml-md-0">
          <li class="nav-item dropdown">
            <a class="nav-link pr-0" href="javascript:void(0);" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
              <div class="media align-items-center">
                <div class="media-body ml-2">
                  <span class="mb-0 text-sm text-white font-weight-bold">{{ name }}</span>
                </div>
              </div>
            </a>
            <div class="dropdown-menu dropdown-menu-right">
              <div class="dropdown-header noti-title">
                <h6 class="text-overflow m-0">Welcome!</h6>
              </div>
              <router-link :to="{ name:'profile' }" class="dropdown-item">
                <i class="ni ni-single-02"></i>
                <span>My profile</span>
              </router-link>
              <form @submit.prevent="submit" method="POST" class="dropdown-item">
                <i class="ni ni-user-run"></i>
                <button class="btn btn-link text-dark" style="padding: 0; font-weight: normal;" type="submit">Logout</button>
              </form>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script>
import axios from "axios";
import authHeader from "@/service/auth-header";
import getRoles from "@/service/get-roles";

export default {
  data() {
    return {
      name: "",
    }
  },
  mounted() {
    this.fetchData()
  },
  methods: {
    submit() {
      axios.delete(`${process.env.VUE_APP_SERVICE_URL}/auth/logout`,{ headers: authHeader() })
          .then(response => {
            if(response.data.code === 200) {
              localStorage.removeItem("token")
              localStorage.removeItem("refresh-token")
              alert(response.data.status)
              this.$router.push({
                name: 'auth.login'
              })
            }
          }).catch(error => {
        console.log(error.response.data.status)
      })
    },
    async fetchData() {
      let getRole = await getRoles();
      this.name = getRole.name
    },
  }
}
</script>