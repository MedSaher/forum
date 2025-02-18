// Function to go back to previous page
function goBack() {
    window.location.href = "/";;
}
window.onload = function() {
    // Retrieve stored error details
    const title = sessionStorage.getItem("errorTitle");
    const message = sessionStorage.getItem("errorMessage");
    const status = sessionStorage.getItem("errorStatus");

    // Update the error page elements
    if (title && message) {
        document.querySelector('.error-title').textContent = title;
        document.querySelector('.error-message').textContent = message;
    }
    showModal("Error Status:", status); // Debugging
};