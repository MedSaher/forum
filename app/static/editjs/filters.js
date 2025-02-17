function filterPosts(selected) {
    const posts = document.getElementsByClassName("post");
    Array.from(posts).forEach(post => {
        const itemCategory = post.getAttribute("data-category");
        if (selected === "All" || itemCategory.split(",").includes(selected)) {
            post.style.display = "block";
        } else {
            post.style.display = "none";
        }
    });
}