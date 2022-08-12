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
              <div class="card-body pb-0">
                <form @submit.prevent="submitSearch" method="GET">
                  <div class="input-group">
                    <input type="text" class="form-control" name="search" v-model="search" placeholder="contoh: Bantal Polkadot">
                    <div class="input-group-append">
                      <button class="btn btn-primary"><i class="fas fa-search"></i></button>
                    </div>
                  </div>
                </form>
              </div>
              <div class="row">
                <div class="col-12 col-lg-3">
                  <!-- Card header -->
                  <div class="card-header">
                    <!-- Title -->
                    <h5 class="h3 mb-0">Semua Kategori</h5>
                  </div>
                  <div class="card-body">
                    <div class="form-check mb-1" v-for="category in categories" :key="category.id_category">
                      <input class="form-check-input" type="radio" name="inlineRadioOptions" :id="category.id_category">
                      <label class="form-check-label" :for="category.id_category">{{ category.name }}</label>
                    </div>
                  </div>
                </div>
                <div class="col-12 col-lg-9">
                  <!-- Card header -->
                  <div class="card-header">
                    <!-- Title -->
                    <h5 class="h3 mb-0"><span class="text-primary">{{ allProducts.length }}</span> Produk ditemukan</h5>
                  </div>
                  <!-- Card body -->
                  <div class="card-body" v-if="isLoading">
                    <p class="mt-2 text-center">Loading...</p>
                  </div>
                  <div class="card-body" v-else>
                    <!-- List group -->
                    <div class="row" v-if="products.length > 0">
                      <div class="col-12 col-md-6 col-lg-4 col-xl-3" v-for="item in products" :key="item.id_product">
                        <div class="card">
                          <div class="embed-responsive embed-responsive-16by9">
                            <img class="card-img-top embed-responsive-item" :src="getImage(item.image)" alt="Preview Image">
                          </div>
                          <div class="card-body">
                            <router-link :to="{ name: 'guest.product.detail', params: { code: item.code } }" class="card-title mb-3 text-dark">{{ item.name }}</router-link>
                            <p class="font-weight-bold">Rp {{ numberWithCommas(item.price) }}</p>
                            <small><i class="ni ni-pin-3 text-primary"></i> Tanggerang, Banten</small>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div class="row" v-else>
                      <div class="col-12">
                        <div class="alert alert-secondary">
                          <h5 class="alert-heading">Oops!</h5>
                          <p>Tidak ada produk yang tersedia.</p>
                        </div>
                      </div>
                    </div>
                    <button @click="loadMore()" v-if="products.length !== this.totalData" class="btn btn-secondary d-block mx-auto px-5">
                      Lihat lainnya <i class="ni ni-bold-down"></i>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="card">
            <!-- Card header -->
            <div class="card-header">
              <!-- Title -->
              <h5 class="h3 mb-0">Rekomendasi Paket Diskon Barang</h5>
            </div>
            <div class="card-body" v-if="isLoading2">
              <p class="mt-2 text-center">Loading...</p>
            </div>
            <div class="card-body" v-else>
              <div class="row" v-if="recommendation.length > 0">
                <div class="col-12 col-md-6 col-lg-3" v-for="item in recommendation" :key="item.apriori_id">
                  <div class="card card-pricing border-0 mb-4">
                    <div class="embed-responsive embed-responsive-16by9">
                      <img class="card-img-top embed-responsive-item" :src="getImage(item.image)" alt="Preview Image">
                    </div>
                    <div class="card-body pb-3">
                      <router-link :to="{ name: 'guest.product.recommendation', params: { code: item.code, id: item.id_apriori } }" class="text-dark d-none d-lg-block">
                        Paket Barang {{ item.item.length > 20 ? UpperWord(item.item.slice(0, 20)) + "..." : UpperWord(item.item) }}
                      </router-link>
                      <router-link :to="{ name: 'guest.product.recommendation', params: { code: item.code, id: item.id_apriori } }" class="text-dark d-block d-lg-none">
                        Paket Barang {{ UpperWord(item.item) }}
                      </router-link>
                      <ul class="list-unstyled">
                        <li v-for="(value,i) in item.item.split(', ')" :key="i">
                          <div class="d-flex align-items-center">
                            <div class="icon icon-xs icon-shape bg-gradient-primary text-white shadow rounded-circle">
                              <i class="ni ni-basket"></i>
                            </div>
                            <span class="pl-2 text-sm">{{ UpperWord(value) }}</span>
                          </div>
                        </li>
                      </ul>
                      <div class="card-footer p-0 pt-2 text-center m-0">
                        <router-link :to="{ name: 'guest.product.recommendation', params: { code: item.code, id: item.id_apriori } }" class="font-weight-bold">
                          {{ item.discount }}%
                        </router-link>
                        <small class="text-muted mb-0">diskon</small>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="row" v-else>
                <div class="col-12">
                  <div class="alert alert-secondary">
                    <h5 class="alert-heading">Oops!</h5>
                    <p>Tidak ada rekomendasi yang tersedia.</p>
                  </div>
                </div>
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
  .card-img-top {
    width: 100%;
    object-fit: cover;
  }
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
  mounted() {
    this.fetchData()
    this.fetchDataRecommendation()
    this.fetchCategory()
    document.getElementsByTagName("body")[0].classList.remove("bg-default");
  },
  data: function () {
    return {
      products: [],
      allProducts: [],
      categories: [],
      carts: [],
      totalCart: 0,
      recommendation: [],
      limitData: 8,
      totalData: 0,
      isLoading: true,
      isLoading2: true,
      search: ""
    };
  },
  methods: {
    submitSearch(){
      let search = ""
      if(this.search !== ""){
        this.$router.replace({ name: "guest.product", query: { search: this.search } })
        search = "?search=" + this.search
      } else {
        this.$router.replace({ name: "guest.product" })
      }
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/products${search}`,{ headers: authHeader() }).then((response) => {
        if(response.data.data != null) {
          this.totalData = response.data.data.length;
          this.allProducts = response.data.data;
          this.products = response.data.data.slice(0, this.limitData);
        } else {
          this.totalData = 0;
          this.allProducts = [];
          this.products = [];
        }
      });
    },
    async fetchCategory() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/categories`,{ headers: authHeader() }).then((response) => {
        this.categories = response.data.data;
      });
    },
    async fetchData() {
      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      let search = this.$route.query.search
      if (search === undefined) {
        search = ""
      }
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/products?search=${search}`,{ headers: authHeader() }).then((response) => {
        if(response.data.data != null) {
          this.totalData = response.data.data.length;
          this.allProducts = response.data.data;
          this.products = response.data.data.slice(0, this.limitData);
        }
      });

      if(this.carts.length > 0){
        this.totalCart = JSON.parse(localStorage.getItem('my-carts')).reduce((total, item) => {
          return total + item.quantity
        }, 0)
      }

      this.isLoading = false
    },
    async fetchDataRecommendation() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/apriori/actives`,{ headers: authHeader() }).then((response) => {
        if(response.data.data != null) {
          this.recommendation = response.data.data;
        }
      }).catch((error) => {
        console.log(error)
      });

      this.isLoading2 = false
    },
    loadMore(){
      this.products = this.allProducts.slice(0, this.limitData += 8);
    },
    numberWithCommas(x) {
      return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".");
    },
    getImage(image) {
      return image;
    },
    UpperWord(str) {
      return str.toLowerCase().replace(/\b[a-z]/g, function (letter) {
        return letter.toUpperCase();
      })
    },
  }
}
</script>