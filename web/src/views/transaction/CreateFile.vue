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
                <h3 class="mb-0">Create Transaction By CSV File</h3>
              </div>
              <!-- Card body -->
              <div class="card-body">
                <form @submit.prevent="submit" method="POST">
                <label class="form-control-label">Upload</label>  <small class="text-danger">*</small>
                  <div class="custom-file mb-3">
                    <input type="file" class="custom-file-input" @change="submitFile" required>
                    <label class="custom-file-label">Select file</label>
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

export default {
    components: {
        Footer,
        Sidebar,
        Header,
        Topbar
    },
    data: function () {
      return {
        products: [],
        transaction: {
          file: null
        }
      };
    },
    methods: {
      submit() {
        const config = {
          headers: {
            'Content-Type': 'multipart/form-data',
            ...authHeader(),
          }
        }
        const formData = new FormData()
        formData.append("file", this.transaction.file)

        axios.post(`${process.env.VUE_APP_SERVICE_URL}/transactions/csv`, formData, config)
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
      submitFile(e){
        this.transaction.file = e.target.files[0]
      }
    }
}
</script>