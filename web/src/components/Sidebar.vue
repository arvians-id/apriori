<template>
  <nav class="sidenav navbar navbar-vertical fixed-left navbar-expand-xs navbar-light bg-white" id="sidenav-main">
    <div class="scrollbar-inner">
      <!-- Brand -->
      <div class="sidenav-header d-flex align-items-center">
        <router-link class="navbar-brand" :to="{ name: 'admin' }">
          <img src="/frontend/img/brand/blue.png" class="navbar-brand-img" alt="...">
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
              <router-link class="nav-link" :to="{ name: 'admin' }">
                <i class="ni ni-shop text-primary"></i>
                <span class="nav-link-text">Home</span>
              </router-link>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="#navbar-examples" data-toggle="collapse" role="button" aria-expanded="false" aria-controls="navbar-examples">
                <i class="ni ni-cart text-orange"></i>
                <span class="nav-link-text">Transaction</span>
              </a>
              <div class="collapse" id="navbar-examples">
                <ul class="nav nav-sm flex-column">
                  <li class="nav-item">
                    <router-link :to="{ name: 'transaction' }" class="nav-link">Data</router-link>
                  </li>
                  <li class="nav-item">
                    <router-link :to="{ name: 'transaction.create' }" class="nav-link">Create</router-link>
                  </li>
                  <li class="nav-item">
                    <router-link :to="{ name: 'transaction.create.csv' }" class="nav-link">Create CSV</router-link>
                  </li>
                </ul>
              </div>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="#navbar-components" data-toggle="collapse" role="button" aria-expanded="false" aria-controls="navbar-components">
                <i class="ni ni-app text-info"></i>
                <span class="nav-link-text">Product</span>
              </a>
              <div class="collapse" id="navbar-components">
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
              <a class="nav-link" href="#navbar-forms" data-toggle="collapse" role="button" aria-expanded="false" aria-controls="navbar-forms">
                <i class="ni ni-ungroup text-pink"></i>
                <span class="nav-link-text">Apriori</span>
              </a>
              <div class="collapse" id="navbar-forms">
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
          </ul>
          <!-- Divider -->
          <hr class="my-3">
          <!-- Heading -->
          <h6 class="navbar-heading p-0 text-muted">Others</h6>
          <!-- Navigation -->
          <ul class="navbar-nav mb-md-3">
            <li class="nav-item">
              <router-link class="nav-link" :to="{ name: 'profile' }">
                <i class="ni ni-single-02"></i>
                <span class="nav-link-text">My Account</span>
              </router-link>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="https://demos.creative-tim.com/argon-dashboard/docs/foundation/colors.html">
                <i class="ni ni-curved-next"></i>
                <span class="nav-link-text">Logout</span>
              </a>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </nav>
</template>