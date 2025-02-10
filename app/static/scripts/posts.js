document.addEventListener("DOMContentLoaded", async function() {
    var logged = false;
    let userProfile = document.getElementById("user-profile");

    (async function fetchUserProfile() {
        try {
            const response = await fetch('http://localhost:8080/profile');
            const user = await response.json();

            // Create user container
            const userInfo = document.createElement("div");
            userInfo.classList.add("user-info");

            if (user) {
                logged = true;

                // Profile Image
                const profileImg = document.createElement("img");
                profileImg.src = "./app/uploads/" + user.profile_pic;
                profileImg.alt = "Profile Picture";
                profileImg.classList.add("profile-pic");

                // User Name
                const nameSpan = document.createElement("span");
                nameSpan.classList.add("user-name");
                nameSpan.textContent = `${user.first_name} ${user.last_name}`;

                // Logout Button
                const logoutBtn = document.createElement("button");
                logoutBtn.classList.add("logout-btn");
                logoutBtn.textContent = "Logout";
                logoutBtn.onclick = async function() {
                    try {
                        const logoutResponse = await fetch('http://localhost:8080/logout', {
                            method: 'POST',
                            credentials: 'include'
                        });

                        if (logoutResponse.ok) {
                            window.location.reload();
                        } else {
                            console.error("Logout failed");
                        }
                    } catch (error) {
                        console.error("Logout error:", error);
                    }
                };

                // Append elements
                userInfo.appendChild(profileImg);
                userInfo.appendChild(nameSpan);
                userInfo.appendChild(logoutBtn);
            } else {
                // Login Button for non-logged-in users
                const loginBtn = document.createElement("button");
                loginBtn.classList.add("login-btn");
                loginBtn.textContent = "Login";
                loginBtn.onclick = function() {
                    window.location.href = '/register';
                };

                userInfo.appendChild(loginBtn);
            }

            userProfile.appendChild(userInfo);

        } catch (error) {
            console.error('Error fetching profile:', error);
        }
    })();

    // Fetch all existing categories:
    const side_bar = document.getElementById("navbar");

    try {
        const response = await fetch('http://localhost:8080/all_categories');
        const categories = await response.json();

        categories.forEach(category => {
            let category_item = document.createElement("button");
            category_item.classList.add("category");
            category_item.textContent = category.name;

            // Add an event listener for filtering posts
            category_item.addEventListener("click", () => filterPosts(category.name));

            // Append each category button to the navbar
            side_bar.appendChild(category_item);
        });

    } catch (error) {
        console.error('Error:', error);
    }


    // Fetch all existing posts:
    const main_content = document.getElementById("main-content");

    try {
        const response = await fetch('http://localhost:8080/all_posts');
        const posts = await response.json();

        posts.forEach(post => {
            const card = document.createElement("article");
            card.classList.add("post");
            card.setAttribute("data-category", post.categoryName);

            // Post Header (Author, Category, Time)
            const postHeader = document.createElement("div");
            postHeader.classList.add("post-header");

            const author = document.createElement("span");
            author.classList.add("post-author");
            author.textContent = `${post.authorFirstName} ${post.authorLastName}`;

            const category = document.createElement("span");
            category.classList.add("post-category");
            category.textContent = post.categoryName;

            const timestamp = document.createElement("span");
            timestamp.classList.add("post-time");
            timestamp.textContent = new Date(post.time).toLocaleString();

            // Append header elements
            postHeader.appendChild(author);
            postHeader.appendChild(category);
            postHeader.appendChild(timestamp);

            // Title
            const title = document.createElement("h2");
            title.classList.add("post-title");
            title.textContent = post.title;

            // Content
            const text = document.createElement("p");
            text.classList.add("post-content");
            text.textContent = post.content;

            // Footer (Likes, Dislikes, Comments)
            const postFooter = document.createElement("div");
            postFooter.classList.add("post-footer");

            // Like button
            const likeBtn = document.createElement("span");
            likeBtn.classList.add("like-btn");

            const likeIcon = document.createElement("i");
            likeIcon.classList.add("fa-solid", "fa-thumbs-up");

            const likeCounter = document.createElement("span");
            likeCounter.classList.add("like-count");
            likeCounter.textContent = post.likeCount;

            likeBtn.appendChild(likeIcon);
            likeBtn.appendChild(likeCounter);

            // Dislike button
            const dislikeBtn = document.createElement("span");
            dislikeBtn.classList.add("dislike-btn");

            const dislikeIcon = document.createElement("i");
            dislikeIcon.classList.add("fa-solid", "fa-thumbs-down");

            const dislikeCounter = document.createElement("span");
            dislikeCounter.classList.add("dislike-count");
            dislikeCounter.textContent = post.dislikeCount;

            dislikeBtn.appendChild(dislikeIcon);
            dislikeBtn.appendChild(dislikeCounter);

            // Comments button
            const cmtBtn = document.createElement("span");
            cmtBtn.classList.add("comment-btn");

            const cmtIcon = document.createElement("i");
            cmtIcon.classList.add("fa-solid", "fa-comment");

            cmtBtn.appendChild(cmtIcon);

            // Attach Event Listener **INSIDE THE LOOP**
            likeBtn.addEventListener("click", async() => {
                if (!logged) {
                    alert("You need to log in to like posts!");
                    return;
                }

                try {
                    const response = await fetch(`http://localhost:8080/vote_for_post`, {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json"
                        },
                        body: JSON.stringify({ postId: post.id, value: 1 })
                    });

                    if (response.ok) {
                        window.location.reload();
                    } else {
                        console.error("Failed to like post");
                    }
                } catch (error) {
                    console.error("Error:", error);
                }
            });

            // Handle the dislike btn behavior:
            dislikeBtn.addEventListener("click", async() => {
                try {
                    const response = await fetch(`http://localhost:8080/vote_for_post`, { // Use backticks
                        method: "POST", // Set HTTP method to POST
                        headers: {
                            "Content-Type": "application/json"
                        },
                        body: JSON.stringify({ postId: post.id, value: -1 }) // Send data if required by backend
                    });
                    if (response.ok) {
                        window.location.reload(); // Refresh page to update dislike count
                    }
                } catch (error) {
                    console.error("Error:", error);
                }
            });

            // handle comments

            cmtBtn.addEventListener("click", async function() {
                closeBtn = document.getElementById("close-btn")
                modalPostTitle.textContent = post.title;
                modalPostContent.textContent = post.content;
                // commentBox.value = "";
                // commentsList.innerHTML = "";
                commentModal.style.display = "block";
                // await loadComments(postId);
                closeBtn.addEventListener("click", () => commentModal.style.display = "none");
                window.addEventListener("click", event => {
                    if (event.target === commentModal) commentModal.style.display = "none";
                });
                submitCommentBtn.onclick = () => submitComment(postId);
            });
            // Append all elements to footer
            postFooter.appendChild(likeBtn);
            postFooter.appendChild(dislikeBtn);
            postFooter.appendChild(cmtBtn);

            // Append everything to the post card
            card.appendChild(postHeader);
            card.appendChild(title);
            card.appendChild(text);
            card.appendChild(postFooter);

            // Append post to main content
            main_content.appendChild(card);
        });

    } catch (error) {
        console.error('Error:', error);
    }

});


// Create a function to filter based on the post category:
function filterPosts(selected) {
    const cards = document.getElementsByClassName("card")
    console.log(cards);
    Array.from(cards).forEach(card => {
        const itemCategory = card.getAttribute("data-category");
        if (itemCategory === selected) {
            card.style.display = "block";
        } else {
            card.style.display = "none";
        }
    });
}

// comments handler