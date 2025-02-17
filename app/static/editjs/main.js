document.addEventListener("DOMContentLoaded", async function() {
    let logged = await fetchUserProfile();
    await fetchCategories();
    await fetchPosts();
});