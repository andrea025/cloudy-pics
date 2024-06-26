<script>
import PhotoCard from '../components/PhotoCard.vue';

export default {
  name: 'HomeView',
  components: {PhotoCard},
  data: function() {
    return {
      loading: false,
      errormsg: null,
      id: localStorage.getItem('id'),
      photos: [],
      photoShowComments: {},
      photoHasLike: {},
      commentText: '',
    };
  },
  methods: {
    async getMyStream() {
      this.loading = true;
      this.errormsg = null;
      try {
        const response = await this.$axios.get('/users/' + this.id + '/stream');
        this.photos = response.data.data;
      } catch (e) {
        switch (e.response.status) {
          case 400:
            this.errormsg = 'Ops, there was something wrong with your request.';
            break;
          case 401:
            this.errormsg = 'You need to login in order to perform this action.';
            break;
          case 403:
            this.errormsg = 'Action forbidden.';
            break;
          case 500:
            this.errormsg = 'Ops, there was an internal problem with the server.';
            break;
          default:
            this.errormsg = e.toString();
        }
      }
      this.loading = false;
    },
  },
  mounted() {
    this.getMyStream();
  },
};
</script>


<template>
  <div class="section page">
    <div class="loader-overlay" v-if="this.loading">
      <div class="loader" >
        <LoadingSpinner :loading="this.loading" />
      </div>
    </div>
    <div class="inner">
      <div class="center">
        <ErrorMsg v-if="this.errormsg" :msg="this.errormsg" />
        <div class="feeds-ctn" v-else-if="!this.errormsg">
          <div class="posts-ctn">
            <div v-for="photo in this.photos" :key="photo" class="post-ctn">
              <PhotoCard :photoObj="photo" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


<style scoped>
.section > .inner {
  display: flex;
  justify-content: center;
}
.center {
  max-width: 615px;
  margin-right: 25px;
}
.feeds-ctn {
  max-width: 615px;
}
.post-ctn {
  margin-top: 20px;
  width: 100%;
  border: 1px solid var(--border-color);
}
.loader-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.9);
  z-index: 9999;
}
.loader{
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  width: 100%;
}
</style>
