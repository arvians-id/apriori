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
    <div class="header bg-primary pb-7">
      <div class="container">
        <div class="header-body">
          <div class="row align-items-center">
            <div class="col-lg-6">
              <div class="pr-5">
                <h1 class="display-2 text-white font-weight-bold mb-0">Toko Ryzy Olshop</h1>
                <h2 class="display-4 text-white font-weight-light">Produk Original Berkualitas dan Terpercaya.</h2>
                <p class="text-white mt-4">Beli produk Ryzy Olshop online, produk terlengkap dan harga terbaik. Dapatkan berbagai promo menarik. Belanja aman dan nyaman hanya di Toko Ryzy.</p>
                <div class="mt-5">
                  <router-link :to="{ name: 'guest.product' }" class="btn btn-default my-2">Beli Sekarang!</router-link>
                </div>
              </div>
            </div>
            <div class="col-lg-6">
              <div class="row pt-5">
                <div class="col-md-6">
                  <div class="card">
                    <div class="card-body">
                      <div class="icon icon-shape bg-gradient-red text-white rounded-circle shadow mb-4">
                        <i class="ni ni-active-40"></i>
                      </div>
                      <h5 class="h3">Best Quality</h5>
                      <p>Kualitas yang sangat baik, terpercaya dan tahan lama.</p>
                    </div>
                  </div>
                  <div class="card">
                    <div class="card-body">
                      <div class="icon icon-shape bg-gradient-info text-white rounded-circle shadow mb-4">
                        <i class="ni ni-active-40"></i>
                      </div>
                      <h5 class="h3">Shopping Convenience</h5>
                      <p>Pemesanan yang sangat nyaman dan aman.</p>
                    </div>
                  </div>
                </div>
                <div class="col-md-6 pt-lg-5 pt-4">
                  <div class="card mb-4">
                    <div class="card-body">
                      <div class="icon icon-shape bg-gradient-success text-white rounded-circle shadow mb-4">
                        <i class="ni ni-active-40"></i>
                      </div>
                      <h5 class="h3">Fast Response</h5>
                      <p>Tingkat kecepatan respon 100% kepada pembeli.</p>
                    </div>
                  </div>
                  <div class="card mb-4">
                    <div class="card-body">
                      <div class="icon icon-shape bg-gradient-warning text-white rounded-circle shadow mb-4">
                        <i class="ni ni-active-40"></i>
                      </div>
                      <h5 class="h3">Good Service</h5>
                      <p>Pelayanan yang tidak akan membuat kecewa.</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <section class="py-6 pb-9 bg-default">
        <div class="container">
          <div class="row justify-content-center text-center">
            <div class="col-md-6">
              <h2 class="display-3 text-white">About Us</h2>
              <p class="lead text-white">
                Toko Ryzy merupakan toko yang menjual berbagai macam produk rumah tangga yang berkualitas. Berdiri sejak tahun 2021 dan sudah memiliki
                ribuan pelanggan hingga saat ini!
              </p>
            </div>
          </div>
        </div>
      </section>
      <section class="section section-lg pt-lg-0 mt--7">
        <div class="container">
          <div class="row justify-content-center">
            <div class="col-lg-12">
              <div class="row">
                <div class="col-lg-4" v-for="item in apriories" :key="item.id_apriori">
                  <div class="card card-lift--hover shadow border-0">
                    <div class="card-body py-5">
                      <router-link :to="{ name: 'guest.product.recommendation', params: { code: item.code, id: item.id_apriori } }" class="display-2 text-dark">{{ item.discount }}%</router-link>
                      <p class="text-muted">discount</p>
                      <ul class="list-unstyled my-4">
                        <li v-for="(value,i) in item.item.split(', ')" :key="i">
                          <div class="d-flex align-items-center mb-2">
                            <div class="icon icon-xs icon-shape bg-gradient-primary text-white shadow rounded-circle">
                              <i class="ni ni-basket"></i>
                            </div>
                            <span class="pl-2 text-sm">{{ UpperWord(value) }}</span>
                          </div>
                        </li>
                      </ul>
                      <div>
                        <span class="badge badge-pill badge-primary">#TOP RATED</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
      <section class="py-6">
        <div class="container">
          <div class="row row-grid align-items-center">
            <div class="col-md-6 order-md-2">
              <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/guling.png" class="img-fluid">
            </div>
            <div class="col-md-6 order-md-1">
              <div class="pr-md-5">
                <h1 class="text-white">Guling Berkualitas</h1>
                <p class="text-white">Guling merek Inoac adalah yang terbaik. Meskipun merek lokal, namun Inoac mampu menghadirkan produk yang kualitasnya bersaing dengan produk-produk luar.
                  Bahkan, busa yang dipakai dalam setiap produk Inoac dijamin memberikan kenyamanan tingkat tinggi dan sangat awet.
                  Inoac menjamin dapat memberikan kenyamanan beristirahat yang sangat baik bagi setiap penggunanya</p>
                <ul class="list-unstyled">
                  <li class="py-2">
                    <div class="d-flex align-items-center">
                      <div>
                        <div class="badge badge-circle badge-success mr-3">
                          <i class="ni ni-like-2"></i>
                        </div>
                      </div>
                      <div>
                        <h4 class="mb-0 text-secondary">Merek Inoac</h4>
                      </div>
                    </div>
                  </li>
                  <li class="py-2">
                    <div class="d-flex align-items-center">
                      <div>
                        <div class="badge badge-circle badge-success mr-3">
                          <i class="ni ni-paper-diploma"></i>
                        </div>
                      </div>
                      <div>
                        <h4 class="mb-0 text-secondary">Isian Foam</h4>
                      </div>
                    </div>
                  </li>
                  <li class="py-2">
                    <div class="d-flex align-items-center">
                      <div>
                        <div class="badge badge-circle badge-success mr-3">
                          <i class="ni ni-satisfied"></i>
                        </div>
                      </div>
                      <div>
                        <h4 class="mb-0 text-secondary">Material dari Microtex</h4>
                      </div>
                    </div>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </section>
      <section>
        <div class="container">
          <div class="row row-grid align-items-center">
            <div class="col-md-6">
              <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/bantal.png" class="img-fluid">
            </div>
            <div class="col-md-6">
              <div class="pr-md-5">
                <h1 class="text-white">Bantal Berkualitas</h1>
                <p class="text-white">Bantal merek Inoac adalah yang terbaik. Meskipun merek lokal, namun Inoac mampu menghadirkan produk yang kualitasnya bersaing dengan produk-produk luar.
                  Bahkan, busa yang dipakai dalam setiap produk Inoac dijamin memberikan kenyamanan tingkat tinggi dan sangat awet.
                  Inoac menjamin dapat memberikan kenyamanan beristirahat yang sangat baik bagi setiap penggunanya</p>
                <ul class="list-unstyled">
                  <li class="py-2">
                    <div class="d-flex align-items-center">
                      <div>
                        <div class="badge badge-circle badge-success mr-3">
                          <i class="ni ni-like-2"></i>
                        </div>
                      </div>
                      <div>
                        <h4 class="mb-0 text-secondary">Merek Inoac</h4>
                      </div>
                    </div>
                  </li>
                  <li class="py-2">
                    <div class="d-flex align-items-center">
                      <div>
                        <div class="badge badge-circle badge-success mr-3">
                          <i class="ni ni-paper-diploma"></i>
                        </div>
                      </div>
                      <div>
                        <h4 class="mb-0 text-secondary">Isian Foam</h4>
                      </div>
                    </div>
                  </li>
                  <li class="py-2">
                    <div class="d-flex align-items-center">
                      <div>
                        <div class="badge badge-circle badge-success mr-3">
                          <i class="ni ni-satisfied"></i>
                        </div>
                      </div>
                      <div>
                        <h4 class="mb-0 text-secondary">Material dari Microtex</h4>
                      </div>
                    </div>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </section>
      <section class="py-7 section-nucleo-icons bg-white overflow-hidden">
        <div class="container">
          <div class="text-center">
            <h4 class="display-4 mb-5 mt-5">Also Available on these E-Commerce</h4>
            <div class="row justify-content-center">
              <div class="w-10 mx-2 mb-2">
                <a href="https://www.lazada.co.id/shop/ryzy-olshop/" target="_blank">
                  <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/lazada.png" width="75" class="img-fluid">
                </a>
              </div>
              <div class="w-10 mx-2 mb-2">
                <a href="https://shopee.co.id/ryzy.olshop" target="_blank">
                  <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/shopee.png" width="60" class="img-fluid">
                </a>
              </div>
              <div class="w-10 mx-2 mb-2">
                <a href=" https://www.creative-tim.com/product/argon-dashboard" target="_blank">
                  <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/tokopedia.png" width="60" class="img-fluid">
                </a>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
    <footer class="py-5 bg-primary" id="footer-main">
      <div class="container">
        <div class="row align-items-center justify-content-xl-between">
          <div class="col-lg-6">
            <div class="copyright text-center text-xl-left text-muted text-white">
              &copy; 2019 <a href="https://www.creative-tim.com" class="font-weight-bold ml-1 text-secondary" target="_blank">Creative Tim</a>
            </div>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
import Sidebar from "@/components/guest/Sidebar.vue"
import Topbar from "@/components/guest/Topbar.vue"
import Header from "@/components/guest/Header.vue"
import axios from "axios";
import authHeader from "@/service/auth-header";

export default {
  components: {
    Sidebar,
    Header,
    Topbar
  },
  mounted() {
    this.fetchData()
    if(authHeader()["Authorization"] !== undefined) {
      this.fetchNotification()
    }
    document.getElementsByTagName("body")[0].classList.remove("bg-default");
  },
  data: function () {
    return {
      apriories: [],
      carts: [],
      totalCart: 0,
      totalNotification: 0,
      notifications: []
    };
  },
  methods: {
    fetchData() {
      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      axios.get(`${process.env.VUE_APP_SERVICE_URL}/apriori/actives`,{ headers: authHeader() }).then((response) => {
        this.apriories = response.data.data.slice(0,3);
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });

      if(this.carts.length > 0){
        this.totalCart = JSON.parse(localStorage.getItem('my-carts')).reduce((total, item) => {
          return total + item.quantity
        }, 0)
      }
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
    UpperWord(str) {
      return str.toLowerCase().replace(/\b[a-z]/g, function (letter) {
        return letter.toUpperCase();
      })
    }
  }
}
</script>