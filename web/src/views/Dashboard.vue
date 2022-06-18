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
                <h3 class="mb-0">List of Recommendation Packages </h3>
                <h3 class="mb-0">{{ getDate }}</h3>
              </div>
              <!-- Card body -->
              <div class="card-body">
                <div class="row">
                  <div class="col-12 col-md-6 col-lg-4 col-xl-3" v-for="item in apriories" :key="item.id_apriori">
                    <div class="card card-pricing border-0 text-center mb-4">
                      <div class="card-header bg-transparent">
                        <h4 class="text-uppercase ls-1 text-primary py-3 mb-0">Recommendation pack</h4>
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
                        <router-link :to="{ name: 'apriori.detail', params: { code: item.code } }" class=" text-muted">{{ item.code }}</router-link>
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
  created() {
    this.fetchData()
  },
  data: function () {
    return {
      apriories: [],
      getDate: "No date found"
    };
  },
  methods: {
    fetchData() {
      axios.get("http://localhost:3000/api/apriori/actives", { headers: authHeader() }).then((response) => {
        this.apriories = response.data.data;
        this.getDate = `${this.apriories[0].range_date}`
      });
    },
    UpperWord(str) {
      return str.toLowerCase().replace(/\b[a-z]/g, function (letter) {
        return letter.toUpperCase();
      })
    }
  }
}
</script>