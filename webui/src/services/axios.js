import axios from 'axios';

const instance = axios.create({
  baseURL: __API_URL__,
  timeout: 1000 * 5,
});

instance.interceptors.request.use(function(req) {
  req.headers.authorization = 'Bearer ' + localStorage.getItem('id');
  return req;
});

export default instance;
