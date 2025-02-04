document.addEventListener("DOMContentLoaded", function () {
    post_form = document.getElementById("postForm")
    post_form.addEventListener("submit", async function (e) {
        e.preventDefault()
        data = new FormData(post_form)
        try {
            const response = await fetch("http://localhost:8080/add_post", {
                method: "POST",
                body: data,
            })
            if (response.ok) {
                const result = await response.json();
                alert("Registration successful: " + result.message);
                
            } else {
                console.log("post creation has failed")
            }
        } catch (error) {
            console.error("Error: " + error)
        }
    })
})