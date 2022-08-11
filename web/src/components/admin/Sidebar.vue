<template>
  <nav class="sidenav navbar navbar-vertical fixed-left navbar-expand-xs navbar-light bg-white" id="sidenav-main">
    <div class="scrollbar-inner">
      <!-- Brand -->
      <div class="sidenav-header d-flex align-items-center">
        <router-link class="navbar-brand" :to="{ name: 'admin' }" style="font-family: 'Fugaz One', cursive; color: #6772e5">
          RYZYOLSHOP
        </router-link>
        <div class="ml-auto">
          <!-- Sidenav toggler -->
          <div class="sidenav-toggler d-none d-xl-block" data-action="sidenav-unpin" data-target="#sidenav-main">
            <div class="sidenav-toggler-inner">
              <i class="sidenav-toggler-line"></i>
              <i class="sidenav-toggler-line"></i>
              <i class="sidenav-toggler-line"></i>
            </div>
          </div>
        </div>
      </div>
      <div class="navbar-inner">
        <!-- Collapse -->
        <div class="collapse navbar-collapse" id="sidenav-collapse-main">
          <!-- Nav items -->
          <ul class="navbar-nav">
            <li class="nav-item">
              <router-link :class="getActiveNavLink('admin')" :to="{ name: 'admin' }">
                <i class="ni ni-shop text-primary"></i>
                <span class="nav-link-text">Home</span>
              </router-link>
            </li>
            <li class="nav-item">
              <a :class="getActiveNavLink('transaction')" href="#navbar-examples" data-toggle="collapse" role="button" :aria-expanded="getAreaExpandedBool('transaction')" aria-controls="navbar-examples">
                <i class="ni ni-cart text-orange"></i>
                <span class="nav-link-text">Transaction</span>
              </a>
              <div :class="getCollapseNavLink('transaction')" id="navbar-examples">
                <ul class="nav nav-sm flex-column">
                  <li class="nav-item">
                    <router-link :to="{ name: 'transaction' }" class="nav-link">Data</router-link>
                  </li>
                  <li class="nav-item">
                    <router-link :to="{ name: 'transaction.create' }" class="nav-link">Create</router-link>
                  </li>
                  <li class="nav-item">
                    <router-link :to="{ name: 'transaction.create.csv' }" class="nav-link">Import</router-link>
                  </li>
                </ul>
              </div>
            </li>
            <li class="nav-item">
              <a :class="getActiveNavLink('product')" href="#navbar-components" data-toggle="collapse" role="button" :aria-expanded="getAreaExpandedBool('product')" aria-controls="navbar-components">
                <i class="ni ni-app text-info"></i>
                <span class="nav-link-text">Product</span>
              </a>
              <div :class="getCollapseNavLink('product')" id="navbar-components">
                <ul class="nav nav-sm flex-column">
                  <li class="nav-item">
                    <router-link :to="{ name: 'product' }" class="nav-link">Data</router-link>
                  </li>
                  <li class="nav-item">
                    <router-link :to="{ name: 'product.create' }" class="nav-link">Create</router-link>
                  </li>
                </ul>
              </div>
            </li>
            <li class="nav-item">
              <a :class="getActiveNavLink('category')" href="#navbar-category" data-toggle="collapse" role="button" :aria-expanded="getAreaExpandedBool('category')" aria-controls="navbar-category">
                <i class="ni ni-tag text-dark"></i>
                <span class="nav-link-text">Category</span>
              </a>
              <div :class="getCollapseNavLink('category')" id="navbar-category">
                <ul class="nav nav-sm flex-column">
                  <li class="nav-item">
                    <router-link :to="{ name: 'category' }" class="nav-link">Data</router-link>
                  </li>
                  <li class="nav-item">
                    <router-link :to="{ name: 'category.create' }" class="nav-link">Create</router-link>
                  </li>
                </ul>
              </div>
            </li>
            <li class="nav-item">
              <a :class="getActiveNavLink('apriori')" href="#navbar-forms" data-toggle="collapse" role="button" :aria-expanded="getAreaExpandedBool('apriori')" aria-controls="navbar-forms">
                <i class="ni ni-ungroup text-pink"></i>
                <span class="nav-link-text">Apriori</span>
              </a>
              <div :class="getCollapseNavLink('apriori')" id="navbar-forms">
                <ul class="nav nav-sm flex-column">
                  <li class="nav-item">
                    <router-link :to="{ name: 'apriori' }" class="nav-link">Data</router-link>
                  </li>
                  <li class="nav-item">
                    <router-link :to="{ name: 'apriori.create' }" class="nav-link">Create</router-link>
                  </li>
                </ul>
              </div>
            </li>
            <li class="nav-item">
              <a :class="getActiveNavLink('user')" href="#navbar-users" data-toggle="collapse" role="button" :aria-expanded="getAreaExpandedBool('user')" aria-controls="navbar-users">
                <i class="ni ni-circle-08 text-danger"></i>
                <span class="nav-link-text">User</span>
              </a>
              <div :class="getCollapseNavLink('user')" id="navbar-users">
                <ul class="nav nav-sm flex-column">
                  <li class="nav-item">
                    <router-link :to="{ name: 'user' }" class="nav-link">Data</router-link>
                  </li>
                  <li class="nav-item">
                    <router-link :to="{ name: 'user.create' }" class="nav-link">Create</router-link>
                  </li>
                </ul>
              </div>
            </li>
            <li class="nav-item">
              <a :class="getActiveNavLink('user-order')" href="#navbar-user-order" data-toggle="collapse" role="button" :aria-expanded="getAreaExpandedBool('user-order')" aria-controls="navbar-user-order">
                <i class="ni ni-spaceship text-blue"></i>
                <span class="nav-link-text">Order</span>
              </a>
              <div :class="getCollapseNavLink('user-order')" id="navbar-user-order">
                <ul class="nav nav-sm flex-column">
                  <li class="nav-item">
                    <router-link :to="{ name: 'user-order' }" class="nav-link">Data</router-link>
                  </li>
                </ul>
              </div>
            </li>
          </ul>
          <!-- Divider -->
          <hr class="my-3">
          <!-- Heading -->
          <h6 class="navbar-heading p-0 text-muted">Others</h6>
          <!-- Navigation -->
          <ul class="navbar-nav mb-md-3">
            <li class="nav-item">
              <router-link class="nav-link" :to="{ name: 'guest.index' }">
                <i class="ni ni-basket"></i>
                <span class="nav-link-text">Back to Shop</span>
              </router-link>
            </li>
            <li class="nav-item">
              <router-link :class="getActiveNavLink('profile')" :to="{ name: 'profile' }">
                <i class="ni ni-single-02"></i>
                <span class="nav-link-text">My Account</span>
              </router-link>
            </li>
            <li class="nav-item">
              <form @submit.prevent="submit" method="POST" class="nav-link">
                <i class="ni ni-curved-next"></i>
                <button class="btn btn-link" style="padding: 0; font-weight: normal; color: rgba(0, 0, 0, .6)" type="submit">Logout</button>
              </form>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </nav>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Fugaz+One&display=swap');

.navbar-brand {
  font-size: 1.40rem;
}
</style>

<script>
import axios from "axios";
import authHeader from "@/service/auth-header";

export default {
  methods: {
    submit() {
      axios.delete(`${process.env.VUE_APP_SERVICE_URL}/auth/logout`,{ headers: authHeader() })
          .then(response => {
            if(response.data.code === 200) {
              localStorage.removeItem("token")
              localStorage.removeItem("refresh-token")
              alert(response.data.status)
              this.$router.push({
                name: 'auth.login'
              })
            }
          }).catch(error => {
        console.log(error.response.data.status)
      })
    },
    getActiveNavLink(name) {
      let classString = "nav-link "

      let routeName = this.$route.name
      let route = routeName.split(".")
      if (route[0] === name) {
        classString += "active"
      }
      return classString;
    },
    getCollapseNavLink(name) {
      let classString = "collapse "

      let routeName = this.$route.name
      let route = routeName.split(".")
      if (route[0] === name) {
        classString += "show"
      }
      return classString;
    },
    getAreaExpandedBool(name) {
      let classString = "false"

      let routeName = this.$route.name
      let route = routeName.split(".")
      if (route[0] === name) {
        classString = "true"
      }
      return classString;
    }
  }
}
</script>