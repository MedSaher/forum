document.addEventListener("DOMContentLoaded", async function () {
    document.getElementById("login_btn").addEventListener("click", function () {
        alert("Button clicked!");
        // Redirect to the dashboard on successful login
        window.location.href = '/register';
    });
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
            category_list.appendChild(category_item)
        })
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

            const title = document.createElement("div");
            title.classList.add("card-title");
            title.textContent = post.title;

            const content = document.createElement("div");
            content.classList.add("card-content");
            content.textContent = post.content;

            card.appendChild(title);
            card.appendChild(content);
            main_content.appendChild(card);
        });
    } catch (error) {
        console.error('Error:', error);
    }

    // Fetch the loged in user based on the session:
    try {
        const response = await fetch('http://localhost:8080/profile');
        const user = await response.json();
        console.log(user)
    } catch (error) {
        console.error('Error !!!:', error);
    }
});
