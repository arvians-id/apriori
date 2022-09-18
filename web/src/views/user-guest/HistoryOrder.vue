<template>
  <!-- Sidenav -->
  <Sidebar :totalNotification="totalNotification" />
  <!-- Main content -->
  <div class="main-content" id="panel">
    <!-- Topnav -->
    <Topbar :totalCart="totalCart" :carts="carts" :totalNotification="totalNotification" :notifications="notifications" />
    <!-- Header -->
    <Header />
    <!-- Page content -->
    <div class="container-fluid mt--6">
      <div class="row">
        <div class="col-12">
          <div class="card-wrapper">
            <!-- Custom form validation -->
            <div class="card">
              <!-- Card header -->
              <div class="card-header">
                <!-- Title -->
                <h5 class="h3 mb-0">Riwayat Penilaian Anda</h5>
              </div>
              <!-- Card body -->
              <div class="card-body p-0" v-if="isLoading">
                <div class="loading-skeleton m-3">
                  <div v-for="item in 5" :key="item">
                    <p class="w-25 mb-0 mb-2">Title</p>
                    <p class="m-0 p-0 py-3">Description</p>
                    <hr class="mb-3 mt-3">
                  </div>
                </div>
              </div>
              <div class="card-body p-0" v-else>
                <!-- List group -->
                <ul class="list-group list-group-flush" data-toggle="checklist" v-if="orders.length > 0">
                  <li class="checklist-entry list-group-item flex-column align-items-start py-4 px-4" v-for="(item,i) in orders" :key="i">
                    <div class="checklist-item checklist-item-success">
                      <div class="checklist-info">
                        <h5 class="checklist-title mb-0">{{ item.name }}</h5>
                        <small>Rp. {{ numberWithCommas(item.price) }} - {{ item.quantity }}x</small>
                        <small class="d-block">Total Pesanan: Rp. {{ numberWithCommas(item.total_price_item) }}</small>
                      </div>
                      <router-link class="btn btn-sm btn-success" :to="{ name: 'member.history.rate', params: { id_order: item.id_order } }">
                        Beri penilaian
                      </router-link>
                    </div>
                  </li>
                </ul>
                <ul class="list-group list-group-flush list" v-else>
                  <li class="list-group-item">
                    <div class="alert alert-secondary">
                      <h5 class="alert-heading">Oops!</h5>
                      <p>Pesanan masih kosong nih. <router-link :to="{ name: 'guest.product' }" >Beli produk disini!</router-link></p>
                    </div>
                  </li>
                </ul>
                <button @click="loadMore" v-if="orders.length !== this.totalData" class="my-3 btn btn-secondary d-block mx-auto px-5">
                  Lihat lainnya <i class="ni ni-bold-down"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      <!-- Footer -->
      <Footer />
    </div>
  </div>
</template>

<style scoped>
@import '../../assets/skeleton.css';
</style>

<script>
import Sidebar from "@/components/guest/Sidebar.vue"
import Topbar from "@/components/guest/Topbar.vue"
import Header from "@/components/guest/Header.vue"
import Footer from "@/components/guest/Footer.vue"
import axios from "axios";
import authHeader from "@/service/auth-header";

export default {
  components: {
    Footer,
    Sidebar,
    Header,
    Topbar
  },
  data(){
    return {
      carts: [],
      totalCart: 0,
      orders: [],
      limitData: 5,
      totalData: 0,
      isLoading: true,
      totalNotification: 0,
      notifications: []
    }
  },
  mounted() {
    if(authHeader()["Authorization"] === undefined) {
      this.$router.push({ name: 'auth.login' })
    }
    this.fetchData()
    if(authHeader()["Authorization"] !== undefined) {
      this.fetchNotification()
    }
    document.getElementsByTagName("body")[0].classList.remove("bg-default");
  },
  methods: {
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/user-order/user`, { headers: authHeader() }).then(response => {
        if(response.data.data != null) {
          this.totalData = response.data.data.length;
          this.orders = response.data.data.slice(0, this.limitData);
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      })

      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      if(this.carts.length > 0){
        this.totalCart = JSON.parse(localStorage.getItem('my-carts')).reduce((total, item) => {
          return total + item.quantity
        }, 0)
      }

      this.isLoading = false
    },
    async fetchNotification() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/notifications/user`, { headers: authHeader() }).then(response => {
        if(response.data.data != null) {
          this.totalNotification = response.data.data.filter(e => e.is_read === false).length
          this.notifications = response.data.data
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      })
    },
    loadMore(){
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/user-order/user`,{ headers: authHeader() }).then((response) => {
        this.orders = response.data.data.slice(0, this.limitData += 5);
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });
    },
    numberWithCommas(x) {
      return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".");
    }
  }
}
</script>