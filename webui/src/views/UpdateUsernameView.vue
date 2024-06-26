<script>
export default {
  name: 'UpdateUsernameView',
  data: function() {
    return {
      errormsg: null,
      username: '',
    };
  },
  methods: {
    async setMyUsername() {
      this.errormsg = null;
      try {
        await this.$axios.patch('/users/' + localStorage.getItem('id'), {username: this.username});
      } catch (e) {
        switch (e.response.status) {
          case 400:
            this.errormsg = 'Ops, there was something wrong with your request.';
            break;
          case 401:
            this.errormsg = 'You need to login in order to perform this action.';
            break;
          case 403:
            this.errormsg = 'Action forbidden: you cannot change the username of another user.';
            break;
          case 409:
            this.errormsg = 'Username already taken.';
            break;
          case 500:
            this.errormsg = 'Ops, there was an internal problem with the server.';
            break;
          default:
            this.errormsg = e.toString();
        }
      }
      this.$router.replace('/users/' + localStorage.getItem('id'));
    },
    checkAuth() {
      if (!(localStorage.getItem('id'))) {
        this.errormsg = 'You need to login in order to perform this action.';
      }
    },
  },
  mounted() {
    this.checkAuth();
  },
};
</script>


<template>
<div class="section page">
  <ErrorMsg v-if="this.errormsg" :msg="this.errormsg" />
  <form @submit.prevent="setMyUsername" v-else-if="!this.errormsg">
    <label for="username">Enter your new username here:</label><br>
    <input type="text" id="username" name="username" placeholder="New Username..." minlength="3" maxlength="16" v-model="username"><br>
    <button type="submit">Update</button>
  </form>
</div>
</template>


<style scoped>
form {
  width: 400px;
  margin: 0 auto;
  text-align: center;
  font-family: 'Helvetica Neue', sans-serif;
}

label {
  font-size: 18px;
  font-weight: bold;
  color: #262626;
  display: block;
  margin-top: 100px;
  margin-bottom: 8px;
}

input[type="text"] {
  width: 60%;
  padding: 12px 20px;
  margin: 8px 0;
  box-sizing: border-box;
  border: 2px solid #262626;
  border-radius: 8px;
  font-size: 16px;
}

button[type="submit"] {
  width: 27%;
  background-color: #e1306c;
  color: white;
  padding: 14px 20px;
  margin: 8px 0;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  font-weight: bold;
}

button[type="submit"]:hover {
  background-color: #cf255b;
}
</style>
