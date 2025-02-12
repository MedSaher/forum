document.addEventListener("DOMContentLoaded", async function(e){
    document.getElementById("commentForm").addEventListener("submit", async function(e){
        e.preventDefault()
        const comment_tag = document.getElementById("commentBox")
        const post_id = parseInt(comment_tag.getAttribute("post_id"), 10);
        const comment = {
            postId: post_id,
            content: comment_tag.value.trim(),
        }
        if(!comment.postId || !comment.content) {
            alert("Comment cannot be empty.");
            return;
        }
        console.log(comment)
        try{
            const response = await fetch(`http://localhost:8080/post_comments`,{
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(comment)
            }
             )
             if(response.ok) {
                
                alert("comment wass added successfully")
             }
        }catch(error) {
            console.error(error)
        }
    })
})

