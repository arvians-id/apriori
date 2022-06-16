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
                     <label class="form-control-label">Minimum Support</label> <small class="text-danger">*</small>
                     <div class="input-group input-group-merge">
                       <input type="number" step="0.1" min="0.1" class="form-control" v-model="minimum_support" placeholder="example: 0.3 or 20" required>
                       <div class="input-group-append">
                         <span class="input-group-text">%</span>
                       </div>
                     </div>
                   </div>
                   <div class="form-group">
                     <label class="form-control-label">Minimum Confidence</label> <small class="text-danger">*</small>
                     <div class="input-group input-group-merge">
                       <input type="number" step="0.1" min="0.1" class="form-control" v-model="maximum_confidence" placeholder="example: 0.3 or 20" required>
                       <div class="input-group-append">
                         <span class="input-group-text">%</span>
                       </div>
                     </div>
                   </div>
                   <div class="form-group">
                     <label class="form-control-label">Minimum Discount</label> <small class="text-danger">*</small>
                     <div class="input-group input-group-merge">
                       <input type="number" min="1" class="form-control" v-model="minimum_discount" placeholder="example: 10" required>
                       <div class="input-group-append">
                         <span class="input-group-text">%</span>
                       </div>
                     </div>
                   </div>
                   <div class="form-group">
                     <label class="form-control-label">Maximum Discount</label> <small class="text-danger">*</small>
                     <div class="input-group input-group-merge">
                       <input type="number" min="1" max="100" class="form-control" v-model="maximum_discount" placeholder="example: 15" required>
                       <div class="input-group-append">
                         <span class="input-group-text">%</span>
                       </div>
                     </div>
                   </div>
                   <div class="form-group">
                     <label class="form-control-label">Start Date</label> <small class="text-danger">*</small>
                     <input type="date" class="form-control" v-model="start_date" required>
                   </div>
                   <div class="form-group">
                     <label class="form-control-label">End Date</label> <small class="text-danger">*</small>
                     <input type="date" class="form-control" v-model="end_date" required>
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