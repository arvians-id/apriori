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
            <div class="card-header d-flex justify-content-between">
              <h3 class="mb-0">Data Transaction</h3>
              <form @submit.prevent="truncate()" method="POST">
                <button class="btn btn-danger btn-sm" type="submit">Clear Data</button>
              </form>
            </div>
            <div class="table-responsive py-3" v-if="isLoading">
              <p class="mt-2 text-center">Loading...</p>
            </div>
            <div class="table-responsive py-4" v-else>
              <table class="table table-flush" id="datatable">
                <thead class="thead-light">
                  <tr>
                    <th>No</th>
                    <th>No Transaction</th>
                    <th>Customer Name</th>
                    <th>Product Name</th>
                    <th>Created At</th>
                    <th>Last Modified</th>
                    <th class="text-center">Action</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(item,i) in transactions" :key="item.id_transaction">
                    <td>{{ (i++) + 1 }}</td>
                    <td>{{ item.no_transaction }}</td>
                    <td>{{ item.customer_name }}</td>
                    <td>{{ item.product_name }}</td>
                    <td>{{ item.created_at }}</td>
                    <td>{{ item.updated_at }}</td>
                    <td class="text-center">
                      <router-link :to="{ name: 'transaction.edit', params: { no_transaction: item.no_transaction } }" class="btn btn-primary btn-sm">Edit</router-link>
                      <form @submit.prevent="submit(item.no_transaction)" method="POST" class="d-inline">
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
import "jquery/dist/jquery.min.js";
import "datatables.net-dt/js/dataTables.dataTables";
import "datatables.net-dt/css/jquery.dataTables.min.css";
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
      transactions: [],
      isLoading: true
    };
  },
  methods: {
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/transactions`, { headers: authHeader() }).then((response) => {
        this.transactions = response.data.data;
        setTimeout(function(){
          $('#datatable').DataTable();
        }, 0);
      });

      this.isLoading = false;
    },
    submit(no_transaction) {
      if(confirm("Are you sure to delete this data?")) {
        axios.delete(`${process.env.VUE_APP_SERVICE_URL}/transactions/` + no_transaction, { headers: authHeader() })
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
    truncate() {
      axios.delete(`${process.env.VUE_APP_SERVICE_URL}/transactions/truncate`, { headers: authHeader() })
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