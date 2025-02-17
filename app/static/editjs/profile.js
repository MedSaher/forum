let logged = false;

async function fetchUserProfile() {
    const userProfile = document.getElementById("user-profile");
    const createPost = document.getElementById("createPost");
    
    // Create a login button dynamically
    const loginBtn = document.createElement("button");
    loginBtn.classList.add("login-btn");
    loginBtn.textContent = "Login";
    loginBtn.onclick = function() {
        window.location.href = '/register';
    };
    
    try {
        const response = await fetch("http://localhost:8080/profile", {
            credentials: "include" // Add this to ensure cookies are sent
        });
        
        // Check for 404 status first
        if (response.status === 404) {
            sessionStorage.setItem("errorStatus", response.status);
            window.location.href = "/errors.html";
            return false;
        }
        
        // Then check if response is ok
        if (!response.ok) {
            userProfile.appendChild(loginBtn);
            return false;
        }
        
        const user = await response.json();
        if (!user) {
            userProfile.appendChild(loginBtn);
            return false;
        }
        
        // User is logged in, create profile elements
        logged = true;
        
        const userInfo = document.createElement("div");
        userInfo.classList.add("user-info");
        
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
        logoutBtn.onclick = async function() {
            try {
                const logoutResponse = await fetch("http://localhost:8080/logout", {
                    method: "POST",
                    credentials: "include"
                });
                
                if (logoutResponse.ok) {
                    window.location.reload();
                } else {
                    console.error("Logout failed");
                }
            } catch (error) {
                console.error("Logout error:", error);
            }
        };
        
        // Append all elements to userInfo
        userInfo.appendChild(profileImg);
        userInfo.appendChild(nameSpan);
        userInfo.appendChild(logoutBtn);
        userProfile.appendChild(userInfo);
        
        // Show create post button if logged in
        if (createPost) {
            createPost.style.display = "block";
        }
        
        return true;
        
    } catch (error) {
        console.error("Error fetching profile:", error);
        userProfile.appendChild(loginBtn);
        return false;
    }
}