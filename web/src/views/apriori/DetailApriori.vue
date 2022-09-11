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
              <h3 class="mb-0">Detail Recommendation</h3>
            </div>
            <div class="card-body">
              <div v-if="isLoading">
                <div class="loading-skeleton row d-flex justify-content-center">
                  <div class="col-12 col-lg-4 mb-2">
                    <img src="https://my-apriori-bucket.s3.ap-southeast-1.amazonaws.com/assets/no-image.png" class="img-fluid mb-2">
                    <p class="pt-4 mb-0">this is button</p>
                  </div>
                  <div class="col-12 col-lg-6">
                    <div class="text-left">
                      <h5 class="h2 p-0 m-0 mb-1">This is title of product</h5>
                      <p class="p-0 m-0 mb-1 w-50">This is title of product</p>
                      <p class="w-50">Pricing</p>
                      <hr class="mt-0 mb-1">
                      <div class="nav-wrapper">
                        <p class="pt-4 mb-0">this is button</p>
                      </div>
                      <div class="card shadow">
                        <div class="card-body">
                          <div class="tab-content" id="skeleton-myTabContent">
                            <div class="tab-pane fade show active" id="skeleton-detail" role="tabpanel" aria-labelledby="detail-tab">
                              <p class="font-weight-bold mb-0 mt-2">Kondisi</p>
                              <p>Original baru</p>
                              <p class="font-weight-bold mb-0 mt-2">Kategori</p>
                              <p class="font-weight-bold">Text</p>
                              <p class="font-weight-bold mb-0 mt-2">Berat Satuan</p>
                              <p>gram</p>
                              <p class="font-weight-bold mb-0 mt-2">Deskripsi</p>
                              <p>Description</p>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div class="row d-flex justify-content-center" v-else>
                <div class="col-12 col-lg-4 text-center mb-2">
                  <img :src="apriori.apriori_image" class="img-fluid mb-2" width="500">
                  <router-link :to="{ name: 'apriori.edit', params: { code: apriori.apriori_code, id: this.$route.params.id } }"  class="btn btn-primary btn-block mb-2">Edit</router-link>
                </div>
                <div class="col-12 col-lg-6">
                  <div class="text-left">
                    <h5 class="h2 text-uppercase p-0 m-0">Paket rekomendasi {{ apriori.apriori_item }}</h5>
                    <p class="p-0 m-0">code : {{ apriori.apriori_code }}</p>
                    <div class="h1 font-weight-bold">Rp. {{ apriori.price_discount }}</div>
                    <small class="p-1 d-inline rounded font-weight-bold" style="background-color: #ffeaef;; color: #ff5c84">{{ apriori.apriori_discount }}%</small>
                    <p class="d-inline ml-2" style="text-decoration: line-through">Rp. {{ apriori.product_total_price }}</p>
                    <hr class="mt-0 mb-1">
                    <div class="nav-wrapper">
                      <ul class="nav nav-pills nav-fill flex-column flex-md-row" id="tabs-icons-text" role="tablist">
                        <li class="nav-item">
                          <a class="nav-link mb-sm-3 mb-md-0 active" id="detail-tab" data-toggle="tab" href="#detail" role="tab" aria-controls="detail" aria-selected="true"><i class="ni ni-tag mr-2"></i>Detail</a>
                        </li>
                      </ul>
                    </div>
                    <div class="card shadow">
                      <div class="card-body">
                        <div class="tab-content" id="myTabContent">
                          <div class="tab-pane fade show active" id="detail" role="tabpanel" aria-labelledby="detail-tab">
                            <p class="font-weight-bold mb-0 mt-2">Kondisi</p>
                            <p>Original baru</p>
                            <p class="font-weight-bold mb-0 mt-2">Jumlah Barang</p>
                            <p>{{ apriori.apriori_item.split(", ").length }}</p>
                            <p class="font-weight-bold mb-0 mt-2">Nama Barang</p>
                            <p>{{ UpperWord(apriori.apriori_item) }}</p>
                            <p class="font-weight-bold mb-0 mt-2">Berat Barang</p>
                            <p>{{ apriori.mass == undefined ? 0 : apriori.mass }} gram</p>
                            <p class="font-weight-bold mb-0 mt-2">Deskripsi</p>
                            <p>{{ apriori.apriori_description == undefined ? "Tidak ada deskripsi" : apriori.apriori_description }}</p>
                          </div>
                        </div>
                      </div>
                    </div>
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
  },
  data: function () {
    return {
      apriori: [],
      isLoading: true
    };
  },
  methods: {
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/apriori/${this.$route.params.code}/detail/${this.$route.params.id}`, { headers: authHeader() }).then((response) => {
        this.apriori = response.data.data;
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });

      this.isLoading = false;
    },
    UpperWord(str) {
      return str.toLowerCase().replace(/\b[a-z]/g, function (letter) {
        return letter.toUpperCase();
      })
    }
  }
}
</script>