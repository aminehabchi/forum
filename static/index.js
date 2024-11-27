// I change this beacuse querySelectorAll gets All buttons at first on loading time 
// which we have only 3 posts when I added the more posts like and dislike action 
// do not work because I don't have them 
// this is a better approach ^_^
document.addEventListener('click', (e) => {
    const button = e.target.closest(".likedislikebtn")
    if (button) {
        const id = button.dataset.id;
        const action = button.dataset.action;
        const type = button.dataset.type

        likeDislike(id, action, type);
    }
});

const likeDislike = (CommentId, action, type) => {

    fetch(`/like-dislike?action=${action}&commentid=${CommentId}&type=${type}`)
        .then(res => {

            if (!res.ok) {
                window.location.replace(`/error?s=${res.status}&st=${res.statusText}`)
                return
            }
            if (res.url.includes('/login')) {
                window.location.replace('/login');
                return;
            }

            return res.json()

        })
        .then((data) => {

            document.querySelector(`#like_post-${CommentId}`).textContent = data.like
            document.querySelector(`#dislike_post-${CommentId}`).textContent = data.dislike
            document.querySelectorAll(`.likedislikebtn[data-id="${CommentId}"]`).forEach(button => {
                button.classList.remove('liked-btn');
                if (data.interaction === 1) {
                    document.querySelector(`.likedislikebtn[data-id="${CommentId}"][data-action="like"]`).classList.add('liked-btn');
                } else if (data.interaction === -1) {
                    document.querySelector(`.likedislikebtn[data-id="${CommentId}"][data-action="dislike"]`).classList.add('liked-btn');
                }
            });

        })
        .catch(err => {
            console.log(err);
        })
}


// ----- Load more Posts ------ //
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
