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
              <p class="mt-2 text-center">Loading...</p>
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
                  <td>{{ item.is_active == 0 ? "Non Active" : "Active" }}</td>
                  <td>{{ item.created_at }}</td>
                  <td class="text-center">
                    <router-link :to="{ name: 'apriori.detail', params: { code: item.code } }" class="btn btn-secondary btn-sm">Detail</router-link>
                    <form @submit.prevent="activate(item.code)" method="POST" class="d-inline mr-1">
                      <button class="btn btn-primary btn-sm">{{ item.is_active == 0 ? "Activate" : "Deactivate" }}</button>
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
          console.log(error.response.data.status)
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
        console.log(error.response.data.status)
      })
    }
  }
}
</script>