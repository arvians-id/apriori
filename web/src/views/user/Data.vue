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
              <h3 class="mb-0">Data User</h3>
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
                  <th>Full Name</th>
                  <th>Email</th>
                  <th>Created At</th>
                  <th>Last Modified</th>
                  <th class="text-center">Action</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(item,i) in users" :key="item.id_user">
                  <td>{{ (i++) + 1 }}</td>
                  <td>{{ item.name }}</td>
                  <td>{{ item.email }}</td>
                  <td>{{ item.created_at }}</td>
                  <td>{{ item.updated_at }}</td>
                  <td class="text-center">
                    <router-link :to="{ name: 'user.edit', params: { id: item.id_user } }" class="btn btn-primary btn-sm">Edit</router-link>
                    <form @submit.prevent="submit(item.id_user)" method="POST" class="d-inline">
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
      users: [],
      isLoading: true
    };
  },
  methods: {
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/users`, { headers: authHeader() }).then((response) => {
        this.users = response.data.data;
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
    submit(id) {
      if(confirm("Are you sure to delete this data?")) {
        axios.delete(`${process.env.VUE_APP_SERVICE_URL}/users/` + id, { headers: authHeader() })
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
}
</script>