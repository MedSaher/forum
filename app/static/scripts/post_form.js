document.addEventListener("DOMContentLoaded", async function () {

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
                alert("Post addes successfully: " + result.message);
                // Redirect to the dashboard on successful login
                window.location.href = '/';
            } else {
                console.log("post creation has failed")
            }
        } catch (error) {
            console.error("Error: " + error)
        }
    })

    // fill the categories dynamically:
    const select_container = document.getElementById("select_categories")
    const select_element = document.createElement("select")
    select_element.name = "chosen_category"
    select_element.id = "chosen_category"
    select_element.classList.add("form-control")
    try {
        const response = await fetch('http://localhost:8080/all_categories');
        const categories = await response.json();
        categories.forEach(category => {
            let category_option = document.createElement("option")
            category_option.value = category.name
            category_option.textContent = category.name
            select_element.appendChild(category_option)
        })
        select_container.appendChild(select_element)
    } catch (error) {
        console.error('Error: ', error)
    }
})

