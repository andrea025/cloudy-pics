import {createRouter, createWebHashHistory} from 'vue-router';
import HomeView from '../views/HomeView.vue';
import LoginView from '../views/LoginView.vue';
import ProfileView from '../views/ProfileView.vue';
import UpdateUsernameView from '../views/UpdateUsernameView.vue';
import FollowersView from '../views/FollowersView.vue';
import FollowingView from '../views/FollowingView.vue';
import BannedView from '../views/BannedView.vue';
import SearchUserView from '../views/SearchUserView.vue';
import PhotoLikesView from '../views/PhotoLikesView.vue';
import PhotoView from '../views/PhotoView.vue';
import NotFoundView from '../views/NotFoundView.vue';


const routes = [
  {
    path: '/',
    redirect: (to) => {
      return '/login';
    },
  },
  {
    path: '/home',
    name: 'home',
    component: HomeView,
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView,
  },
  {
    path: '/users/:id',
    name: 'profile',
    component: ProfileView,
  },
  {
    path: '/users',
    name: 'search-user',
    component: SearchUserView,
  },
  {
    path: '/users/:id/followers',
    name: 'profile-followers',
    component: FollowersView,
  },
  {
    path: '/users/:id/following',
    name: 'profile-following',
    component: FollowingView,
  },
  {
    path: '/users/:id/banned',
    name: 'profile-banned',
    component: BannedView,
  },
  {
    path: '/profile/username',
    name: 'profile-username',
    component: UpdateUsernameView,
  },
  {
    path: '/photos/:id',
    name: 'photo',
    component: PhotoView,
  },
  {
    path: '/photos/:id/likes',
    name: 'photo-likes',
    component: PhotoLikesView,
  },
  {
    path: '/:catchAll(.*)',
    name: 'not-found',
    component: NotFoundView,
  },
];

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;
