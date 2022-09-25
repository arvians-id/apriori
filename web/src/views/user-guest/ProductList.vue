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
              <div class="card-body pb-0">
                <form @submit.prevent="submitSearch" method="GET">
                  <div class="input-group">
                    <input type="text" class="form-control" name="search" v-model="search" placeholder="contoh: Bantal Polkadot">
                    <div class="input-group-append">
                      <button class="btn btn-primary" :disabled="search == ''"><i class="fas fa-search"></i></button>
                    </div>
                  </div>
                </form>
              </div>
              <div class="row">
                <div class="col-12 col-lg-3">
                  <!-- Card header -->
                  <div class="card-header">
                    <!-- Title -->
                    <h5 class="h3 mb-0 mb-1">Semua Kategori</h5>
                  </div>
                  <div class="card-body" v-if="isLoading3">
                    <div class="loading-skeleton">
                      <p class="p-3 mb-2">This is title of category</p>
                      <p class="p-3 mb-2">This is title of category</p>
                      <p class="p-3 mb-2">This is title of category</p>
                      <p class="p-3 mb-2">This is title of category</p>
                      <p class="p-3 mb-2">This is title of category</p>
                    </div>
                  </div>
                  <div class="card-body" v-else>
                    <ul class="list-group" v-if="categories.length > 0">
                      <li class="list-group-item d-flex justify-content-between align-items-center" v-for="category in categories" :key="category.id_category">
                        <router-link :to="{ name: 'guest.product', query: { category: category.name } }">{{ category.name }}</router-link>
                      </li>
                    </ul>
                    <div class="row" v-else>
                      <div class="col-12">
                        <div class="alert alert-secondary">
                          <h5 class="alert-heading">Oops!</h5>
                          <p>Tidak ada kategori yang tersedia.</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="col-12 col-lg-9">
                  <!-- Card header -->
                  <div class="card-header d-lg-none justify-content-between">
                    <!-- Title -->
                    <h5 class="h3 mb-1"><span class="text-primary">{{ allProducts.length }}</span> Produk ditemukan</h5>
                    <div v-if="Object.keys(this.$route.query).length > 0">
                      <button class="btn btn-white btn-sm mb-2" v-for="(query, i) in this.$route.query" :key="i">{{ UpperWord(i) + " : " + UpperWord(query) }}</button>
                      <router-link class="btn btn-danger btn-sm mb-2" :to="{ name: 'guest.product' }" @click="search = ''">Reset Pencarian</router-link>
                    </div>
                  </div>
                  <!-- Card header -->
                  <div class="card-header d-lg-flex d-none justify-content-between">
                    <!-- Title -->
                    <h5 class="h3 mb-1"><span class="text-primary">{{ allProducts.length }}</span> Produk ditemukan</h5>
                    <div v-if="Object.keys(this.$route.query).length > 0">
                      <button class="btn btn-white btn-sm" v-for="(query, i) in this.$route.query" :key="i">{{ UpperWord(i) + " : " + UpperWord(query) }}</button>
                      <router-link class="btn btn-danger btn-sm" :to="{ name: 'guest.product' }" @click="search = ''">Reset Pencarian</router-link>
                    </div>
                  </div>
                  <!-- Card body -->
                  <div class="card-body" v-if="isLoading">
                    <div class="loading-skeleton">
                      <div class="row">
                        <div class="col-12 col-md-6 col-lg-4">
                          <div class="card">
                            <div class="embed-responsive embed-responsive-16by9">
                              <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                            </div>
                            <div class="card-body">
                              <p class="card-title mb-1">This is title</p>
                              <p class="font-weight-bold">This is price of product</p>
                              <p class="w-50 mb-0">this is icon</p>
                            </div>
                          </div>
                        </div>
                        <div class="col-12 col-md-6 col-lg-4">
                          <div class="card">
                            <div class="embed-responsive embed-responsive-16by9">
                              <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                            </div>
                            <div class="card-body">
                              <p class="card-title mb-1">This is title</p>
                              <p class="font-weight-bold">This is price of product</p>
                              <p class="w-50 mb-0">this is icon</p>
                            </div>
                          </div>
                        </div>
                        <div class="col-12 col-md-6 col-lg-4">
                          <div class="card">
                            <div class="embed-responsive embed-responsive-16by9">
                              <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                            </div>
                            <div class="card-body">
                              <p class="card-title mb-1">This is title</p>
                              <p class="font-weight-bold">This is price of product</p>
                              <p class="w-50 mb-0">this is icon</p>
                            </div>
                          </div>
                        </div>
                        <div class="col-12 col-md-6 col-lg-4">
                          <div class="card">
                            <div class="embed-responsive embed-responsive-16by9">
                              <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                            </div>
                            <div class="card-body">
                              <p class="card-title mb-1">This is title</p>
                              <p class="font-weight-bold">This is price of product</p>
                              <p class="w-50 mb-0">this is icon</p>
                            </div>
                          </div>
                        </div>
                        <div class="col-12 col-md-6 col-lg-4">
                          <div class="card">
                            <div class="embed-responsive embed-responsive-16by9">
                              <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                            </div>
                            <div class="card-body">
                              <p class="card-title mb-1">This is title</p>
                              <p class="font-weight-bold">This is price of product</p>
                              <p class="w-50 mb-0">this is icon</p>
                            </div>
                          </div>
                        </div>
                        <div class="col-12 col-md-6 col-lg-4">
                          <div class="card">
                            <div class="embed-responsive embed-responsive-16by9">
                              <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                            </div>
                            <div class="card-body">
                              <p class="card-title mb-1">This is title</p>
                              <p class="font-weight-bold">This is price of product</p>
                              <p class="w-50 mb-0">this is icon</p>
                            </div>
                          </div>
                        </div>
                        <button class="btn btn-secondary d-block mx-auto px-5">This is button</button>
                      </div>
                    </div>
                  </div>
                  <div class="card-body" v-else>
                    <!-- List group -->
                    <div class="row" v-if="products.length > 0">
                      <div class="col-12 col-md-6 col-lg-4" v-for="item in products" :key="item.id_product">
                        <div class="card">
                          <div class="embed-responsive embed-responsive-16by9">
                            <img class="card-img-top embed-responsive-item" :src="item.image" alt="Preview Image">
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
              <h5 class="h3 mb-0">Rekomendasi</h5>
            </div>
            <div class="card-body" v-if="isLoading2">
              <div class="loading-skeleton">
                <div class="row">
                  <div class="col-12 col-md-6 col-lg-4">
                    <div class="card card-pricing border-0 mb-4">
                      <div class="embed-responsive embed-responsive-16by9">
                        <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                      </div>
                      <div class="card-body pb-3">
                        <p>This is title of products</p>
                        <p class="card-title mb-1">This is title</p>
                        <p class="font-weight-bold">This is price of product</p>
                        <div class="card-footer p-0 pt-2 m-0 text-center">
                          <p class="w-25 d-inline">this is discount</p>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="col-12 col-md-6 col-lg-4">
                    <div class="card card-pricing border-0 mb-4">
                      <div class="embed-responsive embed-responsive-16by9">
                        <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                      </div>
                      <div class="card-body pb-3">
                        <p>This is title of products</p>
                        <p class="card-title mb-1">This is title</p>
                        <p class="font-weight-bold">This is price of product</p>
                        <div class="card-footer p-0 pt-2 m-0 text-center">
                          <p class="w-25 d-inline">this is discount</p>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="col-12 col-md-6 col-lg-4">
                    <div class="card card-pricing border-0 mb-4">
                      <div class="embed-responsive embed-responsive-16by9">
                        <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                      </div>
                      <div class="card-body pb-3">
                        <p>This is title of products</p>
                        <p class="card-title mb-1">This is title</p>
                        <p class="font-weight-bold">This is price of product</p>
                        <div class="card-footer p-0 pt-2 m-0 text-center">
                          <p class="w-25 d-inline">this is discount</p>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="col-12 col-md-6 col-lg-4">
                    <div class="card card-pricing border-0 mb-4">
                      <div class="embed-responsive embed-responsive-16by9">
                        <img class="card-img-top embed-responsive-item" src="//placekitten.com/300/200" alt="Preview Image">
                      </div>
                      <div class="card-body pb-3">
                        <p>This is title of products</p>
                        <p class="card-title mb-1">This is title</p>
                        <p class="font-weight-bold">This is price of product</p>
                        <div class="card-footer p-0 pt-2 m-0 text-center">
                          <p class="w-25 d-inline">this is discount</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="card-body" v-else>
              <div class="row" v-if="recommendation.length > 0">
                <div class="col-12 col-md-6 col-lg-4" v-for="item in recommendation" :key="item.apriori_id">
                  <div class="card card-pricing border-0 mb-4">
                    <div class="embed-responsive embed-responsive-16by9">
                      <img class="card-img-top embed-responsive-item" :src="item.image" alt="Preview Image">
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
@import '../../assets/skeleton.css';

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
  watch: {
    '$route': function () {
      this.allFetch()
    }
  },
  mounted() {
    this.allFetch()
  },
  data: function () {
    return {
      products: [],
      allProducts: [],
      categories: [],
      carts: [],
      totalCart: 0,
      recommendation: [],
      limitData: 6,
      totalData: 0,
      isLoading: true,
      isLoading2: true,
      isLoading3: true,
      search: "",
      totalNotification: 0,
      notifications: []
    };
  },
  methods: {
    allFetch(){
      this.fetchData()
      this.fetchDataRecommendation()
      this.fetchCategory()
      if(authHeader()["Authorization"] !== undefined) {
        this.fetchNotification()
      }
      document.getElementsByTagName("body")[0].classList.remove("bg-default");
    },
    submitSearch(){
      let search = ""
      if(this.search !== ""){
        this.$router.push({ query: Object.assign({}, this.$route.query, { search: this.search }) })
        search = "?search=" + this.search
        if(this.$route.query.category !== undefined){
          search += "&category=" + this.$route.query.category
        }
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
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });;
    },
    async fetchCategory() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/categories`,{ headers: authHeader() }).then((response) => {
        if(response.data.data != null) {
          this.categories = response.data.data;
        } else {
          this.categories = [];
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });

      this.isLoading3 = false;
    },
    async fetchData() {
      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      let search = ""
      if (this.$route.query.search !== undefined) {
        search = "?search=" + this.$route.query.search
        if(this.$route.query.category !== undefined){
          search += "&category=" + this.$route.query.category
        }
      }else {
        if(this.$route.query.category !== undefined){
          search += "?category=" + this.$route.query.category
        }
      }

      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/products${search}`,{ headers: authHeader() }).then((response) => {
        if(response.data.data != null) {
          this.totalData = response.data.data.length;
          this.allProducts = response.data.data;
          this.products = response.data.data.slice(0, this.limitData);
        } else {
          this.totalData = 0;
          this.allProducts = [];
          this.products = [];
        }
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

      this.isLoading = false
    },
    async fetchDataRecommendation() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/apriori/actives`,{ headers: authHeader() }).then((response) => {
        if(response.data.data != null) {
          this.recommendation = response.data.data;
        }
      }).catch((error) => {
        console.log(error)
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });

      this.isLoading2 = false
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
      });
    },
    loadMore(){
      this.products = this.allProducts.slice(0, this.limitData += 6);
    },
    numberWithCommas(x) {
      return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ".");
    },
    UpperWord(str) {
      return str.toLowerCase().replace(/\b[a-z]/g, function (letter) {
        return letter.toUpperCase();
      })
    },
  }
}
</script>