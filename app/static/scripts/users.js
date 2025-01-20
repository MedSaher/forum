// Form validation and submission handling
document.addEventListener('DOMContentLoaded', function () {
    // Toggle between login and registration forms
    const showLogin = document.getElementById('showLogin');
    const showRegister = document.getElementById('showRegister');
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');

    showLogin.addEventListener('click', function (e) {
        e.preventDefault();
        loginForm.classList.add('active');
        registerForm.classList.remove('active');
    });

    showRegister.addEventListener('click', function (e) {
        e.preventDefault();
        registerForm.classList.add('active');
        loginForm.classList.remove('active');
    });

    // Registration form handling
    const registrationForm = document.getElementById('registrationForm');
    registrationForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const formData = new FormData(registrationForm);

        try {
            const response = await fetch("http://localhost:8080/register", {
                method: "POST",
                body: formData,
            });

            if (response.ok) {
                const result = await response.json();
                alert("Registration successful: " + result.message);
            } else {
                alert("Registration failed");
            }
        } catch (error) {
            console.error("Error:", error);
            alert("An error occurred");
        }
    });

});
// Login form handling
const loginFormElement = document.getElementById('loginFormElement');
loginFormElement.addEventListener('submit', async function (e) {
    e.preventDefault(); // Prevent default form submission

    // Clear any existing error messages
    clearErrors();

    // Get user input and trim whitespace
    const email = document.getElementById('loginEmail').value.trim();
    const password = document.getElementById('loginPassword').value.trim();
    console.log(email)
    let isValid = true;

    // Validate email
    if (!isValidEmail(email)) {
        showError('loginEmailError', 'Please enter a valid email address');
        isValid = false;
    }

    // Validate password
    if (!password) {
        showError('loginPasswordError', 'Password is required');
        isValid = false;
    }

    // Proceed only if inputs are valid
    if (isValid) {
        const loginData = { email, password }; // Prepare login data

        try {
            // Send a POST request to the login endpoint
            const response = await fetch('/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json', // Indicate JSON payload
                },
                credentials: 'include', // Include cookies (e.g., CSRF tokens)
                body: JSON.stringify(loginData), // Convert login data to JSON
            });

            if (response.ok) {
                // Redirect to the dashboard on successful login
                window.location.href = '/all_posts';
            } else {
                // Display error message from the server
                const errorText = await response.text();
                showError('loginError', errorText);
            }
        } catch (error) {
            // Handle any network or unexpected errors
            console.error('Error during login:', error.message);
            showError('loginError', 'An unexpected error occurred');
        }
    }
});

// Helper function: Validate email format
function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/; // Simple email regex
    return emailRegex.test(email);
}

// Helper function: Show error messages
function showError(elementId, errorMessage) {
    const element = document.getElementById(elementId);
    if (element) {
        element.textContent = errorMessage; // Display error message
        element.style.display = 'block'; // Ensure the element is visible
    } else {
        alert(errorMessage); // Fallback for missing elements
    }
}

// Helper function: Clear all error messages
function clearErrors() {
    const errorElements = document.querySelectorAll('.error-message');
    errorElements.forEach(element => {
        element.textContent = ''; // Clear text
        element.style.display = 'none'; // Hide element
    });
}
