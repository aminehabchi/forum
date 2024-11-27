let offset = 3
let loading = false
let noMorePosts = false

const filterParams = {
    category: document.querySelector('meta[name="filter-category"]')?.content || '',
    created: document.querySelector('meta[name="filter-created"]')?.content || '',
    liked: document.querySelector('meta[name="filter-liked"]')?.content || ''
};

function loadMorePosts() {
    if (loading || noMorePosts) return;

    loading = true;
    loadingContainer.style.visibility = "visible"

    const url = new URL('/load-more-posts', window.location.origin);
    url.searchParams.set('offset', offset);
    url.searchParams.set('category', filterParams.category);
    url.searchParams.set('created', filterParams.created);
    url.searchParams.set('liked', filterParams.liked);

    fetch(url)
        .then(resp => resp.json())
        .then(posts => {
            if (posts == null) {
                noMorePosts = true
                loadingContainer.textContent = 'No more posts to load';
                return;
            }

            posts.forEach(post => {
                const postElement = createPostElement(post);
                postsContainer.appendChild(postElement);
            })

            offset += posts.length
            loading = false;
            loadingContainer.style.visibility = "hidden"
        })
        .catch(err => {
            console.error('Error loading more posts:', err);
            loading = false;
            loadingContainer.style.visibility = "hidden"
        })
}


function createPostElement(post) {
    const div = document.createElement('div');
    div.className = "post"
    div.innerHTML = `
        <div class="profilInfo">
            <img class="profileImg" src="static/profil.png">
            <div class="profile-details">
                <span>${post.Name}</span>
                <span class="time">${post.CreatedAt}</span>
            </div>
        </div>
        <h4>${post.Title}</h4>
        ${post.Category.map(cat => `
            <a class="category" href="/filter?category=${cat}">
                <i class="fas fa-tag"></i> ${cat}
            </a>
        `).join('')}
        <p class="content">${post.Content}</p>
        <div style="margin-top: 15px;">
            <button class="action-btn likedislikebtn ${post.UserInteraction === 1 ? 'liked-btn' : ''}"
                    data-id="${post.ID}" data-action="like" data-type="post">
                <i class="fas fa-thumbs-up"></i>
                <span id="like_post-${post.ID}">${post.Likes}</span>
            </button>
            <button class="action-btn likedislikebtn ${post.UserInteraction === -1 ? 'liked-btn' : ''}"
                    data-id="${post.ID}" data-action="dislike" data-type="post">
                <i class="fas fa-thumbs-down"></i>
                <span id="dislike_post-${post.ID}">${post.Dislikes}</span>
            </button>
            <a class="action-btn" href="/Comment?post_id=${post.ID}">
                <i class="fas fa-comment"></i>${post.NbComment}
            </a>
        </div>
    `;
    return div;
}

const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            loadMorePosts()
        }
    })
}, { threshold: 1.0 })

observer.observe(loadingContainer)
