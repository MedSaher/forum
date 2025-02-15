document.addEventListener("DOMContentLoaded", async function () {
    let logged = false;
    let userProfile = document.getElementById("user-profile");


    // Fill the categories dynamically with checkboxes:
    const select_container = document.getElementById("select_categories");
    await FetchTheSelectCategories(select_container)

    // Add a logic to create a new post:
    post_form = document.getElementById("postForm")
    post_form.addEventListener("submit", async function (e) {
        e.preventDefault()
        data = new FormData(post_form)
        // Add a new post:
        await addPost(data)
    })


    // Ensure user profile is fetched and `logged` is set BEFORE proceeding
    logged = await fetchUserProfile(userProfile, logged);

    // Now, fetch categories and only add "Liked Posts" if the user is logged in
    const side_bar = document.getElementById("navbar");
    // Fill the nav bar dynamically:
    await filNavFilter(side_bar, logged)

    //  Fetch all existing posts:
    await HandlePosts(logged)

    // Add an event listener to create a new comment:
    document.getElementById("commentForm").addEventListener("submit", async function (e) {
        e.preventDefault()
        const comment_tag = document.getElementById("commentBox")
        const post_id = parseInt(comment_tag.getAttribute("post_id"), 10);
        const comment = {
            postId: post_id,
            content: comment_tag.value.trim(),
        }
        comment_tag.value = ""
        if (!comment.postId || !comment.content) {
            alert("Comment cannot be empty.");
            return;
        }
        // Create a comment
        await CreateComment(comment)
    })
});


// End of the main function...
async function HandlePosts(logged) {
    // Fetch all existing posts:
    const main_content = document.getElementById("main-content");

    const posts = await FetchAllPosts();
    if (posts) {
        posts.forEach(async post => {
            console.log(post)
            const card = document.createElement("article");
            card.classList.add("post");
            card.setAttribute("post_id", post.id)

            // Post Header (Author, Category, Time)
            const postHeader = document.createElement("div");
            postHeader.classList.add("post-header");

            const author = document.createElement("span");
            author.classList.add("post-author");
            author.textContent = `${post.authorFirstName} ${post.authorLastName}`;


            // Handling the posts in relation to a category:
            const category = document.createElement("span");
            category.classList.add("post-category");
            await fetchPostCategories(category, post.id)
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
            likeCounter.setAttribute("liked_post_counter", post.id)
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
            dislikeCounter.setAttribute("disliked_post_counter", post.id)
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
            likeBtn.addEventListener("click", async () => {
                if (!logged) {
                    window.location.href = '/login';
                    return;
                }
                await VoteForPost(post.id, 1)
            });

            // Handle the dislike btn behavior:
            dislikeBtn.addEventListener("click", async () => {
                if (!logged) {
                    window.location.href = '/login';
                    return;
                }
                await VoteForPost(post.id, -1)
            });

            // handle comments
            cmtBtn.addEventListener("click", async function () {
                closeBtn = document.getElementById("close-btn")
                comments_container = document.getElementById("commentsList")
                modalPostTitle.textContent = post.title;
                modalPostContent.textContent = post.content;
                commentModal.style.display = "block";
                const comment = document.getElementById("commentBox")
                comment.setAttribute("post_id", post.id)
                commentModal.setAttribute("postId", post.id)
                closeBtn.addEventListener("click", () => commentModal.style.display = "none");
                window.addEventListener("click", event => {
                    while (comments_container.firstChild) {
                        comments_container.removeChild(comments_container.firstChild);
                    }
                    if (event.target === commentModal) commentModal.style.display = "none";
                });
                await GetPostComments(post.id, logged)
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
    }
}

// Fetch all posts:
async function FetchAllPosts() {
    try {
        const response = await fetch('http://localhost:8080/all_posts');
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        const posts = await response.json();
        return posts;
    } catch (error) {
        console.error('Error fetching posts:', error);
        return null;
    }
}

// Filter posts based on the post id that is retrieved from the category:
async function filterPosts(selected) {
    console.log("The chosen category is:", selected)
    try {
        // Fetch the list of post IDs belonging to the selected category
        const response = await fetch(`http://localhost:8080/posts_by_category?category=${encodeURIComponent(selected)}`);
        if (!response.ok) throw new Error("Failed to fetch posts");
        const postIds = await response.json(); // Assuming backend returns a map of post IDs
        console.log(postIds);  // Check the structure of postIds

        // Get all posts in the DOM
        const posts = document.getElementsByClassName("post");

        // Loop through each post and check if its ID is in the retrieved postIds
        Array.from(posts).forEach(post => {
            const postId = post.getAttribute("post_id");
            // Check if the postId is a key in the postIds map, or if "All" is selected
            if (selected === "All" || postIds.hasOwnProperty(postId)) {
                post.style.display = "block";
            } else {
                post.style.display = "none";
            }
        });
    } catch (error) {
        console.error("Error fetching or filtering posts:", error);
    }
}



// A function to filter owned and created posts:
async function FilterOwnedAndCreated(required) {
    try {
        const response = await fetch(`http://localhost:8080/${required}`);
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
    } catch (error) {
        console.error("Error fetching liked posts:", error);
    }
}

// Add a new post:
async function addPost(data) {
    try {
        const response = await fetch("http://localhost:8080/add_post", {
            method: "POST",
            body: data,
        })
        if (response.ok) {
            const result = await response.json();
            console.log("Post was added successfully: -------------->>>" + result.post_id);
          await fetchSpecificPost(result.post_id)
        } else {
            console.log("post creation has failed")
        }
    } catch (error) {
        console.error("Error: " + error)
    }
}

// Fill the select categories:
async function FetchTheSelectCategories(select_container) {
    try {
        const response = await fetch('http://localhost:8080/all_categories');
        const categories = await response.json();

        categories.forEach(category => {
            // Create a div wrapper for better styling
            const checkboxWrapper = document.createElement("div");
            checkboxWrapper.classList.add("form-check");

            // Create the checkbox input
            const checkbox = document.createElement("input");
            checkbox.type = "checkbox";
            checkbox.name = "chosen_categories[]"; // Array notation to handle multiple values
            checkbox.value = category.name;
            checkbox.id = `category_${category.name}`;
            checkbox.classList.add("form-check-input");

            // Create the label for the checkbox
            const label = document.createElement("label");
            label.htmlFor = `category_${category.name}`;
            label.textContent = category.name;
            label.classList.add("form-check-label");

            // Append checkbox and label to wrapper
            checkboxWrapper.appendChild(checkbox);
            checkboxWrapper.appendChild(label);

            // Append wrapper to container
            select_container.appendChild(checkboxWrapper);
        });
    } catch (error) {
        console.error('Error: ', error);
    }

}

// Vote for post:
async function VoteForPost(post_id, vote) {
    try {
        const response = await fetch(`http://localhost:8080/vote_for_post`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ postId: post_id, value: vote })
        });

        if (response.ok) {
            // Select elements using custom attributes
            const likeCounter = document.querySelector(`[liked_post_counter="${post_id}"]`);
            const dislikeCounter = document.querySelector(`[disliked_post_counter="${post_id}"]`);
            if (likeCounter && dislikeCounter) {
                if (vote === 1) {
                    likeCounter.textContent = Number(likeCounter.textContent) + 1;
                    dislikeCounter.textContent = Math.max(0, Number(dislikeCounter.textContent) - 1); // Ensure non-negative
                } else {
                    dislikeCounter.textContent = Number(dislikeCounter.textContent) + 1;
                    likeCounter.textContent = Math.max(0, Number(likeCounter.textContent) - 1);
                }
            }

        } else {
            console.error("Failed to like post");
        }
    } catch (error) {
        console.error("Error:", error);
    }
}

// Fetch post categories:
async function fetchPostCategories(category, postId) {
    try {
        const response = await fetch("http://localhost:8080/post_categories", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ post_id: postId })  // Send JSON data
        });
        console.log("The post id is: ", postId)
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const categories = await response.json();

        categories.forEach(cat => {
            if (category.textContent == "") {
                category.textContent = cat.name
            } else {
                category.textContent += ", " + cat.name;
            }
        });
    } catch (error) {
        console.error("Error fetching categories:", error);
    }
};

// Get post comments:
async function GetPostComments(postId, logged) {
    try {
        const response = await fetch(`http://localhost:8080/get_comments?post_id=${postId}`);
        if (!response.ok) throw new Error("Failed to fetch comments");

        const comments = await response.json();
        console.log(comments);

        if (comments != null) {
            comments.forEach(com => {
                console.log(com)
                let comment = document.createElement("div");
                comment.classList.add("comment_item");
                comment.textContent = com.content;

                // Footer (Likes, Dislikes)
                const commentFooter = document.createElement("div");
                commentFooter.classList.add("comment-footer"); // Fixed this line

                // Like button
                const commentLikeBtn = document.createElement("span");
                commentLikeBtn.classList.add("comment_like-btn");

                const commentLikeIcon = document.createElement("i");
                commentLikeIcon.classList.add("fa-solid", "fa-thumbs-up");

                const commentLikeCounter = document.createElement("span");
                commentLikeCounter.classList.add("comment_like-count");
                commentLikeCounter.setAttribute("liked_comment_counter", com.id)
                commentLikeCounter.textContent = com.likeCount; // Use 'com' instead of 'post'

                commentLikeBtn.appendChild(commentLikeIcon);
                commentLikeBtn.appendChild(commentLikeCounter);

                // Dislike button
                const commentdislikeBtn = document.createElement("span");
                commentdislikeBtn.classList.add("comment_dislike-btn");

                const commentdislikeIcon = document.createElement("i");
                commentdislikeIcon.classList.add("fa-solid", "fa-thumbs-down");

                const commentdislikeCounter = document.createElement("span");
                commentdislikeCounter.classList.add("comment_dislike-count");
                commentdislikeCounter.setAttribute("disliked_comment_counter", com.id)
                commentdislikeCounter.textContent = com.dislikeCount; // Use 'com' instead of 'post'

                commentdislikeBtn.appendChild(commentdislikeIcon);
                commentdislikeBtn.appendChild(commentdislikeCounter);

                // Append buttons to footer
                commentFooter.appendChild(commentLikeBtn);
                commentFooter.appendChild(commentdislikeBtn);


                // Attach Event Listener **INSIDE THE LOOP**
                commentLikeBtn.addEventListener("click", async (e) => {
                    e.stopPropagation();
                    if (!logged) {
                        window.location.href = '/login';
                        return;
                    }
                    await VoteForComment(com.id, 1)
                });

                // Handle the dislike btn behavior:
                commentdislikeBtn.addEventListener("click", async (e) => {
                    e.stopPropagation();
                    if (!logged) {
                        window.location.href = '/login';
                        return;
                    }
                    await VoteForComment(com.id, -1)
                });

                // Append footer to comment
                comment.appendChild(commentFooter);

                // Append comment to container
                comments_container.appendChild(comment);
            });
        }
    } catch (error) {
        console.error(error);
    }
}

// Create a new comment:

async function CreateComment(comm) {
    try {
        const response = await fetch(`http://localhost:8080/post_comments`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(comm)
        }
        )
        if (response.ok) {
            alert("comment wass added successfully")
        }
    } catch (error) {
        console.error(error)
    }
}

// fill the nav bar:
async function filNavFilter(side_bar, logged) {
    try {
        const response = await fetch("http://localhost:8080/all_categories");
        const categories = await response.json();

        categories.forEach((category) => {
            let category_item = document.createElement("button");
            category_item.classList.add("category");
            category_item.textContent = category.name;

            // Add an event listener for filtering posts
            category_item.addEventListener("click", () => filterPosts(category.name));

            // Append each category button to the navbar
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

            likedPostsBtn.addEventListener("click", async () => {
                await FilterOwnedAndCreated("liked")
            });


            ownedPost.addEventListener("click", async () => {
                await FilterOwnedAndCreated("owned")
            });

            // Append the "Liked Posts" button after categories
            side_bar.appendChild(likedPostsBtn);
            side_bar.appendChild(ownedPost);
        }
    } catch (error) {
        console.error("Error:", error);
    }
}

// Function to fetch the user profile
async function fetchUserProfile(userProfile, logged) {
    try {
        const response = await fetch("http://localhost:8080/profile");
        const user = await response.json();

        // Create user container
        const userInfo = document.createElement("div");
        userInfo.classList.add("user-info");
        if (user) {
            logged = true; // Ensure `logged` is updated correctly
            console.log("User is logged in");


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
            logoutBtn.onclick = async function () {
                try {
                    const logoutResponse = await fetch("http://localhost:8080/logout", {
                        method: "POST",
                        credentials: "include",
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
            userProfile.appendChild(userInfo);
        } else {
            // Login Button for non-logged-in users
            const loginBtn = document.createElement("button");
            loginBtn.classList.add("login-btn");
            loginBtn.textContent = "Login";
            loginBtn.onclick = function () {
                window.location.href = '/register';
            };
            userInfo.appendChild(loginBtn);
        }
        userProfile.appendChild(userInfo);
    } catch (error) {
        console.error("Error fetching profile:", error);
    }
    return logged
}

// Vote for comment:
async function VoteForComment(comment_id, vote) {
    console.log("The comment id is: ", comment_id)
    try {
        const response = await fetch(`http://localhost:8080/vote_for_comment`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ comment_id: comment_id, value: vote })
        });

        if (response.ok) {
            // Select elements using custom attributes
            const likeCounter = document.querySelector(`[liked_comment_counter="${comment_id}"]`);
            const dislikeCounter = document.querySelector(`[disliked_comment_counter="${comment_id}"]`);
            if (likeCounter && dislikeCounter) {
                if (vote === 1) {
                    likeCounter.textContent = Number(likeCounter.textContent) + 1;
                    dislikeCounter.textContent = Math.max(0, Number(dislikeCounter.textContent) - 1); // Ensure non-negative
                } else {
                    dislikeCounter.textContent = Number(dislikeCounter.textContent) + 1;
                    likeCounter.textContent = Math.max(0, Number(likeCounter.textContent) - 1);
                }
            }

        } else {
            console.error("Failed to like comment");
        }
    } catch (error) {
        console.error("Error:", error);
    }
}

// fetch a specific post:
async function fetchSpecificPost(post_id) {
    try {
       const response = await fetch(`http://localhost:8080/special_post?post_id=${post_id}`)
        new_post = await response.json()
        
    }catch(error) {
        console.error(error)
    }
}