document.addEventListener("DOMContentLoaded", async function(e){
    document.getElementById(commentForm).addEventListener("submit", async function(e){
        e.preventDefault()
        const comment_tag = document.getElementById("commentBox")
        post_id = comment_tag.getAttribute("post_id")
        const comment = {
            postId: post_Id,
            content: comment_tag.ariaValueMax.trim(),
        }
        if(!comment.postId || !comment.content) {
            alert("Comment cannot be empty.");
            return;
        }

        try{
            const response = await fetch(`http://localhost:8080/post_comments`,
                method: "POST",
                data: comment,
             )
        }catch(error) {

        }
    })
})

