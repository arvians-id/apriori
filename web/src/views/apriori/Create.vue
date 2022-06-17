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
                <form @submit.prevent="submit" method="POST">
                   <div class="form-group">
                     <label class="form-control-label">Minimum Support</label> <small class="text-danger">*</small>
                     <div class="input-group input-group-merge">
                       <input type="number" step="0.1" min="0.1" class="form-control" v-model="apriori.minimum_support" placeholder="example: 0.3 or 20" required>
                       <div class="input-group-append">
                         <span class="input-group-text">%</span>
                       </div>
                     </div>
                   </div>
                   <div class="form-group">
                     <label class="form-control-label">Minimum Confidence</label> <small class="text-danger">*</small>
                     <div class="input-group input-group-merge">
                       <input type="number" step="0.1" min="0.1" class="form-control" v-model="apriori.minimum_confidence" placeholder="example: 0.3 or 20" required>
                       <div class="input-group-append">
                         <span class="input-group-text">%</span>
                       </div>
                     </div>
                   </div>
                   <div class="form-group">
                     <label class="form-control-label">Minimum Discount</label> <small class="text-danger">*</small>
                     <div class="input-group input-group-merge">
                       <input type="number" min="1" class="form-control" v-model="apriori.minimum_discount" placeholder="example: 10" required>
                       <div class="input-group-append">
                         <span class="input-group-text">%</span>
                       </div>
                     </div>
                   </div>
                   <div class="form-group">
                     <label class="form-control-label">Maximum Discount</label> <small class="text-danger">*</small>
                     <div class="input-group input-group-merge">
                       <input type="number" min="1" max="100" class="form-control" v-model="apriori.maximum_discount" placeholder="example: 15" required>
                       <div class="input-group-append">
                         <span class="input-group-text">%</span>
                       </div>
                     </div>
                   </div>
                   <div class="form-group">
                     <label class="form-control-label">Start Date</label> <small class="text-danger">*</small>
                     <input type="date" class="form-control" v-model="apriori.start_date" required>
                   </div>
                   <div class="form-group">
                     <label class="form-control-label">End Date</label> <small class="text-danger">*</small>
                     <input type="date" class="form-control" v-model="apriori.end_date" required>
                   </div>
                  <button class="btn btn-primary" type="submit">Submit form</button>
                </form>
              </div>
            </div>
            <div class="card">
              <div class="card-header bg-transparent">
                <h3 class="mb-0">Result Apriori</h3>
              </div>
              <div class="card-body">
                <template v-for="(item,i) in result" :key="i">
                  <div class="timeline timeline-one-side mt-3" data-timeline-content="axis" data-timeline-axis-style="dashed">
                    <div class="timeline-block">
                      <template v-if="i + 1 == result.length && item[item.length-1].description == `Rules`">
                        <span class="timeline-step badge-primary">
                          <i class="ni ni-diamond"></i>
                        </span>
                      </template>
                      <template v-else>
                        <span class="timeline-step badge-success">
                          {{ i + 1 }}
                        </span>
                      </template>
                      <div class="timeline-content">
                        <table class="table align-items-center table-flush">
                          <thead class="thead-light">
                          <tr>
                            <th>Item Set</th>
                            <th>Transaction</th>
                            <th>Support</th>
                            <template v-if="i + 1 == result.length && item[item.length-1].description == `Rules`">
                              <th>Confidence</th>
                              <th>Discount</th>
                            </template>
                            <th>Status</th>
                          </tr>
                          </thead>
                          <tbody>
                            <template v-for="(value,z) in item" :key="z">
                                <tr>
                                  <td>{{ value.item_set.join(", ") }}</td>
                                  <td>{{ value.transaction }}</td>
                                  <td>{{ value.support }}%</td>
                                  <template v-if="i + 1 == result.length && value.description == `Rules`">
                                    <th>{{ value.confidence }}</th>
                                    <th>{{ value.discount }}</th>
                                  </template>
                                  <td>
                                    <template v-if="value.description == `Eligible` || value.description == `Rules`">
                                      <span class="badge badge-pill badge-success"><i class="ni ni-check-bold" style="font-size: 1.5em;"></i></span>
                                    </template>
                                    <template v-else>
                                      <span class="badge badge-pill badge-danger"><i class="ni ni-fat-remove" style="font-size: 1.5em;"></i></span>
                                    </template>
                                  </td>
                                </tr>
                            </template>
                          </tbody>
                        </table>
                        <div class="mt-3">
                          <template v-for="(value,z) in item" :key="z">
                            <template v-if="value.description == `Eligible`">
                              <span class="badge badge-pill badge-success">{{ value.item_set.join(", ") }}</span>
                            </template>
                          </template>
                        </div>
                      </div>
                    </div>
                  </div>
                </template>
                <template v-if="result">

                </template>
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

<style scoped>
  .timeline-one-side .timeline-content {
    max-width: 300rem;
  }
</style>

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
  data: function () {
    return {
      result: [],
      test: [],
      apriori: {
        minimum_support: 0,
        minimum_confidence: 0,
        minimum_discount: 0,
        maximum_discount: 0,
        start_date: "",
        end_date: ""
      }
    };
  },
  methods: {
    submit() {
      this.apriori = {
            minimum_support: 30,
            minimum_confidence: 70,
            minimum_discount: 10,
            maximum_discount: 15,
            start_date: "2022-05-21",
            end_date: "2022-05-21"
      }

      axios.post("http://localhost:3000/api/apriori/generate", this.apriori)
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              // Logic for clean data aproiri
              let counter = 1
              let totalIterate = response.data.data[response.data.data.length-1].iterate
              for(let i = 0; i < totalIterate; i++){
                this.result.push([])
              }
              for(let i = 0; i < response.data.data.length; i++) {
                if(response.data.data[i].iterate === counter) {
                  this.result[response.data.data[i].iterate-1].push(response.data.data[i])
                } else {
                  counter++
                  this.result[response.data.data[i].iterate-1].push(response.data.data[i])
                }
              }
            }
          }).catch(error => {
            if (error.response == undefined) {
              alert("data is empty")
            } else {
              alert(error.response.data.status)
            }
      })
    }
  }
}
</script>