let logged = false;
async function fetchUserProfile() {
    let userProfile = document.getElementById("user-profile");
    const createPost = document.getElementById("createPost");

    try {
        const response = await fetch("http://localhost:8080/profile");
        const user = await response.json();

        const userInfo = document.createElement("div");
        userInfo.classList.add("user-info");

        if (user) {
            logged = true;

            const profileImg = document.createElement("img");
            profileImg.src = "./app/uploads/" + user.profile_pic;
            profileImg.alt = "Profile Picture";
            profileImg.classList.add("profile-pic");

            const nameSpan = document.createElement("span");
            nameSpan.classList.add("user-name");
            nameSpan.textContent = `${user.first_name} ${user.last_name}`;

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
                    if (response.status === 404) {
                        // Store the status in sessionStorage
                        sessionStorage.setItem("errorStatus", response.status);

                        // Redirect to error page
                        window.location.href = "/errors.html";
                    }
                } catch (error) {
                    console.error("Logout error:", error);
                }
            };

            createPost.style.display = "block";
            userInfo.appendChild(profileImg);
            userInfo.appendChild(nameSpan);
            userInfo.appendChild(logoutBtn);
        } else {
            const loginBtn = document.createElement("button");
            loginBtn.classList.add("login-btn");
            loginBtn.textContent = "Login";
            loginBtn.onclick = function () {
                window.location.href = '/register';
            };
            userInfo.appendChild(loginBtn);
        }

        userProfile.appendChild(userInfo);
        if (response.status === 404) {
            // Store the status in sessionStorage
            sessionStorage.setItem("errorStatus", response.status);
    
            // Redirect to error page
            window.location.href = "/errors.html";
          }
    } catch (error) {
        console.error("Error fetching profile:", error);
    }

    return logged;
}