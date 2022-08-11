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
            <div class="card-body">
              <div class="row align-items-center mx-auto" v-if="isLoading">
                <p class="p-3 mt-2 text-center">Loading...</p>
              </div>
              <div class="row d-flex justify-content-center" v-else>
                <div class="col-12 col-lg-4 text-center mb-2">
                  <img :src="getImage()" class="img-fluid mb-2" width="500">
                  <router-link :to="{ name: 'product.edit', params: { code: this.$route.params.code } }"  class="btn btn-primary btn-block mb-2">Edit</router-link>
                  <form @submit.prevent="submit(this.$route.params.code)" method="POST" class="d-inline">
                    <button class="btn btn-danger btn-block">Delete</button>
                  </form>
                </div>
                <div class="col-12 col-lg-4">
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
                        <p>{{ product.category }}</p>
                        <p class="font-weight-bold mb-0 mt-2">Berat Satuan</p>
                        <p>{{ product.mass }} gram</p>
                        <p class="font-weight-bold mb-0 mt-2">Deskripsi</p>
                        <p>{{ product.description }}</p>
                      </div>
                      <div class="tab-pane fade" id="other-detail" role="tabpanel" aria-labelledby="other-detail-tab">
                        <p class="font-weight-bold mb-0 mt-2">Tanggal Dibuat</p>
                        <p>{{ product.created_at }}</p>
                        <p class="font-weight-bold mb-0 mt-2">Terakhir Diubah</p>
                        <p>{{ product.updated_at }}</p>
                        <p class="font-weight-bold mb-0 mt-2">Status Produk</p>
                        <p>{{ product.is_empty == 0 ? "Produk aktif" : "Produk tidak aktif" }}</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="card">
            <!-- Card header -->
            <div class="card-header d-flex justify-content-between">
              <h3 class="mb-0">Recommendation Packages</h3>
            </div>
            <!-- Card body -->
            <div class="card-body" v-if="isLoading2">
              <p class="mt-2 text-center">Loading...</p>
            </div>
            <div class="card-body" v-else>
              <div class="row">
                <div class="col-12 col-md-6 col-lg-4 col-xl-3" v-for="item in recommendation" :key="item.apriori_id">
                  <div class="card card-pricing border-0 text-center mb-4">
                    <div class="card-header bg-transparent">
                      <router-link :to="{ name: 'apriori.code-detail', params: { code: item.apriori_code, id: item.apriori_id } }" class="text-uppercase h4 ls-1 text-primary py-3 mb-0">Recommendation pack</router-link>
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
import Sidebar from "@/components/admin/Sidebar.vue"
import Topbar from "@/components/admin/Topbar.vue"
import Header from "@/components/admin/Header.vue"
import Footer from "@/components/admin/Footer.vue"
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
      recommendation: [],
      isLoading: true,
      isLoading2: true,
    };
  },
  methods: {
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/products/${this.$route.params.code}`, { headers: authHeader() }).then((response) => {
        this.product = response.data.data;
      });

      this.isLoading = false
    },
    async fetchDataRecommendation() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/products/${this.$route.params.code}/recommendation`, { headers: authHeader() }).then((response) => {
        if(response.data.data != null) {
          this.recommendation = response.data.data;
        }
      }).catch((error) => {
        console.log(error)
      })

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
      if(confirm("Are you sure to delete this product?")) {
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
}
</script>