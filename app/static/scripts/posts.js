document.addEventListener("DOMContentLoaded", async function () {
    let user_profile = document.getElementById("user-profile");
    // Fetch the logged-in user based on the session
    const profile_img = document.createElement("img");
    const nameSpan = document.createElement("span");
    const btn = document.createElement("button");

    try {
        const response = await fetch('http://localhost:8080/profile');
        const user = await response.json();
        console.log(user);

        if (user) {
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
                        alert("Logged out successfully!");
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
});
