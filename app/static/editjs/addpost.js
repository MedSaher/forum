// Get modal elements
const modal = document.getElementById("registrationModal");
const modalMessage = document.getElementById("modalMessage");
const closeModalBtn = document.getElementById("closeModal");
const modalOkBtn = document.getElementById("modalOkBtn");

// Function to show the modal
function showModal(message) {
    modalMessage.textContent = message; // Set message
    modal.style.display = "block"; // Show modal
}

// Close modal function
function closeModal() {
    modal.style.display = "none"; // Hide modal
}

// Close modal on button click
closeModalBtn.addEventListener("click", closeModal);
modalOkBtn.addEventListener("click", closeModal);

// Add event listener to form submission
const post_form = document.getElementById("postForm");
post_form.addEventListener("submit", async function (e) {
    e.preventDefault();
    const data = new FormData(post_form);

    try {
        const response = await fetch("http://localhost:8080/add_post", {
            method: "POST",
            body: data,
        });

        if (!response.ok) {

            // Handle errors based on response status
            if (response.status === 404) {
                sessionStorage.setItem("errorStatus", response.status);
                window.location.href = "/errors.html";
            } else {
                window.location.href = "/login"
            }
        } else {
            const post = await response.json(); // Get post details
            console.log("the new created post: ", post);
            
            // Dynamically create and add the new post to the page
            const postId = document.createElement("input")
            postId.type = "hidden"
            postId.id = "postId"
            postId.value = post.id

            const card = document.createElement("article");
            card.classList.add("post");
            card.setAttribute("data-category", post.categoryName);
            card.setAttribute("post_id", post.id);

            const postHeader = document.createElement("div");
            postHeader.classList.add("post-header");

            const author = document.createElement("span");
            author.classList.add("post-author");
            author.textContent = `${post.authorFirstName} ${post.authorLastName}`

            const category = document.createElement("span");
            category.classList.add("post-category");
            category.textContent = post.categoryName;
            console.log(category);
            
            const timestamp = document.createElement("span");
            timestamp.classList.add("post-time");
            timestamp.textContent = new Date(post.time).toLocaleString();
            postHeader.appendChild(author);
            postHeader.appendChild(category);
            postHeader.appendChild(timestamp);
            postHeader.appendChild(postId);

            const title = document.createElement("h2");
            title.classList.add("post-title");
            title.textContent = post.title;

            const text = document.createElement("p");
            text.classList.add("post-content");
            text.textContent = post.content;

            const postFooter = document.createElement("div");
            postFooter.classList.add("post-footer");

            const likeBtn = createLikeButton(post);
            const dislikeBtn = createDislikeButton(post);
            const cmtBtn = createCommentButton(post);

            postFooter.appendChild(likeBtn);
            postFooter.appendChild(dislikeBtn);
            postFooter.appendChild(cmtBtn);

            card.appendChild(postHeader);
            card.appendChild(title);
            card.appendChild(text);
            card.appendChild(postFooter);
            main_content.appendChild(card);
            main_content.insertBefore(card, main_content.children[1]);
            post_form.reset()
        }
    } catch (error) {
        console.error("Error:", error.message);
        showModal("An error occurred: Creating a post.");
    }
});
