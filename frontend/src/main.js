import { createApp } from 'vue'
import App from './App.vue'
import userProfileSinglePage from './components/userProfileSinglePage.vue'
import { createWebHistory, createRouter } from 'vue-router'
import PostFeed from './components/PostFeed.vue'
import { createStore } from 'vuex'
import AddPost from './components/AddPost.vue'
import axios from 'axios';
import SinglePost from './components/SinglePost.vue';
import EditPost from './components/EditPost.vue';
import LoginPage from './components/LoginPage.vue';
import UserConversation from './components/UserConversation.vue'
import ConversationsList from './components/ConversationsList.vue'
import RegisterPage from './components/RegisterPage.vue'
import CategoryList from './components/CategoryList.vue'
import CategoryPosts from './components/CategoryPosts.vue'

axios.defaults.baseURL = 'http://localhost:5050';
axios.defaults.headers.common['Content-Type'] = 'application/json';
axios.defaults.withCredentials = true;

const routes = [
    { path: '/', component: PostFeed },
    { path: '/user', component: userProfileSinglePage },
    { path: '/post/:id', component: SinglePost }, 
    { path: '/add_post', component: AddPost },
    { path: '/edit/:id', component: EditPost },
    { path: '/login', component: LoginPage },
    { path: '/dm/:id',  component: UserConversation },
    { path: '/dms',  component: ConversationsList },
    { path: '/register', component: RegisterPage },
    {path : '/categories', component: CategoryList},
    { path: '/category/:category', component: CategoryPosts },
    { path: '/:notFound(.*)', redirect: '/' },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

const store = createStore({
    state() {
        return {
            userId: -1,
            currentUser: null,
            users : []
        }
    },
    mutations: {
        setUserId(state, userId) {
            state.userId = userId
        },
        setCurrentUser(state, user) {
            state.currentUser = user
        },
        logout(state) {
            state.currentUser = null
            state.userId = -1
        }
    },
    getters: {
        getCurrentUser: (state) => state.currentUser
    }
})

const app = createApp(App)
app.use(router)
app.use(store)

// enable CORS
app.config.globalProperties.$http = axios


// Provide the store to all components
app.config.globalProperties.$store = store

app.mount('#app')