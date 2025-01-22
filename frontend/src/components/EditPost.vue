<template>
  <div class="edit-post">
    <NavBar :user="getUserWithId($store.state.userId)"></NavBar>
    <div v-if="loading">Loading...</div>
    <div v-else-if="error">{{ error }}</div>
    <div v-else class="edit-form">
      <h2>Edit Post</h2>
      
      <label for="contentText">Content:</label>
      <textarea
        id="contentText"
        v-model="postContent"
        rows="5"
        placeholder="Edit your post..."
      ></textarea>

      <label for="imageURL">Image URL:</label>
      <input 
        id="imageURL" 
        v-model="imageURL" 
        type="url" 
        placeholder="Enter image URL"
      >

      <!-- Image Preview -->
      <div v-if="imageURL" class="image-preview-container">
        <img 
          :src="imageURL" 
          alt="Preview" 
          @error="handleImageError"
          class="image-preview"
        >
        <button 
          type="button" 
          class="clear-image" 
          @click="clearImage"
        >
          Clear Image
        </button>
      </div>


      <div class="category-section">
        <label for="category">Category:</label>
        <select v-if="!showNewCategoryForm" id="category" v-model="selectedCategory" class="category-select">
          <option value="">Select a category</option>
          <option v-for="category in categories" :key="category.idCategory" :value="category.idCategory">
            {{ category.name }}
          </option>
        </select>
        
        <button type="button" @click="toggleNewCategory" class="toggle-category-btn">
          {{ showNewCategoryForm ? 'Select Existing Category' : 'Add New Category' }}
        </button>

        <!-- New Category Form -->
        <div v-if="showNewCategoryForm" class="new-category-form">
          <input 
            v-model="newCategory.name" 
            placeholder="Category Name"
            class="category-input"
          />
          <textarea 
            v-model="newCategory.description" 
            placeholder="Category Description"
            class="category-description"
          ></textarea>
          <button type="button" @click="addNewCategory" class="add-category-btn">
            Create Category
          </button>
        </div>
      </div>

      <div class="buttons">
        <button @click="savePost" class="save-btn">Save Changes</button>
        <button @click="cancel" class="cancel-btn">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import NavBar from "./NavBar.vue";

export default {
  name: "EditPost",
  components: {
    NavBar,
  },

  data() {
    return {
      postContent: "",
      imageURL: "",
      isValidImage: true,
      loading: true,
      error: null,
      users: [],
      baseUrl: "http://localhost:5050",
      categories: [],
      selectedCategory: "",
      showNewCategoryForm: false,
      newCategory: {
        name: '',
        description: ''
      }
    };
  },

  async created() {
    if (this.$store.state.userId == -1) {
      this.$router.push({ path: "/" });
    }
    try {
      // Get post ID from route params
      const postId = this.$route.params.id;
      console.log(postId);
      // Fetch post data
      const response = await axios.get(`${this.baseUrl}/post`, {
        params: {
          id: postId,
        },
      });
      this.postContent = response.data.content_text;
      this.imageURL = response.data.imageURL || "";
      console.log(response.data);
      await this.fetchUsers();
      await this.fetchCategories();
      this.selectedCategory = response.data.categoryID;
    } catch (error) {
      console.error("Error fetching post data:", error);
      this.error = "Failed to load post";
    } finally {
      this.loading = false;
    }
  },

  methods: {
    handleImageError(e) {
      e.target.classList.add('image-error');
      this.isValidImage = false;
    },

    clearImage() {
      this.imageURL = "";
      this.isValidImage = true;
    },

    async fetchUsers() {
      try {
        const response = await axios.get(`${this.baseUrl}/users`);
        this.users = response.data;
      } catch (error) {
        console.error("Error fetching users:", error);
      }
    },

    async fetchCategories() {
      try {
        const response = await axios.get(`${this.baseUrl}/categories`);
        this.categories = response.data;
      } catch (error) {
        console.error("Error fetching categories:", error);
      }
    },

    toggleNewCategory() {
      this.showNewCategoryForm = !this.showNewCategoryForm;
    },

    async addNewCategory() {
      try {
        const response = await axios.post(`${this.baseUrl}/addCategory`, {
          name: this.newCategory.name,
          description: this.newCategory.description
        });

        if (response.status === 200) {
          await this.fetchCategories();
          this.selectedCategory = response.data.categoryId;
          this.showNewCategoryForm = false;
          this.newCategory = { name: '', description: '' };
        }
      } catch (error) {
        console.error('Error adding category:', error);
        this.error = "Failed to add category";
      }
    },

    getUserWithId(id) {
      return this.users.find((user) => user.idUser === id);
    },

    async savePost() {
      try {
        const response = await axios.put(`${this.baseUrl}/editPost`, {
          postID: String(this.$route.params.id),
          contentText: this.postContent,
          imageURL: this.imageURL,
          categoryID: String(this.selectedCategory)
        });

        if (response.status === 200) {
          this.$router.push("/");
        }
      } catch (error) {
        console.error("Error updating post:", error);
        this.error = "Failed to update post";
      }
    },

    cancel() {
      this.$router.push("/");
    },
  },
};
</script>

<style scoped>
.edit-post {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.edit-form {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

label {
  display: block;
  font-weight: bold;
  margin-top: 1em;
}

textarea, input {
  width: 100%;
  padding: 0.5em;
  margin-top: 0.5em;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.image-preview-container {
  margin: 1rem 0;
  position: relative;
  max-width: 100%;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.image-preview {
  width: 100%;
  max-height: 300px;
  object-fit: cover;
  display: block;
}

.image-preview.image-error {
  display: none;
}

.clear-image {
  position: absolute;
  top: 10px;
  right: 10px;
  background: rgba(0,0,0,0.6);
  color: white;
  border: none;
  border-radius: 4px;
  padding: 5px 10px;
  cursor: pointer;
}

.buttons {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 20px;
}

.save-btn {
  background-color: #4caf50;
  color: white;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.cancel-btn {
  background-color: #f44336;
  color: white;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.save-btn:hover, .cancel-btn:hover {
  opacity: 0.9;
}

.category-section {
  margin: 1rem 0;
}

.category-select {
  width: 100%;
  padding: 0.5rem;
  margin-bottom: 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.toggle-category-btn {
  margin-top: 0.5rem;
  padding: 0.5rem;
  background: #4caf50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.new-category-form {
  margin-top: 1rem;
  padding: 1rem;
  background: #f8f8f8;
  border-radius: 4px;
}

.category-input,
.category-description {
  width: 100%;
  padding: 0.5rem;
  margin-bottom: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.category-description {
  min-height: 80px;
}

.add-category-btn {
  background: #4CAF50;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
}
</style>
