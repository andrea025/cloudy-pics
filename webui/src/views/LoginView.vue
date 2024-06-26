<script>
export default {
  name: 'LoginView',
  data: function() {
    return {
      errormsg: null,
	  loading: false,
      username: '',
      id: '',
    };
  },
  methods: {
    async login() {
      this.loading = true;
	  this.errormsg = null;
	  try {
        const response = await this.$axios.post('/session', {username: this.username});
	    this.id = response.data.id;
	  } catch (e) {
	    switch (e.response.status) {
          case 400:
            this.errormsg = 'Ops, there was something wrong with your request: send a valid username.';
            break;
          case 500:
            this.errormsg = 'Ops, there was an internal problem with the server.';
            break;
          default:
            this.errormsg = e.toString();
        }
        this.loading = false;
        return;
	  }
	  this.loading = false;
      localStorage.setItem('id', this.id);
      this.$router.replace('/home');
    },
    checkRedirect() {
      if (localStorage.getItem('id')) {
        this.$router.replace('/home');
      }
    },
  },
  mounted() {
    this.checkRedirect();
  },
};
</script>


<template>
<div>
  <LoadingSpinner v-if="this.loading" />
  <section class="vh-100 gradient-custom">
    <div class="container py-5 h-100">
      <div class="row d-flex justify-content-center align-items-center h-100">
        <ErrorMsg v-if="this.errormsg" :msg="errormsg" />
        <div class="col-12 col-md-8 col-lg-6 col-xl-5">
          <div class="card bg-dark text-white" style="border-radius: 1rem;">
            <div class="card-body p-5 text-center">
              <div class="mb-md-5 mt-md-4 pb-5">
                <h2 class="fw-bold mb-2 text-uppercase">Welcome!</h2>
                <p class="text-white-50 mb-5">Please enter your username to login</p>
                <form @submit.prevent="login">
                  <div class="form-outline form-white mb-4">
                    <label class="form-label">Your Username:</label>
                    <input type="username" class="form-control form-control-lg" placeholder="John" v-model="username" />
                  </div>
                  <button class="btn btn-outline-light btn-lg px-5" type="submit">Login</button>
                </form>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</div>
</template>


<style scoped>
form {
  font-family: 'Helvetica Neue', sans-serif;
}

.gradient-custom {
  /* fallback for old browsers */
  background: #e1306c;
  /* Chrome 10-25, Safari 5.1-6 */
  background: linear-gradient(to right, #e1306c, #a81d49);
  /* W3C, IE 10+/ Edge, Firefox 16+, Chrome 26+, Opera 12+, Safari 7+ */
  background: linear-gradient(to right, #e1306c, #a81d49);
}
</style>
