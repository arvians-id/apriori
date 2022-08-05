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
          <div class="card">
            <div class="card-body">
              <form @submit.prevent="submitSearch" method="GET">
                <div class="form-group">
                  <input type="text" class="form-control" name="search" v-model="search" placeholder="contoh: Bantal Polkadot">
                </div>
                <button class="btn btn-primary">Search</button>
              </form>
            </div>
          </div>
          <div class="card-wrapper">
            <!-- Custom form validation -->
            <div class="card">
              <!-- Card header -->
              <div class="card-header">
                <!-- Title -->
                <h5 class="h3 mb-0">Semua Produk</h5>
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
                      <ul class="list-group list-group-flush">
                        <li class="list-group-item text-danger font-weight-bold">Rp {{ numberWithCommas(item.price) }}</li>
                      </ul>
                      <div class="card-body">
                        <router-link :to="{ name: 'guest.product.detail', params: { code: item.code } }" class="card-title mb-3 text-dark">{{ item.name }}</router-link>
                        <p class="card-text mb-4">{{ item.description.length > 35 ? item.description.slice(0, 35) + "..." : item.description }}</p>
                        <div class="d-flex justify-content-between">
                          <button type="button" @click="min(item)" class="btn btn-danger btn-sm">Kurangi</button>
                          <button type="button" @click="add(item)" class="btn btn-primary btn-sm">Tambah</button>
                        </div>
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
          <!-- Custom form validation -->
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
                <div class="col-12 col-md-6 col-lg-4 col-xl-3" v-for="item in recommendation" :key="item.apriori_id">
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
                      <br>
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
    document.getElementsByTagName("body")[0].classList.remove("bg-default");
  },
  data: function () {
    return {
      products: [],
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
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/products?search=${this.search}`,{ headers: authHeader() }).then((response) => {
        if(response.data.data != null) {
          this.totalData = response.data.data.length;
          this.products = response.data.data.slice(0, this.limitData);
        }
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
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/products?search=${this.search}`,{ headers: authHeader() }).then((response) => {
        this.products = response.data.data.slice(0, this.limitData += 8);
      });
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
    add(item) {
      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      let productItem = this.carts.find(product => product.code === item.code);
      if (productItem) {
        productItem.quantity += 1
        productItem.totalPricePerItem = productItem.price * productItem.quantity
        this.totalCart += 1
      } else {
        this.carts.push({
          id_product: item.id_product,
          code: item.code,
          name: item.name,
          price: item.price,
          image: item.image,
          quantity: 1,
          totalPricePerItem: item.price
        });
        this.totalCart += 1
      }

      localStorage.setItem('my-carts', JSON.stringify(this.carts));
    },
    min(item){
      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      let productItem = this.carts.find(product => product.code === item.code);
      if(productItem !== undefined) {
        if (productItem.quantity > 1) {
          productItem.quantity -= 1
          productItem.totalPricePerItem -= productItem.price
          this.totalCart -= 1
        } else {
          this.carts.splice(this.carts.indexOf(productItem), 1);
          this.totalCart -= 1
        }
        localStorage.setItem('my-carts', JSON.stringify(this.carts));
      }
    }
  }
}
</script>