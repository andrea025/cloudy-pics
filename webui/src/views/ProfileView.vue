<script>
import NewPostIcon from '../components/NewPostIcon.vue';

export default {
  name: 'ProfileView',
  components: {NewPostIcon},
  data: function() {
    return {
      errormsg: null,
      loading: false,
      reqId: localStorage.getItem('id'),
      id: '',
      username: '',
      followersCount: null,
      followingCount: null,
      uploadedPhotos: null,
      photos: [],
      followers: [],
      following: [],
      banned: [],
      isFollowing: null,
      isBanned: null,
      photoBaseUrl: __API_URL__,
    };
  },
  methods: {
    async getUserProfile() {
      this.loading = true;
      this.errormsg = null;
      try {
        const response = await this.$axios.get('/users/' + this.$route.params.id);
        this.id = response.data.id;
        this.username = response.data.username;
        this.followersCount = response.data.followers;
        this.followingCount = response.data.following;
        this.uploadedPhotos = response.data.uploaded_photos;
        this.photos = response.data.photos;
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
          case 404:
            this.errormsg = 'Ops, user not found.';
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
      if (this.reqId != this.id) {
        await this.checkFollowing();
        await this.checkBanned();
      }
      this.loading = false;
    },
    async getFollowers() {
      this.loading = true;
      this.errormsg = null;
      try {
        const response = await this.$axios.get('/users/' + this.$route.params.id + '/followers');
        this.followers = response.data.data;
      } catch (e) {
        switch (e.response.status) {
          case 400:
            this.errormsg = 'Ops, there was something wrong with your request.';
            break;
          case 401:
            this.errormsg = 'You need to login in order to perform this action.';
            break;
          case 404:
            this.errormsg = 'Ops, user not found.';
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
    async getBanned(id) {
      this.loading = true;
      this.errormsg = null;
      try {
        const response = await this.$axios.get('/users/' + id + '/banned');
        this.banned = response.data.data;
      } catch (e) {
        switch (e.response.status) {
          case 400:
            this.errormsg = 'Ops, there was something wrong with your request.';
            break;
          case 401:
            this.errormsg = 'You need to login in order to perform this action.';
            break;
          case 403:
            this.errormsg = 'Action forbidden: you cannot see the list of banned users of another user.';
          case 404:
            this.errormsg = 'Ops, user not found.';
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
    async checkFollowing() {
      await this.getFollowers();
      for (const follower of this.followers) {
        if (follower.id == this.reqId) {
          this.isFollowing = true;
        }
      }
    },
    async checkBanned() {
      await this.getBanned(this.reqId);
      for (const ban of this.banned) {
        if (ban.id == this.id) {
          this.isBanned = true;
        }
      }
    },
    async followAction() {
      this.loading = true;
      this.errormsg = null;
      try {
        if (this.isFollowing) {
          await this.$axios.delete('/users/' + this.reqId + '/following/' + this.id);
          this.followersCount--;
        } else {
          await this.$axios.put('/users/' + this.reqId + '/following/' + this.id);
          this.followersCount++;
        }
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
          case 404:
            this.errormsg = 'Ops, user not found.';
            break;
          case 500:
            this.errormsg = 'Ops, there was an internal problem with the server.';
            break;
          default:
            this.errormsg = e.toString();
        }
      }
      this.loading = false;
      this.isFollowing = !this.isFollowing;
    },
    async banAction() {
      this.loading = true;
      this.errormsg = null;
      try {
        if (this.isBanned) {
          await this.$axios.delete('/users/' + this.reqId + '/banned/' + this.id);
        } else {
          await this.$axios.put('/users/' + this.reqId + '/banned/' + this.id);
        }
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
          case 404:
            this.errormsg = 'Ops, user not found.';
            break;
          case 500:
            this.errormsg = 'Ops, there was an internal problem with the server.';
            break;
          default:
            this.errormsg = e.toString();
        }
      }
      this.loading = false;
      this.isBanned = !this.isBanned;
    },
    getPhoto(photoId) {
      this.$router.push('/photos/' + photoId);
    },
    onPickFile() {
      this.$refs.fileInput.click();
    },
    async onFilePicked(event) {
      const files = event.target.files;
      const fileReader = new FileReader();
      fileReader.readAsDataURL(files[0]);
      const image = files[0];

      if (!['image/jpeg', 'image/png'].includes(image.type)) {
        console.log(image.type);
        this.errormsg = 'File type not supported, please upload a png or jpeg image.';
        return;
      }

      this.loading = true;
      try {
        await this.$axios.post('/users/' + this.id + '/photos', image, {
          headers: {
            'Content-Type': image.type,
          },
        });
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
          case 404:
            this.errormsg = 'Ops, user not found.';
            break;
          case 413:
            this.errormsg = 'The selected image is too large.';
            break;
          case 500:
            this.errormsg = 'Ops, there was an internal problem with the server.';
            break;
          default:
            this.errormsg = e.toString();
        }
      }
      this.$router.go(0);
      this.loading = false;
    },
  },
  mounted() {
		  this.getUserProfile();
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
    <ErrorMsg v-if="this.errormsg" :msg="this.errormsg" />
    <div class="inner" v-if="!this.errormsg">
      <div class="header">
        <div class="ls" ></div>
        <div class="user-info">
          <div class="block-1">
            <div class="username">
              {{ this.username }}
            </div>
            <router-link class="edit-btn" to="/profile/username" v-if="this.reqId == this.id">
              Edit Username
            </router-link>
            <button class="clear-btn" @click="onPickFile">
              <NewPostIcon />
            </button>
            <input type="file" style="display: none" ref="fileInput" accept=".png,.jpeg" @change="onFilePicked" />
          </div>
          <div class="block-2">
            <div class="user-stat">
              <b>{{ this.uploadedPhotos }}</b> photos
            </div>
            <div class="user-stat">
              <router-link class="user-stats-text" :to="`${this.$route.path}/followers`">
                <b>{{ this.followersCount }}</b> followers
              </router-link>
            </div>
            <div class="user-stat">
              <router-link class="user-stats-text" :to="`${this.$route.path}/following`">
                <b>{{ this.followingCount }}</b> following
              </router-link>
            </div>
            <div class="user-stat" v-if="this.reqId == this.id">
              <router-link class="user-stats-text" :to="`${this.$route.path}/banned`">
                banned users
              </router-link>
            </div>
          </div>
          <div v-if="this.reqId != this.id" class="relationship-btn">
            <button class="follow-btn" type="submit" @click="followAction">{{ this.isFollowing ? 'Unfollow' : 'Follow' }}</button>
            <button class="ban-btn" type="submit" @click="banAction">{{ this.isBanned ? 'Unban' : 'Ban' }}</button>
          </div>
        </div>
      </div>
      <div class="user-content">
        <div class="tab-layout">
          <div :class="{'clear-btn': true, active: $route.name === 'profile'}">
            <svg aria-label="" class="_8-yf5 " color="#8e8e8e" fill="#8e8e8e" height="12" role="img" viewBox="0 0 24 24" width="12"><rect fill="none" height="18" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" width="18" x="3" y="3"></rect><line fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" x1="9.015" x2="9.015" y1="3" y2="21"></line><line fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" x1="14.985" x2="14.985" y1="3" y2="21"></line><line fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" x1="21" x2="3" y1="9.015" y2="9.015"></line><line fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" x1="21" x2="3" y1="14.985" y2="14.985"></line></svg>
            PHOTOS
          </div>
        </div>
        <div class="grid" v-if="$route.name === 'profile'">
          <div v-for="photo in this.photos" :key="photo" class="grid-item">
            <div :class="{ 'post-grid-ctn': true }">
              <div class="image">
                <img :src="this.photoBaseUrl+photo.photo_url" alt="">
              </div>
              <div class="image-overlay" @click="this.getPhoto(photo.id)">
                <div class="post-count">
                  <div class="likes">
                    <span>
                      <svg aria-label="Activity Feed" class="_8-yf5 " color="#FFF" fill="#FFF" height="24" role="img" viewBox="0 0 48 48" width="24"><path d="M34.6 3.1c-4.5 0-7.9 1.8-10.6 5.6-2.7-3.7-6.1-5.5-10.6-5.5C6 3.1 0 9.6 0 17.6c0 7.3 5.4 12 10.6 16.5.6.5 1.3 1.1 1.9 1.7l2.3 2c4.4 3.9 6.6 5.9 7.6 6.5.5.3 1.1.5 1.6.5s1.1-.2 1.6-.5c1-.6 2.8-2.2 7.8-6.8l2-1.8c.7-.6 1.3-1.2 2-1.7C42.7 29.6 48 25 48 17.6c0-8-6-14.5-13.4-14.5z"></path></svg>          </span>
                    <span>
                      {{ photo.likes.count }}
                    </span>
                  </div>
                  <div class="comments">
                    <span>
                      <svg aria-label="Comment" class="_8-yf5 " color="none" fill="none" height="24" role="img" viewBox="0 0 24 24" width="24"><path d="M20.656 17.008a9.993 9.993 0 10-3.59 3.615L22 22z" fill="#FFFFFF" stroke="currentColor" stroke-linejoin="round" stroke-width="2"></path></svg>
                    </span>
                    <span>
                      {{ photo.comments.count }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

</template>


<style scoped>
  .header {
    display: flex;
    align-items: flex-start;
    margin-bottom: 45px;
  }
  .ls {
    width: 290px;
    margin-right: 30px;
    display: flex;
    justify-content: center;
  }
  b {
    font-weight: 600;
  }
  .block-1 {
    margin-bottom: 24px;
    display: flex;
    align-items: center;
  }
  .username {
    font-size: 1.6rem;
    font-weight: 600;
    margin-right: 20px;
  }
  .edit-btn {
    text-decoration: none;
    height: 30px;
    width: 125px;
    display: grid;
    place-items: center;
    border: 1px solid var(--border-color);
    font-size: 0.875rem;
    color: #000000;
    border-radius: 4px;
    margin-right: 12px;
    font-weight: 600;
  }
  .block-2 {
    margin-bottom: 24px;
    display: flex;
  }
  .user-stat {
    margin-right: 40px;
  }

.follow-btn, .ban-btn {
  font-family: 'Helvetica Neue', sans-serif;
  font-size: 14px;
  font-weight: 600;
  color: white;
  background-color: #3897f0;
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  cursor: pointer;
}

.follow-btn:hover, .ban-btn:hover {
  background-color: #2b88d9;
}

.follow-btn:active, .ban-btn:active {
  background-color: #2b88d9;
  transform: scale(0.95);
}

.ban-btn {
  background-color: #e1306c;
  margin-left: 10px;
}

.ban-btn:hover {
  background-color: #cf255b;
}

.ban-btn:active {
  background-color: #cf255b;
  transform: scale(0.95);
}

  /* User Content */
  .user-content {
    width: 100%;
  }
  .tab-layout {
    border-top: 1px solid var(--border-color);
    width: 100%;
    display: flex;
    justify-content: center;
  }
  .tab-layout .clear-btn {
    height: 53px;
    width: 73px;
    display: flex;
    justify-content: space-around;
    align-items: center;
    margin-right: 60px;
    font-size: 0.75rem;
    font-weight: 600;
    transition: 0.3s;
    position: relative;
    bottom: 1px;
  }
  .tab-layout .clear-btn:active {
    opacity: 0.5;
  }
  .tab-layout .clear-btn.active {
    color: #000000;
    border-top: 1px solid #000000;
  }
  .tab-layout .clear-btn.active svg {
    color: #000000;
  }
  .tab-layout .clear-btn:nth-last-child(1) {
    margin-right: 0;
  }
  .flex {
    display: flex;
    justify-content: flex-start;
    flex-wrap: wrap;
  }
  .post-grid-ctn {
    flex-basis: 32%;
    margin-right: 2%;
    margin-bottom: 20px;
    position: relative;
    transition: 0.3s;
  }
  .post-grid-ctn.bigger-grid {
    flex-basis: 66.2%;
    margin-right: 0;
  }
  .post-grid-ctn.third {
    margin-right: 0;
  }
  .image {
    width: 100%;
    cursor: pointer;
    z-index: 999;
  }
  .image > img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    cursor: pointer;
    display: block;
  }
  .image:hover + .image-overlay {
    display: grid;
    opacity: 1;
  }

  .image-overlay:hover {
    display: grid;
  }
  .image-overlay {
    cursor: pointer;
    background: #00000044;
    color: #FFFFFF;
    place-items: center;
    position: absolute;
    inset: 0;
    bottom: 3px;
    transition: 1000s;
    opacity: 0;
    display: none;
  }
  .post-grid-ctn.bigger-grid .image-overlay {
    bottom: 3px;
  }
  .post-count {
    display: flex;
  }
  .post-count > div {
    display: flex;
    align-items: center;
  }
  .post-count > div > span:nth-child(1) {
    margin-right: 4px;
  }
  .likes {
    margin-right: 20px;
  }

.user-stats-text{
  color: black;
  text-decoration: none;
}

.delete-button {
  background: transparent;
  border: none;
  height: 24px;
  width: 24px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
}

.grid-item {
  padding: 10px;
}

.image {
  width: 100%;
  height: 285px;
  object-fit: cover;
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
