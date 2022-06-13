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
                <h3 class="mb-0">Buat Pengguna</h3>
              </div>
              <!-- Card body -->
              <div class="card-body">
                 <form>
                  <div class="form-group">
                    <label class="form-control-label">Nama Lengkap</label>
                    <input type="text" class="form-control" v-model="name">
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Email</label>
                    <input type="text" class="form-control" v-model="email">
                  </div>
                   <div class="form-group">
                     <label class="form-control-label">Password</label>
                     <input type="password" class="form-control" v-model="password">
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
import Sidebar from "@/components/Sidebar.vue"
import Topbar from "@/components/Topbar.vue"
import Header from "@/components/Header.vue"
import Footer from "@/components/Footer.vue"

export default {
  components: {
    Footer,
    Sidebar,
    Header,
    Topbar
  }
}
</script>