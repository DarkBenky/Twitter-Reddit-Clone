<template>
    <div class="category-list">
        <NavBar :user="user"></NavBar>
        <div class="categories-container">
            <h2>Categories</h2>
            <div class="categories-grid">
                <div v-for="category in categories" 
                     :key="category.idCategory" 
                     class="category-card"
                     @click="goToCategoryPosts(category.name)">
                    <h3>{{ category.name }}</h3>
                    <p>{{ category.description }}</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import NavBar from './NavBar.vue';

export default {
    name: 'CategoryList',
    components: {
        NavBar
    },
    data() {
        return {
            categories: [],
            baseUrl: "http://localhost:5050"
        };
    },
    async created() {
        await this.fetchCategories();
    },
    methods: {
        async fetchCategories() {
            try {
                const response = await axios.get(`${this.baseUrl}/categories`);
                this.categories = response.data;
            } catch (error) {
                console.error('Error fetching categories:', error);
            }
        },
        goToCategoryPosts(categoryName) {
            this.$router.push(`/category/${categoryName}`);
        }
    }
};
</script>

<style scoped>
.category-list {
    padding: 20px;
}

.categories-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
}

.categories-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 20px;
    margin-top: 20px;
}

.category-card {
    background: white;
    border-radius: 8px;
    padding: 20px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    transition: transform 0.2s;
    cursor: pointer; /* Add this to show it's clickable */
}

.category-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0,0,0,0.2);
}

.category-card h3 {
    margin: 0 0 10px 0;
    color: #333;
}

.category-card p {
    margin: 0;
    color: #666;
    font-size: 0.9em;
}
</style>