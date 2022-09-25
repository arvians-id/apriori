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
              <div class="card-header">
                <h3 class="mb-0">Edit Transaction</h3>
              </div>
              <!-- Card body -->
              <div class="card-body" v-if="isLoading">
                <p class="mt-2 text-center">Loading...</p>
              </div>
              <div class="card-body" v-else>
                <form @submit.prevent="submit" method="POST">
                  <div class="form-group">
                    <label class="form-control-label">Product Name</label> <small class="text-danger">*use ctrl for selecting the product</small>
                    <select multiple class="form-control" v-model="transaction.product_name">
                      <option v-for="(item) in products" :key="item.id_product" :selected="transaction.product_name.includes(item.name.toLowerCase())">{{ item.name.toLowerCase() }}</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Customer Name</label> <small class="text-danger">*</small>
                    <input type="text" class="form-control" v-model="transaction.customer_name">
                  </div>
                  <button class="btn btn-primary" type="submit">Submit form</button>
                </form>
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
import gql from 'graphql-tag'

export default {
  apollo: {
    products: gql`query {
      products: ProductFindAllByAdmin {
        id_product
        name
        created_at
        updated_at
      }
    }`,
  },
  components: {
    Footer,
    Sidebar,
    Header,
    Topbar
  },
  async mounted() {
    await this.fetchData()

    this.isLoading = false;
  },
  data: function () {
    return {
      products: [],
      transaction: {
        product_name: [],
        customer_name: ""
      },
      isLoading: true
    };
  },
  methods: {
    submit() {
      if (this.transaction.product_name.length > 0) {
        let productName = this.transaction.product_name
        this.transaction.product_name = productName.join(", ")
      }

      axios.patch(`${process.env.VUE_APP_SERVICE_URL}/transactions/${this.$route.params.no_transaction}`, this.transaction, { headers: authHeader() })
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              this.$router.push({
                name: 'transaction'
              })
            }
          }).catch(error => {
            if (error.response.status === 400 || error.response.status === 404) {
              alert(error.response.data.status)
            }
          })
    },
    async fetchData() {
      await axios.get(`${process.env.VUE_APP_SERVICE_URL}/transactions/${this.$route.params.no_transaction}`, { headers: authHeader() }).then(response => {
        let productName = response.data.data.product_name
        this.transaction = {
          product_name: productName.split(", "),
          customer_name: response.data.data.customer_name,
        };
      }).catch(error => {
        if (error.response.status === 400 || error.response.status === 404) {
          console.log(error.response.data.status)
        }
      });
    }
  }
}
</script>