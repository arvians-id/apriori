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
                <p class="mt-2 text-center">Loading...</p>
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
        this.apriories = response.data.data;
      });

      this.isLoading = false;
    }
  }
}
</script>