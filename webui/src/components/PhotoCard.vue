<script>
import CommentIcon from './CommentIcon.vue';
import EmojiIcon from './EmojiIcon.vue';
import LikeIcon from './LikeIcon.vue';
import DeleteIcon from './DeleteIcon.vue';

export default {
  name: 'PhotoCard',
  components: {CommentIcon, EmojiIcon, LikeIcon, DeleteIcon},
  props: {
    photoObj: {
      type: Object,
    },
  },
  data: function() {
    return {
      loading: false,
      errormsg: null,
      photo: this.photoObj,
      reqId: localStorage.getItem('id'),
      hasLike: false,
      commentsVisibility: false,
      commentText: '',
      photoBaseUrl: __API_URL__,
    };
  },
  methods: {
    showComments() {
      if (this.photo.comments.count > 0) {
        this.commentsVisibility = true;
      } else {
        this.commentsVisibility = false;
      }
    },
    async likeAction() {
      if (this.reqId == this.photo.owner.id) {
        return;
      }
      this.loading = true;
      this.errormsg = null;
      try {
        if (this.hasLike) {
          await this.$axios.delete('/photos/' + this.photo.id + '/likes/' + this.reqId);
          this.photo.likes.count--;
        } else {
          await this.$axios.put('/photos/' + this.photo.id + '/likes/' + this.reqId);
          this.photo.likes.count++;
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
            this.errormsg = 'Ops, photo not found.';
            break;
          case 500:
            this.errormsg = 'Ops, there was an internal problem with the server.';
            break;
          default:
            this.errormsg = e.toString();
        }
      }
      this.loading = false;
      this.hasLike = !this.hasLike;
    },
    checkLike() {
      if (this.reqId != this.photo.owner.id) {
        for (const user of this.photo.likes.users) {
          if (user.id == this.reqId) {
            this.hasLike = true;
          }
          return;
        }
      }
      this.hasLike = false;
    },
    async commentPhoto() {
      this.loading = true;
      this.errormsg = null;
      try {
        const response = await this.$axios.post('/photos/' + this.photo.id + '/comments', {text: this.commentText});
        this.photo.comments.count++;
        this.photo.comments.comments.push(response.data);
        this.commentsVisibility = true;
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
            this.errormsg = 'Ops, photo not found.';
            break;
          case 500:
            this.errormsg = 'Ops, there was an internal problem with the server.';
            break;
          default:
            this.errormsg = e.toString();
        }
      }
      this.loading = false;
      this.commentText = '';
    },
    async uncommentPhoto(commentId) {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.delete('/photos/' + this.photo.id + '/comments/' + commentId);
        this.photo.comments.comments = this.photo.comments.comments.filter(function(comment) {
          return commentId != comment.id;
        });
        this.photo.comments.count--;
        this.commentsVisibility = this.photo.comments.count > 0 ? this.commentVisibility : false;
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
            this.errormsg = 'Ops, photo or comment not found.';
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
    async deletePhoto() {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.delete('/photos/' + this.photo.id);
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
            this.errormsg = 'Ops, uphoto not found';
            break;
          case 500:
            this.errormsg = 'Ops, there was an internal problem with the server.';
            break;
          default:
            this.errormsg = e.toString();
        }
      }
      this.loading = false;
      this.$router.go(-1);
    },
  },
  mounted() {
    this.checkLike();
  },
};
</script>

<template>
  <ErrorMsg v-if="this.errormsg" :msg="this.errormsg" />
      <div v-if="!this.errormsg">
        <div class="post-profile">
          <div class="user-info-icon">
            <div class="not-image">
              <router-link :to="'/users/' + this.photo.owner.id" class="username">
                {{ this.photo.owner.username }}
              </router-link>
            </div>
          </div>
          <button class="delete-btn" v-if="this.photo.owner.id == this.reqId" @click="this.deletePhoto()">
            <DeleteIcon />
          </button>
        </div>
        <div class="post-media">
          <div class="post-image-ctn">
            <img :src="this.photoBaseUrl+this.photo.photo_url" id="photo" @dblclick="this.likeAction()">
          </div>
        </div>
        <div class="post-desc">
          <div class="post-actions">
            <div class="actions-left">
              <button @click="likeAction()">
                <LikeIcon :active="this.hasLike" />
              </button>
              <button @click="showComments(this.photo.id)">
                <CommentIcon />
              </button>
            </div>
          </div>
          <div class="likes-count">
            <router-link :to="'/photos/' + this.photo.id + '/likes'" v-bind:class="{disabled: this.photo.likes.count == 0}" style="color: black; text-decoration: none;">
              {{ this.photo.likes.count }} likes
            </router-link>
          </div>
          <div class="">
            <button class="clear-btn" @click="showComments(this.photo.id)">
              View all {{ this.photo.comments.count }} comments
            </button>
            <div v-if="this.commentsVisibility" class="comments">
              <div class="comment" v-for="comment in this.photo.comments.comments" :key="comment">
                <div class="comment-content">
                  <div class="comment-header">
                    <div>
                      <span class="username">{{ comment.user.username }}</span>
                      <span class="timestamp">{{ comment.created_datetime }}</span>
                    </div>
                    <button type="submit" class="delete-btn" v-if="comment.user.id == this.reqId" @click="uncommentPhoto(comment.id)">
                      <svg width="24" height="24" fill="none"><path fill-rule="evenodd" clip-rule="evenodd" d="M6 7a1 1 0 0 1 1 1v11a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V8a1 1 0 1 1 2 0v11a3 3 0 0 1-3 3H8a3 3 0 0 1-3-3V8a1 1 0 0 1 1-1z" fill="#000"/><path fill-rule="evenodd" clip-rule="evenodd" d="M10 8a1 1 0 0 1 1 1v8a1 1 0 1 1-2 0V9a1 1 0 0 1 1-1zM14 8a1 1 0 0 1 1 1v8a1 1 0 1 1-2 0V9a1 1 0 0 1 1-1zM4 5a1 1 0 0 1 1-1h14a1 1 0 1 1 0 2H5a1 1 0 0 1-1-1zM8 3a1 1 0 0 1 1-1h6a1 1 0 1 1 0 2H9a1 1 0 0 1-1-1z" fill="#000"/></svg>
                    </button>
                  </div>
                  <p class="comment-text" style="text-align: left;">{{ comment.text }}</p>
                </div>
              </div>
            </div>
          </div>
          <div class="post-time">
            {{ this.photo.created_datetime }}
          </div>
        </div>
        <div class="post-comment">
          <EmojiIcon class="clear-btn"/>
          <input v-model.trim="commentText" placeholder="Add a comment..." type="text">
          <button :class="{'clear-btn': true, 'post-btn': true, fade: commentText.length < 3}" @click="commentPhoto()">
            Post
          </button>
        </div>
      </div>
</template>


<style scoped>
  .post-profile {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 12px;
  }
  .post-media {
    width: 613px;
  }
  .post-desc {
    padding: 0 12px;
  }
  .post-desc > * {
    margin: 4px 0;
    font-size: 0.875rem;
  }
  .post-actions {
    display: flex;
    justify-content: space-between;
    margin: 12px 0;
    margin-top: 16px;
  }
  .post-actions button {
    outline: none;
    background: none;
    border: none;
    padding: 0;
  }
  .post-actions .actions-left button {
    margin-right: 12px;
  }
  .post-actions .actions-right button {
    margin-right: 0;
  }
  .post-desc > .clear-btn {
    display: block;
  }
  .likes-count {
    font-weight: 600;
  }
  .post-time {
    font-size: 0.65rem;
    color: var(--grey-text);
    text-transform: uppercase;
    margin-top: 8px;
  }
  .post-comment {
    border-top: 1px solid var(--border-color);
    padding: 0 12px;
    height: 55px;
    margin-top: 20px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .post-comment > input {
    width: 82%;
    height: 90%;
    outline: none;
    border: none;
    font-size: 0.875rem;
  }
  .post-comment > input::placeholder {
    font-size: 0.875rem;
  }
  .post-btn {
    padding: 0 8px;
    color: #0095F6;
    transition: 0.3s;
    font-weight: 600;
  }

  #photo {
    width: 100%;
  }

  .user-info-icon {
    display: flex;
    align-items: center;
  }
  .not-image {
    margin-left: 12px;
    font-weight: 500;
    font-size: 0.875rem;
  }
  .username {
    text-decoration: none;
    color: #000000;
    font-weight: 600;
  }
  .username:hover {
    text-decoration: none;
  }
.comments {
  min-height: 50px;
  max-height: 200px;
  overflow-y: auto;
  width: 95%;
  margin: 0 auto;
  margin-top: 5px;
  padding: 2px;
}

.comment {
  display: flex;
  margin-bottom: 5px;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 3px;
}

.username {
  font-weight: bold;
  margin-right: 0px;
}

.timestamp {
  color: #999;
  font-size: 12px;
  margin: 40px;
}

.comment-text {
  color: #333;
  font-size: 14px;
  text-align: left;
  font-weight: 100;
}

.delete-btn {
  background: transparent;
  border: none;
}

.disabled {
  pointer-events: none;
}
</style>
