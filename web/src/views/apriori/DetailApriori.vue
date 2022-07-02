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
              <h3 class="mb-0">Recommendation Packages</h3>
            </div>
            <div class="row align-items-center">
              <div class="col-12 col-lg-6 text-center">
                <img :src="getImage()" class="img-fluid my-5" width="500">
              </div>
              <div class="col-12 col-lg-6">
                <div class="card-header text-center border-0 pt-8 pt-md-4 pb-0 pb-md-4">
                  <div class="d-flex justify-content-between">
                    <router-link :to="{ name: 'apriori.edit', params: { code: this.$route.params.code, id: this.$route.params.id } }"  class="btn btn-sm btn-default float-right">Edit</router-link>
                  </div>
                </div>
                <div class="card-pricing border-0 text-center mb-4">
                  <div class="card-body px-lg-7">
                    <div class="display-2">{{ apriori.apriori_discount }}%</div>
                    <span class="text-muted h2" style="text-decoration: line-through">Rp. {{ apriori.product_total_price }}</span>
                    /
                    <span class="text-muted">Rp. {{ apriori.price_discount }}</span>
                    <ul class="list-unstyled my-4" v-if="apriori.apriori_item">
                      <li v-for="(value,i) in apriori.apriori_item.split(', ')" :key="i">
                        <div class="d-flex align-items-center justify-content-center">
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
                    <p>{{ apriori.apriori_description }}</p>
                  </div>
                  <div class="card-footer">
                    <router-link :to="{ name: 'apriori.detail', params: { code: this.$route.params.code } }" class="text-muted">
                      {{ apriori.apriori_code }}</router-link>
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
  },
  data: function () {
    return {
      apriori: [],
    };
  },
  methods: {
    fetchData() {
      axios.get(`${process.env.VUE_APP_SERVICE_URL}/apriori/${this.$route.params.code}/detail/${this.$route.params.id}`, { headers: authHeader() }).then((response) => {
        this.apriori = response.data.data;
      });
    },
    getImage() {
      return this.apriori.apriori_image
    },
    UpperWord(str) {
      return str.toLowerCase().replace(/\b[a-z]/g, function (letter) {
        return letter.toUpperCase();
      })
    }
  }
}
</script>