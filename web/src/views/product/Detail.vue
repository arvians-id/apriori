<template>
  <!-- Sidenav -->
  <Sidebar />
  <!-- Main content -->
  <div class="main-content" id="panel">
    <!-- Topnav -->
    <Topbar />
    <!-- Header -->
    <Header />
    <!-- Page content -->
    <div class="container-fluid mt--6">
      <div class="row">
        <div class="col-xl-12 order-xl-2">
          <div class="card card-profile">
            <!-- Card header -->
            <div class="card-header">
              <h3 class="mb-0">Detail Product</h3>
            </div>
            <div class="row align-items-center">
              <div class="col-12 col-lg-6 text-center">
                <img :src="getImage()" class="img-fluid my-5" width="500">
              </div>
              <div class="col-12 col-lg-6">
                <div class="card-header text-center border-0 pt-8 pt-md-4 pb-0 pb-md-4">
                  <div class="d-flex justify-content-between">
                    <router-link :to="{ name: 'product.edit', params: { code: this.$route.params.code } }"  class="btn btn-sm btn-default float-right">Edit</router-link>
                    <form @submit.prevent="submit(this.$route.params.code)" method="POST" class="d-inline">
                      <button class="btn btn-danger btn-sm">Delete</button>
                    </form>
                  </div>
                </div>
                <div class="card-body pt-0">
                  <div class="row">
                    <div class="col">
                      <div class="card-profile-stats d-flex justify-content-center">
                        <div>
                          <span class="heading">Created At</span>
                          <span class="description">{{ product.created_at }}</span>
                        </div>
                        <div>
                          <span class="heading">Last Modified</span>
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
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="card">
            <!-- Card header -->
            <div class="card-header d-flex justify-content-between">
              <h3 class="mb-0">Recommendation Packages For You</h3>
            </div>
            <!-- Card body -->
            <div class="card-body">
              <div class="row">
                <div class="col-12 col-md-6 col-lg-4 col-xl-3" v-for="item in recommendation" :key="item.apriori_id">
                  <div class="card card-pricing border-0 text-center mb-4">
                    <div class="card-header bg-transparent">
                      <h4 class="text-uppercase ls-1 text-primary py-3 mb-0">Recommendation pack</h4>
                    </div>
                    <div class="card-body mx-auto">
                      <div class="display-2">{{ item.apriori_discount }}%</div>
                      <span class="text-muted h2" style="text-decoration: line-through">Rp. {{ item.product_total_price }}</span>
                      /
                      <span class="text-muted">Rp. {{ item.price_discount }}</span>
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
import Sidebar from "@/components/Sidebar.vue"
import Topbar from "@/components/Topbar.vue"
import Header from "@/components/Header.vue"
import Footer from "@/components/Footer.vue"
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
  },
  data: function () {
    return {
      product: [],
      recommendation: []
    };
  },
  methods: {
    fetchData() {
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/products/${this.$route.params.code}`, { headers: authHeader() }).then((response) => {
        this.product = response.data.data;
      });
    },
    fetchDataRecommendation() {
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/products/${this.$route.params.code}/recommendation`, { headers: authHeader() }).then((response) => {
        if(response.data.data != null) {
          this.recommendation = response.data.data;
        }
      });
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
    }
  }
}
</script>