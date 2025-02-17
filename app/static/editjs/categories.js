async function fetchCategories() {
    const select_container = document.getElementById("select_categories");
    const side_bar = document.getElementById("navbar");

    try {
        const response = await fetch("http://localhost:8080/all_categories");
        const categories = await response.json();

        categories.forEach((category) => {
            const checkboxWrapper = document.createElement("div");
            checkboxWrapper.classList.add("form-check");

            const checkbox = document.createElement("input");
            checkbox.type = "checkbox";
            checkbox.name = "chosen_categories[]";
            checkbox.value = category.name;
            checkbox.id = `category_${category.name}`;
            checkbox.classList.add("form-check-input");

            const label = document.createElement("label");
            label.htmlFor = `category_${category.name}`;
            label.textContent = category.name;
            label.classList.add("form-check-label");

            checkboxWrapper.appendChild(checkbox);
            checkboxWrapper.appendChild(label);
            select_container.appendChild(checkboxWrapper);

            let category_item = document.createElement("button");
            category_item.classList.add("category");
            category_item.textContent = category.name;
            category_item.addEventListener("click", () => filterPosts(category.name));
            side_bar.appendChild(category_item);
        });
        // Ensure "Liked Posts" is only added AFTER the `logged` variable is set
        if (logged) {

            // handle liked posts filter 
            let likedPostsBtn = document.createElement("button");
            likedPostsBtn.classList.add("category");
            likedPostsBtn.textContent = "Liked Posts";

            // Add a filter for owned posts:
            let ownedPost = document.createElement("button");
            ownedPost.classList.add("category");
            ownedPost.textContent = "created posts";

            likedPostsBtn.addEventListener("click", async() => {
                try {
                    const response = await fetch("http://localhost:8080/liked");
                    const likedPostIds = await response.json();

                    // Get all posts currently displayed
                    const allPosts = document.querySelectorAll(".post");

                    allPosts.forEach((post) => {
                        const postId = post.getAttribute("post_id");
                        // Check if post ID exists in the likedPostIds object
                        if (likedPostIds[postId]) {
                            post.style.display = "block"; // Show liked posts
                        } else {
                            post.style.display = "none"; // Hide unliked posts
                        }
                    });
                } catch (error) {
                    console.error("Error fetching liked posts:", error);
                }
            });

            ownedPost.addEventListener("click", async() => {
                try {
                    const response = await fetch("http://localhost:8080/owned");
                    const ownedPosts = await response.json();

                    // Get all posts currently displayed
                    const allPosts = document.querySelectorAll(".post");

                    allPosts.forEach((post) => {
                        const postId = post.getAttribute("post_id");

                        // Check if post ID exists in the likedPostIds object
                        if (ownedPosts[postId]) {
                            post.style.display = "block"; // Show liked posts
                        } else {
                            post.style.display = "none"; // Hide unliked posts
                        }
                    });
                    if (response.status === 404) {
                        // Store the status in sessionStorage
                        sessionStorage.setItem("errorStatus", response.status);
                
                        // Redirect to error page
                        window.location.href = "/error.html";
                      }
                } catch (error) {
                    console.error("Error fetching liked posts:", error);
                }
            });

            // Append the "Liked Posts" button after categories
            side_bar.appendChild(likedPostsBtn);
            side_bar.appendChild(ownedPost);
        }

        if (response.status === 404) {
            // Store the status in sessionStorage
            sessionStorage.setItem("errorStatus", response.status);
    
            // Redirect to error page
            window.location.href = "/error.html";
          }

    } catch (error) {
        console.error("Error fetching categories:", error);
    }
}