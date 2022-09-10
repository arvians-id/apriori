<template>
  <nav class="sidenav navbar navbar-vertical fixed-left navbar-expand-xs navbar-light bg-white" id="sidenav-main">
    <div class="scrollbar-inner">
      <!-- Brand -->
      <div class="sidenav-header d-flex align-items-center">
        <router-link class="navbar-brand" :to="{ name: 'guest.index' }" style="font-family: 'Fugaz One', cursive; color: #6772e5">
          RYZYOLSHOP
        </router-link>
        <div class="ml-auto">
          <!-- Sidenav toggler -->
          <div class="sidenav-toggler d-none d-xl-block" data-action="sidenav-unpin" data-target="#sidenav-main">
            <div class="sidenav-toggler-inner">
              <i class="sidenav-toggler-line"></i>
              <i class="sidenav-toggler-line"></i>
              <i class="sidenav-toggler-line"></i>
            </div>
          </div>
        </div>
      </div>
      <div class="navbar-inner">
        <!-- Collapse -->
        <div class="collapse navbar-collapse" id="sidenav-collapse-main">
          <!-- Nav items -->
          <ul class="navbar-nav">
            <li class="nav-item">
              <router-link :class="getActiveNavLink('guest.index')" :to="{ name: 'guest.index' }">
                <i class="ni ni-shop text-primary"></i>
                <span class="nav-link-text">Home</span>
              </router-link>
            </li>
            <li class="nav-item">
              <router-link :class="getActiveNavLink('guest.product')" :to="{ name: 'guest.product' }">
                <i class="ni ni-box-2 text-info"></i>
                <span class="nav-link-text">Product List</span>
              </router-link>
            </li>
            <li class="nav-item">
              <router-link :class="getActiveNavLink('guest.cart')" :to="{ name: 'guest.cart' }">
                <i class="ni ni-cart text-danger"></i>
                <span class="nav-link-text">My Order</span>
              </router-link>
            </li>
            <li class="nav-item" v-if="isLoggedIn">
              <router-link :class="getActiveNavLink('member.order')" :to="{ name: 'member.order' }">
                <i class="ni ni-delivery-fast text-dark"></i>
                <span class="nav-link-text">History Order</span>
              </router-link>
            </li>
            <li class="nav-item" v-if="isLoggedIn">
              <router-link :class="getActiveNavLink('member.history')" :to="{ name: 'member.history' }">
                <i class="ni ni-satisfied text-success"></i>
                <span class="nav-link-text">Rate the Order</span>
              </router-link>
            </li>
          </ul>
          <template v-if="isLoggedIn">
            <!-- Divider -->
            <hr class="my-3">
            <!-- Heading -->
            <h6 class="navbar-heading p-0 text-muted">Others</h6>
            <!-- Navigation -->
            <ul class="navbar-nav mb-md-3">
              <li class="nav-item">
                <router-link :class="getActiveNavLink('member.notification')" :to="{ name: 'member.notification' }">
                  <i class="ni ni-bell-55"></i>
                  <span class="nav-link-text">Notification</span>
                  <span class="badge badge badge-success ml-2" v-if="totalNotification > 0">{{ totalNotification }}</span>
                </router-link>
              </li>
              <li class="nav-item">
                <router-link :class="getActiveNavLink('member.profile')" :to="route">
                  <i class="ni ni-single-02"></i>
                  <span class="nav-link-text">My Account</span>
                </router-link>
              </li>
              <li class="nav-item">
                <form @submit.prevent="submit" method="POST" class="nav-link">
                  <i class="ni ni-curved-next"></i>
                  <button class="btn btn-link" style="padding: 0; font-weight: normal; color: rgba(0, 0, 0, .6)" type="submit">Logout</button>
                </form>
              </li>
            </ul>
          </template>
        </div>
      </div>
    </div>
  </nav>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Fugaz+One&display=swap');

.navbar-brand {
  font-size: 1.40rem;
}
</style>

<script>
import authHeader from "@/service/auth-header";
import axios from "axios";
import getRoles from "@/service/get-roles";

export default {
  props: {
    totalNotification: {
      type: Number,
      default: 0
    },
  },
  mounted() {
    this.checkLogin()
    this.checkRole()
  },
  data() {
    return {
      isLoggedIn : false,
      route: { name: 'member.profile' },
    }
  },
  methods: {
    checkLogin() {
      if(authHeader()["Authorization"]) {
        this.isLoggedIn = true
      }
    },
    async checkRole() {
      let getRole = await getRoles();
      if(getRole != null) {
        if(getRole.role == "1") {
          this.route = {name: 'profile'}
        }
      }
    },
    getActiveNavLink(name) {
      let classString = "nav-link "

      let routeName = this.$route.name.split('.')[1]
      if (routeName === name.split('.')[1]) {
        classString += "active"
      }
      return classString;
    },
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
  }
}
</script>