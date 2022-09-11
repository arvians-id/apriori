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
      <!-- Table -->
      <div class="row">
        <div class="col">
          <div class="card">
            <!-- Card header -->
            <div class="card-header">
              <h3 class="mb-0">Data Apriori</h3>
            </div>
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
            <div class="table-responsive py-4" v-else>
              <table class="table table-flush" id="datatable">
                <thead class="thead-light">
                <tr>
                  <th>No</th>
                  <th>Code</th>
                  <th>Date Range</th>
                  <th>Status</th>
                  <th>Created At</th>
                  <th class="text-center">Action</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(item,i) in apriories" :key="item.id_apriori">
                  <td>{{ (i++) + 1 }}</td>
                  <td>{{ item.code }}</td>
                  <td>{{ item.range_date }}</td>
                  <td>{{ item.is_active === false ? "Non Active" : "Active" }}</td>
                  <td>{{ item.created_at }}</td>
                  <td class="text-center">
                    <router-link :to="{ name: 'apriori.detail', params: { code: item.code } }" class="btn btn-secondary btn-sm">Detail</router-link>
                    <form @submit.prevent="activate(item.code)" method="POST" class="d-inline mr-1">
                      <button class="btn btn-primary btn-sm">{{ item.is_active === false ? "Activate" : "Deactivate" }}</button>
                    </form>
                    <form @submit.prevent="submit(item.code)" method="POST" class="d-inline">
                      <button class="btn btn-danger btn-sm">Delete</button>
                    </form>
                  </td>
                </tr>
                </tbody>
              </table>
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
import axios from "axios";
import $ from "jquery";
import Sidebar from "@/components/admin/Sidebar.vue"
import Topbar from "@/components/admin/Topbar.vue"
import Header from "@/components/admin/Header.vue"
import Footer from "@/components/admin/Footer.vue"
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
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/apriori`, { headers: authHeader() }).then((response) => {
        this.apriories = response.data.data;
        setTimeout(function(){
          $('#datatable').DataTable();
        }, 0);
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });

      this.isLoading = false;
    },
    submit(code) {
      if(confirm("Are you sure to delete this data?")) {
        axios.delete(`${process.env.VUE_APP_SERVICE_URL}/apriori/` + code, { headers: authHeader() })
            .then(response => {
              if(response.data.code === 200) {
                alert(response.data.status)
                this.fetchData()
              }
            }).catch(error => {
              if (error.response.status === 400 || error.response.status === 404) {
                alert(error.response.data.status)
              }
        })
      }
    },
    activate(code) {
      axios.patch(`${process.env.VUE_APP_SERVICE_URL}/apriori/` + code, null,{ headers: authHeader() })
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              this.fetchData()
            }
          }).catch(error => {
            if (error.response.status === 400 || error.response.status === 404) {
              alert(error.response.data.status)
            }
      })
    }
  }
}
</script>