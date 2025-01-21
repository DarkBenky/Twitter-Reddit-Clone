<template>
  <div>
    <NavBar :user="user"></NavBar>
    <div class="add-post">
      <h2>Add a Post</h2>
      <form @submit.prevent="submitPost">
        <label for="contentText">Content:</label>
        <textarea id="contentText" v-model="contentText" required></textarea>

        <label for="category">Category:</label>
        <select id="category" v-model="selectedCategory" class="category-select">
          <option value="">Select a category</option>
          <option v-for="category in categories" 
                  :key="category.idCategory" 
                  :value="category.idCategory">
            {{ category.name }}
          </option>
        </select>

        <button type="submit">Add Post</button>

        <p v-if="message" :class="{ success: success, error: !success }">
          {{ message }}
        </p>
      </form>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import NavBar from "./NavBar.vue";

export default {
  name: "AddPost",

  components: {
    NavBar,
  },

  data() {
    return {
      user: {},
      baseUrl: "http://localhost:5050",
      categories: [],
      contentText: "",
      selectedCategory: "",
      message: "",
      success: false,
    };
  },

  async created() {
    if (this.$store.state.userId == -1) {
      this.$router.push({ path: "/" });
    }
    this.fetchCategories();
  },

  methods: {
    async fetchCategories() {
      try {
        const response = await axios.get(`${this.baseUrl}/categories`);
        this.categories = response.data;
      } catch (error) {
        console.error("Error fetching categories:", error);
      }
    },

    async submitPost() {
      try {
        const response = await axios.post(`${this.baseUrl}/addPost`, {
          content_text: this.contentText.toString(),
          categoryID: this.selectedCategory.toString(),
          userID: this.$store.state.userId.toString()
        });
        console.log(response.data);
        this.message = "Post added successfully!";
        this.success = true;
        this.contentText = "";
        this.selectedCategory = "";
        
        setTimeout(() => {
          this.$router.push("/");
        }, 1500);
      } catch (error) {
        this.message = error.response?.data?.error || "Failed to add post";
        this.success = false;
      }
    },

    clearForm() {
      this.contentText = "";
    },
  },
};
</script>

<style scoped>
.add-post {
  max-width: 500px;
  margin: auto;
  padding: 1em;
  background-color: #f9f9f9;
  border-radius: 8px;
}

label {
  display: block;
  font-weight: bold;
  margin-top: 1em;
}

input,
textarea {
  width: 100%;
  padding: 0.5em;
  margin-top: 0.5em;
  border: 1px solid #ccc;
  border-radius: 4px;
}

button {
  margin-top: 1em;
  padding: 0.5em 1em;
  background-color: #4caf50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #45a049;
}

.success {
  color: green;
}

.error {
  color: red;
}

.category-select {
  width: 100%;
  padding: 0.5rem;
  margin-bottom: 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.category-select:focus {
  outline: none;
  border-color: #4CAF50;
  box-shadow: 0 0 0 2px rgba(76, 175, 80, 0.2);
}
</style>
