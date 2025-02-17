let page = 1; // Start from page 1
const limit = 8; // Number of posts per page
let isFetching = false; // Prevent duplicate requests
const main_content = document.getElementById("main-content");

async function fetchPosts() {
    if (isFetching) return; // Prevent multiple calls at the same time
    isFetching = true;

    try {
        const response = await fetch(`http://localhost:8080/all_posts?page=${page}&limit=${limit}`);
        const posts = await response.json();

        if (posts == null) {
            window.removeEventListener("scroll", handleScroll); // Stop fetching if no more posts
            return;
        }
        posts.forEach(post => {
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
        });


        page++; // Increment page for next fetch

    } catch (error) {
        console.error("Error fetching posts:", error);
    } finally {
        isFetching = false;
    }
}

// Scroll Event Listener for Infinite Scroll
function handleScroll() {
    if (window.innerHeight + window.scrollY >= document.body.offsetHeight - 500) {
        fetchPosts();
    }
}

// Load initial posts and attach scroll event listener
document.addEventListener("DOMContentLoaded", () => {
    fetchPosts();
    window.addEventListener("scroll", handleScroll);
});