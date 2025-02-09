document.addEventListener("DOMContentLoaded", async function () {
    var loged = false
    let user_profile = document.getElementById("user-profile");
    // Fetch the logged-in user based on the session
    const profile_img = document.createElement("img");
    const nameSpan = document.createElement("span");
    const btn = document.createElement("button");

    try {
        const response = await fetch('http://localhost:8080/profile');
        const user = await response.json();
        if (user) {
            loged = true
            // add a create header:
            let link_to = document.createElement("a")
            link_to.className = "nav-link"
            link_to.textContent = "Create post"
            link_to.href = "/add_post"
            user_profile.appendChild(link_to)
            // Handle the profile image
            profile_img.src = "./app/uploads/" + user.profile_pic;
            profile_img.alt = "Profile Picture";
            profile_img.id = "navProfilePic";

            // Handle the user name
            nameSpan.id = "navUserName";
            nameSpan.textContent = `${user.first_name} ${user.last_name}`;

            // Set button for logout
            btn.id = "logout_btn";
            btn.textContent = "Logout";
            btn.onclick = async function () {
                try {
                    const logoutResponse = await fetch('http://localhost:8080/logout', {
                        method: 'POST',
                        credentials: 'include'
                    });

                    if (logoutResponse.ok) {
                        window.location.reload(); // Refresh page to show login button
                    } else {
                        console.error("Logout failed");
                    }
                } catch (error) {
                    console.error("Logout error:", error);
                }
            };
            user_profile.appendChild(profile_img);
            user_profile.appendChild(nameSpan);
            user_profile.appendChild(btn);
        } else {
            // Set button for login if user is not logged in
            btn.id = "login_btn";
            btn.textContent = "Login";
            user_profile.appendChild(btn);
            // Handle the login logic:
            document.getElementById("login_btn").addEventListener("click", function () {
                alert("Button clicked!");
                // Redirect to the dashboard on successful login
                window.location.href = '/register';
            });
        }
    } catch (error) {
        console.error('Error !!!:', error);
    }

    // Fetch all existing categories:
    const side_bar = document.getElementById("sidebar")
    let category_list = document.createElement("ul")
    category_list.classList.add("category-list")
    try {
        const response = await fetch('http://localhost:8080/all_categories');
        const categories = await response.json();
        categories.forEach(category => {
            let category_item = document.createElement("li")
            category_item.classList.add("category-item")
            category_item.textContent = category.name
            // Add an event listener for filter:
            category_item.addEventListener("click", () => filterPosts(category.name))
            category_list.appendChild(category_item)
        })

        if(loged) {
            // Add a filter for liked posts:
            let likedPosts = document.createElement("li")
            likedPosts.classList.add("other-item")
            likedPosts.textContent = "liked posts"

            // Add a filter for owned posts:
            let ownedPost = document.createElement("li")
            ownedPost.classList.add("other-item")
            ownedPost.textContent = "created posts"

            // Filter based on liked posts:
            likedPosts.addEventListener("click", async () => {
                try {
                    const response = await fetch('http://localhost:8080/liked')

                } catch(error){

                }
            })

            // Filter based on liked posts:
            ownedPost.addEventListener("click", async () => {
                likedPosts.addEventListener("click", async () => {
                    try {
                        const response = await fetch('http://localhost:8080/owned')
                    } catch(error){
    
                    }
                })
            })

            category_list.appendChild(likedPosts)
            category_list.appendChild(ownedPost)
        }

    
        side_bar.appendChild(category_list)
    } catch (error) {
        console.error('Error: ', error)
    }

    // Fetch all existing posts:
    const main_content = document.getElementById("main-content");
    try {
        const response = await fetch('http://localhost:8080/all_posts');
        const posts = await response.json();

        posts.forEach(post => {
            const card = document.createElement("div");
            card.classList.add("card");
            card.setAttribute("data-category", post.categoryName)
            card.setAttribute("post_id", post.id)
            const title = document.createElement("div");
            title.classList.add("card-title");
            title.textContent = post.title;

            const content = document.createElement("div");
            content.classList.add("card-content");
            content.textContent = post.content;

            // Create a vote section:
            const buttonsDiv = document.createElement("div");
            buttonsDiv.classList.add("vote");

            // Create a comment section:
            const commentSection = document.createElement("div");
            commentSection.classList.add("comment-section");

            if (loged) {
                console.log("loded")
                // Create buttons container
                // Create like button
                const likeBtn = document.createElement("button");
                likeBtn.classList.add("like-btn");
                likeBtn.textContent = "👍 ";

                const likeCount = document.createElement("span");
                likeCount.classList.add("like-count");
                likeCount.textContent = post.likeCount;
                likeBtn.appendChild(likeCount);

                // Create dislike button
                const dislikeBtn = document.createElement("button");
                dislikeBtn.classList.add("dislike-btn");
                dislikeBtn.textContent = "👎 ";

                // Comment Button
                const commentBtn = document.createElement("button");
                commentBtn.textContent = "💬 Comment";

                const commentBox = document.createElement("textarea");
                commentBox.classList.add("comment-box");
                commentBox.placeholder = "Write a comment...";

                const submitCommentBtn = document.createElement("button");
                submitCommentBtn.classList.add("comment-btn");
                submitCommentBtn.textContent = "Submit";


                const dislikeCount = document.createElement("span");
                dislikeCount.classList.add("dislike-count");
                dislikeCount.textContent = post.dislikeCount;
                dislikeBtn.appendChild(dislikeCount);

                // Attach event listeners
                // Handle the like btn behavior:
                likeBtn.addEventListener("click", async () => {
                    try {
                        const response = await fetch(`http://localhost:8080/vote_for_post`, { // Use backticks
                            method: "POST", // Set HTTP method to POST
                            headers: {
                                "Content-Type": "application/json"
                            },
                            body: JSON.stringify({ postId: post.id, value: 1 }) // Send data if required by backend
                        });

                        if (response.ok) {
                            window.location.reload(); // Refresh page to update dislike count
                        }
                    } catch (error) {
                        console.error("Error:", error);
                    }
                });

                // Handle the dislike btn behavior:
                dislikeBtn.addEventListener("click", async () => {
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


                // Toggle Comment Section visibility
                commentBtn.addEventListener("click", () => {
                    commentSection.style.display = commentSection.style.display === "none" || !commentSection.style.display ? "block" : "none";
                });

                // Submit Comment Event
                submitCommentBtn.addEventListener("click", () => {
                    const comment = commentBox.value.trim();
                    if (comment) {
                        const commentPara = document.createElement("p");
                        commentPara.textContent = comment;
                        commentSection.appendChild(commentPara);
                        commentBox.value = ""; // Clear comment box after submission
                    }
                });


                // Append buttons to buttonsDiv
                buttonsDiv.appendChild(likeBtn);
                buttonsDiv.appendChild(dislikeBtn);
                buttonsDiv.appendChild(commentBtn)

            }
            card.appendChild(title);
            card.appendChild(content);
            main_content.appendChild(card);
            if (loged) {
                card.appendChild(buttonsDiv);
                card.appendChild(commentSection)
            }
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

// Create a function to fiter based on liked posts:
