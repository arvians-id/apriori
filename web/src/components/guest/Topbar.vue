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
          <li class="nav-item dropdown">
            <a class="nav-link" href="#" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
              {{ totalCart }} <i class="ni ni-cart"></i>
            </a>
            <div class="dropdown-menu dropdown-menu-xl dropdown-menu-right py-0 overflow-hidden">
              <!-- Dropdown header -->
              <div class="px-3 py-3">
                <h6 class="text-sm text-muted m-0">Kamu memiliki <strong class="text-primary">{{ totalCart }}</strong> barang dikeranjang.</h6>
              </div>
              <!-- List group -->
              <div class="list-group list-group-flush overflow-auto" style="height: 400px" v-if="carts.length > 0">
                <template v-for="(item,i) in carts" :key="i">
                  <router-link
                        :to="{ name: 'guest.product.recommendation', params: { code: item.id_product, id: item.code } }"
                        v-if="item.name.includes(`Paket Rekomendasi`)"
                        class="list-group-item list-group-item-action">
                    <div class="row align-items-center">
                      <div class="col-auto">
                        <!-- Avatar -->
                        <img alt="Image placeholder" :src="getImage(item.image)" class="avatar rounded-circle">
                      </div>
                      <div class="col ml--2">
                        <div class="d-flex justify-content-between align-items-center">
                          <div>
                            <h4 class="mb-0 text-sm">{{ item.name.length > 30 ? item.name.slice(0, 30) + "..." : item.name  }}</h4>
                          </div>
                          <div class="text-right text-muted">
                            <small>{{ item.quantity }} item</small>
                          </div>
                        </div>
                        <p class="text-sm mb-0">Rp {{  item.totalPricePerItem }}</p>
                      </div>
                    </div>
                  </router-link>
                  <router-link :to="{ name: 'guest.product.detail', params: { code: item.code } }" class="list-group-item list-group-item-action" v-else>
                    <div class="row align-items-center">
                      <div class="col-auto">
                        <!-- Avatar -->
                        <img alt="Image placeholder" :src="getImage(item.image)" class="avatar rounded-circle">
                      </div>
                      <div class="col ml--2">
                        <div class="d-flex justify-content-between align-items-center">
                          <div>
                            <h4 class="mb-0 text-sm">{{ item.name.length > 30 ? item.name.slice(0, 30) + "..." : item.name  }}</h4>
                          </div>
                          <div class="text-right text-muted">
                            <small>{{ item.quantity }} item</small>
                          </div>
                        </div>
                        <p class="text-sm mb-0">Rp {{  item.totalPricePerItem }}</p>
                      </div>
                    </div>
                  </router-link>
                </template>
              </div>
              <div class="list-group list-group-flush" v-else>
                <div class="list-group-item">
                  <div class="text-center">
                      <h4 class="mb-0 text-sm text-muted">Keranjang kamu masih kosong.</h4>
                  </div>
                </div>
              </div>
              <!-- View all -->
              <router-link :to="{ name: 'guest.cart' }" class="dropdown-item text-center text-primary font-weight-bold py-3">Lihat semua</router-link>
            </div>
          </li>
        </ul>
        <ul class="navbar-nav align-items-center ml-auto ml-md-0" v-if="isLoggedIn">
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
              <router-link :to="route" class="dropdown-item">
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
        <ul class="navbar-nav align-items-center ml-auto ml-md-0" v-else>
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
      </div>
    </div>
  </nav>
</template>

<script>

import authHeader from "@/service/auth-header";
import axios from "axios";

export default {
  props: {
    totalCart: {
      type: Number,
      default: 0
    },
    carts: {
      type: Array,
      default: () => []
    },
  },
  mounted() {
    this.checkLogin()
    if(this.isLoggedIn) {
      this.fetchData()
    }
    this.checkRole()
  },
  data() {
    return {
      name: "",
      isLoggedIn: false,
      route: { name: 'member.profile' },
    }
  },
  methods: {
    submit() {
      axios.delete(`${process.env.VUE_APP_SERVICE_URL}/auth/logout`,{ headers: authHeader() })
          .then(response => {
            if(response.data.code === 200) {
              localStorage.removeItem("token")
              localStorage.removeItem("refresh-token")
              localStorage.removeItem("user")
              localStorage.removeItem("name")
              alert(response.data.status)
              this.$router.push({
                name: 'auth.login'
              })
            }
          }).catch(error => {
        console.log(error.response.data.status)
      })
    },
    checkLogin() {
      if(authHeader()["Authorization"]) {
        this.isLoggedIn = true
      }
    },
    checkRole(){
      if(localStorage.getItem("role") === "1") {
        this.route = { name: 'profile' }
      }
    },
    fetchData() {
      this.name = localStorage.getItem("name")
    },
    getImage(image) {
      return image;
    }
  }
}
</script>