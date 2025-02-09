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
            // Create a container for each checkbox
            let checkboxContainer = document.createElement("div");
    
            // Create the checkbox input
            let checkbox = document.createElement("input");
            checkbox.type = "checkbox";
            checkbox.id = category.name;
            checkbox.name = "categories"; // Use a common name to send multiple selections
            checkbox.value = category.name;
    
            // Create the label
            let label = document.createElement("label");
            label.htmlFor = category.name;
            label.textContent = category.name;
    
            // Append checkbox and label to the container
            checkboxContainer.appendChild(checkbox);
            checkboxContainer.appendChild(label);
    
            // Append the container to the parent element
            select_container.appendChild(checkboxContainer);
        });
    } catch (error) {
        console.error('Error: ', error);
    }    
})

