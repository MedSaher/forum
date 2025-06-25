document.addEventListener('DOMContentLoaded', function () {
    // Toggle between login and registration forms
    const showLogin = document.getElementById('showLogin');
    const showRegister = document.getElementById('showRegister');
    const go_home = document.getElementById('go_home');
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

    go_home.addEventListener("click", (e) => {
        e.preventDefault();
        window.location.href = "/";
    });


    // Registration form handling
    const registrationForm = document.getElementById('registrationForm');
    const modal = document.getElementById('registrationModal');
    const modalMessage = document.getElementById('modalMessage');
    const closeModal = document.getElementById('closeModal');
    const modalOkBtn = document.getElementById('modalOkBtn');

    // Form submission event listener
    registrationForm.addEventListener("submit", async (e) => {
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
                showModal("Registration successful: " + result.message, true);
            } else {
                showModal("Registration failed", false);
            }
            if (response.status === 404) {
                // Store the status in sessionStorage
                sessionStorage.setItem("errorStatus", response.status);

                // Redirect to error page
                window.location.href = "/error.html";
            }
        } catch (error) {
            console.error("Error:", error);
            showModal("An error occurred", false);
        }
    });

    // Function to display modal
    function showModal(message, success) {
        modalMessage.textContent = message;
        modal.style.display = "flex";

        modalOkBtn.onclick = () => {
            modal.style.display = "none";
            if (success) window.location.href = '/'; // Redirect only on success
        };
    }

    // Close modal on clicking "X" button
    closeModal.onclick = () => {
        modal.style.display = "none";
    };

    // Close modal if user clicks outside the modal content
    window.onclick = (e) => {
        if (e.target === modal) modal.style.display = "none";
    };

    // Login form handling
    const loginFormElement = document.getElementById('loginFormElement');

    loginFormElement.addEventListener('submit', async function (e) {
        e.preventDefault();
        clearErrors(); // Clear previous errors

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
                    window.location.href = '/'; // Redirect on success
                } else {
                    const errorText = await response.text();
                    showModal(errorText, false); // Show error in modal
                }
                if (response.status === 404) {
                    // Store the status in sessionStorage
                    sessionStorage.setItem("errorStatus", response.status);

                    // Redirect to error page
                    window.location.href = "/error.html";
                }
            } catch (error) {
                console.error('Error during login:', error.message);
                showModal('An unexpected error occurred', false);
            }
        }
    });

    // Function to display modal
    function showModal(message, success) {
        modalMessage.textContent = message;
        modal.style.display = "flex";

        modalOkBtn.onclick = () => {
            modal.style.display = "none";
            if (success) window.location.href = '/'; // Redirect only on success
        };
    }

    // Close modal on clicking "X" button
    closeModal.onclick = () => {
        modal.style.display = "none";
    };

    // Close modal if user clicks outside the modal content
    window.onclick = (e) => {
        if (e.target === modal) modal.style.display = "none";
    };
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
        showModal(errorMessage);
    }
}

// Helper: Clear all error messages
function clearErrors() {
    document.querySelectorAll('.error-message').forEach(element => {
        element.textContent = '';
        element.style.display = 'none';
    });
}