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
            <div class="card-header">
              <h3 class="mb-0">Detail Produk</h3>
            </div>
            <div class="row align-items-center mx-auto" v-if="isLoading">
              <p class="p-3 mt-2 text-center">Loading...</p>
            </div>
            <div class="row align-items-center" v-else>
              <div class="col-12 col-lg-6 text-center">
                <img :src="getImage()" class="img-fluid my-5" width="500">
              </div>
              <div class="col-12 col-lg-6">
                <div class="card-body pt-0">
                  <div class="row">
                    <div class="col">
                      <div class="card-profile-stats d-flex justify-content-center">
                        <div>
                          <span class="heading">Dibuat</span>
                          <span class="description">{{ product.created_at }}</span>
                        </div>
                        <div>
                          <span class="heading">Terakhir diubah</span>
                          <span class="description">{{ product.updated_at }}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div class="text-center">
                    <h5 class="h3 text-uppercase">
                      {{ product.name }} - {{ product.code }}
                    </h5>
                    <div class="h5 font-weight-300">
                      <i class="ni location_pin mr-2"></i>Rp. {{ product.price }}
                    </div>
                    <div>
                      <i class="ni education_hat mr-2"></i> {{ product.description }}
                    </div>
                    <button type="button" @click="min(product)" class="btn btn-danger btn-sm mt-3">-</button>
                    <button class="btn disabled mt-3">{{ quantity }} item</button>
                    <button type="button" @click="add(product)" class="btn btn-primary btn-sm mt-3">+</button>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="card">
            <!-- Card header -->
            <div class="card-header d-flex justify-content-between">
              <h3 class="mb-0">Rekomendasi Paket Diskon Untuk Kamu</h3>
            </div>
            <div class="card-body" v-if="isLoading2">
              <p class="mt-2 text-center">Loading...</p>
            </div>
            <!-- Card body -->
            <div class="card-body" v-else>
              <div class="row">
                <div class="col-12 col-md-6 col-lg-4 col-xl-3" v-for="item in recommendation" :key="item.apriori_id">
                  <div class="card card-pricing border-0 text-center mb-4">
                    <div class="card-header bg-transparent">
                      <router-link :to="{ name: 'guest.product.recommendation', params: { code: item.apriori_code, id: item.apriori_id } }" class="text-uppercase h4 ls-1 text-primary py-3 mb-0">Paket Rekomendasi</router-link>
                    </div>
                    <div class="card-body mx-auto">
                      <div class="display-2">{{ item.apriori_discount }}%</div>
                      <span class="text-muted h2" style="text-decoration: line-through">Rp {{ item.product_total_price }}</span>
                      /
                      <span class="text-muted">Rp {{ item.price_discount }}</span>
                      <ul class="list-unstyled my-4">
                        <li v-for="(value,i) in item.apriori_item.split(', ')" :key="i">
                          <div class="d-flex align-items-center">
                            <div>
                              <div class="icon icon-xs icon-shape bg-gradient-primary text-white shadow rounded-circle">
                                <i class="ni ni-basket"></i>
                              </div>
                            </div>
                            <div>
                              <span class="pl-2 text-sm">{{ UpperWord(value) }}</span>
                            </div>
                          </div>
                        </li>
                      </ul>
                    </div>
                    <div class="card-footer">
                      <h3>Rp. {{ item.price_discount }}</h3>
                    </div>
                  </div>
                </div>
                <p v-if="recommendation.length === 0" class="mx-auto">No Recommendation Found</p>
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
  mounted() {
    this.fetchData()
    this.fetchDataRecommendation()
    document.getElementsByTagName("body")[0].classList.remove("bg-default");
  },
  data: function () {
    return {
      product: [],
      recommendation: [],
      carts: [],
      totalCart: 0,
      quantity: 0,
      isLoading: true,
      isLoading2: true,
    };
  },
  methods: {
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
    submit(no_product) {
      axios.delete(`${process.env.VUE_APP_SERVICE_URL}/products/` + no_product, { headers: authHeader() })
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              this.$router.push({
                name: 'product'
              })
            }
          }).catch(error => {
        console.log(error.response.data.status)
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