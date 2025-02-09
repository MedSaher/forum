document.addEventListener('DOMContentLoaded', function() {
    // Toggle between login and registration forms
    const showLogin = document.getElementById('showLogin');
    const showRegister = document.getElementById('showRegister');
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');

    showLogin.addEventListener('click', function(e) {
        e.preventDefault();
        loginForm.classList.add('active');
        registerForm.classList.remove('active');
    });

    showRegister.addEventListener('click', function(e) {
        e.preventDefault();
        registerForm.classList.add('active');
        loginForm.classList.remove('active');
    });

    // Registration form handling
    const registrationForm = document.getElementById('registrationForm');
    registrationForm.addEventListener("submit", async(e) => {
        e.preventDefault();

        clearErrors(); // Clear previous errors

        const formData = new FormData(registrationForm);
        let isValid = validateRegistration(formData);

        if (!isValid) return;

        try {
            const response = await fetch("http://localhost:8080/register", {
                method: "POST",
                body: formData,
            });

            if (response.ok) {
                const result = await response.json();
                alert("Registration successful: " + result.message);
                window.location.href = '/'; // Redirect after registration
            } else {
                showError('registrationError', 'Registration failed');
            }
        } catch (error) {
            console.error("Error:", error);
            showError('registrationError', 'An error occurred');
        }
    });

    // Login form handling
    const loginFormElement = document.getElementById('loginFormElement');
    loginFormElement.addEventListener('submit', async function(e) {
        e.preventDefault();

        clearErrors();

        const email = document.getElementById('loginEmail').value.trim();
        const password = document.getElementById('loginPassword').value.trim();
        let isValid = true;

        if (!isValidEmail(email)) {
            showError('loginEmailError', 'Please enter a valid email address');
            isValid = false;
        }
        if (!password) {
            showError('loginPasswordError', 'Password is required');
            isValid = false;
        }

        if (isValid) {
            try {
                const response = await fetch('/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    credentials: 'include',
                    body: JSON.stringify({ email, password }),
                });

                if (response.ok) {
                    window.location.href = '/';
                } else {
                    const errorText = await response.text();
                    showError('loginError', errorText);
                }
            } catch (error) {
                console.error('Error during login:', error.message);
                showError('loginError', 'An unexpected error occurred');
            }
        }
    });
});

// Validate registration inputs
function validateRegistration(formData) {
    let isValid = true;

    const firstName = formData.get('firstName').trim();
    const lastName = formData.get('lastName').trim();
    const email = formData.get('email').trim();
    const password = formData.get('password').trim();

    if (firstName.length < 2) {
        showError('firstNameError', 'First name must be at least 2 characters');
        isValid = false;
    }
    if (lastName.length < 2) {
        showError('lastNameError', 'Last name must be at least 2 characters');
        isValid = false;
    }
    if (!isValidEmail(email)) {
        showError('emailError', 'Invalid email address');
        isValid = false;
    }
    if (password.length < 6) {
        showError('passwordError', 'Password must be at least 6 characters');
        isValid = false;
    }

    return isValid;
}

// Helper: Validate email format
function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

// Helper: Show error messages
function showError(elementId, errorMessage) {
    const element = document.getElementById(elementId);
    if (element) {
        element.textContent = errorMessage;
        element.style.display = 'block';
    } else {
        alert(errorMessage);
    }
}

// Helper: Clear all error messages
function clearErrors() {
    document.querySelectorAll('.error-message').forEach(element => {
        element.textContent = '';
        element.style.display = 'none';
    });
}