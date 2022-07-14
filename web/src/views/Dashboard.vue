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
        <div class="col-12">
          <div class="card-wrapper">
            <!-- Custom form validation -->
            <div class="card">
              <!-- Card header -->
              <div class="card-header d-flex justify-content-between">
                <h3 class="mb-0">All Recommendation Packages </h3>
              </div>
              <!-- Card body -->
              <div class="card-body" v-if="isLoading">
                <p class="mt-2 text-center">Loading...</p>
              </div>
              <div class="card-body" v-else>
                <div class="row">
                  <div class="col-12 col-md-6 col-lg-4 col-xl-3" v-for="item in apriories" :key="item.id_apriori">
                    <div class="card card-pricing border-0 text-center mb-4">
                      <div class="card-header bg-transparent">

                        <h4 class="text-uppercase ls-1 text-primary py-3 mb-0">
                          <router-link :to="{ name: 'apriori.detail', params: { code: item.code } }">Recommendation pack</router-link>
                        </h4>
                      </div>
                      <div class="card-body mx-auto">
                        <div class="display-2">{{ item.discount }}%</div>
                        <span class="text-muted">discount</span>
                        <ul class="list-unstyled my-4">
                          <li v-for="(value,i) in item.item.split(', ')" :key="i">
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
                        {{ getDate }}
                      </div>
                    </div>
                  </div>
                  <p v-if="apriories.length == 0" class="mx-auto">No Recommendation Found</p>
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
    document.getElementsByTagName("body")[0].classList.remove("bg-default");
  },
  data: function () {
    return {
      apriories: [],
      getDate: "No date found",
      isLoading: true,
    };
  },
  methods: {
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/apriori/actives`, { headers: authHeader() }).then((response) => {
        this.apriories = response.data.data;
        this.getDate = `${this.apriories[0].range_date}`
      });

      this.isLoading = false
    },
    UpperWord(str) {
      return str.toLowerCase().replace(/\b[a-z]/g, function (letter) {
        return letter.toUpperCase();
      })
    }
  }
}
</script>