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
              <div class="table-responsive">
                <table class="table align-items-center table-flush">
                  <thead class="thead-light">
                    <tr>
                      <th class="text-center">Item Set</th>
                      <th class="text-center">Iterate</th>
                      <th class="text-center">Discount</th>
                      <th class="text-center">Support</th>
                      <th class="text-center">Confidence</th>
                      <th class="text-center">Range Date</th>
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
import Sidebar from "@/components/Sidebar.vue"
import Topbar from "@/components/Topbar.vue"
import Header from "@/components/Header.vue"
import Footer from "@/components/Footer.vue"
import axios from "axios";

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
    };
  },
  methods: {
    fetchData() {
      axios.get(`http://localhost:3000/api/apriori/${this.$route.params.code}`).then((response) => {
        this.apriories = response.data.data;
      });
    }
  }
}
</script>