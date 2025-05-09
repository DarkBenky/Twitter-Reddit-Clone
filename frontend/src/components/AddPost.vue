<template>
  <div>
    <NavBar :user="user"></NavBar>
    <div class="add-post">
      <h2>Add a Post</h2>
      <form @submit.prevent="submitPost">
        <label for="contentText">Content:</label>
        <textarea id="contentText" v-model="contentText" required></textarea>

        <label for="imageURL">Image URL:</label>
        <input id="imageURL" v-model="imageURL" type="url" placeholder="Enter image URL">

        <!-- Image Preview -->
        <div v-if="imageURL" class="image-preview-container">
          <img :src="imageURL" alt="Preview" @error="handleImageError" class="image-preview">
          <button type="button" class="clear-image" @click="clearImage">
            Clear Image
          </button>
        </div>

        <!-- add secondary images -->
        <div class="secondary-images-section">
          <label for="secondaryImageURL">Secondary Images:</label>
          <div class="secondary-image-input">
            <input id="secondaryImageURL" v-model="secondaryImageURL" type="url" placeholder="Enter image URL">
            <button type="button" class="add-image-btn" @click="addSecondaryImage">
              Add Image
            </button>
          </div>

          <!-- preview of secondary images -->
          <div class="secondary-images-grid" v-if="secondaryImages.length > 0">
            <div v-for="(secondaryImage, index) in secondaryImages" :key="index" class="secondary-image-container">
              <img :src="secondaryImage" alt="Secondary image" @error="handleSecondaryImageError($event, index)" class="secondary-image-preview">
              <button type="button" class="remove-secondary-image" @click="removeSecondaryImage(index)">
                ✖
              </button>
            </div>
          </div>
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
            <input v-model="newCategory.name" placeholder="Category Name" class="category-input" />
            <textarea v-model="newCategory.description" placeholder="Category Description"
              class="category-description"></textarea>
            <button type="button" @click="addNewCategory" class="add-category-btn">
              Create Category
            </button>
          </div>
        </div>

        <button type="submit">Add Post</button>

        <p v-if="message" :class="{ success: success, error: !success }">
          {{ message }}
        </p>
      </form>
    </div>
  </div>
</template>

<script>
// import axios from "axios";
import NavBar from "./NavBar.vue";
import api from "../services/api.js";
import config from "../config.js";

export default {
  name: "AddPost",

  components: {
    NavBar,
  },

  data() {
    return {
      user: {},
      baseUrl: config.apiUrl,
      categories: [],
      contentText: "",
      selectedCategory: "",
      message: "",
      success: false,
      imageURL: "",
      isValidImage: false,
      showNewCategoryForm: false,
      newCategory: {
        name: '',
        description: ''
      },
      secondaryImages: [],
      secondaryImageURL: "",
    };
  },

  async created() {
    if (this.$store.state.userId == -1) {
      this.$router.push({ path: "/" });
    }
    this.fetchCategories();
  },

  methods: {
    addSecondaryImage() {
      if (this.secondaryImageURL) {
        this.secondaryImages.push(this.secondaryImageURL);
        this.secondaryImageURL = "";
      }
    },

    removeSecondaryImage(index) {
      this.secondaryImages.splice(index, 1);
    },

    async fetchCategories() {
      try {
        const response = await api.get(`${this.baseUrl}/categories`);
        this.categories = response.data;
      } catch (error) {
        console.error("Error fetching categories:", error);
      }
    },

    handleImageError(e) {
      e.target.classList.add('image-error');
      this.isValidImage = false;
    },

    clearImage() {
      this.imageURL = "";
      this.isValidImage = false;
    },

    async submitPost() {
      try {
        const response = await api.post(`${this.baseUrl}/addPost`, {
          content_text: this.contentText.toString(),
          categoryID: this.selectedCategory.toString(),
          userID: this.$store.state.userId.toString(),
          imageURL: this.imageURL,
          secondaryImages: this.secondaryImages
        });
        console.log(response.data);
        console.log("number of secondary images: ", this.secondaryImages.length);
        this.message = "Post added successfully!";
        this.success = true;
        this.contentText = "";
        this.selectedCategory = "";
        this.imageURL = "";

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

    toggleNewCategory() {
      this.showNewCategoryForm = !this.showNewCategoryForm;
    },

    async addNewCategory() {
      try {
        const response = await api.post(`${this.baseUrl}/addCategory`, {
          name: this.newCategory.name,
          description: this.newCategory.description
        });

        if (response.status === 200) {
          // Update categories list
          await this.fetchCategories();
          // Select the new category
          this.selectedCategory = response.data.categoryId;
          // Reset form
          this.showNewCategoryForm = false;
          this.newCategory = { name: '', description: '' };
          // Show success message
          this.message = "Category added successfully";
          this.success = true;
        }
      } catch (error) {
        console.error('Error adding category:', error);
        this.message = error.response?.data?.error || "Failed to add category";
        this.success = false;
      }
    }
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

.image-preview-container {
  margin: 1rem 0;
  position: relative;
  max-width: 100%;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
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
  background: rgba(0, 0, 0, 0.6);
  color: white;
  border: none;
  border-radius: 4px;
  padding: 5px 10px;
  cursor: pointer;
}

.clear-image:hover {
  background: rgba(0, 0, 0, 0.8);
}

.category-section {
  margin: 1rem 0;
}

.toggle-category-btn {
  margin-top: 0.5rem;
  padding: 0.5rem;
  background: #4caf50;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
}

.new-category-form {
  margin-top: 1rem;
  padding: 1rem;
  background: #f8f8f8;
  border-radius: 4px;
}

.category-input {
  width: 100%;
  padding: 0.5rem;
  margin-bottom: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.category-description {
  width: 100%;
  padding: 0.5rem;
  margin-bottom: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
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

.add-category-btn:hover {
  background: #45a049;
}

/* Secondary Images Section */
.secondary-images-section {
  margin: 1.5rem 0;
}

.secondary-image-input {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.secondary-image-input input {
  flex: 1;
}

.add-image-btn {
  margin-top: 0;
  white-space: nowrap;
  background-color: #4361ee;
}

.add-image-btn:hover {
  background-color: #3a56d4;
}

.secondary-images-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 1rem;
  margin-top: 1rem;
}

.secondary-image-container {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  aspect-ratio: 1/1;
}

.secondary-image-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.secondary-image-preview.image-error {
  display: none;
}

.remove-secondary-image {
  position: absolute;
  top: 5px;
  right: 5px;
  background: rgba(255, 0, 0, 0.7);
  color: white;
  border: none;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  margin: 0;
  cursor: pointer;
}

.remove-secondary-image:hover {
  background: rgba(255, 0, 0, 0.9);
}
</style>
