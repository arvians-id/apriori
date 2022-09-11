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
                <div :class="{ hide: !submitted }" id="hide">
                  <ul class="list-group list-group-flush mb-4" data-toggle="checklist">
                    <li class="checklist-entry list-group-item flex-column align-items-start py-4 px-4">
                      <div class="checklist-item checklist-item-success">
                        <div class="checklist-info">
                          <h5 class="checklist-title mb-0">Minimum Support</h5>
                          <small>{{ apriori.minimum_support }}%</small>
                        </div>
                      </div>
                    </li>
                    <li class="checklist-entry list-group-item flex-column align-items-start py-4 px-4">
                      <div class="checklist-item checklist-item-warning">
                        <div class="checklist-info">
                          <h5 class="checklist-title mb-0">Minimum Confidence</h5>
                          <small>{{ apriori.minimum_confidence }}%</small>
                        </div>
                      </div>
                    </li>
                    <li class="checklist-entry list-group-item flex-column align-items-start py-4 px-4">
                      <div class="checklist-item checklist-item-info">
                        <div class="checklist-info">
                          <h5 class="checklist-title mb-0">Discount</h5>
                          <small>{{ `${apriori.minimum_discount}% - ${apriori.maximum_discount}%` }}</small>
                        </div>
                      </div>
                    </li>
                    <li class="checklist-entry list-group-item flex-column align-items-start py-4 px-4">
                      <div class="checklist-item checklist-item-danger">
                        <div class="checklist-info">
                          <h5 class="checklist-title mb-0">Range Date</h5>
                          <small>{{ `${apriori.start_date} - ${apriori.end_date}` }}</small>
                        </div>
                      </div>
                    </li>
                  </ul>
                </div>
                <template v-for="(item,i) in result" :key="i">
                  <div class="timeline timeline-one-side mt-3" data-timeline-content="axis" data-timeline-axis-style="dashed">
                    <div class="timeline-block">
                      <template v-if="i + 1 == result.length && item[item.length-1].description == `Rules`">
                        <span class="timeline-step badge-primary">
                          <i class="ni ni-diamond"></i>
                        </span>
                      </template>
                      <template v-else>
                        <span class="timeline-step badge-info">
                          {{ i + 1 }}
                        </span>
                      </template>
                      <div class="timeline-content table-responsive">
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
                                    <td>{{ value.confidence }}%</td>
                                    <td>{{ value.discount }}%</td>
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
                            <template v-if="value.description == `Eligible` || value.description == `Rules`">
                              <span class="badge badge-pill badge-success mr-2">{{ value.item_set.join(", ") }}</span>
                            </template>
                          </template>
                        </div>
                      </div>
                    </div>
                  </div>
                  <form @submit.prevent="save" method="POST" class="mt-3 text-center" v-if="item[item.length-1].description == `Rules`">
                    <button class="btn btn-primary" type="submit">Save</button>
                  </form>
                </template>
                <template v-if="result.length == 0">
                  <p class="text-center">No data available</p>
                </template>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="card">
        <div class="card-header">
          <h3 class="mb-0">Guide</h3>
        </div>
        <div class="card-body">
          <div class="row">
            <div class="col-md-4">
              <button type="button" class="btn btn-block btn-primary mb-3" data-toggle="modal" data-target="#modal-rules">What is Association Rules?</button>
              <div class="modal fade" id="modal-rules" tabindex="-1" role="dialog" aria-labelledby="modal-rules" aria-hidden="true">
                <div class="modal-dialog modal- modal-dialog-centered modal-" role="document">
                  <div class="modal-content">
                    <div class="modal-header">
                      <h6 class="modal-title" id="modal-title-rules">What is Association Rules?</h6>
                      <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true">×</span>
                      </button>
                    </div>
                    <div class="modal-body">
                      <p>Association rules atau aturan asosiasi adalah teknik data mining untuk menemukan aturan asosiasi antara suatu kombinasi item</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-md-4">
              <button type="button" class="btn btn-block btn-primary mb-3" data-toggle="modal" data-target="#modal-apriori">What is Apriori?</button>
              <div class="modal fade" id="modal-apriori" tabindex="-1" role="dialog" aria-labelledby="modal-apriori" aria-hidden="true">
                <div class="modal-dialog modal- modal-dialog-centered modal-" role="document">
                  <div class="modal-content">
                    <div class="modal-header">
                      <h6 class="modal-title" id="modal-title-apriori">What is Apriori?</h6>
                      <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true">×</span>
                      </button>
                    </div>
                    <div class="modal-body">
                      <p>Algoritma apriori adalah suatu metode untuk mencari pola hubungan antar satu atau lebih item dalam suatu dataset</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-md-4">
              <button type="button" class="btn btn-block btn-primary mb-3" data-toggle="modal" data-target="#modal-support">What is Support?</button>
              <div class="modal fade" id="modal-support" tabindex="-1" role="dialog" aria-labelledby="modal-support" aria-hidden="true">
                <div class="modal-dialog modal- modal-dialog-centered modal-" role="document">
                  <div class="modal-content">
                    <div class="modal-header">
                      <h6 class="modal-title" id="modal-title-support">What is Support?</h6>
                      <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true">×</span>
                      </button>
                    </div>
                    <div class="modal-body">
                      <p>Support (nilai penunjang) adalah persentase kombinasi item tersebut dalam database</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-md-4">
              <button type="button" class="btn btn-block btn-primary mb-3" data-toggle="modal" data-target="#modal-confidence">What is Confidence?</button>
              <div class="modal fade" id="modal-confidence" tabindex="-1" role="dialog" aria-labelledby="modal-confidence" aria-hidden="true">
                <div class="modal-dialog modal- modal-dialog-centered modal-" role="document">
                  <div class="modal-content">
                    <div class="modal-header">
                      <h6 class="modal-title" id="modal-title-confidence">What is Confidence?</h6>
                      <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true">×</span>
                      </button>
                    </div>
                    <div class="modal-body">
                      <p>Confidence (nilai kepastian) yaitu kuatnya hubungan antar item dalam aturan asosiatif yang terbentuk oleh metode asosiasi dalam data mining</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-md-4">
              <button type="button" class="btn btn-block btn-primary mb-3" data-toggle="modal" data-target="#modal-discount">About Discount</button>
              <div class="modal fade" id="modal-discount" tabindex="-1" role="dialog" aria-labelledby="modal-discount" aria-hidden="true">
                <div class="modal-dialog modal- modal-dialog-centered modal-" role="document">
                  <div class="modal-content">
                    <div class="modal-header">
                      <h6 class="modal-title" id="modal-title-discount">About Discount</h6>
                      <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true">×</span>
                      </button>
                    </div>
                    <div class="modal-body">
                      <p>Anda dapat menetapkan minimal diskon dan maksimal diskon agar perhitungan apriori tidak melebihi diskon yang tidak anda inginkan.</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-md-4">
              <button type="button" class="btn btn-block btn-primary mb-3" data-toggle="modal" data-target="#modal-date">About Date</button>
              <div class="modal fade" id="modal-date" tabindex="-1" role="dialog" aria-labelledby="modal-date" aria-hidden="true">
                <div class="modal-dialog modal- modal-dialog-centered modal-" role="document">
                  <div class="modal-content">
                    <div class="modal-header">
                      <h6 class="modal-title" id="modal-title-date">About Date</h6>
                      <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                        <span aria-hidden="true">×</span>
                      </button>
                    </div>
                    <div class="modal-body">
                      <p>Anda dapat mengatur tanggal yang anda inginkan kedalam perhitungan apriori.</p>
                    </div>
                  </div>
                </div>
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
  .hide {
    display:none;
   }
</style>

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
      result: [],
      test: [],
      submitted: false,
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
      this.result = []
      axios.post(`${process.env.VUE_APP_SERVICE_URL}/apriori/generate`, this.apriori, { headers: authHeader() })
          .then(response => {
            this.submitted = true
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
              if (error.response.status === 400 || error.response.status === 404) {
                alert(error.response.data.status)
              }
            }
      })
    },
    save() {
      axios.post(`${process.env.VUE_APP_SERVICE_URL}/apriori`, this.result[this.result.length-1], { headers: authHeader() })
          .then(response => {
            if(response.data.code === 200) {
              alert(response.data.status)
              this.$router.push({
                name: 'apriori'
              })
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