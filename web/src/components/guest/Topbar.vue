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
                <a href="javascript:void(0);" class="text-muted text-sm" @click="reloadPage">Refresh halaman</a>
              </div>
              <!-- List group -->
              <div class="list-group list-group-flush" v-if="cart.length > 0">
                <template v-for="(item,i) in cart" :key="i">
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
    this.fetchCart()
  },
  data() {
    return {
      name: "",
      isLoggedIn: false,
      cart: [],
      totalCart: 0,
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
    getImage(image) {
      return image;
    },
    fetchCart(){
      localStorage.getItem('my-carts') ? this.cart = JSON.parse(localStorage.getItem('my-carts')).slice(0,5) : this.cart = []

      this.totalCart = JSON.parse(localStorage.getItem('my-carts')).reduce((total, item) => {
        return total + item.quantity
      }, 0)
    },
    reloadPage() {
      window.location.reload();
    }
  }
}
</script>