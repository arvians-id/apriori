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
        <div class="col-xl-12 order-xl-2">
          <div class="card card-profile">
            <!-- Card header -->
            <div class="card-body mt-5">
              <div class="row align-items-center mx-auto" v-if="isLoading">
                <p class="p-3 mt-2 text-center">Loading...</p>
              </div>
              <div class="row d-flex justify-content-center" v-else>
                <div class="col-12 col-lg-3 mb-2">
                  <img :src="getImage()" class="img-fluid mb-2" width="500">
                </div>
                <div class="col-12 col-lg-6">
                  <div class="text-left">
                    <h5 class="h2 text-uppercase p-0 m-0">{{ product.name }}</h5>
                    <p class="p-0 m-0">code : {{ product.code }}</p>
                    <div class="h1 font-weight-bold">Rp. {{ product.price }}</div>
                    <hr class="m-2">
                    <ul class="nav nav-tabs" id="myTab" role="tablist">
                      <li class="nav-item">
                        <a class="nav-link active" id="detail-tab" data-toggle="tab" href="#detail" role="tab" aria-controls="detail" aria-selected="true">Detail</a>
                      </li>
                      <li class="nav-item">
                        <a class="nav-link" id="other-detail-tab" data-toggle="tab" href="#other-detail" role="tab" aria-controls="other-detail" aria-selected="false">Lainnya</a>
                      </li>
                    </ul>
                    <div class="tab-content" id="myTabContent">
                      <div class="tab-pane fade show active" id="detail" role="tabpanel" aria-labelledby="detail-tab">
                        <p class="font-weight-bold mb-0 mt-2">Kondisi</p>
                        <p>Original baru</p>
                        <p class="font-weight-bold mb-0 mt-2">Kategori</p>
                        <p>
                          <router-link :to="{ name: 'guest.product', query: { category: category } }" v-for="(category, i) in product.category.split(', ')" :key="i" class="text-primary font-weight-bold">
                            {{ category }}{{ product.category.split(', ').length - 1 != i ? ', ' : '' }}
                          </router-link>
                        </p>
                        <p class="font-weight-bold mb-0 mt-2">Berat Satuan</p>
                        <p>{{ product.mass }} gram</p>
                        <p class="font-weight-bold mb-0 mt-2">Deskripsi</p>
                        <p>{{ product.description == "" ? "Tidak ada deskripsi" : product.description }}</p>
                        <hr class="m-0 mb-3">
                        <div class="media">
                          <img src="https://my-apriori.s3.ap-southeast-1.amazonaws.com/assets/ryzy.jpg" width="53" class="mr-3" alt="...">
                          <div class="media-body">
                            <h3 class="mt-0 mb-0">Toko Ryzy Olshop</h3>
                            <p>Produk Original Berkualitas dan Terpercaya..</p>
                          </div>
                        </div>
                        <hr class="m-0 mb-3">
                        <p class="font-weight-bold mb-0 mt-2">Pengiriman</p>
                        <p class="mb-1"><i class="ni ni-pin-3"></i> Dikirim dari Tanggerang, Banten</p>
                        <p><i class="ni ni-delivery-fast"></i> Tersedia pengiriman dengan TIKI, JNE dan POS Indonesia</p>
                      </div>
                      <div class="tab-pane fade" id="other-detail" role="tabpanel" aria-labelledby="other-detail-tab">
                        <p class="font-weight-bold mb-0 mt-2">Tanggal Dibuat</p>
                        <p>{{ product.created_at }}</p>
                        <p class="font-weight-bold mb-0 mt-2">Terakhir Diubah</p>
                        <p>{{ product.updated_at }}</p>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="col-12 col-lg-2">
                  <div class="border p-3 rounded text-center" style="color: #525f7f">
                    <p class="mb-0">Atur jumlah yang pembelian</p>
                    <div>
                      <button type="button" @click="min(product)" class="btn btn-danger btn-sm">-</button>
                      <button class="btn disabled">{{ quantity }} item</button>
                      <button type="button" @click="add(product)" class="btn btn-primary btn-sm">+</button>
                    </div>
                  </div>
                  <h3 class="mb-0 mt-3">Produk yang serupa</h3>
                  <hr class="mb-3 p-0">
                  <div class="card card-pricing border shadow-none" v-for="item in productSimilarCategory" :key="item.id_product">
                    <div class="embed-responsive embed-responsive-16by9">
                      <img class="card-img-top embed-responsive-item" :src="getImage(item.image)" alt="Preview Image">
                    </div>
                    <div class="card-body">
                      <router-link :to="{ name: 'guest.product.detail', params: { code: item.code } }" class="card-title m-0">{{ item.name }}</router-link>
                      <p class="card-text p-0 m-0">Rp. {{ item.price }}</p>
                    </div>
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
                      <router-link :to="{ name: 'guest.product.recommendation', params: { code: item.apriori_code, id: item.apriori_id } }" class="text-dark d-none d-lg-block">
                        Paket Barang {{ item.apriori_item.length > 20 ? UpperWord(item.apriori_item.slice(0, 20)) + "..." : UpperWord(item.apriori_item) }}
                      </router-link>
                      <router-link :to="{ name: 'guest.product.recommendation', params: { code: item.apriori_code, id: item.apriori_id } }" class="text-dark d-block d-lg-none">
                        Paket Barang {{ UpperWord(item.apriori_item) }}
                      </router-link>
                      <ul class="list-unstyled">
                        <li v-for="(value,i) in item.apriori_item.split(', ')" :key="i">
                          <div class="d-flex align-items-center">
                            <div class="icon icon-xs icon-shape bg-gradient-primary text-white shadow rounded-circle">
                              <i class="ni ni-basket"></i>
                            </div>
                            <span class="pl-2 text-sm">{{ UpperWord(value) }}</span>
                          </div>
                        </li>
                      </ul>
                      <div class="card-footer p-0 pt-2 text-center m-0">
                        <router-link :to="{ name: 'guest.product.recommendation', params: { code: item.apriori_code, id: item.apriori_id } }" class="font-weight-bold">
                          {{ item.apriori_discount }}%
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
      product: [],
      recommendation: [],
      productSimilarCategory: [],
      carts: [],
      totalCart: 0,
      quantity: 0,
      isLoading: true,
      isLoading2: true,
      categories: [],
    };
  },
  methods: {
    allFetch(){
      this.fetchCategories()
      this.fetchData()
      this.fetchDataRecommendation()
      this.fetchSimilarCategory()
      document.getElementsByTagName("body")[0].classList.remove("bg-default");
    },
    async fetchCategories(){
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/categories`, { headers: authHeader() })
        .then(response => {
          this.categories = response.data.data
        })
        .catch(error => {
          console.log(error)
        })
    },
    async fetchSimilarCategory(){
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/products/${this.$route.params.code}/category`, { headers: authHeader() })
          .then(response => {
            this.productSimilarCategory = response.data.data
          })
          .catch(error => {
            console.log(error)
          })
    },
    async fetchData() {
      localStorage.getItem("my-carts")
          ? (this.carts = JSON.parse(localStorage.getItem("my-carts")))
          : (this.carts = []);

      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/products/${this.$route.params.code}`, { headers: authHeader() }).then((response) => {
        this.product = response.data.data;
      });

      if(this.carts.length > 0){
        this.totalCart = JSON.parse(localStorage.getItem('my-carts')).reduce((total, item) => {
          return total + item.quantity
        }, 0)
      }

      let productItem = this.carts.find(product => product.code === this.$route.params.code);
      this.quantity = productItem ? productItem.quantity : 0;

      this.isLoading = false
    },
    async fetchDataRecommendation() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/products/${this.$route.params.code}/recommendation`, { headers: authHeader() }).then((response) => {
        if(response.data.data != null) {
          this.recommendation = response.data.data;
        } else {
          this.recommendation = []
        }
      }).catch((error) => {
        console.log(error);
      });

      this.isLoading2 = false
    },
    getImage() {
      return this.product.image
    },
    UpperWord(str) {
      return str.toLowerCase().replace(/\b[a-z]/g, function (letter) {
        return letter.toUpperCase();
      })
    },
    add(item) {
      let productItem = this.carts.find(product => product.code === item.code);
      if (productItem) {
        productItem.quantity += 1
        productItem.totalPricePerItem = productItem.price * productItem.quantity
        this.quantity += 1
        this.totalCart += 1
      } else {
        this.carts.push({
          id_product: item.id_product,
          code: item.code,
          name: item.name,
          price: item.price,
          mass: item.mass,
          image: item.image,
          quantity: 1,
          totalPricePerItem: item.price
        });
        this.quantity = 1
        this.totalCart += 1
      }

      localStorage.setItem('my-carts', JSON.stringify(this.carts));
    },
    min(item){
      let productItem = this.carts.find(product => product.code === item.code);
      if (productItem !== undefined) {
        if (productItem.quantity > 1) {
          productItem.quantity -= 1
          productItem.totalPricePerItem -= productItem.price
          this.quantity -= 1
          this.totalCart -= 1
        } else {
          this.quantity = 0
          this.totalCart -= 1
          this.carts.splice(this.carts.indexOf(productItem), 1);
        }
        localStorage.setItem('my-carts', JSON.stringify(this.carts));
      }
    }
  }
}
</script>