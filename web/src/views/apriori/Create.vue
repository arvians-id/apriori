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
                <h3 class="mb-0">Generate Apriori</h3>
              </div>
              <!-- Card body -->
              <div class="card-body">
                 <form>
                  <div class="form-group">
                    <label class="form-control-label">Minimum Pendukung (Support)</label> <small class="text-danger">*%</small>
                    <input type="text" class="form-control" v-model="minimum_support" placeholder="contoh: 0.3 atau 20">
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Maximum Kepercayaan (Confidence)</label> <small class="text-danger">*%</small>
                    <input type="text" class="form-control" v-model="maximum_confidence" placeholder="contoh: 0.3 atau 20">
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Minimum Diskon</label> <small class="text-danger">*%</small>
                    <input type="text" class="form-control" v-model="minimum_discount" placeholder="contoh: 10">
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Maximum Diskon</label> <small class="text-danger">*%</small>
                    <input type="text" class="form-control" v-model="maximum_discount" placeholder="contoh: 15">
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Dimulai</label>
                    <input type="date" class="form-control" v-model="start_date">
                  </div>
                  <div class="form-group">
                    <label class="form-control-label">Hingga</label>
                    <input type="date" class="form-control" v-model="end_date">
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