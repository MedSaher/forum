// Get modal elements
const commentModal = document.getElementById("commentModal");
var commentBox = document.getElementById("commentBox");
const submitCommentBtn = document.getElementById("submitCommentBtn");
const commentsList = document.getElementById("commentsList");

function createLikeButton(post) {
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

    // Attach Event Listener **INSIDE THE LOOP**
    likeBtn.addEventListener("click", async () => {
        if (!logged) {
            window.location.href = '/login';
            return;
        }
        await VoteForPost(post.id, 1)
    });
    return likeBtn;
}

function createDislikeButton(post) {
    const dislikeBtn = document.createElement("span");
    dislikeBtn.classList.add("dislike-btn");

    const dislikeBtnIcon = document.createElement("i");
    dislikeBtnIcon.classList.add("fa-solid", "fa-thumbs-down");

    const dislikeBtnCounter = document.createElement("span");
    dislikeBtnCounter.classList.add("dislike-count");
    dislikeBtnCounter.setAttribute("disliked_post_counter", post.id)
    dislikeBtnCounter.textContent = post.dislikeCount;

    dislikeBtn.appendChild(dislikeBtnIcon);
    dislikeBtn.appendChild(dislikeBtnCounter);


    // Handle the dislike btn behavior:
    dislikeBtn.addEventListener("click", async () => {
        if (!logged) {
            window.location.href = '/login';
            return;
        }
        await VoteForPost(post.id, -1)
    });

    return dislikeBtn;
}

function createCommentButton(post) {
    const commentBtn = document.createElement("span");
    commentBtn.classList.add("comment-btn");

    const commentIcon = document.createElement("i");
    commentIcon.classList.add("fa-solid", "fa-comment"); // Fix missing class on icon
    commentBtn.appendChild(commentIcon);

    commentBtn.addEventListener("click", async function () {
        closeBtn = document.getElementById("close-btn")
        modalPostTitle.textContent = post.title;
        modalPostContent.textContent = post.content;

        currentPostId = post.id; // Ensure correct property

        // Set post ID in a hidden input field (to be used later)
        document.getElementById("postId").value = currentPostId;

        commentModal.style.display = "block";
        await loadComments(currentPostId)
        // await GetPostComments(post.id, logged)
        closeBtn.addEventListener("click", () => commentModal.style.display = "none");
        window.addEventListener("click", event => {
            if (event.target === commentModal) commentModal.style.display = "none";
        });
    });
    return commentBtn;
}


document.addEventListener("DOMContentLoaded", function () {
    // Submit comment function
    submitCommentBtn.addEventListener("click", async function () {
        const currentPostId = parseInt(document.getElementById("postId").value, 10);
        if (!currentPostId) {
            showModal("Post ID is missing")
            return;
        }

        const comment = {
            postId: currentPostId,
            content: commentBox.value.trim(),
        }
        const commentText = commentBox.value.trim();
        commentBox.value = ""
        if (commentText === "") {
            showModal("comment cannot be empty")
            return;
        }
        // Create a comment
        added = await CreateComment(comment)
        let no_comment = document.getElementById("noCommentsMessage");
        if (no_comment) {
            no_comment.innerText = "";
        }
    });
});

function createLikeCmtButton(comment) {
    const likeBtn = document.createElement("span");
    likeBtn.classList.add("like-btn");

    const likeIcon = document.createElement("i");
    likeIcon.classList.add("fa-solid", "fa-thumbs-up");

    const likeCounter = document.createElement("span");
    likeCounter.classList.add("like-count");
    likeCounter.textContent = comment.likeCount;
    likeCounter.setAttribute("liked_comment_counter", comment.id)
    likeBtn.appendChild(likeIcon);
    likeBtn.appendChild(likeCounter);
    likeBtn.addEventListener("click", async () => {
        if (!logged) {
            window.location.href = '/login';
            return;
        }
        await VoteForComment(comment.id, 1)
    });

    return likeBtn;
}

function createDislikeCmtButton(comment) {
    const dislikeBtn = document.createElement("span");
    dislikeBtn.classList.add("dislike-btn");

    const dislikeIcon = document.createElement("i");
    dislikeIcon.classList.add("fa-solid", "fa-thumbs-down");

    const dislikeCounter = document.createElement("span");
    dislikeCounter.classList.add("dislike-count");

    dislikeCounter.textContent = comment.dislikeCount;
    dislikeCounter.setAttribute("disliked_comment_counter", comment.id)

    dislikeBtn.appendChild(dislikeIcon);
    dislikeBtn.appendChild(dislikeCounter);

    dislikeBtn.addEventListener("click", async () => {
        if (!logged) {
            window.location.href = '/login';
            return;
        }
        await VoteForComment(comment.id, -1)
    });

    return dislikeBtn;
}

// Load all comments to a post:
async function loadComments(postId) {
    try {
        // Fetch comments for the specific post
        const response = await fetch(`/get_comments?post_id=${postId}`);
        if (!response.ok) {
            throw new Error("Failed to fetch comments");
        }
        const comments = await response.json();
        const commentsList = document.getElementById("commentsList");
        commentsList.innerHTML = "";

        if (comments === null) {
            let noCommentsMessage = document.createElement("p");
            noCommentsMessage.id = "noCommentsMessage";
            noCommentsMessage.textContent = "No comments yet. Be the first to comment!";
            // Append it to the comments list dynamically
            commentsList.appendChild(noCommentsMessage);
        } else {

            comments.forEach(comment => {
                const commentElement = document.createElement("div");
                commentElement.classList.add("comment");

                commentElement.innerHTML = `
                    <span class="author">${comment.firstName} ${comment.lastName}</span>
                    <span class="timestamp">${new Date(comment.timestamp).toLocaleString()}</span>
                    <p>${comment.content}</p>
                `;
                // Create buttons separately and append them
                const likeBtn = createLikeCmtButton(comment);
                const dislikeBtn = createDislikeCmtButton(comment);

                const hrline = document.createElement("hr")
                commentElement.appendChild(likeBtn);
                commentElement.appendChild(dislikeBtn);
                commentElement.appendChild(hrline);
                commentsList.appendChild(commentElement);
            });
        }
    } catch (error) {
        showModal("Error loading comments!")
        document.getElementById("commentsList").innerHTML = "<p>Error loading comments.</p>";
    }
}

// Vote for post:
// Frontend JavaScript (vote.js)
async function VoteForPost(post_id, vote) {
    try {
        const response = await fetch(`http://localhost:8080/vote_for_post`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ postId: post_id, value: vote })
        });

        const data = await response.json();

        if (response.ok) {
            const likeCounter = document.querySelector(`[liked_post_counter="${post_id}"]`);
            const dislikeCounter = document.querySelector(`[disliked_post_counter="${post_id}"]`);

            if (likeCounter && dislikeCounter) {
                likeCounter.textContent = data.likeCount;
                dislikeCounter.textContent = data.dislikeCount;
            }
        } else {
            showModal("Error : Failed to vote")
        }
    } catch (error) {
        showModal("error " + error)
    }
}

// Create a new comment:
async function CreateComment(comm) {
    let added = false
    try {
        const response = await fetch(`http://localhost:8080/post_comment`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(comm)
        }
        )
        if (response.status == 502) {
            window.location.href = "/login";
            return
        }
        if (!response.ok) {
            showModal("Error creating a comment!")
        } else {
            let comment = await response.json()
            const commentElement = document.createElement("div");
            commentElement.classList.add("comment");

            commentElement.innerHTML = `
                <span class="author">${comment.firstName} ${comment.lastName}</span>
                <span class="timestamp">${new Date(comment.timestamp).toLocaleString()}</span>
                <p>${comment.content}</p>
            `;
            // Create buttons separately and append them
            const likeBtn = createLikeCmtButton(comment);
            const dislikeBtn = createDislikeCmtButton(comment);

            const hrline = document.createElement("hr")
            commentElement.appendChild(likeBtn);
            commentElement.appendChild(dislikeBtn);
            commentElement.appendChild(hrline);
            commentsList.insertBefore(commentElement, commentsList.firstChild);
            added = true
        }
    } catch (error) {
        showModal("Error creating a comment!")
    }
    return added
}

// Vote for comment:
async function VoteForComment(comment_id, vote) {
    try {
        const response = await fetch(`http://localhost:8080/vote_for_comment`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ comment_id: comment_id, value: vote })
        });

        const data = await response.json();

        if (response.ok) {
            const likeCounter = document.querySelector(`[liked_comment_counter="${comment_id}"]`);
            const dislikeCounter = document.querySelector(`[disliked_comment_counter="${comment_id}"]`);

            if (likeCounter && dislikeCounter) {
                likeCounter.textContent = data.likeCount;
                dislikeCounter.textContent = data.dislikeCount;
            }
        } else {
            showModal("Failed to vote!")
        }
    } catch (error) {
        showModal("error " + error)
    }
}

// Error message
function showError(message) {
    // Create a new <p> element
    let errorMsg = document.createElement("p");

    // Add text content
    errorMsg.textContent = message;

    // Apply styles
    errorMsg.style.color = "red";
    errorMsg.style.fontWeight = "bold";

    // // Append to the container
    return errorMsg
}