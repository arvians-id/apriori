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
              <div class="card-header border-0">
                <h3 class="mb-0">{{ this.$route.params.code }}</h3>
              </div>
              <!-- Light table -->
              <div class="table-responsive py-3" v-if="isLoading">
                <div class="loading-skeleton">
                  <div class="d-flex justify-content-between mx-4 mt-3">
                    <p>Show Entries Data</p>
                    <p>Search Data In Table</p>
                  </div>
                  <hr class="mt-1 mb-4">
                  <div class="mx-4">
                    <p>Column of data</p>
                    <p>Column of data</p>
                    <p>Column of data</p>
                    <p>Column of data</p>
                    <p>Column of data</p>
                    <p>Column of data</p>
                    <p>Column of data</p>
                    <p>Column of data</p>
                    <p>Column of data</p>
                    <p>Column of data</p>
                  </div>
                  <hr class="mt-4 mb-1">
                  <div class="d-flex justify-content-between mx-4 mt-4">
                    <p>Show Entries Data</p>
                    <div class="d-none d-lg-block">
                      <p class="d-inline mx-2">Search</p>
                      <p class="d-inline mx-2">Search Data In Table</p>
                      <p class="d-inline mx-2">Search</p>
                    </div>
                  </div>
                </div>
              </div>
              <div class="table-responsive" v-else>
                <table class="table align-items-center table-flush">
                  <thead class="thead-light">
                    <tr>
                      <th class="text-center">Item Set</th>
                      <th class="text-center">Iterate</th>
                      <th class="text-center">Discount</th>
                      <th class="text-center">Support</th>
                      <th class="text-center">Confidence</th>
                      <th class="text-center">Range Date</th>
                      <th class="text-center">Action</th>
                    </tr>
                  </thead>
                  <tbody class="list">
                    <tr v-for="item in apriories" :key="item.id_apriori">
                      <td>{{ item.item }}</td>
                      <td class="text-center">{{ item.item.split(",").length }}</td>
                      <td class="text-center">
                        <div class="d-flex align-items-center">
                          <span class="completion mr-2">{{ item.discount }}%</span>
                          <div>
                            <div class="progress">
                              <div class="progress-bar bg-warning" role="progressbar" :aria-valuenow="item.discount" aria-valuemin="0" aria-valuemax="100" :style="`width: ${item.discount}%;`"></div>
                            </div>
                          </div>
                        </div>
                      </td>
                      <td class="text-center">
                        <div class="d-flex align-items-center">
                          <span class="completion mr-2">{{ item.support }}%</span>
                          <div>
                            <div class="progress">
                              <div class="progress-bar bg-warning" role="progressbar" :aria-valuenow="item.support" aria-valuemin="0" aria-valuemax="100" :style="`width: ${item.support}%;`"></div>
                            </div>
                          </div>
                        </div>
                      </td>
                      <td class="text-center">
                        <div class="d-flex align-items-center">
                          <span class="completion mr-2">{{ item.confidence }}%</span>
                          <div>
                            <div class="progress">
                              <div class="progress-bar bg-warning" role="progressbar" :aria-valuenow="item.confidence" aria-valuemin="0" aria-valuemax="100" :style="`width: ${item.confidence}%;`"></div>
                            </div>
                          </div>
                        </div>
                      </td>
                      <td class="text-center">{{ item.range_date }}</td>
                      <td class="text-center">
                        <router-link :to="{ name: 'apriori.code-detail', params: { code: item.code, id: item.id_apriori } }" class="btn btn-secondary btn-sm">Detail</router-link>
                        <router-link :to="{ name: 'apriori.edit', params: { code: item.code, id: item.id_apriori } }" class="btn btn-primary btn-sm">Edit</router-link>
                      </td>
                    </tr>
                  </tbody>
                </table>
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
@import '../../assets/skeleton.css';
</style>

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
      apriories: [],
      isLoading: true
    };
  },
  methods: {
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/apriori/${this.$route.params.code}`, { headers: authHeader() }).then((response) => {
        if (response.data.data !== undefined) {
          this.apriories = response.data.data;
        }
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });

      this.isLoading = false;
    }
  }
}
</script>