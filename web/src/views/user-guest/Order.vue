<template>
  <!-- Sidenav -->
  <Sidebar />
  <!-- Main content -->
  <div class="main-content" id="panel">
    <!-- Topnav -->
    <Topbar :totalCart="totalCart" :carts="carts" />
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
                <h5 class="h3 mb-0">Riwayat Pesanan Saya</h5>
              </div>
              <!-- Card body -->
              <div class="card-body p-0" v-if="isLoading">
                <p class="p-3 mt-2 text-center">Loading...</p>
              </div>
              <div class="card-body p-0" v-else>
                <!-- List group -->
                <ul class="list-group list-group-flush" data-toggle="checklist" v-if="orders.length > 0">
                  <li class="checklist-entry list-group-item flex-column align-items-start py-4 px-4" v-for="(item,i) in orders" :key="i">
                    <div :class="getColor(`checklist-item checklist-item-`, item.transaction_status)">
                      <div class="checklist-info">
                        <router-link :to="{ name: 'guest.order.detail', params: { order_id: item.order_id } }">
                          <h5 class="checklist-title mb-0">Order ID #{{ item.order_id }}</h5>
                        </router-link>
                        <small>{{ item.transaction_time }}</small>
                      </div>
                      <div :class="getColor(`btn disabled btn-sm btn-`,item.transaction_status)">
                        {{ item.transaction_status }}
                      </div>
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
    }
  },
  mounted() {
    if(authHeader()["Authorization"] === undefined) {
      this.$router.push({ name: 'auth.login' })
    }
    this.fetchData()
    document.getElementsByTagName("body")[0].classList.remove("bg-default");
  },
  methods: {
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/user-order`, { headers: authHeader() }).then(response => {
        if(response.data.data != null) {
          this.totalData = response.data.data.length;
          this.orders = response.data.data.slice(0, this.limitData);
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
    loadMore(){
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/user-order`,{ headers: authHeader() }).then((response) => {
        this.orders = response.data.data.slice(0, this.limitData += 5);
      });
    },
    getColor(classNames, status) {
      let className = classNames
      if (status === "settlement") {
        className += "success"
      } else if (status === "pending") {
        className += "info"
      } else {
        className += "danger"
      }

      return className
    }
  }
}
</script>