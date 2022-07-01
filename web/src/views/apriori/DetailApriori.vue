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
              <h3 class="mb-0">Detail Apriori</h3>
            </div>
            <div class="row align-items-center">
              <div class="col-12 col-lg-6 text-center">
                <img :src="getImage()" class="img-fluid my-5" width="500">
              </div>
              <div class="col-12 col-lg-6">
                <div class="card-pricing border-0 text-center mb-4">
                  <div class="card-body px-lg-7">
                    <div class="display-2">{{ apriori.discount }}%</div>
                    <span class="text-muted">discount</span>
                    <ul class="list-unstyled my-4" v-if="apriori.item">
                      <li v-for="(value,i) in apriori.item.split(', ')" :key="i">
                        <div class="d-flex align-items-center justify-content-center">
                          <div>
                            <div class="icon icon-xs icon-shape bg-gradient-primary text-white shadow rounded-circle">
                              <i class="fas fa-terminal"></i>
                            </div>
                          </div>
                          <div>
                            <span class="pl-2 text-sm">{{ UpperWord(value) }}</span>
                          </div>
                        </div>
                      </li>
                    </ul>
                    <p>{{ apriori.description }}</p>
                  </div>
                  <div class="card-footer">
                    <router-link :to="{ name: 'apriori.detail', params: { code: this.$route.params.code } }" class="text-muted">
                      {{ apriori.code }}</router-link>
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
      axios.get(`http://localhost:3000/api/apriori/${this.$route.params.code}/detail/${this.$route.params.id}`, { headers: authHeader() }).then((response) => {
        this.apriori = response.data.data;
      });
    },
    getImage() {
      return this.apriori.image
    },
    UpperWord(str) {
      return str.toLowerCase().replace(/\b[a-z]/g, function (letter) {
        return letter.toUpperCase();
      })
    }
  }
}
</script>